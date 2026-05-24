FROM alpine:latest
RUN apk --no-cache add ca-certificates
ARG TARGETPLATFORM
ENTRYPOINT ["/usr/bin/pgxcli"]
COPY $TARGETPLATFORM/pgxcli /usr/bin/
