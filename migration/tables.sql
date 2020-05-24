CREATE TABLE IF NOT EXISTS domains (
    id UUID PRIMARY KEY,
    servers_changed BOOLEAN NOT NULL,
    ssl_grade TEXT NOT NULL,
    previous_ssl_grade TEXT NOT NULL,
    logo TEXT NOT NULL,
    title TEXT NOT NULL,
    is_down BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS servers (
    id UUID PRIMARY KEY,
    domain_id UUID NOT NULL REFERENCES domains (id) ON DELETE CASCADE,
    address TEXT NOT NULL,
    ssl_grade TEXT NOT NULL,
    country TEXT NOT NULL,
    owner TEXT NOT NULL
);