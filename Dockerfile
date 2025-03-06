FROM golang:1.23.2-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY .env .

RUN go mod tidy
RUN go build -o main ./cmd/web/

EXPOSE 4000

CMD ["./main"]
