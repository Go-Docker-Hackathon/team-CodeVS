FROM golang:1.3
RUN go get github.com/cyrus-and/gdb
RUN go get github.com/go-martini/martini
RUN go get github.com/gorilla/websocket
ADD . /code
WORKDIR /code
EXPOSE 8080
CMD ["go", "run", "main.go"]
