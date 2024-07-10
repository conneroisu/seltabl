// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package master

import (
	"context"
)

type Querier interface {
	//DeleteFileByID
	//
	//  DELETE FROM
	//      files
	//  WHERE
	//      id = ?
	DeleteFileByID(ctx context.Context, arg DeleteFileByIDParams) error
	//DeleteHTMLByID
	//
	//  DELETE FROM
	//      htmls
	//  WHERE
	//      id = ?
	DeleteHTMLByID(ctx context.Context, arg DeleteHTMLByIDParams) error
	//DeleteLineByID
	//
	//  DELETE FROM
	//      lines
	//  WHERE
	//      id = ?
	DeleteLineByID(ctx context.Context, arg DeleteLineByIDParams) error
	//DeleteSelectorByID
	//
	//  DELETE FROM
	//  	selectors
	//  WHERE
	//  	id = ?
	DeleteSelectorByID(ctx context.Context, arg DeleteSelectorByIDParams) error
	//DeleteStructByID
	//
	//  DELETE FROM
	//      structs
	//  WHERE
	//      id = ?
	DeleteStructByID(ctx context.Context, arg DeleteStructByIDParams) error
	//DeleteTagByID
	//
	//  DELETE FROM
	//      tags
	//  WHERE
	//      id = ?
	DeleteTagByID(ctx context.Context, arg DeleteTagByIDParams) error
	//DeleteURL
	//
	//  DELETE FROM
	//      urls
	//  WHERE
	//      id = ?
	DeleteURL(ctx context.Context, arg DeleteURLParams) error
	//GetFileByID
	//
	//  SELECT
	//      id, uri, updated_at, created_at
	//  FROM
	//      files
	//  WHERE
	//      id = ? LIMIT 1
	GetFileByID(ctx context.Context, arg GetFileByIDParams) (*File, error)
	//GetFileByURI
	//
	//  SELECT
	//      id, uri, updated_at, created_at
	//  FROM
	//      files
	//  WHERE
	//      uri = ?
	GetFileByURI(ctx context.Context, arg GetFileByURIParams) (*File, error)
	//GetLineByID
	//
	//  SELECT
	//      id, value, number
	//  FROM
	//      lines
	//  WHERE
	//      id = ?
	GetLineByID(ctx context.Context, arg GetLineByIDParams) (*GetLineByIDRow, error)
	//GetSelectorByID
	//
	//  SELECT
	//  	id, value, url_id, occurances, context
	//  FROM
	//  	selectors
	//  WHERE
	//  	id = ?
	GetSelectorByID(ctx context.Context, arg GetSelectorByIDParams) (*Selector, error)
	//GetSelectorByValue
	//
	//  SELECT
	//  	id, value, url_id, occurances, context
	//  FROM
	//  	selectors
	//  WHERE
	//  	value = ?
	GetSelectorByValue(ctx context.Context, arg GetSelectorByValueParams) (*Selector, error)
	//GetSelectorsByContext
	//
	//  SELECT
	//  	id, value, url_id, occurances, context
	//  FROM
	//  	selectors
	//  WHERE
	//  	context = ?
	GetSelectorsByContext(ctx context.Context, arg GetSelectorsByContextParams) ([]*Selector, error)
	//GetSelectorsByURL
	//
	//  SELECT
	//  	selectors.id, selectors.value, selectors.url_id, selectors.occurances, selectors.context
	//  FROM
	//  	selectors
	//  	JOIN urls ON urls.id = selectors.url_id
	//  WHERE
	//  	urls.value = ?
	GetSelectorsByURL(ctx context.Context, arg GetSelectorsByURLParams) ([]*Selector, error)
	//GetStructByID
	//
	//  SELECT
	//      id, value, url_id, start_line_id, end_line_id, file_id, context
	//  FROM
	//      structs
	//  WHERE
	//      id = ?
	GetStructByID(ctx context.Context, arg GetStructByIDParams) (*Struct, error)
	//GetStructByValue
	//
	//  SELECT
	//      id, value, url_id, start_line_id, end_line_id, file_id, context
	//  FROM
	//      structs
	//  WHERE
	//      value = ?
	GetStructByValue(ctx context.Context, arg GetStructByValueParams) (*Struct, error)
	//GetStructsByEndLineID
	//
	//  SELECT
	//      id, value, url_id, start_line_id, end_line_id, file_id, context
	//  FROM
	//      structs
	//  WHERE
	//      end_line_id = ?
	GetStructsByEndLineID(ctx context.Context, arg GetStructsByEndLineIDParams) ([]*Struct, error)
	//GetStructsByFileID
	//
	//  SELECT
	//      id, value, url_id, start_line_id, end_line_id, file_id, context
	//  FROM
	//      structs
	//  WHERE
	//      file_id = ?
	GetStructsByFileID(ctx context.Context, arg GetStructsByFileIDParams) ([]*Struct, error)
	//GetStructsByFileIDAndEndLineID
	//
	//  SELECT
	//      id, value, url_id, start_line_id, end_line_id, file_id, context
	//  FROM
	//      structs
	//  WHERE
	//      file_id = ? AND end_line_id = ?
	GetStructsByFileIDAndEndLineID(ctx context.Context, arg GetStructsByFileIDAndEndLineIDParams) ([]*Struct, error)
	//GetStructsByFileIDAndStartLineID
	//
	//  SELECT
	//      id, value, url_id, start_line_id, end_line_id, file_id, context
	//  FROM
	//      structs
	//  WHERE
	//      file_id = ? AND start_line_id = ?
	GetStructsByFileIDAndStartLineID(ctx context.Context, arg GetStructsByFileIDAndStartLineIDParams) ([]*Struct, error)
	//GetStructsByStartLineEndlineRange
	//
	//  SELECT
	//      structs.id, structs.value, url_id, start_line_id, end_line_id, structs.file_id, context, lines.id, lines.file_id, lines.value, lines.number, lines2.id, lines2.file_id, lines2.value, lines2.number
	//  FROM
	//      structs
	//  JOIN lines ON lines.id = structs.start_line_id
	//  JOIN lines AS lines2 ON lines2.id = structs.end_line_id
	//  WHERE
	//      lines.number <= ? AND lines2.number >= ?
	GetStructsByStartLineEndlineRange(ctx context.Context, arg GetStructsByStartLineEndlineRangeParams) ([]*GetStructsByStartLineEndlineRangeRow, error)
	//GetStructsByStartLineID
	//
	//  SELECT
	//      id, value, url_id, start_line_id, end_line_id, file_id, context
	//  FROM
	//      structs
	//  WHERE
	//      start_line_id = ?
	GetStructsByStartLineID(ctx context.Context, arg GetStructsByStartLineIDParams) ([]*Struct, error)
	//GetStructsByValue
	//
	//  SELECT
	//      id, value, url_id, start_line_id, end_line_id, file_id, context
	//  FROM
	//      structs
	//  WHERE
	//      value = ?
	GetStructsByValue(ctx context.Context, arg GetStructsByValueParams) ([]*Struct, error)
	//GetTagByFieldIDAndValue
	//
	//  SELECT
	//      id, value, start, "end", line_id, field_id
	//  FROM
	//      tags
	//  WHERE
	//      field_id = ? AND value = ?
	GetTagByFieldIDAndValue(ctx context.Context, arg GetTagByFieldIDAndValueParams) (*Tag, error)
	//GetTagByID
	//
	//  SELECT
	//      id, value, start, "end", line_id, field_id
	//  FROM
	//      tags
	//  WHERE
	//      id = ?
	GetTagByID(ctx context.Context, arg GetTagByIDParams) (*Tag, error)
	//GetTagByValue
	//
	//  SELECT
	//      id, value, start, "end", line_id, field_id
	//  FROM
	//      tags
	//  WHERE
	//      value = ?
	GetTagByValue(ctx context.Context, arg GetTagByValueParams) (*Tag, error)
	//GetTagsByFieldID
	//
	//  SELECT
	//      id, value, start, "end", line_id, field_id
	//  FROM
	//      tags
	//  WHERE
	//      field_id = ?
	GetTagsByFieldID(ctx context.Context, arg GetTagsByFieldIDParams) ([]*Tag, error)
	//GetURLByValue
	//
	//  SELECT
	//      id, value, html_id
	//  FROM
	//      urls
	//  WHERE
	//      value = ?
	GetURLByValue(ctx context.Context, arg GetURLByValueParams) (*Url, error)
	//InsertFile
	//
	//  INSERT INTO
	//      files (uri)
	//  VALUES
	//      (?) RETURNING id, uri, updated_at, created_at
	InsertFile(ctx context.Context, arg InsertFileParams) (*File, error)
	//InsertHTML
	//
	//  INSERT INTO
	//      htmls (value)
	//  VALUES
	//      (?) RETURNING id, value, updated_at, created_at
	InsertHTML(ctx context.Context, arg InsertHTMLParams) (*Html, error)
	//InsertLine
	//
	//  INSERT INTO
	//      lines (value)
	//  VALUES
	//      (?) RETURNING id, value, number
	InsertLine(ctx context.Context, arg InsertLineParams) (*InsertLineRow, error)
	//****************************************************************************
	//****************************************************************************
	//
	//
	//  /*
	//   ** selectors.sql
	//   ** Description: This file contains the SQLite queries for the selectors table
	//   ** Dialect: sqlite3
	//   */
	//  INSERT INTO
	//  	selectors (value, url_id, context, occurances)
	//  VALUES
	//  	(?, ?, ?, ?) RETURNING id, value, url_id, occurances, context
	InsertSelector(ctx context.Context, arg InsertSelectorParams) (*Selector, error)
	//InsertStruct
	//
	//  INSERT INTO
	//      structs (file_id, start_line_id, end_line_id, value)
	//  VALUES
	//      (?, ?, ?, ?) RETURNING id, value, url_id, start_line_id, end_line_id, file_id, context
	InsertStruct(ctx context.Context, arg InsertStructParams) (*Struct, error)
	//InsertTag
	//
	//  INSERT INTO
	//      tags (value, start, end, line_id, field_id)
	//  VALUES
	//      (?, ?, ?, ?, ?)
	InsertTag(ctx context.Context, arg InsertTagParams) error
	//InsertURL
	//
	//  INSERT INTO
	//      urls (value, html_id)
	//  VALUES
	//      (?, ?) RETURNING id, value, html_id
	InsertURL(ctx context.Context, arg InsertURLParams) (*Url, error)
	//ListAll
	//
	//  SELECT
	//      urls.id,
	//      urls.value,
	//      htmls.value as html,
	//      selectors.value as selector
	//  FROM
	//      urls
	//      JOIN htmls ON urls.html_id = htmls.id
	//      JOIN selectors ON urls.id = selectors.url_id
	ListAll(ctx context.Context) ([]*ListAllRow, error)
	//****************************************************************************
	//
	//  /*
	//  ** File: files.sql
	//  ** Description: This file contains the SQLite queries for the files table
	//  ** Dialect: sqlite3
	//  */
	//
	//  SELECT
	//      id, uri, updated_at, created_at
	//  from
	//      files
	ListFiles(ctx context.Context) ([]*File, error)
	//****************************************************************************
	//****************************************************************************
	//
	//
	//  /*
	//   ** File: htmls.sql
	//   ** Description: This file contains the SQLite queries for the htmls table
	//   ** Dialect: sqlite3
	//   */
	//  SELECT
	//      id, value, updated_at, created_at
	//  from
	//      htmls
	ListHTMLs(ctx context.Context) ([]*Html, error)
	//****************************************************************************
	//****************************************************************************
	//
	//
	//  /*
	//  ** File: lines.sql
	//  ** Description: This file contains the SQLite queries for the lines table
	//  ** Dialect: sqlite3
	//  */
	//
	//  SELECT
	//      id, file_id, value, number
	//  from
	//      lines
	ListLines(ctx context.Context) ([]*Line, error)
	//****************************************************************************
	//****************************************************************************
	//
	//
	//  /*
	//  ** File: structs.sql
	//  ** Description: This file contains the SQLite queries for the structs table
	//  ** Dialect: sqlite3
	//  */
	//
	//  SELECT
	//      id, value, url_id, start_line_id, end_line_id, file_id, context
	//  from
	//      structs
	ListStructs(ctx context.Context) ([]*Struct, error)
	//****************************************************************************
	//****************************************************************************
	//
	//
	//  /*
	//  ** File: tags.sql
	//  ** Description: This file contains the SQLite queries for the tags table
	//  ** Dialect: sqlite3
	//  */
	//
	//  SELECT
	//      id, value, start, "end", line_id, field_id
	//  from
	//      tags
	ListTags(ctx context.Context) ([]*Tag, error)
	//****************************************************************************
	//
	//  /*
	//   ** File: urls.sql
	//   ** Description: This file contains the SQLite queries for the urls table
	//   ** Dialect: sqlite3
	//   */
	//  SELECT
	//      id, value, html_id
	//  from
	//      urls
	ListURLs(ctx context.Context) ([]*Url, error)
	//UpdateFileByID
	//
	//  UPDATE
	//      files
	//  SET
	//      uri = ?
	//  WHERE
	//      id = ? RETURNING id, uri, updated_at, created_at
	UpdateFileByID(ctx context.Context, arg UpdateFileByIDParams) (*File, error)
	//UpdateHTMLByID
	//
	//  UPDATE
	//      htmls
	//  SET
	//      value = ?
	//  WHERE
	//      id = ? RETURNING id, value, updated_at, created_at
	UpdateHTMLByID(ctx context.Context, arg UpdateHTMLByIDParams) (*Html, error)
	//UpdateLineByID
	//
	//  UPDATE
	//      lines
	//  SET
	//      value = ?
	//  WHERE
	//      id = ? RETURNING id, value, number
	UpdateLineByID(ctx context.Context, arg UpdateLineByIDParams) (*UpdateLineByIDRow, error)
	//UpdateSelectorByID
	//
	//  UPDATE
	//  	selectors
	//  SET
	//  	value = ?,
	//  	url_id = ?,
	//  	context = ?,
	//  	occurances = ?
	//  WHERE
	//  	id = ?
	UpdateSelectorByID(ctx context.Context, arg UpdateSelectorByIDParams) error
	//UpdateStructByID
	//
	//  UPDATE
	//      structs
	//  SET
	//      value = ?,
	//      start_line_id = ?,
	//      end_line_id = ?
	//  WHERE
	//      id = ? RETURNING id, value, url_id, start_line_id, end_line_id, file_id, context
	UpdateStructByID(ctx context.Context, arg UpdateStructByIDParams) (*Struct, error)
	//UpdateTagByID
	//
	//  UPDATE
	//      tags
	//  SET
	//      value = ?,
	//      start = ?,
	//      end = ?,
	//      line_id = ?,
	//      field_id = ?
	//  WHERE
	//      id = ?
	UpdateTagByID(ctx context.Context, arg UpdateTagByIDParams) error
	//UpdateURL
	//
	//  UPDATE
	//      urls
	//  SET
	//      value = ?,
	//      html_id = ?
	//  WHERE
	//      id = ?
	UpdateURL(ctx context.Context, arg UpdateURLParams) error
	//UpsertURL
	//
	//  INSERT INTO
	//      urls (value, html_id)
	//  VALUES
	//      (?, ?)
	//  ON CONFLICT (value)
	//  DO UPDATE
	//      SET
	//          html_id = excluded.html_id RETURNING id, value, html_id
	UpsertURL(ctx context.Context, arg UpsertURLParams) (*Url, error)
}

var _ Querier = (*Queries)(nil)
