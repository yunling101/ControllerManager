ARG GOLANG_BUILDER="1.21.9-alpine"

FROM golang:${GOLANG_BUILDER} as go_builder

ENV GOPROXY=https://goproxy.cn
ENV GO111MODULE=on
ENV GOPATH=/go

ARG NAME="alertmanager-plug"
ARG OUTPUT="bin/${NAME}"
ARG BUILDER="cmd/${NAME}/${NAME}.go"

ARG ARCH="amd64"
ARG OS="linux"

WORKDIR /go/src/github.com/yunling101/ControllerManager
COPY . .

RUN CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -ldflags "-w -s" -o ${OUTPUT} ${BUILDER}

FROM alpine:3.19

WORKDIR /opt
ENV LANG zh_CN.UTF-8
ENV SECRET_KEY Need to specify

COPY --from=go_builder /go/src/github.com/yunling101/ControllerManager/bin/alertmanager-plug .
USER nobody

COPY ./cmd/alertmanager-plug/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["--alertmanager-base-dir", "/opt/alertmanager"]

EXPOSE 9098