/*
** File: schemas/logs.sql
** Description: This file contains the SQLite schema for the logs table
** Dialect: sqlite
*/
/******************************************************************************/

/* Table: logs */
CREATE TABLE IF NOT EXISTS logs (
	id              INTEGER   NOT NULL PRIMARY KEY AUTOINCREMENT,
	request_id      INTEGER,
	response_id     INTEGER,
	notification_id INTEGER,
	value           TEXT      NOT NULL,
	created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (id, value),
	FOREIGN KEY (request_id) REFERENCES requests (id),
	FOREIGN KEY (response_id) REFERENCES responses (id),
	FOREIGN KEY (notification_id) REFERENCES notifications (id)
);

/* Index: idx_logs_value */
CREATE INDEX IF NOT EXISTS idx_logs_value ON logs (value);

/******************************************************************************/
