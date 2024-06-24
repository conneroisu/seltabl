-- Table to store log levels
CREATE TABLE IF NOT EXISTS log_levels (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    level_name TEXT NOT NULL UNIQUE
);

-- Table to store log messages
CREATE TABLE IF NOT EXISTS logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    level_id INTEGER,
    message TEXT NOT NULL,
    source TEXT,
    user_id INTEGER,
    FOREIGN KEY (level_id) REFERENCES log_levels (id)
);

-- Table to store errors
CREATE TABLE IF NOT EXISTS errors (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    error_code INTEGER,
    error_message TEXT NOT NULL,
    stack_trace TEXT,
    log_id INTEGER,
    FOREIGN KEY (log_id) REFERENCES logs (id)
);

-- Table to store user interactions
CREATE TABLE IF NOT EXISTS user_interactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER,
    interaction_type TEXT NOT NULL,
    details TEXT,
    log_id INTEGER,
    FOREIGN KEY (log_id) REFERENCES logs (id)
);

-- Insert common log levels
INSERT INTO log_levels (level_name) VALUES ('DEBUG');
INSERT INTO log_levels (level_name) VALUES ('INFO');
INSERT INTO log_levels (level_name) VALUES ('WARN');
INSERT INTO log_levels (level_name) VALUES ('ERROR');
INSERT INTO log_levels (level_name) VALUES ('FATAL');
