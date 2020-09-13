FROM golang:latest
WORKDIR /app

RUN apt-get update \
 && apt-get install -y \
    mariadb-client \
    less \
    git

COPY go.mod .
COPY go.sum .
RUN go mod download

CMD /bin/sh -c "while sleep 1000; do :; done"
