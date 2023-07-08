FROM node:18-alpine AS node-builder

WORKDIR /app/ui
COPY ui ./
RUN npm install && npm run build

# switch to bullseye due to https://github.com/GoogleContainerTools/distroless/issues/1342
FROM golang:1.20-bullseye AS builder

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
RUN apt-get update \
  && apt-get install -y --no-install-recommends minify libwebp-dev \
  && find ./public -type f \( \
  -name "*.html" \
  -o -name '*.js' \
  -o -name '*.css' \
  \) \
  -print0 | \
  xargs -0  -I '{}' sh -c 'minify -o "{}" "{}"'

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o ./bin/song-stitch cmd/*.go

# hadolint ignore=DL3006
FROM gcr.io/distroless/base-debian11 AS build-release-stage
ARG ARCH=aarch64

# Copy dependency for webp
COPY --from=builder /usr/lib/${ARCH}-linux-gnu/libwebp.so* /usr/lib/${ARCH}-linux-gnu/

WORKDIR /app

USER nonroot:nonroot
COPY --chown=nonroot:nonroot --from=builder /app/bin/song-stitch /app/song-stitch
COPY --chown=nonroot:nonroot --from=builder /app/public /app/public
COPY --chown=nonroot:nonroot assets/ /app/assets

ENTRYPOINT ["/app/song-stitch"]
