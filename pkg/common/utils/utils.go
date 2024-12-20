package utils

import (
	"crypto/rand"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	mathRand "math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	pkgError "tts-poc-service/pkg/common/error"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/nacl/box"
)

//go:embed *
var files embed.FS

// Atoi This function is used for convert string to int
func Atoi(s string) int {
	d, _ := strconv.Atoi(s)
	return d
}

// ReadReader function for read body response client
func ReadReader(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

// Errorf wrap string as error
func Errorf(msg string) error {
	return fmt.Errorf(msg)
}

func Pointer[T any](object T) *T {
	return &object
}

// Now is a function to get time.Now with format 2006-01-02 15:04:05.000
func Now() string {
	return time.Now().Format("2006-01-02 15:04:05.000")
}

func TimeParse(date string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", date)
}

func TodayUnix() int64 {
	today, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	return today.Unix()
}

func NowWithAdd(duration time.Duration) string {
	return time.Now().Add(duration).Format("2006-01-02 15:04:05.000")
}

// IfThenElse is a function for ternary operator
func IfThenElse[T any](cond bool, this, that T) T {
	if cond {
		return this
	}
	return that
}

// Unmarshal is a function to casting object
func Unmarshal(source, target any) error {
	bytes, _ := json.Marshal(source)
	return json.Unmarshal(bytes, &target)
}

// BindRequestAndValidate is a function to validate request body, header and param.
// This function use go-playground and echo.DefaultBinder
func BindRequestAndValidate(c echo.Context, request any) error {
	if err := c.Bind(request); err != nil {
		return Errorf(pkgError.BAD_REQUEST)
	}

	if err := (&echo.DefaultBinder{}).BindHeaders(c, request); err != nil {
		return err
	}

	if err := (&echo.DefaultBinder{}).BindPathParams(c, request); err != nil {
		return err
	}

	if err := c.Validate(request); err != nil {
		validationError := err.(validator.ValidationErrors)
		errorString := ""
		for _, v := range validationError {
			valid := ""
			switch {
			case v.Tag() == "required":
				valid = " is required"
			case v.Tag() == "required_if":
				valid = fmt.Sprintf(" is required if %s", v.Param())
			case v.Tag() == "number":
				valid = " should only contain number"
			case v.Tag() == "email":
				valid = " should be email format"
			case v.Tag() == "special-character":
				valid = " should not contain special character"
			case v.Tag() == "oneof":
				valid = fmt.Sprintf(" should only contain one of %s", v.Param())
			case v.Tag() == "len":
				valid = fmt.Sprintf(" should be %s digit", v.Param())
			case v.Tag() == "max":
				valid = fmt.Sprintf(" should have maximum %s digit", v.Param())
			case v.Tag() == "min":
				valid = fmt.Sprintf(" should have minimum %s digit", v.Param())
			case v.Tag() == "password":
				valid = fmt.Sprintf(" should at least 8 character, 1 lower case, 1 upper case, 1 symbol and 1 number")
			}
			if errorString == "" {
				errorString += v.Field() + valid
			} else {
				errorString += ";" + v.Field() + valid
			}
		}
		return fmt.Errorf(errorString)
	}
	return nil
}

// SliceContain is a function to check slice contain key or not
func SliceContain[T any](slice []T, key T) bool {
	for _, o := range slice {
		if any(o) == any(key) {
			return true
		}
	}
	return false
}

func UniqSliceString(slice []string) []string {
	keys := make(map[string]bool)
	var result []string

	for _, item := range slice {
		if !keys[item] {
			result = append(result, item)
			keys[item] = true
		}
	}
	return result
}

func OpenCustomKey(chipertext []byte, pub, priv *[32]byte) (plaintext []byte, ok bool) {
	rest, nonce := ExtractNonce(chipertext)
	return box.Open(nil, rest, &nonce, pub, priv)
}

func SealCustomKey(plain []byte, pub, priv *[32]byte) (chipertext []byte) {
	nonce := GenerateNonce()
	return box.Seal(nonce[:], plain, &nonce, pub, priv)
}

func GenerateNonce() [24]byte {
	var nonce [24]byte
	_, _ = io.ReadFull(rand.Reader, nonce[:])
	return nonce
}

func ExtractNonce(ciphertext []byte) ([]byte, [24]byte) {
	var nonce [24]byte

	_, rest := copy(nonce[:], ciphertext[:24]), ciphertext[24:]

	return rest, nonce
}

func GetDiffDateYYYYMMDD(start, end, units string) int64 {
	first, _ := time.Parse("2006-01-02", start)
	second, _ := time.Parse("2006-01-02", end)

	diff := second.Sub(first)

	if units == "year" {
		return int64(diff.Hours() / 24 / 365)
	} else if units == "hour" {
		return int64(diff.Hours())
	} else if units == "day" {
		return int64(diff.Hours() / 24)
	}

	return int64(diff.Hours())
}
func GenerateRandomDigit(n int) string {
	mathRand.Seed(time.Now().UnixNano())
	id := ""
	for i := 0; i < n; i++ {
		id += strconv.Itoa(mathRand.Intn(10))
	}

	return id
}

func GenerateFile(dir, name string) (*os.File, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0777); err != nil {
			return nil, err
		}
	}

	filename := dir + name
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func ObjectToString(in interface{}) string {
	out, err := json.Marshal(in)
	if err != nil {
		return ""
	}
	return string(out)
}

func RandInt(min, max int) int {
	return min + mathRand.Intn(max-min)
}

func SpaceMap(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func IsNumeric(word string) bool {
	return regexp.MustCompile(`\d`).MatchString(word)
}

func ValidateDates(dateFrom, dateTo string) error {
	const dateFormat = "2006-01-02T15:04:05Z"

	// Parse date_from
	from, err := time.Parse(dateFormat, dateFrom)
	if err != nil {
		return fmt.Errorf(pkgError.INVALID_DATE_FORMAT)
	}

	// Parse date_to
	to, err := time.Parse(dateFormat, dateTo)
	if err != nil {
		return fmt.Errorf(pkgError.INVALID_DATE_FORMAT)
	}

	// Check if date_from is less than or equal to date_to
	if !from.Before(to) && !from.Equal(to) {
		return fmt.Errorf(pkgError.INVALID_DATE_FROM_VALUE)
	}

	// Check if date_to is greater than or equal to date_from
	if !to.After(from) && !to.Equal(from) {
		return fmt.Errorf(pkgError.INVALID_DATE_TO_VALUE)
	}

	return nil
}

func ValidateInput(input string, allowedList []string) bool {
	if input == "" {
		return true
	}
	input = strings.ToUpper(input)
	for _, validValue := range allowedList {
		if input == validValue {
			return true
		}
	}
	return false
}
