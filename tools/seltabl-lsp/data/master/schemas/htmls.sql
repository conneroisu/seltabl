/*
 ** File: htmls.sql
 ** Description: This file contains the SQLite schema for the htmls table
 ** Dialect: sqlite3
 */
/******************************************************************************/
/* Table: htmls */
CREATE TABLE IF NOT EXISTS htmls (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	value TEXT NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (id, value)
);

/* Index: idx_htmls_html */
CREATE INDEX IF NOT EXISTS idx_htmls_html ON htmls (value);

/******************************************************************************/
