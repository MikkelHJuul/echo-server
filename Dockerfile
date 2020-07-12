FROM golang:1.14 as builder

RUN go build github.com/MikkelHJuul/echo-server


FROM scratch
COPY --from=builder echo-server /
ENV PORT 8080
EXPOSE 8080
ENTRYPOINT ["/echo-server"]
