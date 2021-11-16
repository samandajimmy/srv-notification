FROM alpine:3.14

ARG ARG_LOG_LEVEL=error
ARG ARG_LOG_FORMAT=json
ARG ARG_BIN_FILE=app
ARG ARG_APP_ENV=1
ARG ARG_PORT=3000

RUN apk add --no-cache ca-certificates

LABEL maintainer="Saggaf Arsyad <saggaf@nusantarabetastudio.com>"

COPY . /app

ENV LOG_LEVEL ${ARG_LOG_LEVEL}
ENV LOG_FORMAT ${ARG_LOG_FORMAT}
ENV BIN_FILE ${ARG_BIN_FILE}
ENV APP_ENV ${ARG_APP_ENV}

WORKDIR /app

RUN chmod +x ${ARG_BIN_FILE}

EXPOSE ${ARG_PORT}

ENTRYPOINT /app/${BIN_FILE} -env=${APP_ENV}