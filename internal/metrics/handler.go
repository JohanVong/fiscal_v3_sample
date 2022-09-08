package metrics

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/JohanVong/fiscal_v3_sample/internal/db"
	"github.com/JohanVong/fiscal_v3_sample/internal/models"
	"github.com/JohanVong/fiscal_v3_sample/internal/services/calculations"
)

var Stage = os.Getenv("STAGE")

const (
	OFDStatusOK int = iota
	OFDStatusUnknownID
	OFDStatusInvalidToken
	OFDStatusProtocolError
	OFDStatusUnknownCommand
	OFDStatusUnsupportedCommand
	OFDStatusInvalidConfiguration
	OFDStatusSSLIsNotAllowed
	OFDStatusInvalidRequestNumber
	OFDStatusInvalidRetryRequest
	OFDStatusCantCancelTicket
	OFDStatusOpenShiftTimeout
	OFDStatusInvalidLoginPassword
	OFDStatusIncorrectRequestData
	OFDStatusNotEnoughCash
	OFDStatusBlocked
	OFDStatusServiceTemporarilyUnavailable = 254
	OFDStatusUnknownError                  = 255
	OFDStatusTimeout                       = 500
)

var (
	ofdStatusKT = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "ofd_status_kt",
			Help:      "Статус ОФД Казактелеком",
		})

	averageOperationTime = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "oper_avrg_time",
			Help:      "Среднее время операции",
		})

	ofdStatusTK = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "ofd_status_tk",
			Help:      "Статус ОФД Транстелеком",
		})

	kkmActiveNumber = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "kkm_active_number",
			Help:      "Количество активных ккм",
		})

	kkmInactiveNumber = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "kkm_inactive_number",
			Help:      "Количество заблокированных ккм",
		})

	kkmOfflineModeNumber = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "kkm_offline_mode_number",
			Help:      "Количество ккм в автономном режиме",
		})

	kkmOfdStatusOK = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "kkm_ofd_status_ok",
			Help:      "Количество ккм с успешными операциями в ОФД",
		})

	kkmOfdStatusUnknownID = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "kkm_ofd_status_unknown_id",
			Help:      "Количество ккм с ошибкой от ОФД 'неизвестный ID устройства'",
		})

	kkmOfdStatusInvalidToken = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "kkm_ofd_status_invalid_token",
			Help:      "Количество ккм с ошибкой от ОФД 'неверный токен'",
		})

	kkmOfdStatusProtocolError = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "kkm_ofd_status_protocol_error",
			Help:      "Количество ккм с ошибкой от ОФД 'ошибка протокола'",
		})

	kkmOfdStatusUnknownCommand = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "kkm_ofd_status_unknown_command",
			Help:      "Количество ккм с ошибкой от ОФД 'неизвестная команда'",
		})

	kkmOfdStatusUnsupportedCommand = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "kkm_ofd_status_unsupported_command",
			Help:      "Количество ккм с ошибкой от ОФД 'команда не поддерживается'",
		})

	kkmOfdStatusInvalidConfiguration = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "kkm_ofd_status_invalid_configuration",
			Help:      "Количество ккм с ошибкой от ОФД 'неверные натройки устройства'",
		})

	kkmOfdStatusSSLIsNotAllowed = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "kkm_ofd_status_ssl_is_not_allowed",
			Help:      "Количество ккм с ошибкой от ОФД 'использование SSL не разрешено'",
		})

	kkmOfdStatusInvalidRequestNumber = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "kkm_ofd_status_invalid_request_number",
			Help:      "Количество ккм с ошибкой от ОФД 'неправильный номер запроса'",
		})

	kkmOfdStatusInvalidRetryRequest = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "kkm_ofd_status_invalid_retry_request",
			Help:      "Количество ккм с ошибкой от ОФД 'неправильная попытка отправки повторного запроса'",
		})

	kkmOfdStatusOpenShiftTimeout = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "kkm_ofd_status_open_shift_timeout",
			Help:      "Количество ккм с ошибкой от ОФД 'время открытой смены истекло'",
		})

	kkmOfdStatusIncorrectRequestData = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "kkm_ofd_status_incorrect_request_data",
			Help:      "Количество ккм с ошибкой от ОФД 'неверные входные данные'",
		})

	kkmOfdStatusNotEnoughCash = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "kkm_ofd_status_not_enough_cash",
			Help:      "Количество ккм с ошибкой от ОФД 'недостточно наличных'",
		})

	kkmOfdStatusBlocked = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "kkm_ofd_status_blocked",
			Help:      "Количество ккм с ошибкой от ОФД 'касса заблокирована'",
		})

	kkmOfdStatusServiceTemporarilyUnavailable = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "kkm_ofd_status_service_temporarily_unavailable",
			Help:      "Количество ккм с ошибкой от ОФД 'сервис временно недоступен'",
		})

	kkmOfdStatusUnknownError = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "kkm_ofd_status_unknown_error",
			Help:      "Количество ккм с ошибкой от ОФД 'неизвесттная ошибка'",
		})

	kkmOfdStatusTimeout = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "kkm_ofd_status_timeout",
			Help:      "Количество ккм, соединение с которым прервалось по таймауту",
		})

	documentOfflineModeNumber = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "document_offline_mode_number",
			Help:      "Количество операций в автономном режиме",
		})

	operationsPerMinuteTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "operations_per_minute_total",
			Help:      "Общее количество операций в ОФД в минуту",
		})

	operationsPerMinuteSuccess = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "operations_per_minute_success",
			Help:      "Количество успешных операций в ОФД в минуту",
		})

	operationsPerMinuteError = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "fiscal2",
			Name:      "operations_per_minute_error",
			Help:      "Количество операций в ОФД с ошибкой в минуту",
		})
)

