CREATE TABLE IF NOT EXISTS domains (
    id UUID PRIMARY KEY,
    servers_changed BOOLEAN NOT NULL,
    ssl_grade TEXT NOT NULL,
    previous_ssl_grade TEXT NOT NULL,
    logo TEXT NOT NULL,
    title TEXT NOT NULL,
    is_down BOOLEAN NOT NULL,
    host TEXT NOT NULL
);