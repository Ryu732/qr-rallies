# Dockerfile.heroku（Heroku Container Registry用）
FROM golang:1.24-alpine AS builder

# 必要なパッケージをインストール
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# 依存関係をコピーしてダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# バイナリをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 本番用の軽量イメージ
FROM alpine:latest

# SSL証明書を追加
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# ビルドしたバイナリをコピー
COPY --from=builder /app/main .

# Herokuが動的にポートを設定
CMD ["./main"]
