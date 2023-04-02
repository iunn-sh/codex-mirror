FROM golang:alpine3.17

WORKDIR /app
ADD . /app

#RUN ls -al /app
RUN cd /app
RUN go mod init github.com/iunn-sh/codex-mirror
RUN go build -o app

ENTRYPOINT ./app