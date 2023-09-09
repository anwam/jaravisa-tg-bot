FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/...

FROM debian:buster-slim AS runner
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/main /app/main

ARG _TELEGRAM_BOT_TOKEN
ARG _NOTION_SECRET
ARG _NOTION_DATABASE_ID
ARG _PORT

ENV TELEGRAM_BOT_TOKEN=$_TELEGRAM_BOT_TOKEN
ENV NOTION_SECRET=$_NOTION_SECRET
ENV NOTION_DATABASE_ID=$_NOTION_DATABASE_ID
ENV PORT=$_PORT

CMD ["/app/main"]
