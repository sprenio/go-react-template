package domain

import (
	"backend/pkg/apperror"
	"database/sql/driver"
	"fmt"
)

type ExpenseType string

const (
	ExpenseTypeAnnual    ExpenseType = "ANNUAL"
	ExpenseTypeFixed     ExpenseType = "FIXED"
	ExpenseTypeIrregular ExpenseType = "IRREGULAR"
	ExpenseTypeOneOff    ExpenseType = "ONEOFF"
)

func (y *ExpenseType) Scan(value any) error {
	if value == nil {
		return apperror.New("invalid ExpenseType value")
	}
	switch v := value.(type) {
	case []byte:
		*y = ExpenseType(v)
	case string:
		*y = ExpenseType(v)
	default:
		return apperror.New(fmt.Sprintf("invalid ExpenseType value: %T", value))
	}

	if *y != ExpenseTypeAnnual && *y != ExpenseTypeFixed && *y != ExpenseTypeIrregular && *y != ExpenseTypeOneOff {
		return apperror.New(fmt.Sprintf("invalid ExpenseType enum: %s", *y))
	}
	return nil
}

func (y ExpenseType) Value() (driver.Value, error) {
	if y != ExpenseTypeAnnual && y != ExpenseTypeFixed && y != ExpenseTypeIrregular && y != ExpenseTypeOneOff {
		return nil, apperror.New(fmt.Sprintf("invalid ExpenseType enum: %s", y))
	}
	return string(y), nil
}
