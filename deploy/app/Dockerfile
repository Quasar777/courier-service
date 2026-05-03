FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o service ./cmd/app/main.go

FROM curlimages/curl:8.17.0
WORKDIR /
COPY --from=builder /app/service /courier-service
EXPOSE 8080
USER nobody:nobody
ENTRYPOINT ["/courier-service"]