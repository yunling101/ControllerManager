FROM alpine:3.19

WORKDIR /opt

ARG VERSION="0.27.0"
ARG PACKAGE="alertmanager-${VERSION}.linux-amd64"

RUN wget https://github.com/prometheus/alertmanager/releases/download/v${VERSION}/${PACKAGE}.tar.gz \
    && tar -zxvf ${PACKAGE}.tar.gz \
    && rm -f ${PACKAGE}.tar.gz \
    && mv ${PACKAGE} alertmanager \
    && mkdir -p /opt/alertmanager/bin \
    && mkdir -p /opt/alertmanager/config \
    && mkdir -p /opt/alertmanager/data \
    && mkdir -p /opt/alertmanager/template \
    && mv /opt/alertmanager/alertmanager /opt/alertmanager/bin \
    && mv /opt/alertmanager/amtool /opt/alertmanager/bin \
    && mv /opt/alertmanager/alertmanager.yml /opt/alertmanager/config

WORKDIR /opt/alertmanager

COPY ./third_party/alertmanager/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["--config.file", "/opt/alertmanager/config/alertmanager.yml", \
    "--storage.path", "/opt/alertmanager/data"]

EXPOSE 9093
EXPOSE 9094