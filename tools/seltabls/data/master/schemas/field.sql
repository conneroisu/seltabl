/*
** File: field.sql
** Description: This file contains the SQLite schema for the field table
** Dialect: sqlite3
*/
/******************************************************************************/

/* Table: field */
CREATE TABLE IF NOT EXISTS field (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	line INTEGER NOT NULL,
	start INTEGER NOT NULL,
	end INTEGER NOT NULL,
	struct_id INTEGER NOT NULL,
	UNIQUE (name),
	FOREIGN KEY (struct_id) REFERENCES structs (id)
);

/* Index: idx_fields_name */
CREATE INDEX IF NOT EXISTS idx_fields_name ON field (name);

/* Index: idx_fields_type */
CREATE INDEX IF NOT EXISTS idx_fields_type ON field (type);

/******************************************************************************/
