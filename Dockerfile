FROM golang:1.12 as build
WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

FROM gcr.io/distroless/base
COPY --from=build /go/bin/app /
COPY server.crt /
COPY server.key /
CMD ["/app"]
