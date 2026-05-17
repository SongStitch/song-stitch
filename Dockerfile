FROM golang:1.26-bookworm@sha256:252599aeb51ad60b83e4d8821802068127c528c707cb7dd7afd93be057c6011c AS builder

WORKDIR /app

COPY go.mod go.sum Makefile ./
COPY vendor ./vendor
COPY cmd ./cmd
COPY internal ./internal
COPY assets ./assets
COPY ui ./ui

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

RUN rm -rf public \
      && mkdir -p public/assets \
      && cp -R ui/public/. public/ \
      && cp -R ui/assets/. public/assets/ \
      && cp ui/*.html ui/*.css ui/*.js public/

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

FROM gcr.io/distroless/static-debian12:nonroot@sha256:d093aa3e30dbadd3efe1310db061a14da60299baff8450a17fe0ccc514a16639 AS build-release-stage

WORKDIR /app

COPY --chown=nonroot:nonroot --from=builder /app/bin/song-stitch /app/song-stitch
COPY --chown=nonroot:nonroot --from=builder /app/public /app/public
COPY --chown=nonroot:nonroot assets/NotoSans-Bold.ttf assets/NotoSans-Regular.ttf /app/assets/

ENTRYPOINT ["/app/song-stitch"]
