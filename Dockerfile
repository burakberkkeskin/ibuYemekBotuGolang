FROM golang:1.18-alpine3.15 as builder
RUN apk --no-cache add ca-certificates
RUN apk add git
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /app/ibuYemekBot

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/ibuYemekBot /
ENTRYPOINT ["/ibuYemekBot"]