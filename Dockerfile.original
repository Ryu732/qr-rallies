# Dockerfile（air不要版）
FROM golang:1.24-alpine AS builder

# 必要なパッケージをインストール
RUN apk add --no-cache git

WORKDIR /app

# 依存関係をコピーしてダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# バイナリをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 本番用の軽量イメージ
FROM alpine:latest AS production

# SSL証明書を追加
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# ビルドしたバイナリをコピー
COPY --from=builder /app/main .

# 本番環境用のコマンド
CMD ["./main"]

# 開発用のステージ（airなし）
FROM golang:1.24-alpine AS development

RUN apk add --no-cache git

WORKDIR /app

# 依存関係をコピー
COPY go.mod go.sum ./
RUN go mod download

# 開発環境用のコマンド（go runで直接実行）
CMD ["go", "run", "main.go"]
