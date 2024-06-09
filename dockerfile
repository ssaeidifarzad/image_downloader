FROM golang:1.21.0

WORKDIR /app

COPY . .


RUN go mod download

RUN go build -o main .

CMD ["./cmd/main"]

