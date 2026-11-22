package domain

type TableAlias struct {
	Name         string
	Alias        string
	ColumnPrefix string
}

func NewTableAlias(name, alias string) TableAlias {
	columnPrefix := alias
	if alias != "" {
		columnPrefix = alias + "."
	}
	return TableAlias{Name: name, Alias: alias, ColumnPrefix: columnPrefix}
}
