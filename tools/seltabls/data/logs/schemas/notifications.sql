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
