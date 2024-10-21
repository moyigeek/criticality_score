FROM golang:1.22 AS builder

ENV GO111MODULE=on

# predownload go modules to speed up incremental builds
RUN mkdir -p /build
WORKDIR /build
COPY go.mod go.sum /build/
RUN go mod download && rm -rf /build

# build the app
COPY . /build

RUN make all && \
    mkdir -p /app && \
    cp bin/* /app

# final stage
FROM ubuntu:24.04

RUN apt-get update && apt-get install -y ca-certificates && \
    curl -sSfL http://mirrors.hust.edu.cn/get | sh -s -- deploy ubuntu && \
    apt-get update && apt-get install -y git cron make

ENV TZ=Asia/Shanghai
ENV PATH=/app:$PATH

COPY ./workflow /workflow
 
ENV CFG_FILE=/config/config.json
ENV APP_BIN=/app
ENV STORAGE_DIR=/storage

RUN /workflow/gen_crontab.sh -o /data/log/app.log -r /data/rec > /tmp/update.cron && \
    crontab /tmp/update.cron && \
    rm /tmp/update.cron

RUN echo '#!/bin/bash' > /gitlink.sh && \
    echo 'APP_BIN=/app CFG_FILE=/config/config.json /workflow/update.sh -C /data/rec gitlink' >> /gitlink.sh && \
    chmod +x /gitlink.sh

COPY --from=builder /app /app

CMD ["cron", "-f"]


