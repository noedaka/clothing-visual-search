package mlclient

import (
	"context"
	"time"

	pb "github.com/noedaka/clothing-visual-search/backend/internal/mlpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.MLServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:   conn,
		client: pb.NewMLServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetEmbedding(
	ctx context.Context,
	image []byte,
	imageFormat string,
) ([]float32, error) {
	var req pb.ImageRequest

	req.SetImageData(image)
	req.SetImageFormat(imageFormat)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

	resp, err := c.client.GetEmbedding(ctx, &req)
	if err != nil {
		return nil, err
	}

	return resp.GetEmbedding(), nil
}
