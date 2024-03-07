package schema

type Column struct {
	OrdinalPosition int    `json:"ordinalPosition"`
	ColumnName      string `json:"columnName"`
	DataType        string `json:"dataType"`
}

type QueryOption struct {
	Take  int
	Skip  int64
	Order *string
}

type QueryFilter struct {
	Equal         []string
	NotEqual      []string
	Like          []string
	Greater       []string
	Lesser        []string
	In            []string
	NotIn         []string
	ObjectContain []string
	ArrayContain  []string
	ObjectLike    []string
}
