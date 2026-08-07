FROM golang:1.26.5 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN go build -trimpath -ldflags='-s -w' -o /out/http-server-projeto-korp ./cmd/http-server-projeto-korp

FROM scratch

EXPOSE 8080

COPY --from=builder /out/http-server-projeto-korp /http-server-projeto-korp

ENTRYPOINT ["/http-server-projeto-korp"]
