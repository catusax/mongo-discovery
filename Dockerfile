FROM golang:alpine AS builder
ENV CGO_ENABLED=0 GOOS=linux
RUN apk --update --no-cache add ca-certificates gcc libtool make musl-dev protoc
WORKDIR /go/src/mongo-discovery
COPY . .
RUN make tidy build

FROM scratch
ENV CONTAINER=docker
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /go/src/mongo-discovery/mongo-discovery /mongo-discovery
ENTRYPOINT ["/mongo-discovery"]
CMD []