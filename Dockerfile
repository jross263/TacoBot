FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o tacobot .

FROM alpine:latest AS runner
WORKDIR /root/
COPY --from=builder /app/tacobot .
CMD ["./tacobot"]