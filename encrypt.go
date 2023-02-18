package datatype

import (
	"database/sql/driver"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type EncryptService interface {
	Encode(value string) string
	Decode(value string) string
	Mask(value string) string
}

// 默认加密服务
type DefaultEncryptService struct{}

func (es DefaultEncryptService) Encode(value string) string {
	var data []string

	for _, pair := range RegexpSplit(`[0-9]{8,20}`, strings.ToUpper(value)) {
		if pair.Match {
			size := len(pair.Value)
			ch := pair.Value[size-1:]

			pair.Value = pair.Value[:size-1]

			offset, _ := strconv.Atoi(ch)

			if offset == 0 {
				offset = 10
			}

			if offset%2 == 1 {
				pair.Value = ReverseString(pair.Value)
			}

			v := big.NewInt(0)
			v.SetString(pair.Value, 10)

			pair.Value = ""

			for _, c := range v.Text(16) {
				i, _ := strconv.ParseInt(string(c), 16, 10)
				pair.Value += fmt.Sprintf("%x", (i+int64(offset))%16)
			}

			pair.Value += ch
			pair.Value += "!"

			_size := len(pair.Value)

			if _size < size {
				pair.Value = pair.Value[:_size/2] + strings.Repeat("*", size-_size) + pair.Value[_size/2:]
			}
		}

		data = append(data, pair.Value)
	}

	return strings.Join(data, "")
}

func (es DefaultEncryptService) Decode(value string) string {
	var data []string

	for _, pair := range RegexpSplit(`[0-9a-f*]{6,18}[0-9]!`, value) {
		if pair.Match {
			size := len(pair.Value) - 1
			pair.Value = strings.ReplaceAll(strings.ReplaceAll(pair.Value, "*", ""), "!", "")

			_size := len(pair.Value)
			ch := pair.Value[_size-1:]
			pair.Value = pair.Value[:_size-1]

			offset, _ := strconv.Atoi(ch)

			if offset == 0 {
				offset = 10
			}

			value := pair.Value
			pair.Value = ""

			for _, c := range value {
				i, _ := strconv.ParseInt(string(c), 16, 10)
				pair.Value += fmt.Sprintf("%x", (i-int64(offset)+16)%16)
			}

			v := big.NewInt(0)
			v.SetString(pair.Value, 16)

			pair.Value = v.String()

			if len(pair.Value) < size {
				pair.Value = strings.Repeat("0", size-len(pair.Value)) + pair.Value
			}

			if offset%2 == 1 {
				pair.Value = ReverseString(pair.Value)
			}

			pair.Value += ch
		}

		data = append(data, pair.Value)
	}

	return strings.Join(data, "")
}

func (es DefaultEncryptService) Mask(value string) string {
	var data []string

	for _, pair := range RegexpSplit(`[0-9]{8,20}`, value) {
		if pair.Match {
			size := len(pair.Value)

			if size >= 10 {
				pair.Value = pair.Value[:4] + strings.Repeat("*", size-8) + pair.Value[size-4:]
			} else {
				pair.Value = pair.Value[:3] + strings.Repeat("*", size-6) + pair.Value[size-3:]
			}
		}

		data = append(data, pair.Value)
	}

	return strings.Join(data, "")
}

type EncryptConfig struct {
	// 加密服务
	Service EncryptService
}

var EncryptOptions = EncryptConfig{
	Service: DefaultEncryptService{},
}

type Encrypt string

// GORM
func (e *Encrypt) Scan(value any) error {
	if v, ok := value.(string); ok {
		*e = Encrypt(EncryptOptions.Service.Decode(v))
	}

	return nil
}

func (e Encrypt) Value() (driver.Value, error) {
	return EncryptOptions.Service.Encode(string(e)), nil
}

func (e Encrypt) Mask() string {
	return EncryptOptions.Service.Mask(string(e))
}

// String
func (e Encrypt) String() string {
	return e.Mask()
}
