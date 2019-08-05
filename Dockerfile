FROM ubuntu:18.04

LABEL maintainer="a.diachenko@outlook.com"

WORKDIR /app

EXPOSE 1337

RUN apt-get update 
RUN apt-get upgrade -y
RUN apt-get install hugo -y
RUN mkdir db

COPY . /app

CMD [ "./wiki" ]

