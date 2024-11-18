FROM mysql:8.0

COPY ./initdb.d/01-create_db.sql /docker-entrypoint-initdb.d/
COPY ./my.cnf /etc/mysql/conf.d/my.cnf
