FROM golang:1.26.0-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o scheduler .

FROM alpine:3.22
RUN apk --no-cache add ca-certificates && \
    addgroup -S -g 10001 scheduler && \
    adduser -S -D -H -u 10001 -G scheduler scheduler && \
    mkdir -p /data && chown scheduler:scheduler /data

WORKDIR /app
COPY --from=builder /build/scheduler .
COPY --from=builder /build/web ./web

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

USER scheduler
EXPOSE 7540
CMD ["./scheduler"]
