package domain

import (
	"fmt"
	"strings"
)

type Table struct {
	Headers []string   `json:"headers"`
	Rows    [][]string `json:"rows"`
}

func (table *Table) String() string {
	return fmt.Sprintf("Table[%s] (%d rows)", strings.Join(table.Headers, ", "), len(table.Rows))
}
