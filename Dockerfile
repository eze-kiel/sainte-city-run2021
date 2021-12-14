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

RUN setcap cap_net_raw+ep server

FROM gcr.io/distroless/base-debian10

WORKDIR /app

COPY --from=builder /build/server .
COPY --from=builder /build/views .

EXPOSE 8080

ENTRYPOINT [ "/app/server" ]