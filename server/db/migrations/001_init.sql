CREATE TABLE IF NOT EXISTS servers (
    name        TEXT PRIMARY KEY,
    version     TEXT NOT NULL DEFAULT '',
    title       TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    data        TEXT NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE VIRTUAL TABLE IF NOT EXISTS servers_fts USING fts5(
    name, title, description,
    content='servers',
    content_rowid='rowid'
);

-- Triggers to keep FTS index in sync
CREATE TRIGGER IF NOT EXISTS servers_ai AFTER INSERT ON servers BEGIN
    INSERT INTO servers_fts(rowid, name, title, description)
    VALUES (new.rowid, new.name, new.title, new.description);
END;

CREATE TRIGGER IF NOT EXISTS servers_ad AFTER DELETE ON servers BEGIN
    INSERT INTO servers_fts(servers_fts, rowid, name, title, description)
    VALUES ('delete', old.rowid, old.name, old.title, old.description);
END;

CREATE TRIGGER IF NOT EXISTS servers_au AFTER UPDATE ON servers BEGIN
    INSERT INTO servers_fts(servers_fts, rowid, name, title, description)
    VALUES ('delete', old.rowid, old.name, old.title, old.description);
    INSERT INTO servers_fts(rowid, name, title, description)
    VALUES (new.rowid, new.name, new.title, new.description);
END;

CREATE TABLE IF NOT EXISTS sync_log (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    synced_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    server_count INTEGER NOT NULL DEFAULT 0,
    duration_ms  INTEGER NOT NULL DEFAULT 0
);
