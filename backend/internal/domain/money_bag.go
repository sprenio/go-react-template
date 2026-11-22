package domain

import (
	"backend/internal/contexthelper"
	"context"
	"encoding/json"
	"maps"
	"sort"
	"strings"
)

type MoneyBag struct {
	amounts map[string]Money
	order   []string
}

type ComparisonBag struct {
	Planned MoneyBag
	Actual  MoneyBag
}
type ItemBag struct {
	Id  uint
	Bag MoneyBag
}

func NewMoneyBag() MoneyBag {
	return MoneyBag{
		amounts: make(map[string]Money),
		order:   make([]string, 0),
	}
}

func NewMoneyBagFromContext(ctx context.Context) MoneyBag {
	return MoneyBag{
		amounts: make(map[string]Money),
		order:   contexthelper.GetCurrencyOrder(ctx),
	}
}

func NewMoneyBagFromSlice(items []MoneyWithCurrencyCode) MoneyBag {
	bag := NewMoneyBag()

	for _, item := range items {
		bag.Add(item.CurrencyCode, item.Amount)
	}

	return bag
}

func normalizeCurrency(currency string) string {
	return strings.ToUpper(strings.TrimSpace(currency))
}

func (b *MoneyBag) SetOrder(order []string) {
	b.order = order
}

func (b MoneyBag) Add(currency string, amount Money) {
	currency = normalizeCurrency(currency)

	b.amounts[currency] += amount

	if b.amounts[currency].IsZero() {
		delete(b.amounts, currency)
	}
}

func (b MoneyBag) Sub(currency string, amount Money) {
	b.Add(currency, -amount)
}

func (b MoneyBag) Mul(multiplier int64) {
	for currency, amount := range b.amounts {
		b.amounts[currency] = amount.Mul(multiplier)

		if b.amounts[currency].IsZero() {
			delete(b.amounts, currency)
		}
	}
}

func (b MoneyBag) Div(divisor int64) {
	if divisor == 0 {
		panic("moneybag: division by zero")
	}

	for currency, amount := range b.amounts {
		b.amounts[currency] = amount.Div(divisor)

		if b.amounts[currency].IsZero() {
			delete(b.amounts, currency)
		}
	}
}

func (b MoneyBag) AddBag(other MoneyBag) {
	for currency, amount := range other.amounts {
		b.Add(currency, amount)
	}
}

func (b MoneyBag) SubBag(other MoneyBag) {
	for currency, amount := range other.amounts {
		b.Sub(currency, amount)
	}
}

func (b MoneyBag) Ceil(step int64) MoneyBag {
	result := b.newBag()
	for currency, amount := range b.amounts {
		result.Add(currency, amount.Ceil(step))
	}
	return result
}

func (b MoneyBag) Floor(step int64) MoneyBag {
	result := b.newBag()
	for currency, amount := range b.amounts {
		result.Add(currency, amount.Floor(step))
	}
	return result
}

func (b MoneyBag) Round(step int64) MoneyBag {
	result := b.newBag()
	for currency, amount := range b.amounts {
		result.Add(currency, amount.Round(step))
	}
	return result
}

func (b MoneyBag) Amount(currency string) Money {
	return b.amounts[normalizeCurrency(currency)]
}

func (b MoneyBag) HasCurrency(currency string) bool {
	_, exists := b.amounts[normalizeCurrency(currency)]
	return exists
}

func (b MoneyBag) IsZero() bool {
	return len(b.amounts) == 0
}

func (b MoneyBag) Clone() MoneyBag {
	result := b.newBag()
	result.amounts = maps.Clone(b.amounts)
	result.order = append([]string{}, b.order...)
	return result
}

func (b MoneyBag) Negate() MoneyBag {
	result := b.newBag()

	for currency, amount := range b.amounts {
		result.amounts[currency] = -amount
	}

	return result
}

func (b MoneyBag) AddSlice(items []MoneyWithCurrencyCode) MoneyBag {
	for _, item := range items {
		b.Add(item.CurrencyCode, item.Amount)
	}
	return b
}

func (b MoneyBag) ToSlice() []MoneyWithCurrencyCode {
	result := make([]MoneyWithCurrencyCode, 0, len(b.amounts))

	// sort for deterministic results
	currencies := make([]string, 0, len(b.amounts))

	for currency := range b.amounts {
		currencies = append(currencies, currency)
	}

	sort.Strings(currencies)

	for _, currency := range currencies {
		result = append(result, MoneyWithCurrencyCode{
			CurrencyCode: currency,
			Amount:       b.amounts[currency],
		})
	}

	return result
}

func (b MoneyBag) ToSortedSlice(order []string) []MoneyWithCurrencyCode {
	result := make([]MoneyWithCurrencyCode, 0, len(b.amounts))

	// map currency positions
	orderIndex := make(map[string]int, len(order))

	for i, currency := range order {
		orderIndex[normalizeCurrency(currency)] = i
	}

	// build slice
	for currency, amount := range b.amounts {
		result = append(result, MoneyWithCurrencyCode{
			CurrencyCode: currency,
			Amount:       amount,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		left := normalizeCurrency(result[i].CurrencyCode)
		right := normalizeCurrency(result[j].CurrencyCode)

		leftPos, leftExists := orderIndex[left]
		rightPos, rightExists := orderIndex[right]

		// both exist in order
		if leftExists && rightExists {
			return leftPos < rightPos
		}

		// only left exists -> earlier
		if leftExists {
			return true
		}

		// only right exists -> later
		if rightExists {
			return false
		}

		// neither exists -> alphabetically
		return strings.Compare(left, right) < 0
	})

	return result
}

func (b MoneyBag) newBag() MoneyBag {
	result := NewMoneyBag()
	if len(b.order) > 0 {
		result.order = append([]string{}, b.order...)
	}
	return result
}

func (b MoneyBag) MarshalJSON() ([]byte, error) {
	if len(b.order) > 0 {
		return json.Marshal(b.ToSortedSlice(b.order))
	}
	return json.Marshal(b.ToSlice())
}
