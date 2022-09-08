package validator

import (
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/ericlagergren/decimal"
	"gopkg.in/asaskevich/govalidator.v9"

	"github.com/JohanVong/fiscal_v3_sample/internal/models"
)

func InitCustomValidator() {
	govalidator.CustomTypeTagMap.Set("isPhoneNumber", govalidator.CustomTypeValidator(func(i interface{}, context interface{}) bool {
		re, _ := regexp.Compile(`^\d{11}$`)
		phoneNumber, ok := i.(string)
		if ok == false {
			return false
		} else {
			return re.Match([]byte(phoneNumber))
		}
	}))

	govalidator.CustomTypeTagMap.Set("twoDecimalPlaces", govalidator.CustomTypeValidator(func(i interface{},
		context interface{}) bool {
		switch v := i.(type) {
		case *decimal.Big:
			if v.Scale() <= 2 {
				return true
			}
		case decimal.Big:
			if v.Scale() <= 2 {
				return true
			}
		case *models.Decimal:
			if v.Scale() <= 2 {
				return true
			}
		case models.Decimal:
			if v.Scale() <= 2 {
				return true
			}
		case string:
			d, _ := decimal.New(0, 0).SetString(v)
			if d.IsFinite() && d.Sign() > 0 && d.Scale() <= 2 {
				return true
			} else {
				return false
			}
		}
		return false
	}))

	govalidator.CustomTypeTagMap.Set("isEscapedDateTime", govalidator.CustomTypeValidator(func(i interface{},
		context interface{}) bool {
		v := strings.TrimSpace(i.(string))
		if len(v) > 0 {
			vEscaped, err := url.QueryUnescape(v)
			if err != nil {
				return false
			}
			_, err = time.Parse("2006-01-02T15:04:05", vEscaped)
			if err != nil {
				return false
			}
		}

		return true
	}))

	govalidator.CustomTypeTagMap.Set("isEscapedDateTimeNoSeconds", govalidator.CustomTypeValidator(func(i interface{},
		context interface{}) bool {
		v := strings.TrimSpace(i.(string))
		if len(v) > 0 {
			vEscaped, err := url.QueryUnescape(v)
			if err != nil {
				return false
			}
			_, err = time.Parse("2006-01-02T15:04", vEscaped)
			if err != nil {
				return false
			}
		}

		return true
	}))

	govalidator.CustomTypeTagMap.Set("isNonNegative", govalidator.CustomTypeValidator(func(i interface{},
		context interface{}) bool {
		switch v := i.(type) {
		case *decimal.Big:
			return !v.Signbit()
		case decimal.Big:
			return !v.Signbit()
		case *models.Decimal:
			return !v.Signbit()
		case models.Decimal:
			return !v.Signbit()
		case int:
			return v >= 0
		case int64:
			return v >= 0
		}
		return false
	}))
}
