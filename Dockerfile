ARG GO_VERSION=1.14.2
FROM golang:${GO_VERSION} as build-env

WORKDIR /go/src/github.com/ustrugany/projectx/

COPY api/vendor api/vendor
COPY api .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o projectx ./cmd/server/main.go
RUN chmod +x ./projectx

ENTRYPOINT ["./projectx"]
CMD ""

#ENTRYPOINT ""
#CMD ["./projectx"]

#FROM alpine:3.7
#WORKDIR /app
#COPY --from=build-env /go/src/github.com/ustrugany/projectx/projectx ./booking/docs
#COPY --from=build-env /go/src/github.com/marcusolsson/goddd/tracking/docs ./tracking/docs
#COPY --from=build-env /go/src/github.com/marcusolsson/goddd/handling/docs ./handling/docs
#COPY --from=build-env /go/src/github.com/marcusolsson/goddd/goapp .
#EXPOSE 8080
#ENTRYPOINT ["./goapp"]