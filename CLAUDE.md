# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

go-elasticsearchパッケージを使用してElasticsearchの様々な機能を試すためのGoプレイグラウンドリポジトリです。以下のサンプルコードが含まれています：
- Single InsertとBulk Insertのパフォーマンス比較
- 並行バルクインサートの実装パターン
- Learning to Rank (LTR)を使用した検索実装

## 開発コマンド

### Elasticsearchの起動
```bash
docker-compose up -d
```

### ベンチマークの実行
```bash
# bulk-insert-vs-single-insertディレクトリで
go test -bench=. -benchmem

# concurrent-bulk-insertディレクトリで  
go test -bench=. -benchmem
```

### テストの実行
```bash
go test ./...
```

### 単一テストの実行
```bash
go test -run TestName
```

## アーキテクチャ

### プロジェクト構成
- **bulk-insert-vs-single-insert/**: 単一ドキュメント挿入とバルク操作のパフォーマンス比較
  - `insert.go`: SingleInsert、BulkInsertメソッド（リフレッシュあり/なし）を含むクライアント実装
  - `insert_test.go`: 異なる挿入戦略を比較するベンチマークテスト

- **concurrent-bulk-insert/**: 並行バルクインデックスパターンの探求
  - 複数実装: BulkInsert、BulkInsertConcurrent、BulkInsertConcurrentV2、BulkInsertConcurrentV3
  - 効率的なバルク操作のために`esutil.BulkIndexer`を使用

- **search-using-ltr/**: Learning to Rank実装
  - 異なる検索エンジン（Elasticsearch、OpenSearch、Solr）用のJupyterノートブック
  - LTR実験用にTMDBデータセットを使用

### 主要な設計パターン

1. **クライアント構造**: 各モジュールはベースクライアントと型付きElasticsearchクライアントの両方をラップするClient構造体を使用
   ```go
   type Client struct {
       baseClient  *elasticsearch.Client
       typedClient *elasticsearch.TypedClient
   }
   ```

2. **インデックス管理**: CreateIndexメソッドはインデックスの削除（存在する場合）と適切なマッピングでの作成を処理

3. **バルク操作**: 効率的なバルクドキュメントインデックスのために`esutil.BulkIndexer`を使用

4. **ベンチマーク構造**: テストは異なるドキュメント数（10、100、1000、10000）でテーブル駆動ベンチマークを使用

## 環境設定

- Go 1.24.2（go.modファイルで指定）
- Elasticsearch 8.14.1（Docker経由）
- デフォルトのElasticsearch URL: http://localhost:9200
- ローカル開発用にセキュリティは無効化（xpack.security.enabled=false）