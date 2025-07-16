package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type Client struct {
	baseClient *elasticsearch.Client
}

func NewClient() (*Client, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}
	return &Client{
		baseClient: es,
	}, nil
}

func (c *Client) BulkInsert(ctx context.Context, index string, docs []map[string]interface{}) error {
	bulkCfg := esutil.BulkIndexerConfig{
		Client: c.baseClient,
		Index:  index,
	}

	indexer, err := esutil.NewBulkIndexer(bulkCfg)
	if err != nil {
		return fmt.Errorf("failed to create bulk indexer: %w", err)
	}
	for i, doc := range docs {
		data, err := json.Marshal(doc)
		if err != nil {
			return fmt.Errorf("failed to marshal document %d: %w", i+1, err)
		}
		err = indexer.Add(
			ctx,
			esutil.BulkIndexerItem{
				Action: "index",
				Body:   strings.NewReader(string(data)),
			},
		)
		if err != nil {
			return fmt.Errorf("failed to add document %d to bulk indexer: %w", i+1, err)
		}
	}
	if err := indexer.Close(ctx); err != nil {
		return fmt.Errorf("failed to close bulk indexer: %w", err)
	}
	return nil
}

func (c *Client) BulkInsertConcurrentWithTimeSleep(ctx context.Context, index string, docs []map[string]interface{}, chunkSize int, sleepSec int) {
	for i := 0; i < len(docs); i += chunkSize {
		end := i + chunkSize
		if end > len(docs) {
			end = len(docs)
		}
		chunk := docs[i:end]
		go func(chunk []map[string]interface{}) {
			if err := c.BulkInsert(ctx, index, chunk); err != nil {
				// log.Printf("failed to bulk insert: %v", err)
			}
		}(chunk)
	}

	// log.Printf("Waiting for all goroutines to finish...")
	time.Sleep(time.Duration(sleepSec) * time.Second)
	// log.Println("Done.")
}

func (c *Client) BulkInsertConcurrent(ctx context.Context, index string, docs []map[string]interface{}, chunkSize int) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	for i := 0; i < len(docs); i += chunkSize {
		end := i + chunkSize
		if end > len(docs) {
			end = len(docs)
		}
		chunk := docs[i:end]

		wg.Add(1)
		go func(chunk []map[string]interface{}) {
			defer wg.Done()
			if err := c.BulkInsert(ctx, index, chunk); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}(chunk)
	}

	wg.Wait()
	close(errCh)

	if err := <-errCh; err != nil {
		return fmt.Errorf("error occurred during concurrent bulk insert: %w", err)
	}

	return nil
}

// BulkInsertConcurrentV2 は、単一のBulkIndexerを複数のgoroutineで共有する、より効率的な並行処理です。
// BulkIndexerが内部的に並行処理を行うため、このアプローチが推奨されます。
func (c *Client) BulkInsertConcurrentV2(ctx context.Context, index string, docs []map[string]interface{}, numWorkers int) error {
	// esutil.BulkIndexerは内部で並行処理をサポートしています。
	// NumWorkersを設定すると、その数だけワーカーgoroutineが起動し、リクエストを並行して送信します。
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:     c.baseClient,
		Index:      index,
		NumWorkers: numWorkers,
	})
	if err != nil {
		return fmt.Errorf("failed to create bulk indexer: %w", err)
	}

	// ドキュメントをチャネルに投入
	docCh := make(chan map[string]interface{})
	go func() {
		for _, doc := range docs {
			docCh <- doc
		}
		close(docCh)
	}()

	// 複数のgoroutineでBulkIndexerにドキュメントを追加
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for doc := range docCh {
				data, err := json.Marshal(doc)
				if err != nil {
					// log.Printf("failed to marshal document: %v", err)
					continue
				}
				err = bi.Add(
					ctx,
					esutil.BulkIndexerItem{
						Action: "index",
						Body:   strings.NewReader(string(data)),
					},
				)
				if err != nil {
					// log.Printf("failed to add document to bulk indexer: %v", err)
				}
			}
		}()
	}

	wg.Wait()

	// BulkIndexerを閉じて、すべてのドキュメントが処理されるのを待つ
	if err := bi.Close(ctx); err != nil {
		return fmt.Errorf("failed to close bulk indexer: %w", err)
	}

	// stats := bi.Stats()
	// log.Printf("V2: Indexed [%d] documents with [%d] workers", stats.NumIndexed, numWorkers)
	return nil
}

// BulkInsertConcurrentV3 は、BulkIndexerの内部並行処理に完全に任せる最もシンプルな実装です。
// クライアント側のオーバーヘッドが最小限になります。
func (c *Client) BulkInsertConcurrentV3(ctx context.Context, index string, docs []map[string]interface{}, numWorkers int) error {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:     c.baseClient,
		Index:      index,
		NumWorkers: numWorkers,
	})
	if err != nil {
		return fmt.Errorf("failed to create bulk indexer: %w", err)
	}

	for _, doc := range docs {
		data, err := json.Marshal(doc)
		if err != nil {
			return fmt.Errorf("failed to marshal document: %w", err)
		}
		err = bi.Add(
			ctx,
			esutil.BulkIndexerItem{
				Action: "index",
				Body:   strings.NewReader(string(data)),
			},
		)
		if err != nil {
			return fmt.Errorf("failed to add document to bulk indexer: %w", err)
		}
	}

	if err := bi.Close(ctx); err != nil {
		return fmt.Errorf("failed to close bulk indexer: %w", err)
	}

	return nil
}
