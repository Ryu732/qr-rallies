# Dockerfile
FROM golang:1.24.1-alpine

# 必要なパッケージをインストール
RUN apk add --no-cache git

# 作業ディレクトリを設定
WORKDIR /app

# Airをインストール（ホットリロード用）
RUN go install github.com/air-verse/air@latest

# go.modとgo.sumをコピーして依存関係をダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ポートを公開
EXPOSE 8080

# Air設定ファイルをコピー
COPY .air.toml .

# tmpディレクトリを作成
RUN mkdir -p tmp

# デフォルトでAirを起動
CMD ["air", "-c", ".air.toml"]
