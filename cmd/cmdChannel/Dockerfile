ARG GOLANG_BUILDER="1.21.9-alpine"

FROM golang:${GOLANG_BUILDER} as go_builder

ENV GOPROXY=https://goproxy.cn
ENV GO111MODULE=on
ENV GOPATH=/go

ARG NAME="cmdChannel"
ARG OUTPUT="bin/${NAME}"
ARG BUILDER="cmd/${NAME}/${NAME}.go"

ARG ARCH="amd64"
ARG OS="linux"

WORKDIR /go/src/github.com/yunling101/ControllerManager
COPY . .

RUN CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -ldflags "-w -s" -o ${OUTPUT} ${BUILDER}

FROM scratch

WORKDIR /opt
COPY --from=go_builder /go/src/github.com/yunling101/ControllerManager/bin/cmdChannel .

CMD ["/opt/cmdChannel"]