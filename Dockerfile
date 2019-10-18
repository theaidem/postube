FROM alpine:3.10
RUN apk add -U --no-cache ca-certificates
WORKDIR /app
ARG appname
ENV APP_NAME=${appname}
ADD ${APP_NAME} /app
CMD ./${APP_NAME}