FROM golang:1.20.2-alpine3.17

WORKDIR /app
ADD . /app

RUN cd /app
RUN go build -o app

ENTRYPOINT ./app

# TODO: minimize container size with multi-stage