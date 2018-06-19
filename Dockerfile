FROM alpine:latest

WORKDIR /app
VOLUME /app
COPY startup.sh /startup.sh

RUN apk add --update --no-cache mysql mysql-client && rm -f /var/cache/apk/*
COPY my.cnf /etc/mysql/my.cnf


EXPOSE 3306
CMD ["mysqld"]