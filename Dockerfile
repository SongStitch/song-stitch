FROM node:21-alpine@sha256:78c45726ea205bbe2f23889470f03b46ac988d14b6d813d095e2e9909f586f93 AS node-builder

WORKDIR /app/ui
COPY ui ./
RUN npm install && npm run build

FROM golang:1.25-bookworm@sha256:09f53deea14d4019922334afe6258b7b776afc1d57952be2012f2c8c4076db05 AS builder

WORKDIR /app

COPY go.mod go.sum Makefile ./
COPY vendor ./vendor
COPY cmd ./cmd
COPY internal ./internal
COPY assets ./assets
COPY --from=node-builder /app/public ./public

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Minify Assets
# hadolint ignore=DL3008
RUN dpkg --add-architecture amd64 && apt-get update && apt-get update \
  && apt-get install -y --no-install-recommends \
  minify=2.12.4-2 \
  libwebp-dev=1.2.4-0.2+deb12u1 \
  && find ./public -type f \( \
  -name "*.html" \
  -o -name '*.js' \
  -o -name '*.css' \
  \) \
  -print0 | \
  xargs -0  -I '{}' sh -c 'minify -o "{}" "{}"'

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w -linkmode 'external' -extldflags '-static'" -o ./bin/song-stitch cmd/*.go

FROM gcr.io/distroless/static-debian12:nonroot@sha256:2b7c93f6d6648c11f0e80a48558c8f77885eb0445213b8e69a6a0d7c89fc6ae4 AS build-release-stage

WORKDIR /app

COPY --chown=nonroot:nonroot --from=builder /app/bin/song-stitch /app/song-stitch
COPY --chown=nonroot:nonroot --from=builder /app/public /app/public
COPY --chown=nonroot:nonroot assets/NotoSans-Bold.ttf assets/NotoSans-Regular.ttf /app/assets/

ENTRYPOINT ["/app/song-stitch"]
