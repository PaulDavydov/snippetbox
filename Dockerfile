FROM golang:1.23.2-alpine

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .
RUN go build -o main ./cmd/web/

EXPOSE 4000

CMD ["./main"]
