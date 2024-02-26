FROM golang:1.21 AS go

RUN mkdir -p /gwarr

COPY ./dist/gwarr /gwarr/gwarr

ENTRYPOINT ["/gwarr/gwarr"]
