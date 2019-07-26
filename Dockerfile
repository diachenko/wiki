FROM golang:latest

LABEL maintainer="a.diachenko@outlook.com"

WORKDIR /app

EXPOSE 1337

RUN apt update 
RUN apt upgrade -y
RUN apt install hugo -y
RUN mkdir db
RUN go get github.com/gorilla/mux
RUN go get github.com/boltdb/bolt
RUN mkdir md


COPY . /app

RUN go build
CMD [ "./wiki" ]

