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
