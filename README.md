# Wiki

## Disclaimer

The goal of this project is to create small personal wiki engine.

Uses Hugo as a static site generator.

Editing is made in separate page (template to be defined)

Go BE is used to manage MD files.

>Q: Why not use already existing one?
>
>A: because [fatal issue (en)](https://en.wikipedia.org/wiki/Not_invented_here)/[rus](http://lurkmore.to/%D0%A4%D0%B0%D1%82%D0%B0%D0%BB%D1%8C%D0%BD%D1%8B%D0%B9_%D0%BD%D0%B5%D0%B4%D0%BE%D1%81%D1%82%D0%B0%D1%82%D0%BE%D0%BA)

Everything here is not supposed to be working.

literally, it's draft of the draft.

## Docker manual

> Temprorary: run ```go build``` after pulling repo
To start docker build it first:

```
docker build --tag=wiki .
```

Then run it:

```
docker run -d -p YOUR_PORT_HERE:1337/tcp -v PATH_TO_YOUR_WEB_SERVER_FOLDER:/app/public/ wiki:latest
example:
docker run -d -p 1337:1337/tcp -v /usr/www/html/wiki2:/app/public/ wiki:latest
```
