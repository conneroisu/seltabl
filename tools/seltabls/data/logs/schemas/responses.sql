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
