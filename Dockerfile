# builder image
FROM golang:1.17.1-alpine3.14 AS pipeline-example-builder

RUN mkdir /build

COPY *.go /build/
COPY go.mod /build/
COPY go.sum /build/

WORKDIR /build

RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .


# final image
FROM alpine:3.14 AS pipeline-example

COPY --from=pipeline-example-builder /build/main .

CMD ./main --help