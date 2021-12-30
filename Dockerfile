FROM happyxhw/alpine:3.15.0

LABEL MAINTAINER="happyxhw"

RUN mkdir -p /app/config
WORKDIR /app

ADD iself iself

ENV TZ=Asia/Shanghai

ENTRYPOINT [ "/app/iself" ]