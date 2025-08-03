# Gin API デモプロジェクト

Gin フレームワークを使用したRESTful APIデモプロジェクトです。ユーザー認証機能と投稿管理機能を提供します。

## 📋 目次

- [機能概要](#機能概要)
- [技術スタック](#技術スタック)
- [プロジェクト構造](#プロジェクト構造)
- [インストール方法](#インストール方法)
- [実行方法](#実行方法)
- [API仕様](#api仕様)
- [Postmanでのテスト方法](#postmanでのテスト方法)
- [開発者向け情報](#開発者向け情報)

## 🚀 機能概要

### 認証機能
- ユーザー登録
- ユーザーログイン
- Bearer Token認証

### ユーザー管理機能
- プロフィール取得
- ユーザーリスト表示（ページネーション対応）
- ユーザー詳細表示
- ユーザー情報更新（本人のみ）
- ユーザー削除（自分以外）

### 投稿管理機能
- 投稿作成（下書き/公開）
- 投稿リスト表示（ステータス別フィルタリング）
- 投稿詳細表示（閲覧数カウント）
- 投稿更新（作成者のみ）
- 投稿削除（作成者のみ）
- マイ投稿一覧

## 🛠 技術スタック

- **言語**: Go 1.24
- **フレームワーク**: Gin
- **データベース**: PostgreSQL 15
- **ORM**: GORM
- **コンテナ**: Docker & Docker Compose
- **認証**: Bearer Token
- **API仕様**: RESTful API

## 📁 プロジェクト構造

```
gin-api-demo/
├── cmd/server/           # アプリケーションエントリーポイント
│   └── main.go
├── internal/             # 内部パッケージ
│   ├── auth/            # 認証関連
│   │   ├── handler/     # HTTPハンドラー
│   │   └── service/     # ビジネスロジック
│   ├── user/            # ユーザー関連
│   │   ├── model/       # データモデル
│   │   └── repository/  # データアクセス層
│   └── post/            # 投稿関連
│       ├── handler/
│       ├── model/
│       ├── repository/
│       └── service/
├── pkg/                 # 共通パッケージ
│   ├── logger/          # ログ管理
│   ├── middleware/      # ミドルウェア
│   └── util/           # ユーティリティ
├── router/              # ルーティング設定
├── docker-compose.yml   # Docker構成
├── Dockerfile          # Dockerイメージ設定
└── README.md           # このファイル
```

## 💻 インストール方法

### 前提条件

以下のソフトウェアがインストールされている必要があります：

- [Docker](https://www.docker.com/) 20.10+
- [Docker Compose](https://docs.docker.com/compose/) 2.0+
- [Git](https://git-scm.com/)

### 1. リポジトリのクローン

```bash
git clone <repository-url>
cd gin-api-demo
```

### 2. 環境設定の確認

`docker-compose.yml`ファイルで以下の設定が正しいことを確認してください：

```yaml
# PostgreSQLの設定
POSTGRES_USER: postgres
POSTGRES_PASSWORD: postgres
POSTGRES_DB: gin_demo

# アプリケーションの設定
DATABASE_URL: "host=postgres user=postgres password=postgres dbname=gin_demo port=5432 sslmode=disable TimeZone=Asia/Shanghai"
PORT: 8080
```

## 🔄 実行方法

### Docker Composeを使用した実行（推奨）

```bash
# アプリケーションとデータベースを同時に起動
docker-compose up --build

# バックグラウンドで実行する場合
docker-compose up --build -d

# ログを確認する場合
docker-compose logs -f app
```

### 個別での実行

```bash
# PostgreSQLのみ起動
docker-compose up postgres -d

# Go アプリケーションを直接実行
export DATABASE_URL="host=localhost user=postgres password=postgres dbname=gin_demo port=5432 sslmode=disable TimeZone=Asia/Shanghai"
go run cmd/server/main.go
```

### サービスの停止

```bash
# すべてのサービスを停止
docker-compose down

# データベースのボリュームも削除する場合
docker-compose down -v
```

## 📚 API仕様

### ベースURL
```
http://localhost:8080
```

### エンドポイント一覧

#### 健康チェック
- `GET /hello` - サーバー状態確認

#### 認証API
- `POST /api/v1/auth/register` - ユーザー登録
- `POST /api/v1/auth/login` - ユーザーログイン

#### ユーザー管理API（認証必須）
- `GET /api/v1/user/profile` - プロフィール取得
- `GET /api/v1/user/list` - ユーザーリスト取得
- `GET /api/v1/user/:id` - ユーザー詳細取得
- `PUT /api/v1/user/:id` - ユーザー情報更新
- `DELETE /api/v1/user/:id` - ユーザー削除

#### 投稿管理API
**公開API（認証不要）**
- `GET /api/v1/posts` - 投稿リスト取得
- `GET /api/v1/posts/:id` - 投稿詳細取得

**認証必須API**
- `POST /api/v1/posts` - 投稿作成
- `PUT /api/v1/posts/:id` - 投稿更新
- `DELETE /api/v1/posts/:id` - 投稿削除
- `GET /api/v1/posts/my` - マイ投稿一覧

### レスポンス形式

すべてのAPIは以下の統一された形式でレスポンスを返します：

```json
{
    "statusCode": 200,
    "message": "成功メッセージ",
    "data": {}
}
```

### 認証方法

認証が必要なAPIには、Authorizationヘッダーにベアラートークンを設定してください：

```
Authorization: Bearer <your-token>
```

## 🧪 Postmanでのテスト方法

### 1. Postman Collectionのインポート

1. Postmanを開く
2. 「Import」ボタンをクリック
3. プロジェクトルートの `gin-api-demo.postman_collection.json` ファイルを選択
4. インポートを完了

### 2. 環境変数の設定

Collectionには以下の環境変数が事前設定されています：

```
base_url: http://localhost:8080
auth_token: (自動設定)
user_id: (自動設定)
post_id: (自動設定)
```

### 3. テスト手順

#### 基本的なテストフロー

1. **サーバー確認**
   ```
   Health Check → Welcome Page
   ```

2. **ユーザー登録とログイン**
   ```
   Authentication → Register User
   Authentication → Login
   ```
   ※ トークンが自動的に保存されます

3. **ユーザー管理のテスト**
   ```
   User Management → Get My Profile
   User Management → Get User List
   User Management → Update User
   ```

4. **投稿機能のテスト**
   ```
   Posts → Create Post (Draft)
   Posts → Publish Draft Post
   Posts → Get Published Posts
   Posts → Get My Posts
   ```

#### 詳細なテストシナリオ

**シナリオ1：新規ユーザー登録からブログ投稿まで**

1. `POST /api/v1/auth/register` - 新規ユーザー登録
2. `GET /api/v1/user/profile` - プロフィール確認
3. `POST /api/v1/posts` - 下書き投稿作成
4. `PUT /api/v1/posts/:id` - 投稿を公開状態に更新
5. `GET /api/v1/posts` - 公開投稿一覧で確認

**シナリオ2：投稿の管理**

1. `POST /api/v1/posts` - 複数の投稿を作成
2. `GET /api/v1/posts/my` - 自分の投稿一覧確認
3. `PUT /api/v1/posts/:id` - 投稿内容を更新
4. `GET /api/v1/posts/:id?view=true` - 閲覧数増加テスト
5. `DELETE /api/v1/posts/:id` - 投稿削除

**シナリオ3：権限テスト**

1. 2つのユーザーアカウントを作成
2. ユーザーAで投稿を作成
3. ユーザーBでユーザーAの投稿編集を試行（失敗することを確認）
4. ユーザーBで自分のアカウント削除を試行（失敗することを確認）

### 4. 自動化機能

Collectionには以下の自動化機能が組み込まれています：

- **自動トークン管理**: ログイン成功時にトークンを自動保存
- **自動IDセット**: 作成されたリソースのIDを自動保存
- **自動認証**: 認証が必要なリクエストに自動でトークンを設定

## 🔧 開発者向け情報

### ローカル開発環境

```bash
# 依存関係のインストール
go mod download

# アプリケーションの実行（PostgreSQLが別途必要）
go run cmd/server/main.go

# テストの実行
go test ./...

# ビルド
go build -o main ./cmd/server
```

### 効率的なDocker開発ワークフロー

#### 初回起動
```bash
# すべてのサービスを初回ビルドして起動
docker-compose up --build -d
```

#### 通常の開発サイクル

**コード変更後（ソースファイルのみ変更）**
```bash
# 通常の起動（自動でビルドが必要かどうか判断）
docker-compose up -d

# または、アプリケーションコンテナのみ再起動
docker-compose restart app
```

**依存関係やDockerfile変更後**
```bash
# ビルドが必要な場合のみ
docker-compose up --build -d
```

**完全リセットが必要な場合**
```bash
# すべてのコンテナとボリュームを削除して再構築
docker-compose down -v
docker-compose up --build -d
```

#### 開発に便利なコマンド

```bash
# アプリケーションのログをリアルタイムで確認
docker-compose logs -f app

# データベースのログを確認
docker-compose logs postgres

# 実行中のコンテナ状態を確認
docker-compose ps

# アプリケーションコンテナに入る
docker-compose exec app sh

# データベースコンテナに入る
docker-compose exec postgres psql -U postgres -d gin_demo
```

#### パフォーマンス最適化のポイント

- **不要なファイルの除外**: `.dockerignore`ファイルで不要なファイルをビルドコンテキストから除外
- **レイヤーキャッシュの活用**: `go.mod`と`go.sum`を先にコピーして依存関係をキャッシュ
- **マルチステージビルド**: 本番環境では軽量なイメージを使用することを推奨

### データベースマイグレーション

アプリケーション起動時に自動的にテーブルが作成されます：

- `users` - ユーザー情報
- `posts` - 投稿情報

### ログ設定

ログは以下の形式で出力されます：

```
INFO: 2024/01/01 12:00:00 logger.go:20: サーバーがポート8080で開始されました
```

### トラブルシューティング

**問題**: データベース接続エラー
```
解決策: 
1. PostgreSQLコンテナが起動していることを確認
2. docker-compose logs postgres でログを確認
3. 必要に応じて docker-compose down -v でボリュームをリセット
```

**問題**: ポート8080が使用中
```
解決策:
1. docker-compose.yml のポート設定を変更
2. または既存のプロセスを停止
```

**問題**: Postmanでトークンが自動設定されない
```
解決策:
1. Testsタブのスクリプトが正しく設定されているか確認
2. レスポンスが正常に返されているか確認
```
