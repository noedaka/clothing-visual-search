package milvusclient

import (
	"context"
	"fmt"
	"log"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

//TODO: Заменить типы id на int64 везде

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

	return &Client{client: c}, nil
}

func (c *Client) CreateproductCollection(ctx context.Context) error {
	has, err := c.client.HasCollection(ctx, "products")
	if err != nil {
		return err
	}

	if has {
		log.Printf("Collection already exists")
		return nil
	}

	schema := &entity.Schema{
		CollectionName: "products",
		Description:    "Product image embeddings",
		AutoID:         false,
		Fields: []*entity.Field{
			{
				Name:       "product_id",
				DataType:   entity.FieldTypeInt64,
				PrimaryKey: true,
				AutoID:     false,
			},
			{
				Name:       "embedding",
				DataType:   entity.FieldTypeFloatVector,
				TypeParams: map[string]string{"dim": "2048"},
			},
		},
	}

	err = c.client.CreateCollection(ctx, schema, 2)
	if err != nil {
		return err
	}

	idx, err := entity.NewIndexIvfFlat(entity.L2, 128)
	if err != nil {
		return err
	}

	err = c.client.CreateIndex(ctx, "products", "embedding", idx, false)
	if err != nil {
		return err
	}

	log.Println("Collection and index created")
	return nil
}

func (c *Client) InsertEmbedding(ctx context.Context, productID int64, embedding []float32) error {
	idColumn := entity.NewColumnInt64("product_id", []int64{productID})
	embedColumn := entity.NewColumnFloatVector("embedding", 2048, [][]float32{embedding})

	_, err := c.client.Insert(ctx, "products", "", idColumn, embedColumn)
	return err
}

func (c *Client) SearchSimilar(ctx context.Context, queryVector []float32, topK int) ([]int64, error) {
	searchParam, err := entity.NewIndexIvfFlatSearchParam(16)
	if err != nil {
		return nil, err
	}

	searchResults, err := c.client.Search(
		ctx,
		"products",             // collectionName
		[]string{},             // partitionNames (пусто)
		"",                     // expr (фильтр)
		[]string{"product_id"}, // outputFields – какие поля вернуть
		[]entity.Vector{entity.FloatVector(queryVector)}, // векторы запроса
		"embedding", // vectorField
		entity.L2,   // metricType
		topK,        // topK
		searchParam, // searchParam
	)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	var productIDs []int64
	for _, result := range searchResults {
		col := result.Fields.GetColumn("product_id")
		if col == nil {
			continue
		}
		if int64Col, ok := col.(*entity.ColumnInt64); ok {
			ids := int64Col.Data()
			productIDs = append(productIDs, ids...)
		} else {
			return nil, fmt.Errorf("unexpected column type for product_id")
		}
	}
	return productIDs, nil
}

func (c *Client) DeleteByProductID(ctx context.Context, productID int64) error {
	expr := fmt.Sprintf("product_id == %d", productID)
	return c.client.Delete(ctx, "products", "", expr)
}
