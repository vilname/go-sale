# Стадия сборки
FROM golang:1.22-alpine as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY .. .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Run stage
FROM alpine:3.19.1
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main /app/.env ./
COPY --from=builder /app/migration ./migration
COPY --from=builder /app/docs ./docs

CMD ["./main"]