package domain

import (
	"backend/pkg/apperror"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

type YesNo string

const (
	Yes YesNo = "YES"
	No  YesNo = "NO"
)

func (y *YesNo) Scan(value any) error {
	if value == nil {
		*y = No
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*y = YesNo(v)
	case string:
		*y = YesNo(v)
	default:
		return apperror.New(fmt.Sprintf("invalid YesNo value: %T", value))
	}

	if *y != Yes && *y != No {
		return apperror.New(fmt.Sprintf("invalid YesNo enum: %s", *y))
	}
	return nil
}

func (y YesNo) Value() (driver.Value, error) {
	if y != Yes && y != No {
		return nil, apperror.New(fmt.Sprintf("invalid YesNo enum: %s", y))
	}
	return string(y), nil
}

func NewYesNoFromString(s string) (YesNo, error) {
	value := YesNo(s)
	if value != Yes && value != No {
		return "", apperror.New(fmt.Sprintf("invalid YesNo enum: %s", value))
	}
	return value, nil
}

func (y YesNo) MarshalJSON() ([]byte, error) {
	if y != Yes && y != No {
		return nil, apperror.New(fmt.Sprintf("invalid YesNo enum: %s", y))
	}

	if y == Yes {
		return []byte("true"), nil
	}
	return []byte("false"), nil
}

func (y *YesNo) UnmarshalJSON(data []byte) error {
	// null → No
	if string(data) == "null" {
		*y = No
		return nil
	}

	// bool
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		if b {
			*y = Yes
		} else {
			*y = No
		}
		return nil
	}

	// string
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		s = strings.ToUpper(strings.TrimSpace(s))
		if s == "" {
			*y = No
			return nil
		}

		switch s {
		case "YES":
			*y = Yes
		case "NO":
			*y = No
		default:
			return apperror.New(fmt.Sprintf("invalid YesNo value: %s", s))
		}
		return nil
	}

	return apperror.New(fmt.Sprintf("invalid YesNo JSON value: %s", string(data)))
}

func (y *YesNo) ToBool() bool {
	return *y == Yes
}
