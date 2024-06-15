
/* 
** File: selectors.sql
** Description: This file contains the SQLite schema for the selectors table
** Dialect: sqlite3
*/
/******************************************************************************/

/* Table: selectors */
CREATE TABLE IF NOT EXISTS selectors (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	selector TEXT NOT NULL,
	url TEXT NOT NULL,
	context TEXT NOT NULL,
	UNIQUE (selector, url)
);

/* Index: idx_selectors_id */
CREATE INDEX IF NOT EXISTS idx_selectors_url ON selectors (url);

/* Index: idx_selectors_selector */
CREATE INDEX IF NOT EXISTS idx_selectors_selector ON selectors (selector);

/******************************************************************************/
