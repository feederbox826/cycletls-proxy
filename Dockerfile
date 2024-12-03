FROM golang:alpine AS builder
WORKDIR /build
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build

FROM scratch
EXPOSE 8080
COPY --from=builder /build/cycletls-proxy /cycletls-proxy

ENTRYPOINT ["/cycletls-proxy"]