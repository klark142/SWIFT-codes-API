FROM golang:1.22-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o swift-codes ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o swift-codes-import ./cmd/import

FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache postgresql-client

COPY --from=builder /app/swift-codes .
COPY --from=builder /app/swift-codes-import .

COPY data/swiftcodes_data.csv data/swiftcodes_data.csv

COPY entrypoint.sh .

RUN chmod +x entrypoint.sh

EXPOSE 8080

CMD ["./entrypoint.sh"]
