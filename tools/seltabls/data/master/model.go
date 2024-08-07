// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package master

import (
	"time"
)

type Field struct {
	ID       int64  `db:"id" json:"id"`
	StructID int64  `db:"struct_id" json:"struct_id"`
	Name     string `db:"name" json:"name"`
	LineID   int64  `db:"line_id" json:"line_id"`
}

type File struct {
	ID        int64      `db:"id" json:"id"`
	Uri       string     `db:"uri" json:"uri"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
}

type Html struct {
	ID        int64      `db:"id" json:"id"`
	Value     string     `db:"value" json:"value"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
}

type Line struct {
	ID     int64  `db:"id" json:"id"`
	FileID int64  `db:"file_id" json:"file_id"`
	Value  string `db:"value" json:"value"`
	Number int64  `db:"number" json:"number"`
}

type Selector struct {
	ID         int64  `db:"id" json:"id"`
	Value      string `db:"value" json:"value"`
	UrlID      int64  `db:"url_id" json:"url_id"`
	Occurances int64  `db:"occurances" json:"occurances"`
	Context    string `db:"context" json:"context"`
}

type Struct struct {
	ID          int64  `db:"id" json:"id"`
	Value       string `db:"value" json:"value"`
	UrlID       int64  `db:"url_id" json:"url_id"`
	StartLineID int64  `db:"start_line_id" json:"start_line_id"`
	EndLineID   int64  `db:"end_line_id" json:"end_line_id"`
	FileID      int64  `db:"file_id" json:"file_id"`
	Context     string `db:"context" json:"context"`
}

type Tag struct {
	ID      int64  `db:"id" json:"id"`
	Value   string `db:"value" json:"value"`
	Start   int64  `db:"start" json:"start"`
	End     int64  `db:"end" json:"end"`
	LineID  int64  `db:"line_id" json:"line_id"`
	FieldID int64  `db:"field_id" json:"field_id"`
}

type Url struct {
	ID     int64  `db:"id" json:"id"`
	Value  string `db:"value" json:"value"`
	HtmlID int64  `db:"html_id" json:"html_id"`
}
