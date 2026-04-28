FROM golang:1.22-alpine AS build

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/tle-syncer ./cmd/syncer

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/tle-syncer /usr/local/bin/tle-syncer
ENTRYPOINT ["/usr/local/bin/tle-syncer"]
