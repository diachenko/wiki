FROM ubuntu:18.04

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
#RUN mkdir md

#VOLUME /var/www/html/wiki2

COPY . /app
#RUN hugo
#RUN go build
CMD [ "./wiki" ]

