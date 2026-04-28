FROM golang:1.22-alpine AS build

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/orbit-engine ./cmd/engine

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/orbit-engine /usr/local/bin/orbit-engine
EXPOSE 7001
ENTRYPOINT ["/usr/local/bin/orbit-engine"]
