CREATE USER IF NOT EXISTS craig WITH PASSWORD 'cockroach';

CREATE DATABASE IF NOT EXISTS checker;

GRANT ALL ON DATABASE checker TO craig;