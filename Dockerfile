FROM golang:1.23-alpine AS base
WORKDIR /app

ENV GO111MODULE="on"
ENV GOOS="linux"
ENV CGO_ENABLED=0

# System dependencies
RUN apk update \
    && apk add --no-cache \
    ca-certificates \
    git \
    && update-ca-certificates


### Development with hot reload and debugger
FROM base AS dev
WORKDIR /app

# Hot reloading mod
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod \
    go install github.com/air-verse/air@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest
    # && go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
EXPOSE 8080
EXPOSE 2345

ENTRYPOINT ["/bin/sh", "-c"]
CMD ["air"]

### Executable builder
FROM base AS builder
WORKDIR /app

# Application dependencies
COPY . /app
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download \
    && go mod verify
    # && go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o my-great-program -a ./cmd/app

# Copy migrate
# COPY --from=base /go/bin/migrate /usr/local/bin/migrate

### Production
FROM alpine:latest

RUN apk update \
    && apk add --no-cache \
    ca-certificates \
    curl \
    tzdata \
    && update-ca-certificates

# Copy executable
COPY --from=builder /app/my-great-program /usr/local/bin/my-great-program

# Copy migrate
# COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

COPY app.env /
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/my-great-program"]
