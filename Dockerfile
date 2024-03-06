FROM golang:1.22-alpine3.19 as builder
WORKDIR /app
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o api ./cmd/v3/main.go

FROM scratch
WORKDIR /
COPY --from=builder /app/api ./
EXPOSE 3000
ENTRYPOINT ["./api"]