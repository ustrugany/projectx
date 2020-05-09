ARG GO_VERSION=1.14.2
FROM golang:${GO_VERSION} as build-env

WORKDIR /go/src/github.com/ustrugany/projectx/

COPY api/vendor api/vendor
COPY api .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o projectx ./cmd/server/main.go
RUN chmod +x ./projectx

ARG PORT=8090
EXPOSE ${PORT}

ENTRYPOINT ["./projectx"]
CMD ""