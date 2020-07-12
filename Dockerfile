FROM golang:1.14 as builder
ENV GO111MODULE=on \
        CGO_ENABLED=0 \
        GOOS=linux \
        GOARCH=amd64
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY main.go .
RUN go build -o echo-server .


FROM scratch
COPY --from=builder /build/echo-server /
ENV PORT 8080
EXPOSE 8080
ENTRYPOINT ["/echo-server"]
