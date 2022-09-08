package cpcr

import (
	"fmt"
	"strings"

	"github.com/JohanVong/fiscal_v3_sample/internal/transport/protobuf/ofd/kkm_proto"
)

var OfdErrorText = map[kkm_proto.ResultTypeEnum]string{
	kkm_proto.ResultTypeEnum_RESULT_TYPE_UNKNOWN_ID:                      "Аппарат не зарегистрирован в системе",
	kkm_proto.ResultTypeEnum_RESULT_TYPE_INVALID_TOKEN:                   "Отправка данных невозможна, необходимо произвести сброс токена",
	kkm_proto.ResultTypeEnum_RESULT_TYPE_PROTOCOL_ERROR:                  "Ошибка протокола, обратитесь в сервисную службу",
	kkm_proto.ResultTypeEnum_RESULT_TYPE_UNKNOWN_COMMAND:                 "Ошибка протокола, обратитесь в сервисную службу",
	kkm_proto.ResultTypeEnum_RESULT_TYPE_UNSUPPORTED_COMMAND:             "Данная команда не поддерживается сервером, обратитесь в сервисную службу",
	kkm_proto.ResultTypeEnum_RESULT_TYPE_INVALID_CONFIGURATION:           "",
	kkm_proto.ResultTypeEnum_RESULT_TYPE_SSL_IS_NOT_ALLOWED:              "Использование защищенного соединиения запрещено. Подключите услугу или используйте открытый канал связи",
	kkm_proto.ResultTypeEnum_RESULT_TYPE_INVALID_REQUEST_NUMBER:          "Порядковый номер запроса REQNUM тот же, что и в предыдущем запросе, но токен другой",
	kkm_proto.ResultTypeEnum_RESULT_TYPE_INVALID_RETRY_REQUEST:           "REQNUM и TOKEN имеют те же значения, что и в предыдущем запросе, но код команды отличается",
	kkm_proto.ResultTypeEnum_RESULT_TYPE_CANT_CANCEL_TICKET:              "Можно отменить только последний чек, и при этом после отменяемого  чека не было отправлено ни одной команды, кроме служебной",
	kkm_proto.ResultTypeEnum_RESULT_TYPE_OPEN_SHIFT_TIMEOUT_EXPIRED:      "Истек период, в течение которого смена может быть открыта",
	kkm_proto.ResultTypeEnum_RESULT_TYPE_INVALID_LOGIN_PASSWORD:          "Неправильный логин или пароль",
	kkm_proto.ResultTypeEnum_RESULT_TYPE_INCORRECT_REQUEST_DATA:          "Данные корректны с точки зрения протокола, но неверны в конкретном контексте",
	kkm_proto.ResultTypeEnum_RESULT_TYPE_NOT_ENOUGH_CASH:                 "Недостаточно наличных в кассе",
	kkm_proto.ResultTypeEnum_RESULT_TYPE_BLOCKED:                         "Касса заблокирована",
	kkm_proto.ResultTypeEnum_RESULT_TYPE_SAME_TAXPAYER_AND_CUSTOMER:      "Совпадает значение ИИН/БИН покупателя и продавца",
	kkm_proto.ResultTypeEnum_RESULT_TYPE_SERVICE_TEMPORARILY_UNAVAILABLE: "Сервис временно недоступен",
	kkm_proto.ResultTypeEnum_RESULT_TYPE_UNKNOWN_ERROR:                   "Неизвестная ошибка",
}

type OfdError struct {
	Code      uint32
	text      string
	extraText string
}

func NewOfdError(code uint32, extraText string) *OfdError {
	//defer recoverpanic.RecoverPanic()
	e := new(OfdError)
	e.Code = code
	e.text = OfdErrorText[kkm_proto.ResultTypeEnum(code)]
	e.extraText = strings.TrimSpace(extraText)
	return e
}

func (e *OfdError) Error() string {
	if e.extraText != "" {
		idInfo := getIdInfo(e.extraText)
		if idInfo != "" {
			return fmt.Sprintf("ofd result error: Code: %d: %s: %s", e.Code, e.text, idInfo)
		}
	}
	return fmt.Sprintf("ofd result error: Code: %d: %s", e.Code, e.text)
}

func getIdInfo(substr string) string {
	var idInfo string

	if strings.ContainsAny(substr, "{}") && strings.Contains(substr, ",") {
		arr := strings.Split(substr, ",")
		for _, val := range arr {
			if !strings.Contains(val, "id:") {
				continue
			}

			val = strings.TrimSpace(val)
			idInfo = val
			break
		}
	}

	return idInfo
}
