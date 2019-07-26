FROM ubuntu:18.04

LABEL maintainer="a.diachenko@outlook.com"

RUN apt-get update 
RUN apt-get upgrade -y
RUN apt-get install hugo -y
RUN apt-get install golang

