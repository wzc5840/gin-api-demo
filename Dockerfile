# 公式Golangイメージの使用
FROM golang:1.24

# 作業ディレクトリを設定
WORKDIR /app

# go.mod と go.sum をコピー
COPY go.mod go.sum ./

# 依存関係をダウンロード
RUN go mod download

# ソースコードをコピー
COPY . .

# 実行ファイルのビルド (新しいmain.goの場所を指定)
RUN go build -o main ./cmd/server

# PostgreSQL接続用の環境変数を設定
ENV DATABASE_URL="host=postgres user=postgres password=postgres dbname=gin_demo port=5432 sslmode=disable TimeZone=Asia/Shanghai"
ENV PORT=8080

# ポート8080を公開
EXPOSE 8080

# アプリケーションを起動
CMD ["./main"]
