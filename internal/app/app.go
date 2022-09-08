package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	"github.com/JohanVong/fiscal_v3_sample/configs"
	"github.com/JohanVong/fiscal_v3_sample/internal/db"
	"github.com/JohanVong/fiscal_v3_sample/internal/metrics"
	"github.com/JohanVong/fiscal_v3_sample/internal/migrations"
	"github.com/JohanVong/fiscal_v3_sample/internal/services/calculations"
	"github.com/JohanVong/fiscal_v3_sample/internal/services/graylog"
	"github.com/JohanVong/fiscal_v3_sample/internal/services/ofd"
	"github.com/JohanVong/fiscal_v3_sample/internal/services/systemconfig"
	"github.com/JohanVong/fiscal_v3_sample/internal/services/validator"
	"github.com/JohanVong/fiscal_v3_sample/internal/transport/protobuf/ofd/cpcr"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/panicwrap"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var _mssqlDbUser = os.Getenv("FISCAL_DB_USER")
var _mssqlDbPassword = os.Getenv("FISCAL_DB_PASSWORD")
var _mssqlDbHost = os.Getenv("FISCAL_DB_HOST")
var _mssqlDbPort = os.Getenv("FISCAL_DB_PORT")
var _mssqlDbName = os.Getenv("FISCAL_DB_NAME")
var _ktOfdAddress = os.Getenv("KT_OFD_ADDRESS")
var _tkOfdAddress = os.Getenv("TK_OFD_ADDRESS")
var _queueInstance = os.Getenv("QUEUE_INSTANCE")

var _redisLocal = "localhost:6379"
var _redisProd = os.Getenv("REDIS_PROD")

func Run(configPath string) {
	runtime.GOMAXPROCS(4)

	validator.InitCustomValidator()
	var err error
	if os.Getenv("PANICWRAP_ENABLED") == "0" {
		exitStatus, err := panicwrap.BasicWrap(panicHandler)
		if err != nil {

			panic(err)
		}
		if exitStatus >= 0 {
			os.Exit(exitStatus)
		}
	}

	if err != nil {
		panic(err)
	}
	// logs.Register("graylog", graylog.NewConsole)
	// logs.SetLogger("graylog", ``)

	mockOFD := cpcr.MockOFD{
		Down:         false,
		NoConnection: false,
		MinDelay:     275 * time.Millisecond,
		MaxDelay:     380 * time.Millisecond,
	}

	ktOFD, err := cpcr.NewOFDHandler(_ktOfdAddress)
	if err != nil {
		panic(fmt.Sprint(err, string(debug.Stack())))
	}
	tkOFD, err := cpcr.NewOFDHandler(_tkOfdAddress)
	if err != nil {
		panic(fmt.Sprint(err, string(debug.Stack())))
	}

	ofd.OfdPool = make(ofd.Pool)

	//ofd.OfdPool[1] = ktTestOFD
	ofd.OfdPool[2] = &mockOFD
	ofd.OfdPool[3] = ktOFD
	ofd.OfdPool[4] = tkOFD
	query := url.Values{}
	query.Add("database", _mssqlDbName)

	dbUrl := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(_mssqlDbUser, _mssqlDbPassword),
		Host:     fmt.Sprintf("%s:%s", _mssqlDbHost, _mssqlDbPort),
		RawQuery: query.Encode(),
	}

	db.Orm, err = gorm.Open("mssql", dbUrl.String())

	db.Orm.SingularTable(true)
	db.Orm.LogMode(false)

	go systemconfig.UpdateSettings()
	if err != nil {
		panic(err)
	}

	graylog.GraylogClient, err = graylog.NewClient(configs.Graylog)
	if err != nil {
		panic(err)
	}

	db.RedisCl = redis.NewClient(&redis.Options{
		Addr:     configs.Redis,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, redis_err := db.RedisCl.Ping().Result()
	if redis_err != nil {
		panic(redis_err)
	}

	// middleware.InitEnforcer()
	// notifier.Notifier, err = notifier.NewBitrixNotifier(beego.AppConfig.String("bitrix_notifier"),
	// 	"Уведомления Фискал 2.0")
	// if err != nil {
	// 	m := graylog.NewGELFMessage("ErrorBitrix")
	// 	*m.SendDate = time.Now()
	// 	m.SendMessage = err.Error()
	// 	_, _ = graylog.SendMessage(m)
	// }
	// operpool.Operations = operpool.NewOperPool()
	// operpool.Operations.Start()
	// models.KKMOFDStatusStorage = make(map[int]*models.KKMOFDStatus)

	if _queueInstance == "1" {
		run := make(chan bool)
		var wg sync.WaitGroup
		go ofd.ProcessOfflineQueue(&wg, time.Second*180)
		wg.Wait()
		<-run
	} else if _queueInstance == "0" {
		err = migrations.LaunchMigration()
		if err != nil {
			panic(err)
		}
		log.Println(metrics.RegisterMetrics())
		StartMetricsServer()
		go calculations.AddOperTime()
		// beego.Run()
	} else {
		var wg sync.WaitGroup
		go ofd.ProcessOfflineQueue(&wg, time.Second*10)
		wg.Wait()
		err = migrations.LaunchMigration()
		if err != nil {
			panic(err)
		}
		log.Println(metrics.RegisterMetrics())
		StartMetricsServer()
		go calculations.AddOperTime()
		// beego.Run()
	}
}

func StartMetricsServer() {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	go r.Run(":3000")
}

func panicHandler(output string) {
	m := graylog.NewGELFMessage("CrashLog")
	*m.SendDate = time.Now()
	if len(output) > 30000 {
		m.SendMessage = output[:30000]
	} else {
		m.SendMessage = output
	}

	payload, _ := json.Marshal(m)

	reader := bytes.NewReader(payload)
	cl := http.DefaultClient
	_, _ = cl.Post("http://192.168.151.110:12222/gelf", "application/json", reader)

	os.Exit(1)
}
