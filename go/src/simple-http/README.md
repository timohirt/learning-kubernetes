## Simple HTTP microservice

This service exposes a `/health` endpoint via HTTP on port `30000`.

It compiles and runs on your machine as easy as this:

```bash
go build -o main ./src/*
./main
2019/01/25 18:27:43 Running server!
```

Building a docker image is possible as well.

```bash
docker build -t simple-http:current .
```
