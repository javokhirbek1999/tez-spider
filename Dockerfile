FROM golang:1.19.2-alpine3.15


WORKDIR /app


COPY go.mod ./
COPY go.sum ./
RUN go mod download


COPY . .

RUN go build -o main .

EXPOSE 4000

CMD ["./main"]


