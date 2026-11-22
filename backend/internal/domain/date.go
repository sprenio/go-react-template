package domain

import (
	"backend/pkg/apperror"
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type Date struct {
	time.Time
}

const dateLayout = "2006-01-02"

func NewDate(t time.Time) Date {
	return Date{t}
}

func (d Date) Value() (driver.Value, error) {
	return d.Time, nil
}

func (d *Date) Scan(value any) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		d.Time = v
		return nil
	case string:
		parsed, err := time.Parse(dateLayout, v)
		if err != nil {
			return err
		}
		d.Time = parsed
		return nil
	default:
		return apperror.New(fmt.Sprintf("cannot scan %T into Date", value))
	}
}

func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + d.Time.Format(dateLayout) + `"`), nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)

	if s == "" || s == "null" {
		d.Time = time.Time{}
		return nil
	}

	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return err
	}

	d.Time = t
	return nil
}
