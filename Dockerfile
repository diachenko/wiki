FROM golang:latest

LABEL maintainer="a.diachenko@outlook.com"

WORKDIR /app

EXPOSE 1337

RUN apt-get update 
RUN apt-get upgrade -y
RUN apt-get install hugo -y
RUN mkdir db
#RUN go get github.com/gorilla/mux
#RUN go get github.com/boltdb/bolt
#RUN go get github.com/gomarkdown/markdown
RUN mkdir md


COPY . /app

#RUN go build
CMD [ "./wiki" ]

