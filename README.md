# go-elasticsearch-playground

Go言語の[go-elasticsearch](https://github.com/elastic/go-elasticsearch)パッケージを使用して、Elasticsearchの様々な機能を試すためのサンプルコード集です。

## 概要

このリポジトリでは、Elasticsearchとの連携において実践的な実装パターンを提供します。パフォーマンス比較、並行処理の最適化、Learning to Rank (LTR)の実装例など、実際の開発で役立つサンプルコードを収録しています。

## プロジェクト構成

### bulk-insert-vs-single-insert/
Single InsertとBulk Insertのパフォーマンス比較。大量データ投入時のベストプラクティスを検証。

### concurrent-bulk-insert/
並行バルクインサートの実装パターンとパフォーマンス比較。`esutil.BulkIndexer`の効果的な使い方を探求。

### search-using-ltr/
Learning to Rank (LTR)を使用した検索の実装例。機械学習を活用した検索結果のランキング改善。

## クイックスタート

```bash
# リポジトリのクローン
git clone https://github.com/yourusername/go-elasticsearch-playground.git
cd go-elasticsearch-playground

# Elasticsearchの起動
docker-compose up -d

# 各プロジェクトの実行例
cd bulk-insert-vs-single-insert
go test -bench=. -benchmem
```

## 環境要件

- Go 1.18以上
- Docker & Docker Compose
- Elasticsearch 8.14.1
