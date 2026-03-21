FROM golang:1.24 AS builder

ARG VERSION
ARG COMMIT

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# VERSION and COMMIT will be taken from env
RUN make build

# -------------------
FROM alpine:3.22.1 AS runtime

ARG PORT=3000
ARG DB_HOST
ARG DB_PORT
ARG DB_USER
ARG DB_PASSWORD
ARG DB_NAME
ARG DB_SSL_MODE=require

RUN apk update && apk add --no-cache ca-certificates curl 
RUN curl -sSf https://atlasgo.sh | sh -s -- --yes

WORKDIR /app

COPY --from=builder /app/bin/server .
COPY --from=builder /app/atlas.hcl .
COPY --from=builder /app/domain/migrations/ ./domain/migrations/
COPY ./docker/entrypoint.sh ./docker/entrypoint.sh

RUN addgroup -S app && adduser -S app -G app \
	&& chmod +x /app/docker/entrypoint.sh \
	&& chown -R app:app /app

USER app

ENV PORT=${PORT}
ENV DB_HOST=${DB_HOST}
ENV DB_PORT=${DB_PORT}
ENV DB_USER=${DB_USER}
ENV DB_PASSWORD=${DB_PASSWORD}
ENV DB_NAME=${DB_NAME}
ENV DB_SSL_MODE=${DB_SSL_MODE}
ENV DB_RUN_MIGRATIONS=true
ENV DB_MIGRATION_MAX_ATTEMPTS=30
ENV DB_MIGRATION_RETRY_DELAY_SECONDS=2
ENV DB_MIGRATION_ALLOW_DIRTY=true
ENV DB_MIGRATION_BASELINE=
EXPOSE ${PORT}

ENTRYPOINT ["/app/docker/entrypoint.sh"]