FROM alpine:3.6
EXPOSE 3000

ENV GODEBUG netdns=go

ADD drone-secret-plugin /bin/
ENTRYPOINT ["/bin/drone-secret-plugin"]
