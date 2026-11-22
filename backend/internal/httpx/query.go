package httpx

import (
	"backend/internal/domain"
	"backend/internal/response"
	"backend/pkg/apperror"
	"backend/pkg/logger"
	"net/http"
	"strconv"
	"time"
)

type AbortRequest struct{}

type QueryParams struct {
	Date          *time.Time
	YearMonth     *domain.YearMonth
	SourceId      *uint
	CategoryId    *uint
	SubcategoryId *uint
	OperationType *domain.CategoryType
	Query         *string
	Expenses      *bool
}

func ParseQueryParams(r *http.Request) (QueryParams, error) {
	var qp QueryParams

	if dateStr := r.FormValue("date"); dateStr != "" {
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return qp, apperror.Wrap(err, "invalid date")
		}
		qp.Date = &t
	}

	if ymStr := r.FormValue("yearMonth"); ymStr != "" {
		ym, err := domain.NewYearMonthFromString(ymStr)
		if err != nil {
			return qp, apperror.Wrap(err, "invalid yearMonth")
		}
		logger.InfoCtx(r.Context(), "Parsed yearMonth", "yearMonth", ym)
		qp.YearMonth = &ym
	}

	if sourceIdStr := r.FormValue("sourceId"); sourceIdStr != "" {
		sourceId, err := strconv.ParseUint(sourceIdStr, 10, 64)
		if err != nil {
			return qp, apperror.Wrap(err, "invalid sourceId")
		}
		var sourceIdUint uint = uint(sourceId)
		qp.SourceId = &sourceIdUint
	}

	if categoryIdStr := r.FormValue("categoryId"); categoryIdStr != "" {
		categoryId, err := strconv.ParseUint(categoryIdStr, 10, 64)
		if err != nil {
			return qp, apperror.Wrap(err, "invalid categoryId")
		}
		var categoryIdUint uint = uint(categoryId)
		qp.CategoryId = &categoryIdUint
	}

	if subcategoryIdStr := r.FormValue("subcategoryId"); subcategoryIdStr != "" {
		subcategoryId, err := strconv.ParseUint(subcategoryIdStr, 10, 64)
		if err != nil {
			return qp, apperror.Wrap(err, "invalid subcategoryId")
		}
		var subcategoryIdUint uint = uint(subcategoryId)
		qp.SubcategoryId = &subcategoryIdUint
	}

	if operationTypeStr := r.FormValue("operationType"); operationTypeStr != "" {
		var operationType domain.CategoryType
		switch operationTypeStr {
		case string(domain.Expense):
			operationType = domain.Expense
		case string(domain.Income):
			operationType = domain.Income
		default:
			return qp, apperror.New("invalid operationType")
		}
		qp.OperationType = &operationType
	}

	if queryStr := r.FormValue("query"); queryStr != "" {
		qp.Query = &queryStr
	}

	if expensesStr := r.FormValue("expenses"); expensesStr != "" {
		expenses, err := strconv.ParseBool(expensesStr)
		if err != nil {
			return qp, apperror.Wrap(err, "invalid expenses")
		}
		qp.Expenses = &expenses
	}

	return qp, nil
}

func (qp *QueryParams) ResolveYearMonth() domain.YearMonth {
	if qp.YearMonth != nil {
		return *qp.YearMonth
	}
	yearMonth, _ := domain.NewYearMonthFromTime(time.Now())
	return yearMonth
}

func (qp *QueryParams) ResolveDate() time.Time {
	if qp.Date != nil {
		return *qp.Date
	}

	if qp.YearMonth != nil {
		return qp.YearMonth.EndOfMonth()
	}

	now := time.Now()
	ym := domain.YearMonth{
		Year:  now.Year(),
		Month: int(now.Month()),
	}

	return ym.EndOfMonth()
}

func (qp *QueryParams) ResolveSourceId() uint {
	if qp.SourceId != nil {
		return *qp.SourceId
	}
	return 0
}

func (qp *QueryParams) ResolveCategoryId() uint {
	if qp.CategoryId != nil {
		return *qp.CategoryId
	}
	return 0
}

func (qp *QueryParams) ResolveSubcategoryId() uint {
	if qp.SubcategoryId != nil {
		return *qp.SubcategoryId
	}
	return 0
}

func (qp *QueryParams) ResolveQuery() string {
	if qp.Query != nil {
		return *qp.Query
	}
	return ""
}

func (qp *QueryParams) ResolveOperationType() *domain.CategoryType {
	if qp.OperationType != nil {
		return qp.OperationType
	}
	return nil
}

func (qp *QueryParams) ResolveExpenses() bool {
	if qp.Expenses != nil {
		return *qp.Expenses
	}
	return false
}

func MustParseQueryParams(w http.ResponseWriter, r *http.Request) QueryParams {
	qp, err := ParseQueryParams(r)
	if err != nil {
		response.BadRequestErrorResponse(w)
		logger.ErrorCtx(r.Context(), "Error parsing query params", "error", err)
		panic(AbortRequest{})
	}
	return qp
}
