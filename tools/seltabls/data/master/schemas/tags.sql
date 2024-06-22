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
