FROM golang:1.19-alpine

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN mkdir -p /tmp/gopath && chmod -R 777 /tmp/gopath
ENV GOPATH=/tmp/gopath

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
