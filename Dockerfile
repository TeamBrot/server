FROM golang:1.15

WORKDIR /go/src/app

RUN go get -v github.com/gorilla/websocket

COPY . .
RUN go build -o server .

CMD [ "./server" ]

