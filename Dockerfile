FROM ubuntu:14.04

MAINTAINER Jose Leon

ADD localhost.crt /app/
ADD server.key /app/
ADD bootstrap.sh /
ADD server /app/

ENTRYPOINT ["/bootstrap.sh"]

