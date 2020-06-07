.PHONY: db tables

db:
	cockroach start-single-node --certs-dir=certs --http-addr=localhost:9090

db-config:
	cockroach sql --certs-dir=certs --host=localhost:26257 < migration/config_db.sql

tables:
	cockroach sql --certs-dir=certs --user=craig --host=localhost:26257 --database=checker < migration/tables.sql
	

