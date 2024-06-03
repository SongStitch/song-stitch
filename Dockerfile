FROM node:21-alpine@sha256:c986eb0b8970240f8d648e524bab46016b78f290f912aac16a4aa6705dde05f4 AS node-builder

WORKDIR /app/ui
COPY ui ./
RUN npm install && npm run build

FROM golang:1.22-bookworm@sha256:5c56bd47228dd572d8a82971cf1f946cd8bb1862a8ec6dc9f3d387cc94136976 AS builder

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

FROM gcr.io/distroless/static-debian12:nonroot@sha256:e9ac71e2b8e279a8372741b7a0293afda17650d926900233ec3a7b2b7c22a246 AS build-release-stage

WORKDIR /app

COPY --chown=nonroot:nonroot --from=builder /app/bin/song-stitch /app/song-stitch
COPY --chown=nonroot:nonroot --from=builder /app/public /app/public
COPY --chown=nonroot:nonroot assets/NotoSans-Bold.ttf assets/NotoSans-Regular.ttf /app/assets/

ENTRYPOINT ["/app/song-stitch"]
