FROM node:21-alpine@sha256:78c45726ea205bbe2f23889470f03b46ac988d14b6d813d095e2e9909f586f93 AS node-builder

WORKDIR /app/ui
COPY ui ./
RUN npm install && npm run build

FROM golang:1.23-bookworm@sha256:1a5326b07cbab12f4fd7800425f2cf25ff2bd62c404ef41b56cb99669a710a83 AS builder

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
    && dpkg --add-architecture amd64 && apt-get update && \
    apt-get install -y --no-install-recommends \
    minify:amd64=2.12.4-2 \
    libwebp-dev:amd64=1.2.4-0.2+deb12u1 \
    && find ./public -type f \( \
    -name "*.html" \
    -o -name '*.js' \
    -o -name '*.css' \
    \) \
    -print0 | \
    xargs -0 \
    -I '{}' sh -c 'minify -o "{}" "{}"'

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w -linkmode 'external' -extldflags '-static'" -o ./bin/song-stitch cmd/*.go

FROM gcr.io/distroless/static-debian12:nonroot@sha256:26f9b99f2463f55f20db19feb4d96eb88b056e0f1be7016bb9296a464a89d772 AS build-release-stage

WORKDIR /app

COPY --chown=nonroot:nonroot --from=builder /app/bin/song-stitch /app/song-stitch
COPY --chown=nonroot:nonroot --from=builder /app/public /app/public
COPY --chown=nonroot:nonroot assets/NotoSans-Bold.ttf assets/NotoSans-Regular.ttf /app/assets/

ENTRYPOINT ["/app/song-stitch"]
