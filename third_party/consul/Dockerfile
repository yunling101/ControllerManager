FROM alpine:3.19

WORKDIR /opt/consul

ARG VERSION="1.9.17"
ARG PACKAGE="consul_${VERSION}_linux_amd64"

RUN wget https://releases.hashicorp.com/consul/${VERSION}/${PACKAGE}.zip \
    && unzip ${PACKAGE}.zip \
    && rm -f ${PACKAGE}.zip \
    && mkdir -p /opt/consul/bin \
    && mkdir -p /opt/consul/config \
    && mkdir -p /opt/consul/data \
    && mv consul /opt/consul/bin

COPY ./third_party/consul/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
ENTRYPOINT ["docker-entrypoint.sh"]

CMD ["agent", "-ui", "-node=node", "-server", "-datacenter=dc", "-bootstrap-expect=1"]

EXPOSE 8300
EXPOSE 8500