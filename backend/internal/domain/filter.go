package domain

import (
	"strconv"
	"strings"
	"time"
)

type ResultLimit struct {
	Offset uint
	Limit  uint
}

type FilterParser interface {
	GetFilterConditions(filter *Filter) (string, []interface{})
}

type Filter struct {
	Id             *uint
	SourceId       *uint
	SourceFromId   *uint
	SourceToId     *uint
	DateFrom       *time.Time
	DateTo         *time.Time
	YearMonth      *YearMonth
	MinYearMonth   *YearMonth
	MaxYearMonth   *YearMonth
	FromAmountFrom *Money
	FromAmountTo   *Money
	ToAmountFrom   *Money
	ToAmountTo     *Money
	AmountFrom     *Money
	AmountTo       *Money
	UserId         *uint
	Enabled        *YesNo
	CategoryId     *uint
	SubcategoryId  *uint
	SavingId       *uint
	ExpenseId      *uint
	ExpenseType    *ExpenseType
	CategoryType   *CategoryType
	CurrencyCode   *string
	Comment        *string
	Name           *string
	ExcludeIds     *[]uint

	ResultLimit *ResultLimit

	OrderBy []string

	extraFilters []Filter

	parser *FilterParser
}

func NewFilter() *Filter {
	return &Filter{}
}

func (f *Filter) SetId(id uint) *Filter {
	f.Id = &id
	return f
}

func (f *Filter) SetUserId(userId uint) *Filter {
	f.UserId = &userId
	return f
}

func (f *Filter) SetEnabled(enabled YesNo) *Filter {
	f.Enabled = &enabled
	return f
}

func (f *Filter) SetCurrencyCode(currencyCode string) *Filter {
	f.CurrencyCode = &currencyCode
	return f
}

func (f *Filter) SetComment(comment string) *Filter {
	f.Comment = &comment
	return f
}

func (f *Filter) SetCategoryType(categoryType CategoryType) *Filter {
	cattype := CategoryType(categoryType)
	f.CategoryType = &cattype
	return f
}

func (f *Filter) SetSourceFromId(sourceFromId uint) *Filter {
	f.SourceFromId = &sourceFromId
	return f
}

func (f *Filter) SetSourceToId(sourceToId uint) *Filter {
	f.SourceToId = &sourceToId
	return f
}

func (f *Filter) SetDateFrom(dateFrom *time.Time) *Filter {
	if dateFrom == nil {
		return f
	}
	d := time.Date(
		dateFrom.Year(),
		dateFrom.Month(),
		dateFrom.Day(),
		0, 0, 0, 0,
		dateFrom.Location(),
	)

	f.DateFrom = &d
	return f
}

func (f *Filter) SetSavingId(savingId uint) *Filter {
	f.SavingId = &savingId
	return f
}

func (f *Filter) SetSourceId(sourceId uint) *Filter {
	f.SourceId = &sourceId
	return f
}

func (f *Filter) SetDateTo(dateTo *time.Time) *Filter {
	if dateTo == nil {
		return f
	}
	d := time.Date(
		dateTo.Year(),
		dateTo.Month(),
		dateTo.Day(),
		23, 59, 59, 999999999,
		dateTo.Location(),
	)
	f.DateTo = &d
	return f
}

func (f *Filter) SetDate(date *time.Time) *Filter {
	if date == nil {
		return f
	}
	return f.SetDateFrom(date).SetDateTo(date)
}
func (f *Filter) SetYearMonth(yearMonth YearMonth) *Filter {
	f.YearMonth = &yearMonth
	return f
}

func (f *Filter) SetMinYearMonth(yearMonth YearMonth) *Filter {
	f.MinYearMonth = &yearMonth
	return f
}

func (f *Filter) SetMaxYearMonth(yearMonth YearMonth) *Filter {
	f.MaxYearMonth = &yearMonth
	return f
}

func (f *Filter) SetFromAmountFrom(fromAmountFrom *Money) *Filter {
	if fromAmountFrom == nil {
		return f
	}
	f.FromAmountFrom = fromAmountFrom
	return f
}

func (f *Filter) SetFromAmountTo(fromAmountTo *Money) *Filter {
	if fromAmountTo == nil {
		return f
	}
	f.FromAmountTo = fromAmountTo
	return f
}

func (f *Filter) SetToAmountFrom(toAmountFrom *Money) *Filter {
	if toAmountFrom == nil {
		return f
	}
	f.ToAmountFrom = toAmountFrom
	return f
}

func (f *Filter) SetToAmountTo(toAmountTo *Money) *Filter {
	if toAmountTo == nil {
		return f
	}
	f.ToAmountTo = toAmountTo
	return f
}

func (f *Filter) SetAmountFrom(amountFrom *Money) *Filter {
	if amountFrom == nil {
		return f
	}
	f.AmountFrom = amountFrom
	return f
}

func (f *Filter) SetAmountTo(amountTo *Money) *Filter {
	if amountTo == nil {
		return f
	}
	f.AmountTo = amountTo
	return f
}

