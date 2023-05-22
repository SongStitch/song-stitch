FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum Makefile ./
COPY vendor ./vendor
COPY cmd ./cmd
COPY assets ./assets

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/song-stitch cmd/*.go

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /app

USER nonroot:nonroot
COPY --chown=nonroot:nonroot --from=builder /app/bin/song-stitch /app/song-stitch
COPY --chown=nonroot:nonroot assets/ /app/assets
COPY --chown=nonroot:nonroot public/ /app/public

ENTRYPOINT ["/app/song-stitch"]
