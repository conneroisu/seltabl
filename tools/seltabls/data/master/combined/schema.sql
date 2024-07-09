
/*
* File: field.sql
** Description: This file contains the SQLite schema for the field table
** Dialect: sqlite3
*/
/*****************************************************************************/

/* Table: field */
CREATE TABLE IF NOT EXISTS field (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	struct_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	line_id INTEGER NOT NULL,
	UNIQUE (name),
	FOREIGN KEY (struct_id) REFERENCES structs (id)
);

/* Index: idx_fields_name */
CREATE INDEX IF NOT EXISTS idx_fields_name ON field (name);

/*****************************************************************************/

/*
** File: files.sql
** Description: This file contains the SQLite schema for the files table
** Dialect: sqlite3
*/
/******************************************************************************/

/* Table: files */
CREATE TABLE IF NOT EXISTS files (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	uri TEXT NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (id, uri)
);

/* Index: idx_files_file */
CREATE INDEX IF NOT EXISTS idx_files_uri ON files (uri);

/******************************************************************************/

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
** File: schemas/lines.sql
** Description: This file contains the SQLite schema for the lines table
** Dialect: sqlite3
*/
/******************************************************************************/

/* Table: lines */
CREATE TABLE IF NOT EXISTS lines (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	file_id INTEGER NOT NULL,
	value TEXT NOT NULL,
	number INTEGER NOT NULL,
	UNIQUE (id, value),
	FOREIGN KEY (file_id) REFERENCES files (id)
);

/* Index: idx_lines_number */
CREATE INDEX IF NOT EXISTS idx_lines_number ON lines (number);

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
	occurances INTEGER NOT NULL,
	context TEXT NOT NULL,
	UNIQUE (id),
	FOREIGN KEY (url_id) REFERENCES urls (id)
);

/* Index: idx_selectors_url_id */
CREATE INDEX IF NOT EXISTS idx_selectors_url ON selectors (url_id);

/******************************************************************************/

/*
** File: structs.sql
** Description: This file contains the SQLite schema for the structs table
** Dialect: sqlite3
*/
/******************************************************************************/

/* Table: structs */
CREATE TABLE IF NOT EXISTS structs (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	value TEXT NOT NULL,
	url_id INTEGER NOT NULL,
	start_line_id INTEGER NOT NULL,
	end_line_id INTEGER NOT NULL,
	file_id INTEGER NOT NULL,
	context TEXT NOT NULL,
	UNIQUE (id),
	FOREIGN KEY (url_id) REFERENCES urls (id),
	FOREIGN KEY (file_id) REFERENCES files (id)
);

/* 
** File: tags.sql
** Description: This file contains the SQLite schema for the tags table
** Dialect: sqlite3
*/
/******************************************************************************/

/* Table: tags */
CREATE TABLE IF NOT EXISTS tags (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	value TEXT NOT NULL,
	start INTEGER NOT NULL,
	end INTEGER NOT NULL,
	line_id INTEGER NOT NULL,
	field_id INTEGER NOT NULL,
	UNIQUE (id),
	FOREIGN KEY (field_id) REFERENCES fields (id),
	FOREIGN KEY (line_id) REFERENCES lines (id)
);

/* Index: idx_tags_field_id */
CREATE INDEX IF NOT EXISTS idx_tags_field_id ON tags (field_id);

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
