package milvusclient

import (
	"context"
	"fmt"
	"log"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type Client struct {
	client client.Client
}

func NewClient(ctx context.Context, addr string) (*Client, error) {
	c, err := client.NewClient(ctx, client.Config{
		Address: addr,
	})
	if err != nil {
		return nil, err
	}

	if err := createCollection(c, ctx); err != nil {
		return nil, err
	}

	return &Client{client: c}, nil
}

// InsertEmbedding вставляет эмбеддинг для конкретного изображения
func (c *Client) InsertEmbedding(ctx context.Context, imageID, productID int64, embedding []float32) error {
	collectionName := "product_images"

	idCol := entity.NewColumnInt64("image_id", []int64{imageID})
	pidCol := entity.NewColumnInt64("product_id", []int64{productID})
	embedCol := entity.NewColumnFloatVector("embedding", 2048, [][]float32{embedding})

	_, err := c.client.Insert(ctx, collectionName, "", idCol, pidCol, embedCol)
	if err != nil {
		return fmt.Errorf("insert embedding failed: %w", err)
	}
	return nil
}

// SearchSimilar ищет до topK ближайших векторов с порогом похожести (Range Search)
// и возвращает уникальные product_id.
//
// similarityThreshold — максимальное L2-расстояние (чем меньше значение — тем строже фильтр).
// Пример хороших стартовых значений для нормализованных эмбеддингов: 0.6 – 1.2 (подбирайте экспериментально).
func (c *Client) SearchSimilar(ctx context.Context, queryVector []float32, topK int, threshold float64) ([]int64, error) {
	collectionName := "product_images"

	searchParam, err := entity.NewIndexIvfFlatSearchParam(16)
	if err != nil {
		return nil, fmt.Errorf("create search param failed: %w", err)
	}

	searchParam.AddRadius(threshold)
	searchParam.AddRangeFilter(0.0)

	searchResults, err := c.client.Search(
		ctx,
		collectionName,         // имя коллекции
		[]string{},             // partitions (пусто = все)
		"",                     // expr — скалярный фильтр (не используем)
		[]string{"product_id"}, // возвращаемые поля
		[]entity.Vector{entity.FloatVector(queryVector)},
		"embedding", // поле с вектором
		entity.L2,   // метрика расстояния
		topK,        // максимальное количество результатов
		searchParam,
	)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	productIDSet := make(map[int64]struct{})
	for _, result := range searchResults {
		col := result.Fields.GetColumn("product_id")
		if int64Col, ok := col.(*entity.ColumnInt64); ok {
			for _, pid := range int64Col.Data() {
				productIDSet[pid] = struct{}{}
			}
		}
	}

	productIDs := make([]int64, 0, len(productIDSet))
	for pid := range productIDSet {
		productIDs = append(productIDs, pid)
	}

	return productIDs, nil
}

func (c *Client) DeleteByImageID(ctx context.Context, imageID int64) error {
	expr := fmt.Sprintf("image_id == %d", imageID)
	err := c.client.Delete(ctx, "product_images", "", expr)
	if err != nil {
		return fmt.Errorf("delete by image_id failed: %w", err)
	}
	return nil
}


func createCollection(c client.Client, ctx context.Context) error {
	collectionName := "product_images"
	has, err := c.HasCollection(ctx, collectionName)
	if err != nil {
		return err
	}

	if !has {
		schema := &entity.Schema{
			CollectionName: collectionName,
			Description:    "Embeddings of product images",
			AutoID:         false,
			Fields: []*entity.Field{
				{
					Name:       "image_id",
					DataType:   entity.FieldTypeInt64,
					PrimaryKey: true,
					AutoID:     false,
				},
				{
					Name:       "product_id",
					DataType:   entity.FieldTypeInt64,
					PrimaryKey: false,
				},
				{
					Name:       "embedding",
					DataType:   entity.FieldTypeFloatVector,
					TypeParams: map[string]string{"dim": "2048"}, // размерность вектора
				},
			},
		}

		if err := c.CreateCollection(ctx, schema, 2); err != nil {
			return fmt.Errorf("create collection failed: %w", err)
		}

		idx, err := entity.NewIndexIvfFlat(entity.L2, 128)
		if err != nil {
			return fmt.Errorf("create index params failed: %w", err)
		}
		if err := c.CreateIndex(ctx, collectionName, "embedding", idx, false); err != nil {
			return fmt.Errorf("create index failed: %w", err)
		}
		log.Println("Collection product_images created with index")
	} else {
		log.Println("Collection product_images already exists")
	}

	err = c.LoadCollection(ctx, "product_images", false)
	if err != nil {
		return fmt.Errorf("failed to load collection: %w", err)
	}

	return nil
}