func (f *Filter) SetCategoryId(categoryId uint) *Filter {
	f.CategoryId = &categoryId
	return f
}

func (f *Filter) SetSubcategoryId(subcategoryId uint) *Filter {
	f.SubcategoryId = &subcategoryId
	return f
}

func (f *Filter) SetExpenseId(expenseId uint) *Filter {
	f.ExpenseId = &expenseId
	return f
}

func (f *Filter) SetExpenseType(expenseType ExpenseType) *Filter {
	f.ExpenseType = &expenseType
	return f
}

func (f *Filter) SetExcludeIds(excludeIds []uint) *Filter {
	f.ExcludeIds = &excludeIds
	return f
}

func (f *Filter) SetResultLimit(limit uint) *Filter {
	if f.ResultLimit == nil {
		f.ResultLimit = &ResultLimit{}
	}
	f.ResultLimit.Limit = limit
	return f
}

func (f *Filter) AddOrderBy(orderBy string) *Filter {
	if orderBy != "" {
		f.OrderBy = append(f.OrderBy, orderBy)
	}
	return f
}

func (f *Filter) SetResultOffset(offset uint) *Filter {
	if f.ResultLimit == nil {
		f.ResultLimit = &ResultLimit{}
	}
	f.ResultLimit.Offset = offset
	return f
}

func (f *Filter) SetName(name string) *Filter {
	f.Name = &name
	return f
}

func (f *Filter) GetLimitStatement() string {
	if f.ResultLimit != nil {
		return f.ResultLimit.GetStatement()
	}
	return ""
}

func (f *Filter) GetYearMonthNumber() uint {
	if f.YearMonth != nil {
		return f.YearMonth.Number()
	}
	yearMonth, _ := NewYearMonthFromTime(time.Now())
	return yearMonth.Number()
}

func (f *Filter) GetFilterDates() (*time.Time, *time.Time) {
	if f == nil {
		now := time.Now()
		return nil, &now
	}
	var from *time.Time
	var to *time.Time
	if f.DateFrom != nil {
		from = f.DateFrom
	}
	if f.DateTo != nil {
		to = f.DateTo
	}
	if f.YearMonth != nil {
		startMonth := f.YearMonth.StartOfMonth()
		endMonth := f.YearMonth.EndOfMonth()
		if f.DateFrom == nil || f.DateFrom.Before(startMonth) {
			from = &startMonth
		}
		if f.DateTo == nil || f.DateTo.After(endMonth) {
			to = &endMonth
		}
	}
	if to == nil {
		now := time.Now()
		to = &now
	}
	if from != nil && from.After(*to) {
		from = nil
	}
	return from, to
}

func (rl *ResultLimit) GetStatement() string {
	if rl != nil && rl.Limit > 0 {
		statement := ""
		if rl.Offset > 0 {
			statement += strconv.FormatUint(uint64(rl.Offset), 10) + ","
		}
		statement += strconv.FormatUint(uint64(rl.Limit), 10)
		return "LIMIT " + statement
	}
	return ""
}

func (f *Filter) AddExtraFilter(other *Filter) {
	if other == nil {
		return
	}
	f.extraFilters = append(f.extraFilters, *other)
}

func (f *Filter) GetConditions() (string, []interface{}) {
	var conditions []string
	var params []interface{}

	if f.parser == nil {
		return "0", params
	}
	myParser := *f.parser
	condition, filterParams := myParser.GetFilterConditions(f)
	conditions = append(conditions, condition)
	params = append(params, filterParams...)
	for _, extraFilter := range f.extraFilters {
		condition, extraParams := myParser.GetFilterConditions(&extraFilter)
		conditions = append(conditions, condition)
		params = append(params, extraParams...)
	}
	if len(conditions) == 0 {
		conditions = append(conditions, "1")
	}
	return strings.Join(conditions, " AND "), params
}

func (f *Filter) SetConditionParser(parser FilterParser) *Filter {
	f.parser = &parser
	return f
}

func (f *Filter) IsCategoryTableUsed() bool {
	if f.CategoryType != nil || f.CategoryId != nil || f.UserId != nil {
		return true
	}
	for _, extraFilter := range f.extraFilters {
		if extraFilter.IsCategoryTableUsed() {
			return true
		}
	}
	return false
}

func (f *Filter) IsSubCategoryTableUsed() bool {
	if f.SubcategoryId != nil {
		return true
	}
	for _, extraFilter := range f.extraFilters {
		if extraFilter.IsSubCategoryTableUsed() {
			return true
		}
	}
	return false
}

func (f *Filter) IsSourceTableUsed() bool {
	if f.SourceId != nil || f.UserId != nil {
		return true
	}
	for _, extraFilter := range f.extraFilters {
		if extraFilter.IsSourceTableUsed() {
			return true
		}
	}
	return false
}

func (f *Filter) IsCommentTableUsed() bool {
	if f.Comment != nil {
		return true
	}
	for _, extraFilter := range f.extraFilters {
		if extraFilter.IsCommentTableUsed() {
			return true
		}
	}
	return false
}
