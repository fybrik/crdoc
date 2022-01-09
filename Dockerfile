FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY crdoc .
USER nonroot:nonroot

ENTRYPOINT ["/crdoc"]
