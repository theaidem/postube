FROM alpine:latest
RUN apk add -U --no-cache ca-certificates
WORKDIR /app
ARG appname
ENV APP_NAME=${appname}
CMD ./${APP_NAME}