var perKKMStatus map[int]int

type OperationOFDSuccessRate struct {
	TotalSent          int64
	SuccessfulResponse int64
	ErroneousResponse  int64
}

var MetricMutex sync.Mutex
var OFDStatusMutex sync.Mutex

var OperationOFDSuccessMetric OperationOFDSuccessRate

func RegisterMetrics() (err error) {

	//defer recoverpanic.RecoverPanic()
	err = prometheus.Register(averageOperationTime)
	err = prometheus.Register(ofdStatusKT)
	err = prometheus.Register(ofdStatusTK)
	err = prometheus.Register(kkmActiveNumber)
	err = prometheus.Register(kkmInactiveNumber)
	err = prometheus.Register(kkmOfflineModeNumber)
	err = prometheus.Register(documentOfflineModeNumber)
	err = prometheus.Register(kkmOfdStatusOK)
	err = prometheus.Register(kkmOfdStatusBlocked)
	err = prometheus.Register(kkmOfdStatusIncorrectRequestData)
	err = prometheus.Register(kkmOfdStatusInvalidConfiguration)
	err = prometheus.Register(kkmOfdStatusInvalidRequestNumber)
	err = prometheus.Register(kkmOfdStatusUnsupportedCommand)
	err = prometheus.Register(kkmOfdStatusUnknownID)
	err = prometheus.Register(kkmOfdStatusUnknownError)
	err = prometheus.Register(kkmOfdStatusUnknownCommand)
	err = prometheus.Register(kkmOfdStatusTimeout)
	err = prometheus.Register(kkmOfdStatusSSLIsNotAllowed)
	err = prometheus.Register(kkmOfdStatusServiceTemporarilyUnavailable)
	err = prometheus.Register(kkmOfdStatusProtocolError)
	err = prometheus.Register(kkmOfdStatusOpenShiftTimeout)
	err = prometheus.Register(kkmOfdStatusNotEnoughCash)
	err = prometheus.Register(kkmOfdStatusInvalidToken)
	err = prometheus.Register(kkmOfdStatusInvalidRetryRequest)
	err = prometheus.Register(operationsPerMinuteTotal)
	err = prometheus.Register(operationsPerMinuteError)
	err = prometheus.Register(operationsPerMinuteSuccess)
	var kkmActiveCount float64 = 0
	var kkmInactiveCount float64 = 0
	var kkmOfflineCount float64 = 0
	var documentOfflineCount float64 = 0
	statusTicker := time.NewTicker(time.Second * 60)
	ofdTicker := time.NewTicker(time.Second * 60)

	go func() {
		tx := db.Orm.Begin()
		tx.Exec("SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED")
		for range statusTicker.C {
			tx.Model(&models.Kkm{}).Where("idStatusKKM=?", 2).Count(&kkmActiveCount)
			tx.Model(&models.Kkm{}).Where("idStatusKKM=?", 3).Count(&kkmInactiveCount)
			tx.Model(&models.Kkm{}).Where("OfflineQueue>0").Count(&kkmOfflineCount)
			tx.Model(&models.Document{}).Where("Offline=?", true).Count(&documentOfflineCount)
			kkmActiveNumber.Set(kkmActiveCount)
			averageOperationTime.Set(float64(calculations.GetAvrgOpertime()))
			kkmInactiveNumber.Set(kkmInactiveCount)
			kkmOfflineModeNumber.Set(kkmOfflineCount)
			documentOfflineModeNumber.Set(documentOfflineCount)
		}
		tx.RollbackUnlessCommitted()
	}()

	go func() {
		for range ofdTicker.C {
			MetricMutex.Lock()
			operationsPerMinuteTotal.Set(float64(OperationOFDSuccessMetric.TotalSent))
			operationsPerMinuteSuccess.Set(float64(OperationOFDSuccessMetric.SuccessfulResponse))
			operationsPerMinuteError.Set(float64(OperationOFDSuccessMetric.ErroneousResponse))
			OperationOFDSuccessMetric.TotalSent = 0
			OperationOFDSuccessMetric.ErroneousResponse = 0
			OperationOFDSuccessMetric.SuccessfulResponse = 0
			MetricMutex.Unlock()
			OFDStatusMutex.Lock()
			st := GetKkmOfdStatusMetric()
			for k, v := range st {
				switch k {
				case OFDStatusOK:
					kkmOfdStatusOK.Set(float64(v))
				case OFDStatusBlocked:
					kkmOfdStatusBlocked.Set(float64(v))
				case OFDStatusIncorrectRequestData:
					kkmOfdStatusIncorrectRequestData.Set(float64(v))
				case OFDStatusInvalidConfiguration:
					kkmOfdStatusInvalidConfiguration.Set(float64(v))
				case OFDStatusInvalidRequestNumber:
					kkmOfdStatusInvalidRequestNumber.Set(float64(v))
				case OFDStatusInvalidRetryRequest:
					kkmOfdStatusInvalidRetryRequest.Set(float64(v))
				case OFDStatusInvalidToken:
					kkmOfdStatusInvalidToken.Set(float64(v))
				case OFDStatusNotEnoughCash:
					kkmOfdStatusNotEnoughCash.Set(float64(v))
				case OFDStatusOpenShiftTimeout:
					kkmOfdStatusOpenShiftTimeout.Set(float64(v))
				case OFDStatusProtocolError:
					kkmOfdStatusProtocolError.Set(float64(v))
				case OFDStatusServiceTemporarilyUnavailable:
					kkmOfdStatusServiceTemporarilyUnavailable.Set(float64(v))
				case OFDStatusSSLIsNotAllowed:
					kkmOfdStatusSSLIsNotAllowed.Set(float64(v))
				case OFDStatusTimeout:
					kkmOfdStatusTimeout.Set(float64(v))
				case OFDStatusUnknownCommand:
					kkmOfdStatusUnknownCommand.Set(float64(v))
				case OFDStatusUnknownError:
					kkmOfdStatusUnknownError.Set(float64(v))
				case OFDStatusUnsupportedCommand:
					kkmOfdStatusUnsupportedCommand.Set(float64(v))
				case OFDStatusUnknownID:
					kkmOfdStatusUnknownID.Set(float64(v))
				}
			}
			OFDStatusMutex.Unlock()
		}
	}()

	return err
}

