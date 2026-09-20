FROM golang:1.27-alpine3.24

WORKDIR /app
ADD . /app

RUN cd /app
RUN go build -o app

ENTRYPOINT ./app

# TODO: minimize container size with multi-stage