FROM golang:1.23-apline AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /proxy

FROM alpine:latest
WORKDIR /app
COPY --from=builder /proxy /app/proxy
EXPOSE 80 443
CMD ["/app/proxy"]