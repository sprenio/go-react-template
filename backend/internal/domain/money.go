package domain

import (
	"backend/pkg/apperror"
	"backend/pkg/logger"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Money represents monetary value stored in minor units (e.g. cents).
// Example: 1234 = 12.34
type Money int64

type MoneyWithCurrencyCode struct {
	CurrencyCode string `json:"currency_code"`
	Amount       Money  `json:"amount"`
}

func NewMoneyFromCents(cents int64) Money {
	return Money(cents)
}

func NewMoneyFromFloat(amount float64) Money {
	return Money(math.Round(amount * 100))
}

func NewMoneyFromString(s string) (Money, error) {
	// Trim leading and trailing whitespace
	s = strings.TrimSpace(s)

	// Empty string is treated as zero
	if s == "" {
		return 0, nil
	}
	// Normalize European decimal separator (comma) to dot
	// Reject ambiguous formats like "1,234.56"
	if strings.Contains(s, ",") {
		if strings.Contains(s, ".") {
			return 0, apperror.New(fmt.Sprintf("invalid money format: %s", s))
		}
		s = strings.Replace(s, ",", ".", 1)
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return NewMoneyFromFloat(f), nil
}

func MapToMoneySlice(m map[string]Money) []MoneyWithCurrencyCode {
	result := make([]MoneyWithCurrencyCode, 0, len(m))

	for currency, amount := range m {
		result = append(result, MoneyWithCurrencyCode{
			CurrencyCode: currency,
			Amount:       amount,
		})
	}

	return result
}

func (m Money) Cents() int64 {
	return int64(m)
}

func (m Money) Float() float64 {
	return float64(m) / 100
}

func (m Money) Add(other Money) Money {
	return m + other
}

func (m Money) Sub(other Money) Money {
	return m - other
}

func (m Money) Mul(multiplier int64) Money {
	return Money(int64(m) * multiplier)
}
func (m Money) Div(divisor int64) Money {
	if divisor == 0 {
		panic("money: division by zero")
	}
	return NewMoneyFromFloat(float64(m) / float64(divisor) / 100)
}

func (m Money) IsNegative() bool {
	return m < 0
}

func (m Money) IsZero() bool {
	return m == 0
}
func (m Money) GreaterThan(other Money) bool {
	return m > other
}

func (m Money) LessThan(other Money) bool {
	return m < other
}

func (m Money) Ceil(step int64) Money {
	if step <= 0 {
		return m
	}

	stepInCents := float64(step * 100)
	roundedCents := math.Ceil(float64(m)/stepInCents) * stepInCents

	return NewMoneyFromCents(int64(roundedCents))
}

func (m Money) Floor(step int64) Money {
	if step <= 0 {
		return m
	}

	stepInCents := float64(step * 100)
	roundedCents := math.Floor(float64(m)/stepInCents) * stepInCents

	return NewMoneyFromCents(int64(roundedCents))
}

func (m Money) Round(step int64) Money {
	if step <= 0 {
		return m
	}

	stepInCents := float64(step * 100)
	roundedCents := math.Round(float64(m)/stepInCents) * stepInCents

	return NewMoneyFromCents(int64(roundedCents))
}

func (m Money) String() string {
	return fmt.Sprintf("%.2f", m.Float())
}

func (m Money) Value() (driver.Value, error) {
	return m.Float(), nil // MySQL DECIMAL(10,2)
}
func (m *Money) Scan(value any) error {
	if value == nil {
		*m = 0
		return nil
	}

	switch v := value.(type) {
	case int64:
		*m = Money(v * 100)
	case float64:
		*m = NewMoneyFromFloat(v)
	case []byte:
		f, err := strconv.ParseFloat(string(v), 64)
		if err != nil {
			return err
		}
		*m = NewMoneyFromFloat(f)
	default:
		return errors.New("unsupported money type")
	}
	return nil
}

func (m Money) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(m.Cents(), 10)), nil
}
func (m *Money) UnmarshalJSON(data []byte) error {
	// null → zero
	if string(data) == "null" {
		*m = 0
		return nil
	}
	// string?
	if len(data) >= 2 && data[0] == '"' && data[len(data)-1] == '"' {
		s := string(data[1 : len(data)-1])
		// empty string → zero
		if s == "" {
			*m = 0
			return nil
		}

		val, err := NewMoneyFromString(s)
		if err != nil {
			return err
		}
		*m = val
		return nil
	}
	// number (eg. 12.34)
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		logger.Error("Error unmarshalling money: %v", err)
		return err
	}
	if math.Trunc(f) == f {
		// integer → cents
		*m = NewMoneyFromCents(int64(f))
	} else {
		// float → amount
		*m = NewMoneyFromFloat(f)
	}
	return nil
}
