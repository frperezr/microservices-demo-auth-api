FROM alpine

RUN apk add --update ca-certificates

WORKDIR /src/auth-api

COPY bin/noken-auth-api /usr/bin/noken-auth-api

EXPOSE 3020

CMD ["/bin/sh", "-l", "-c", "noken-auth-api"]