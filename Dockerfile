FROM golang:1.22 AS builder

ENV GO111MODULE=on

COPY . /build
WORKDIR /build

RUN make all && \
    mkdir -p /app && \
    cp bin/* /app

FROM ubuntu:24.04

RUN apt-get update && apt-get install -y ca-certificates && \
    curl -sSfL http://mirrors.hust.edu.cn/get | sh -s -- deploy ubuntu && \
    apt-get update && apt-get install -y git cron make

ENV TZ=Asia/Shanghai
ENV PATH=/app:$PATH

COPY ./workflow /workflow
 
ENV CFG_FILE=/config/config.json
ENV APP_BIN=/app

RUN /workflow/gen_crontab.sh -o /proc/1/fd/1 -r /data/rec > /tmp/update.cron && \
    crontab /tmp/update.cron && \
    rm /tmp/update.cron

COPY --from=builder /app /app

CMD ["cron", "-f"]


