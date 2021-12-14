FROM devopsworks/golang-upx:1.16 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build \
    -o server . && \
    strip server && \
    /usr/local/bin/upx -9 server

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=builder /build/ .

EXPOSE 8080

ENTRYPOINT [ "/server" ]