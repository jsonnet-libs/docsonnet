FROM alpine:3.12
ENTRYPOINT ["/usr/bin/docsonnet"]
COPY docsonnet /usr/bin/docsonnet
