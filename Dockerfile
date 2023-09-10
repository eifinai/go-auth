FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o bin ./main.go

FROM alpine:3.17.3

WORKDIR /app

COPY --from=builder /app/bin ./

ENTRYPOINT ["/app/bin"]
