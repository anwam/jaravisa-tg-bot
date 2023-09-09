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

ARG TELEGRAM_BOT_TOKEN
ARG NOTION_SECRET
ARG NOTION_DATABASE_ID
ARG PORT

ENV TELEGRAM_BOT_TOKEN=$TELEGRAM_BOT_TOKEN
ENV NOTION_SECRET=$NOTION_SECRET
ENV NOTION_DATABASE_ID=$NOTION_DATABASE_ID
ENV PORT=$PORT

CMD ["/app/main"]
