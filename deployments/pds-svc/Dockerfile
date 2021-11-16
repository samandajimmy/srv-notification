# -----------
# Build Stage
# -----------

FROM golang:1.16 as builder

# Set-up go proxy
ENV GO111MODULE=on \
    GOPROXY="https://artifactory.pegadaian.co.id/repository/go-group-01/"
COPY ssl_certificate.crt /usr/local/share/ca-certificates/ssl_certificate.crt

# Install some dependencies needed to build the project
RUN apt install git make \
    && mkdir /usr/src/app \
    && chmod 644 /usr/local/share/ca-certificates/ssl_certificate.crt \
    && update-ca-certificates \
    && echo ${GO111MODULE} ${GOPROXY}

# Here we copy the rest of the source code
WORKDIR /usr/src/app
COPY Makefile go.mod responses.yml ./
COPY .git ./.git
COPY cmd ./cmd
COPY internal ./internal

# Compiling...
RUN ls -la && make release

# ------------
# Deploy Stage
# ------------

FROM alpine:3.14

ARG ARG_LOG_LEVEL=error
ARG ARG_LOG_FORMAT=json
ARG ARG_BIN_FILE=app
ARG ARG_APP_ENV=1
ARG ARG_PORT=3000

RUN apk add --no-cache tzdata ca-certificates

COPY --from=builder /usr/src/app/bin/release /usr/src/app

ENV LOG_LEVEL ${ARG_LOG_LEVEL}
ENV LOG_FORMAT ${ARG_LOG_FORMAT}
ENV BIN_FILE ${ARG_BIN_FILE}
ENV APP_ENV ${ARG_APP_ENV}


# Running...
WORKDIR /usr/src/app

EXPOSE ${ARG_PORT}

ENTRYPOINT ["./pds-svc"]