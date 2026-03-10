package filesclient

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	filesv1 "github.com/agynio/gateway/gen/agynio/api/files/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const uploadChunkSize = 64 * 1024

type Client struct {
	conn   *grpc.ClientConn
	client filesv1.FilesServiceClient
}

type UploadResult struct {
	ID          string    `json:"id"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"contentType"`
	SizeBytes   int64     `json:"sizeBytes"`
	CreatedAt   time.Time `json:"createdAt"`
}

func NewClient(target string) (*Client, error) {
	if strings.TrimSpace(target) == "" {
		return nil, fmt.Errorf("target is required")
	}

	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:   conn,
		client: filesv1.NewFilesServiceClient(conn),
	}, nil
}

func (c *Client) Upload(ctx context.Context, filename, contentType string, sizeBytes int64, body io.Reader) (UploadResult, error) {
	if c == nil || c.client == nil {
		panic("files client is not initialized")
	}
	if ctx == nil {
		panic("context is required")
	}
	if body == nil {
		panic("upload body is required")
	}

	stream, err := c.client.UploadFile(ctx)
	if err != nil {
		return UploadResult{}, err
	}

	metadata := &filesv1.UploadFileMetadata{
		Filename:    filename,
		ContentType: contentType,
		SizeBytes:   sizeBytes,
	}
	if err := stream.Send(&filesv1.UploadFileRequest{Payload: &filesv1.UploadFileRequest_Metadata{Metadata: metadata}}); err != nil {
		return UploadResult{}, err
	}

	buffer := make([]byte, uploadChunkSize)
	for {
		readBytes, readErr := body.Read(buffer)
		if readBytes > 0 {
			chunk := make([]byte, readBytes)
			copy(chunk, buffer[:readBytes])
			if err := stream.Send(&filesv1.UploadFileRequest{Payload: &filesv1.UploadFileRequest_Chunk{Chunk: &filesv1.UploadFileChunk{Data: chunk}}}); err != nil {
				return UploadResult{}, err
			}
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return UploadResult{}, fmt.Errorf("read upload body: %w", readErr)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return UploadResult{}, err
	}
	file := resp.GetFile()
	if file == nil {
		return UploadResult{}, fmt.Errorf("upload response missing file info")
	}
	createdAt := file.GetCreatedAt()
	if createdAt == nil {
		return UploadResult{}, fmt.Errorf("upload response missing created at")
	}

	return UploadResult{
		ID:          file.GetId(),
		Filename:    file.GetFilename(),
		ContentType: file.GetContentType(),
		SizeBytes:   file.GetSizeBytes(),
		CreatedAt:   createdAt.AsTime().UTC(),
	}, nil
}
