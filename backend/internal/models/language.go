package models

type Language struct {
	Id   uint8    `db:"id" json:"id"`
	Code string `db:"code" json:"code"`
	I18nCode string `db:"i18n_code" json:"i18n_code"`
	Name string `db:"name" json:"name"`
}
