
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
