
FROM golang:1.22-alpine

# 必要なパッケージをインストール
RUN apk add --no-cache git

# 作業ディレクトリを設定
WORKDIR /app

# Goの依存関係をキャッシュするための前段階
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# アプリケーションソースコードをコピー
COPY . .

# アプリケーションをビルド
RUN go build -o /todo-app

# 実行可能ファイルを起動
CMD ["/todo-app"]
