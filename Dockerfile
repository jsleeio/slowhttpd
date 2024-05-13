FROM alpine
COPY slowhttpd /usr/bin/slowhttpd
ENTRYPOINT ["/usr/bin/slowhttpd"]
