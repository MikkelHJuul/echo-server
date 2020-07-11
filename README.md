# Echo Server

A very simple HTTP echo server, forked from jmalloc/echo-server

## Behavior
- GET Requests to any other URL will return the request headers and body
- POST with body will log **only** the body

## Configuration

- The `PORT` environment variable sets the server port, which defaults to `8080`
- Set the `LOG_HTTP_BODY` environment variable to dump request bodies to `/dev/stdout`

## Running the server

The examples below show a few different ways of running the server with the HTTP
server bound to a custom TCP port of `10000`.

### Running locally

```
GO111MODULE=off go get -u github.com/MikkelHJuul/echo-server/...
PORT=10000 echo-server
```

### Running under Docker

To run as a container:

```
docker run --detach -p 10000:8080 jmalloc/echo-server
```

To run as a service:

```
docker service create --publish 10000:8080 jmalloc/echo-server
```
