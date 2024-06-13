/*
** File: urls.sql
** Description: This file contains the schema for the urls table.
** Dialect: sqlite3
*/
/******************************************************************************/

/* Table: urls */
CREATE TABLE IF NOT EXISTS urls (
	id TEXT NOT NULL PRIMARY KEY DEFAULT (uuid()),
	url TEXT NOT NULL,
	encoding TEXT NOT NULL,
	UNIQUE (url)
);

/* Index: idx_urls_id */
CREATE INDEX IF NOT EXISTS idx_urls_url ON urls (url);

/******************************************************************************/
