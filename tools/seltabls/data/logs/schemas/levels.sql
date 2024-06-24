/*
** File: schemas/levels.sql
** Description: This file contains the schema for the levels table.
** Dialect: sqlite
*/
/******************************************************************************/

-- Table to store log levels
CREATE TABLE IF NOT EXISTS log_levels (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

/* Index: idx_log_levels_level_name */
CREATE INDEX IF NOT EXISTS idx_log_levels_name ON log_levels (name);

/******************************************************************************/
