FROM golang:1.22-bookworm AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -trimpath -o /out/majorpath ./cmd/server

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
