FROM golang:1.14-rc-buster

COPY .env.example data-update.sh go.mod go.sum main.go /app/

WORKDIR /app

COPY assets /app/assets
COPY enums /app/enums
COPY src /app/src

RUN go build
ENTRYPOINT [ "bash", "data-update.sh" ]