/* 
 ** File: selectors.sql
 ** Description: This file contains the SQLite schema for the selectors table
 ** Dialect: sqlite3
 */
/******************************************************************************/
/* Table: selectors */
CREATE TABLE IF NOT EXISTS selectors (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	value TEXT NOT NULL,
	url_id INTEGER NOT NULL,
	context TEXT NOT NULL,
	UNIQUE (id),
	FOREIGN KEY (url_id) REFERENCES urls (id)
);

/* Index: idx_selectors_url_id */
CREATE INDEX IF NOT EXISTS idx_selectors_url ON selectors (url_id);

/******************************************************************************/
