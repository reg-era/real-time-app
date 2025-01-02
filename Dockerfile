FROM golang:1.23.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go get github.com/mattn/go-sqlite3
RUN go get github.com/gofrs/uuid/v5
RUN go get golang.org/x/crypto@v0.28.0

RUN go mod download

RUN go mod tidy

COPY . .

RUN go build -o forum cmd/main.go

FROM ubuntu:latest

WORKDIR /app

#Set enviroment variables
ENV PORT=8080
ENV DB_PATH="db/data.db"

RUN mkdir db web


COPY ./web/templates ./web/templates
COPY ./web/assets  ./web/assets

COPY --from=builder /app/forum .

EXPOSE ${PORT}

ENTRYPOINT [ "./forum" ]
