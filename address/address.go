package address

import (
	_ "embed"
)

//go:embed address_template.json
var Template string

const (
	TemplateName = "address"
	IndexName    = "addresses"
	Alias        = IndexName + "_v1"
)

type Address struct {
	ID    int64  `json:"id"`
	Line1 string `json:"line_1"`
	Line2 string `json:"line_2"`
}
