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
