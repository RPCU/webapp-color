FROM golang:1.23 AS builder
ENV APP_DIR=/src
WORKDIR ${APP_DIR}
COPY . ${APP_DIR}

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${APP_DIR}/app .

FROM scratch
LABEL org.opencontainers.image.authors="alanamoyel06@gmail.com"
ENV APP_DIR=/src

COPY --from=builder ${APP_DIR} ${APP_DIR}

WORKDIR ${APP_DIR}
CMD [ "./app" ]
