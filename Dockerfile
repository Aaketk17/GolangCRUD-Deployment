FROM golang:1.21.3-alpine3.18

WORKDIR /app

COPY . /app

RUN go get
RUN go build -o bin .

ENTRYPOINT [ "/app/bin" ]