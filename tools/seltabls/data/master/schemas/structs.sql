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
