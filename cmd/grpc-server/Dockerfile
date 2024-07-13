FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /grpc-server ./cmd/grpc-server

FROM alpine:3.14

COPY --from=builder /grpc-server /grpc-server

EXPOSE 50051

CMD ["/grpc-server"]