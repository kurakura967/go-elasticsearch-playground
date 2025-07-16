package main

import (
	"context"
	"fmt"
	"testing"
)

func generateDocs(num int) []map[string]interface{} {
	docs := make([]map[string]interface{}, num)
	for i := 0; i < num; i++ {
		docs[i] = map[string]interface{}{
			"title":   fmt.Sprintf("Test Document %d", i+1),
			"content": "This is a benchmark document.",
		}
	}
	return docs
}

func BenchmarkConcurrentInsert(b *testing.B) {
	client, err := NewClient()
	if err != nil {
		b.Fatalf("failed to create client: %v", err)
	}

	docCounts := []int{100, 1000, 10000}

	for _, count := range docCounts {
		// BulkInsert のベンチマーク
		b.Run(fmt.Sprintf("BulkInsert/%d_docs", count), func(b *testing.B) {
			docs := generateDocs(count)
			indexName := fmt.Sprintf("benchmark-bulk-%d", count)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				// --- セットアップ ---
				// 既存のインデックスを削除 (エラーは無視)
				_, _ = client.baseClient.Indices.Delete([]string{indexName})
				// 新しいインデックスを作成
				_, err := client.baseClient.Indices.Create(indexName)
				if err != nil {
					b.Fatalf("failed to create index: %v", err)
				}
				// --- セットアップ完了 ---
				b.StartTimer()

				if err := client.BulkInsert(context.Background(), indexName, docs); err != nil {
					b.Fatalf("BulkInsert failed: %v", err)
				}
			}
		})

		// BulkInsertConcurrent のベンチマーク
		b.Run(fmt.Sprintf("BulkInsertConcurrent/%d_docs", count), func(b *testing.B) {
			docs := generateDocs(count)
			indexName := fmt.Sprintf("benchmark-concurrent-%d", count)
			chunkSize := 100

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				// --- セットアップ ---
				// 既存のインデックスを削除 (エラーは無視)
				_, _ = client.baseClient.Indices.Delete([]string{indexName})
				// 新しいインデックスを作成
				_, err := client.baseClient.Indices.Create(indexName)
				if err != nil {
					b.Fatalf("failed to create index: %v", err)
				}
				// --- セットアップ完了 ---
				b.StartTimer()

				if err := client.BulkInsertConcurrent(context.Background(), indexName, docs, chunkSize); err != nil {
					b.Fatalf("BulkInsertConcurrent failed: %v", err)
				}
			}
		})

		// BulkInsertConcurrentV2 のベンチマーク
		b.Run(fmt.Sprintf("BulkInsertConcurrentV2/%d_docs", count), func(b *testing.B) {
			docs := generateDocs(count)
			indexName := fmt.Sprintf("benchmark-concurrent-v2-%d", count)
			numWorkers := 4

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				// --- セットアップ ---
				// 既存のインデックスを削除 (エラーは無視)
				_, _ = client.baseClient.Indices.Delete([]string{indexName})
				// 新しいインデックスを作成
				_, err := client.baseClient.Indices.Create(indexName)
				if err != nil {
					b.Fatalf("failed to create index: %v", err)
				}
				// --- セットアップ完了 ---
				b.StartTimer()

				if err := client.BulkInsertConcurrentV2(context.Background(), indexName, docs, numWorkers); err != nil {
					b.Fatalf("BulkInsertConcurrentV2 failed: %v", err)
				}
			}
		})

		// BulkInsertConcurrentV3 のベンチマーク
		b.Run(fmt.Sprintf("BulkInsertConcurrentV3/%d_docs", count), func(b *testing.B) {
			docs := generateDocs(count)
			indexName := fmt.Sprintf("benchmark-concurrent-v3-%d", count)
			numWorkers := 4

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				// --- セットアップ ---
				// 既存のインデックスを削除 (エラーは無視)
				_, _ = client.baseClient.Indices.Delete([]string{indexName})
				// 新しいインデックスを作成
				_, err := client.baseClient.Indices.Create(indexName)
				if err != nil {
					b.Fatalf("failed to create index: %v", err)
				}
				// --- セットアップ完了 ---
				b.StartTimer()

				if err := client.BulkInsertConcurrentV3(context.Background(), indexName, docs, numWorkers); err != nil {
					b.Fatalf("BulkInsertConcurrentV3 failed: %v", err)
				}
			}
		})
	}
}
