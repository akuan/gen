package model

type DicValue struct {
	ID         int    `gorm:"column:id;primary_key" json:"ID"`
	OptionType string `gorm:"column:option_type" json:"optionType"`
	Name       string `gorm:"column:name" json:"name"`
	SortIndex  int    `gorm:"column:sort_index" json:"sortIndex,omitempty"`
}

// TableName sets the insert table name for this struct type
func (d *DicValue) TableName() string {
	return "dic_value"
}
