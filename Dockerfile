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

# 実行ファイルのビルド
RUN go build -o main .

# アプリケーションを起動
CMD ["./main"]
