package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/elastic/go-elasticsearch/v8"
)

// generateDocs は指定された数のダミードキュメントを生成します。
func generateDocs(num int) []map[string]interface{} {
	docs := make([]map[string]interface{}, num)
	for i := 0; i < num; i++ {
		docs[i] = map[string]interface{}{
			"title":  fmt.Sprintf("Test Document %d", i+1),
			"author": "Test Author",
		}
	}
	return docs
}

// BenchmarkInsert は SingleInsert と BulkInsert のパフォーマンスを比較します。
func BenchmarkInsert(b *testing.B) {
	// Elasticsearch クライアントのセットアップ
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}
	client, err := NewClient(cfg)
	if err != nil {
		b.Fatalf("failed to create client: %v", err)
	}

	// ベンチマーク対象のドキュメント件数
	docCounts := []int{10, 100, 1000, 10000}

	for _, count := range docCounts {
		// SingleInsert のベンチマーク
		b.Run(fmt.Sprintf("SingleInsert/%d_docs", count), func(b *testing.B) {
			docs := generateDocs(count)
			indexName := fmt.Sprintf("test-single-%d", count)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// タイマーを停止してテストごとのセットアップ（インデックス作成）を行う
				b.StopTimer()
				err := client.CreateIndex(context.Background(), indexName)
				if err != nil {
					b.Fatalf("failed to create index: %v", err)
				}
				b.StartTimer()

				// 実際の挿入処理を計測
				err = client.SingleInsert(context.Background(), indexName, docs)
				if err != nil {
					b.Fatalf("SingleInsert failed: %v", err)
				}
			}
		})

		// BulkInsert のベンチマーク
		b.Run(fmt.Sprintf("BulkInsert/%d_docs", count), func(b *testing.B) {
			docs := generateDocs(count)
			indexName := fmt.Sprintf("test-bulk-%d", count)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// タイマーを停止してテストごとのセットアップ（インデックス作成）を行う
				b.StopTimer()
				err := client.CreateIndex(context.Background(), indexName)
				if err != nil {
					b.Fatalf("failed to create index: %v", err)
				}
				b.StartTimer()

				// 実際の挿入処理を計測
				err = client.BulkInsert(context.Background(), indexName, docs)
				if err != nil {
					b.Fatalf("BulkInsert failed: %v", err)
				}
			}
		})

		// SingleInsertWithRefresh のベンチマーク
		b.Run(fmt.Sprintf("SingleInsertWithRefresh/%d_docs", count), func(b *testing.B) {
			docs := generateDocs(count)
			indexName := fmt.Sprintf("test-single-refresh-%d", count)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				err := client.CreateIndex(context.Background(), indexName)
				if err != nil {
					b.Fatalf("failed to create index: %v", err)
				}
				b.StartTimer()

				err = client.SingleInsertWithRefresh(context.Background(), indexName, docs)
				if err != nil {
					b.Fatalf("SingleInsertWithRefresh failed: %v", err)
				}
			}
		})

		// BulkInsertWithRefresh のベンチマーク
		b.Run(fmt.Sprintf("BulkInsertWithRefresh/%d_docs", count), func(b *testing.B) {
			docs := generateDocs(count)
			indexName := fmt.Sprintf("test-bulk-refresh-%d", count)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				err := client.CreateIndex(context.Background(), indexName)
				if err != nil {
					b.Fatalf("failed to create index: %v", err)
				}
				b.StartTimer()

				err = client.BulkInsertWithRefresh(context.Background(), indexName, docs)
				if err != nil {
					b.Fatalf("BulkInsertWithRefresh failed: %v", err)
				}
			}
		})
	}
}
