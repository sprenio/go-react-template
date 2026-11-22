package domain

import (
	"backend/pkg/apperror"
	"database/sql/driver"
)

type CategoryType string

const (
	Expense             CategoryType = "EXPENSE"
	Income              CategoryType = "INCOME"
	defaultCategoryType              = Expense
)

func (y *CategoryType) Scan(value any) error {
	if value == nil {
		*y = defaultCategoryType
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*y = CategoryType(v)
	case string:
		*y = CategoryType(v)
	default:
		*y = CategoryType(defaultCategoryType)
	}

	if *y != Expense && *y != Income {
		*y = CategoryType(defaultCategoryType)
	}
	return nil
}

func (y CategoryType) Value() (driver.Value, error) {
	if y != Expense && y != Income {
		return nil, apperror.New("invalid ExpenseIncome enum", "enum", string(y))
	}
	return string(y), nil
}
