FROM alpine:3.19

WORKDIR /opt

ARG VERSION="2.54.1"
ARG PACKAGE="prometheus-${VERSION}.linux-amd64"

RUN wget https://github.com/prometheus/prometheus/releases/download/v${VERSION}/${PACKAGE}.tar.gz \
    && tar -zxvf ${PACKAGE}.tar.gz \
    && rm -f ${PACKAGE}.tar.gz \
    && mv ${PACKAGE} prometheus \
    && mkdir -p /opt/prometheus/config \
    && mkdir -p /opt/prometheus/data \
    && mkdir -p /opt/prometheus/rules \
    && mv prometheus/prometheus.yml /opt/prometheus/config

WORKDIR /opt/prometheus

COPY ./third_party/prometheus/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
ENTRYPOINT ["docker-entrypoint.sh"]
CMD        ["--config.file=/opt/prometheus/config/prometheus.yml", \
             "--storage.tsdb.path=/opt/prometheus/data" ]
EXPOSE 9090