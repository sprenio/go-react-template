package domain

import (
	"backend/pkg/apperror"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type YearMonth struct {
	Year  int
	Month int
}

func NewYearMonth(year, month int) (YearMonth, error) {
	ym := YearMonth{Year: year, Month: month}
	return ym, ym.Validate()
}
func NewYearMonthFromTime(t time.Time) (YearMonth, error) {
	return NewYearMonth(t.Year(), int(t.Month()))
}
func NewYearMonthFromString(s string) (YearMonth, error) {
	if strings.Contains(s, "-") {
		t, err := time.Parse("2006-01", s)
		if err != nil {
			return YearMonth{}, err
		}
		return NewYearMonthFromTime(t)
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return YearMonth{}, err
	}
	return NewYearMonth(n/100, n%100)
}

func (ym YearMonth) Number() uint {
	if err := ym.Validate(); err != nil {
		return 0
	}
	return uint(ym.Year*100 + ym.Month)
}
func (ym *YearMonth) Value() (driver.Value, error) {
	if err := ym.Validate(); err != nil {
		return nil, err
	}
	return int64(ym.Number()), nil
}

func (ym *YearMonth) Scan(value any) error {
	if value == nil {
		return errors.New("YearMonth cannot be NULL")
	}

	var v int
	switch t := value.(type) {
	case int64:
		v = int(t)
	case int32:
		v = int(t)
	case int:
		v = t
	case []byte:
		n, err := strconv.Atoi(string(t))
		if err != nil {
			return err
		}
		v = n
	case string:
		n, err := strconv.Atoi(t)
		if err != nil {
			return err
		}
		v = n
	default:
		return apperror.New(fmt.Sprintf("cannot scan %T into YearMonth", value))
	}

	ym.Year = v / 100
	ym.Month = v % 100

	return ym.Validate()
}

func (ym YearMonth) ToTime() time.Time {
	return ym.StartOfMonth()
}

func (ym YearMonth) EndOfMonth() time.Time {
	return ym.StartOfMonth().AddDate(0, 1, 0).Add(-time.Second)
}

func (ym YearMonth) StartOfMonth() time.Time {
	return time.Date(ym.Year, time.Month(ym.Month), 1, 0, 0, 0, 0, time.UTC)
}

func (ym *YearMonth) AddMonths(m int) (YearMonth, error) {
	newTime := ym.ToTime().AddDate(0, m, 0)
	return NewYearMonthFromTime(newTime)
}

func (ym *YearMonth) Validate() error {
	if ym.Year < 2020 {
		return apperror.New(fmt.Sprintf("invalid year: %d", ym.Year))
	}
	if ym.Month < 1 || ym.Month > 12 {
		return apperror.New(fmt.Sprintf("invalid month: %d", ym.Month))
	}
	return nil
}

func (ym *YearMonth) Before(other *YearMonth) bool {
	return ym.Number() < other.Number()
}

func (ym *YearMonth) After(other *YearMonth) bool {
	return ym.Number() > other.Number()
}

func (ym *YearMonth) Equal(other *YearMonth) bool {
	return ym.Number() == other.Number()
}

func (ym YearMonth) MonthsBetween(other *YearMonth) uint {
	m := (other.Year-ym.Year)*12 + (other.Month - ym.Month)
	if m < 0 {
		return uint(m * -1)
	}
	return uint(m)
}

func (ym *YearMonth) String() string {
	return fmt.Sprintf("%04d-%02d", ym.Year, ym.Month)
}
func (ym *YearMonth) MarshalJSON() ([]byte, error) {
	if err := ym.Validate(); err != nil {
		return nil, err
	}
	return []byte(strconv.FormatUint(uint64(ym.Number()), 10)), nil
}
func (ym *YearMonth) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return errors.New("YearMonth cannot be null")
	}

	// string
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		s = strings.TrimSpace(s)

		if strings.Contains(s, "-") {
			t, err := time.Parse("2006-01", s)
			if err != nil {
				return err
			}
			parsed, parseErr := NewYearMonthFromTime(t)
			if parseErr != nil {
				return parseErr
			}
			return parsed.Validate()
		}

		// YYYYMM
		n, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		ym.Year = n / 100
		ym.Month = n % 100
		return ym.Validate()
	}

	// number
	var n int
	if err := json.Unmarshal(data, &n); err == nil {
		ym.Year = n / 100
		ym.Month = n % 100
		return ym.Validate()
	}

	return apperror.New(fmt.Sprintf("invalid YearMonth JSON value: %s", string(data)))
}
