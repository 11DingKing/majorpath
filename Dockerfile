FROM --platform=$BUILDPLATFORM golang:1.22-bookworm AS build
ARG TARGETARCH
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN if [ "$TARGETARCH" = "arm64" ]; then apt-get update && apt-get install -y --no-install-recommends gcc-aarch64-linux-gnu libc6-dev-arm64-cross && rm -rf /var/lib/apt/lists/*; fi
RUN if [ "$TARGETARCH" = "arm64" ]; then CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64 go build -trimpath -o /out/majorpath ./cmd/server; else CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -trimpath -o /out/majorpath ./cmd/server; fi

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=build /out/majorpath /app/majorpath
COPY migrations /app/migrations
ENV PORT=8080 DATABASE_PATH=/data/majorpath.db
RUN mkdir -p /data
EXPOSE 8080
HEALTHCHECK --interval=5s --timeout=3s --start-period=5s CMD ["/app/majorpath", "healthcheck"]
ENTRYPOINT ["/app/majorpath"]
