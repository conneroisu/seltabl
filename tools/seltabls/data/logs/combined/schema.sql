
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

/*
** File: schemas/logs.sql
** Description: This file contains the SQLite schema for the logs table
** Dialect: sqlite
*/
/******************************************************************************/

/* Table: logs */
CREATE TABLE IF NOT EXISTS logs (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	value TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	request_id INTEGER,
	response_id INTEGER,
	notification_id INTEGER,
	UNIQUE (id, value),
	FOREIGN KEY (request_id) REFERENCES requests (id),
	FOREIGN KEY (response_id) REFERENCES responses (id),
	FOREIGN KEY (notification_id) REFERENCES notifications (id)
);

/* Index: idx_logs_value */
CREATE INDEX IF NOT EXISTS idx_logs_value ON logs (value);

/******************************************************************************/

/*
** File: schemas/notifications.sql
** Description: This file contains the SQLite schema for the notifications table
** Dialect: sqlite
*/
/******************************************************************************/

/* Table: notifications */
CREATE TABLE IF NOT EXISTS notifications (
    id         INTEGER   NOT NULL PRIMARY KEY AUTOINCREMENT,
    method     TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (id, value)
);

/* Index: idx_notifications_value */
CREATE INDEX IF NOT EXISTS idx_notifications_method ON notifications (method);

/******************************************************************************/

/*
** File: schemas/requests.sql
** Description: This file contains the schema for the requests table.
** Dialect: sqlite
*/
/******************************************************************************/

/* Table: requests */
CREATE TABLE IF NOT EXISTS requests (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    rpc_method TEXT NOT NULL,
    rpc_id INTEGER NOT NULL
)

/* Index: idx_requests_rpc_method */    
CREATE INDEX IF NOT EXISTS idx_requests_rpc_method ON requests (rpc_method);

/******************************************************************************/

/*
** File: schemas/responses.sql
** Description: Table to store log messages
** Dialect: sqlite
*/
/******************************************************************************/

/* Table: responses */
CREATE TABLE IF NOT EXISTS responses (      
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    rpc TEXT NOT NULL,
    rpc_id INTEGER NOT NULL,
    result TEXT,
    error TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (rpc_id) REFERENCES rpcs (id)
);

/* Index: idx_responses_timestamp */
CREATE INDEX IF NOT EXISTS idx_responses_timestamp ON responses (timestamp);

/******************************************************************************/
