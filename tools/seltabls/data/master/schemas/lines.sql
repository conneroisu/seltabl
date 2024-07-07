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
