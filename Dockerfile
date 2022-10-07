FROM golang:1.19.2-alpine3.15 AS builder


WORKDIR /app


COPY go.mod ./
COPY go.sum ./
RUN go mod download


COPY . .

RUN go build -o main .


EXPOSE 9090

CMD ["./main"]


