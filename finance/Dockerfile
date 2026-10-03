FROM golang:1.21-bookworm AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/finance .

FROM debian:bookworm-slim
RUN apt-get update && apt-get install --no-install-recommends -y ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd --gid 1000 finance \
    && useradd --uid 1000 --gid finance --no-create-home finance

WORKDIR /app
COPY --from=builder /out/finance /usr/local/bin/finance
USER finance
EXPOSE 3001/tcp
CMD ["finance"]
