FROM golang:1.23.1-alpine3.20 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

FROM alpine:3.19

WORKDIR /app

COPY ./db /app/db
COPY --from=builder /app/main ./

EXPOSE 8080

CMD ["./main"]
