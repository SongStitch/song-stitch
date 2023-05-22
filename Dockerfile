FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY go.mod go.sum Makefile ./
COPY vendor ./vendor
COPY cmd ./cmd
COPY assets ./assets

# RUN make build
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/song-stitch cmd/*.go

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=builder /app/bin/song-stitch /song-stitch
COPY assets/ ./assets
USER nonroot:nonroot

ENTRYPOINT ["/song-stitch"]
