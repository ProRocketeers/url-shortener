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
FROM alpine:3.22 AS runtime

RUN apk update && apk add --no-cache ca-certificates 

WORKDIR /app

COPY --from=builder /app/bin/server .

ENTRYPOINT ["./server"]