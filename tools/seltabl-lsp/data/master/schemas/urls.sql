/*
 ** File: urls.sql
 ** Description: This file contains the SQLite schema for the urls table
 ** Dialect: sqlite3
 */
/******************************************************************************/
/* Table: urls */
CREATE TABLE IF NOT EXISTS urls (
	id INTEGER UNIQUE NOT NULL PRIMARY KEY AUTOINCREMENT,
	value TEXT UNIQUE NOT NULL,
	html_id INTEGER NOT NULL,
	UNIQUE (id, value),
	FOREIGN KEY (html_id) REFERENCES htmls (id)
);

/* Index: idx_urls_html_id */
CREATE INDEX IF NOT EXISTS idx_urls_html_id ON urls (html_id);

/******************************************************************************/
