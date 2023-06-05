FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum Makefile scripts/hashfiles.sh ./
COPY vendor ./vendor
COPY cmd ./cmd
COPY assets ./assets
COPY public ./public

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Minify Assets
# hadolint ignore=DL3008
RUN apt-get update \
    && apt-get install -y --no-install-recommends minify \
    && ./hashfiles.sh \
    && find ./public -type f \( \
    -name "*.html" \
    -o -name '*.js' \
    -o -name '*.css' \
    \) \
    -print0 | \
    xargs -0  -I '{}' sh -c 'minify -o "{}" "{}"'

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./bin/song-stitch cmd/*.go

# hadolint ignore=DL3006
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /app

USER nonroot:nonroot
COPY --chown=nonroot:nonroot --from=builder /app/bin/song-stitch /app/song-stitch
COPY --chown=nonroot:nonroot --from=builder /app/public /app/public
COPY --chown=nonroot:nonroot assets/ /app/assets

ENTRYPOINT ["/app/song-stitch"]
