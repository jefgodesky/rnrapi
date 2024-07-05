FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o rnrapi .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/rnrapi .
EXPOSE 3000

CMD ["./rnrapi"]
