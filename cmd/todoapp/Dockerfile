FROM golang:1.26.0-bookworm AS builder

ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.org
ENV GO111MODULE=on

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download -x

COPY . .
RUN CGO_ENABLED=0 go build -o /app/exe /app/cmd/todoapp

FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/exe /app
CMD ["/app/exe"]