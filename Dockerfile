FROM golang as build-stage
WORKDIR /go/
COPY ./ /go/src
RUN cd /go/src && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o config-reloader

FROM alpine
COPY --from=build-stage /go/src/config-reloader /
EXPOSE 8080
CMD ["/config-reloader"]