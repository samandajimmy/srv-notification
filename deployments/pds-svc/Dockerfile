# -----------
# Build Stage
# -----------

FROM golang as builder

# Install some dependencies needed to build the project
RUN apt install git make && mkdir /usr/src/app

# Here we copy the rest of the source code
WORKDIR /usr/src/app
COPY Makefile go.mod responses.yml ./
COPY .git ./.git
COPY cmd ./cmd
COPY internal ./internal

# Compiling...
RUN ls -la && make api-release

# ------------
# Deploy Stage
# ------------

FROM alpine:latest

RUN apk add --no-cache tzdata

COPY --from=builder /usr/src/app/bin/release /usr/src/app

# Running...
WORKDIR /usr/src/app
ENTRYPOINT ["./pds-svc"]