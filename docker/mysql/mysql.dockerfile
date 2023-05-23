FROM mysql:latest

ENV MYSQL_ROOT_PASSWORD=root
ENV MYSQL_DATABASE=coursefan

COPY ./init.sql /docker-entrypoint-initdb.d/init.sql