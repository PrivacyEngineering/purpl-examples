FROM golang:1.17-alpine

WORKDIR /app

COPY . /app

RUN go build -o /app/server/pingpong-server /app/playground/pingpong/server/server_two.go

CMD ["/app/server/pingpong-server"]
