package systemconfig

import (
	"strconv"
	"sync"
	"time"

	"github.com/JohanVong/fiscal_v3_sample/internal/db"
	"github.com/JohanVong/fiscal_v3_sample/internal/models"
)

const DefaultConfigReadInterval = 60

type systemConfig struct {
	options map[string]string
	mutex   sync.Mutex
}

func (sc *systemConfig) GetOption(option string) (value string, ok bool) {
	sc.mutex.Lock()
	value, ok = sc.options[option]
	sc.mutex.Unlock()
	return value, ok
}

func (sc *systemConfig) Fetch() {
	sc.mutex.Lock()
	if db.Orm != nil {
		if sc.options == nil {
			sc.options = make(map[string]string)
		}
		var options []models.SystemConfig
		db.Orm.Model(&models.SystemConfig{}).Find(&options)
		for _, o := range options {
			sc.options[o.OptionName] = o.OptionValue
		}
	}
	sc.mutex.Unlock()
}

var SystemConfig systemConfig

func UpdateSettings() {
	var readInterval time.Duration
	for {
		SystemConfig.Fetch()
		v, ok := SystemConfig.GetOption("config_read_interval")
		if ok == true {
			intervalValue, err := strconv.ParseInt(v, 10, 64)
			if err == nil && intervalValue > DefaultConfigReadInterval {
				readInterval = time.Duration(intervalValue) * time.Second
			} else {
				readInterval = time.Second * DefaultConfigReadInterval
			}
		} else {
			readInterval = time.Second * DefaultConfigReadInterval
		}
		time.Sleep(readInterval)

	}
}
