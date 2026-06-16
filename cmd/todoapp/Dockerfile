FROM golang:1.26.0-bookworm AS builder

ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.org
ENV GO111MODULE=on

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download -x

# Копируем только существующие директории
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY public/ ./public/
COPY docs/ ./docs/

RUN CGO_ENABLED=0 go build -o /app/exe /app/cmd/todoapp

FROM alpine:3.21

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/exe /app/exe
COPY --from=builder /app/public /app/public

ENV PROJECT_ROOT=/app

EXPOSE 5050

CMD ["/app/exe"]