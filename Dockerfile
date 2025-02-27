FROM golang:1.24-alpine3.21 AS builder

LABEL org.opencontainers.image.source=https://github.com/yrn-go/yrn

#avoid root
ENV USER=appuser
ENV UID=1000

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOPRIVATE=github.com/yrn-go/*

# Configurar credenciais no .netrc
#ARG GITHUB_APP_TOKEN
#RUN echo -e "machine github.com\n  login x-oauth-basic\n  password ${GITHUB_APP_TOKEN}" > ~/.netrc && \
#    chmod 600 ~/.netrc

#avoid root
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd/ ./cmd/
COPY pkg/ ./pkg/
COPY internal/ ./internal/
COPY module/ ./module/

RUN go build -o api ./cmd/api/main.go
RUN go build -o agent ./cmd/agent/main.go
RUN go build -o connector ./cmd/connector/main.go

FROM scratch

LABEL org.opencontainers.image.source=https://github.com/yrn-go/yrn

#avoid rootless
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

WORKDIR /app
COPY --from=builder /app/api /app/api
COPY --from=builder /app/agent /app/agent
COPY --from=builder /app/connector /app/connector

#avoid rootless
USER appuser:appuser