func SetKkmOfdStatusMetric(idKKM int, status int) {
	var ttl time.Duration
	if status == 0 {
		ttl = time.Hour * 24
	} else {
		ttl = time.Minute * 15
	}
	key := fmt.Sprintf("fiscalv2:%s:ofd_status:kkm:%d", Stage, idKKM)
	db.RedisCl.Set(key, status, ttl)
	//log.Println(models.Orm.Debug().Model(&models.Kkm{}).Where("idKKM=?", idKKM).Updates(map[string]interface{}{"ofd_code": status}).RowsAffected)
}

func getKkmOfdStatusMetricKeys() (keys []string, err error) {
	var curs uint64 = 0
	var v []string
	for {
		key := fmt.Sprintf("fiscalv2:%s:ofd_status*", Stage)
		v, curs, err = db.RedisCl.Scan(curs, key, 1000).Result()
		if err != nil {
			return nil, err
		} else {
			if len(v) > 0 {
				keys = append(keys, v...)
			}
		}
		if curs == 0 {
			break
		}
	}

	return
}

func GetKkmOfdStatusMetric() (statuses map[int]int) {
	keys, _ := getKkmOfdStatusMetricKeys()
	val, _ := db.RedisCl.MGet(keys...).Result()
	perKKMStatus = make(map[int]int)

	for _, s := range val {

		v, ok := s.(string)
		if ok == false {
			continue
		}
		st, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		MetricMutex.Lock()
		perKKMStatus[st] = perKKMStatus[st] + 1
		MetricMutex.Unlock()
	}

	return perKKMStatus
}
