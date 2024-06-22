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
