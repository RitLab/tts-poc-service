package validator

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	Validator *validator.Validate
}

func (x *Validator) Validate(i interface{}) error {
	return x.Validator.Struct(i)
}

func New() *Validator {
	v := &Validator{
		Validator: validator.New(),
	}

	v.Validator.RegisterValidation("alphanumspace", v.ValidateStringWithSpace)
	v.Validator.RegisterValidation("alphaspace", v.ValidateAlphaWithSpace)
	v.Validator.RegisterValidation("password", v.ValidatePassword)
	v.Validator.RegisterValidation("device-address", v.ValidateIPDevice)
	v.Validator.RegisterValidation("search-decode", v.QueryDecode)
	v.Validator.RegisterValidation("comma", v.ValidateComma)
	v.Validator.RegisterValidation("routingnumber", v.ValidateRoutingNumber)
	v.Validator.RegisterValidation("yyyy-mm-dd", v.FormatYYYYMMDD)
	v.Validator.RegisterValidation("default", v.DefaultValue)
	v.Validator.RegisterValidation("special-character", v.ValidateSpecialCharacter)

	// register tag json name
	v.Validator.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return v
}

func (x *Validator) ValidateStringWithSpace(fl validator.FieldLevel) bool {
	return regexp.MustCompile("^[a-zA-Z0-9 ]*$").MatchString(fl.Field().String())
}

func (x *Validator) ValidateAlphaWithSpace(fl validator.FieldLevel) bool {
	name := strings.TrimSpace(fl.Field().String())
	if len(name) == 0 {
		return false
	}
	return regexp.MustCompile("^[a-zA-Z ]*$").MatchString(strings.Trim(fl.Field().String(), " "))
}

func (x *Validator) ValidateRoutingNumber(fl validator.FieldLevel) bool {
	return regexp.MustCompile("^[0-9-]*$").MatchString(fl.Field().String())
}

func (x *Validator) DefaultValue(fl validator.FieldLevel) bool {
	if len(fl.Field().String()) == 0 {
		fl.Field().SetString(fl.Param())
	}

	return true
}

func (x *Validator) ValidatePassword(fl validator.FieldLevel) bool {
	if len(fl.Field().String()) < 8 {
		return false
	}
	done, err := regexp.MatchString("([a-z])+", fl.Field().String())
	if err != nil {
		return false
	}
	if !done {
		return false
	}
	done, err = regexp.MatchString("([A-Z])+", fl.Field().String())
	if err != nil {
		return false
	}
	if !done {
		return false
	}
	done, err = regexp.MatchString("([0-9])+", fl.Field().String())
	if err != nil {
		return false
	}
	if !done {
		return false
	}

	done, err = regexp.MatchString("([!@#$%^&*.?-])+", fl.Field().String())
	if err != nil {
		return false
	}
	if !done {
		return false
	}
	return true
}

func (x *Validator) ValidateIPDevice(fl validator.FieldLevel) bool {
	if strings.Contains(fl.Field().String(), "/") {
		ip := strings.Split(fl.Field().String(), "/")[0]
		if __validate_ip(ip) {
			url := strings.Replace(fl.Field().String(), ip, "", -1)
			if url[0] != '/' {
				return false
			}
			return __validate_url(url)
		}
	}

	return __validate_ip(fl.Field().String())
}

func (x *Validator) QueryDecode(fl validator.FieldLevel) bool {
	value, _ := url.QueryUnescape(fl.Field().String())
	fl.Field().SetString(value)
	return true
}

func (x *Validator) ValidateComma(fl validator.FieldLevel) bool {
	digit := strings.Split(fmt.Sprintf("%v", fl.Field().Float()), ".")
	if len(digit) == 1 {
		return true
	}

	param, _ := strconv.ParseFloat(fl.Param(), 64)
	return len(digit[1]) <= int(param)
}

func (x *Validator) FormatYYYYMMDD(fl validator.FieldLevel) bool {
	if len(fl.Field().String()) == 0 {
		return true
	}
	return regexp.MustCompile(`^\d{4}\-(0[1-9]|1[012])\-(0[1-9]|[12][0-9]|3[01])$`).MatchString(fl.Field().String())
}

func (x *Validator) ValidateSpecialCharacter(fl validator.FieldLevel) bool {
	if len(fl.Field().String()) == 0 {
		return true
	}
	return regexp.MustCompile(`^[^!?@#$%&*(),.]*$`).MatchString(fl.Field().String())
}

func __validate_ip(ip string) bool {
	return regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]):[0-9]+$`).MatchString(ip)
}

func __validate_url(url string) bool {
	if regexp.MustCompile(`.*//.*`).MatchString(url) {
		return false
	}
	return regexp.MustCompile("^[a-zA-Z0-9/_?-]*$").MatchString(url)
}
