FROM registry.access.redhat.com/ubi8-minimal:latest
WORKDIR /
COPY crdoc .
ENTRYPOINT ["/crdoc"]
