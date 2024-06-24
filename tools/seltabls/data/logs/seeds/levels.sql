/*
** File: seeds/levels.sql
** Description: Insert common log levels
** Dialect: sqlite
*/
/******************************************************************************/

-- Insert common log levels
INSERT OR IGNORE INTO log_levels (name) VALUES ('DEBUG');
INSERT OR IGNORE INTO log_levels (name) VALUES ('INFO');
INSERT OR IGNORE INTO log_levels (name) VALUES ('WARN');
INSERT OR IGNORE INTO log_levels (name) VALUES ('ERROR');
INSERT OR IGNORE INTO log_levels (name) VALUES ('FATAL');

/******************************************************************************/
