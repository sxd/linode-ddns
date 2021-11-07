FROM gcr.io/distroless/base
ARG TARGETARCH

USER 1025
COPY  dist/linode-ddns_linux_${TARGETARCH} /linux-ddns