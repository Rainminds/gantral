package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/Rainminds/gantral/internal/artifact"
	"github.com/Rainminds/gantral/pkg/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

// Store implements artifact.ImmutableStore using AWS S3.
type Store struct {
	client *s3.Client
	bucket string
}

// NewStore creates a new S3 backward store.
func NewStore(client *s3.Client, bucket string) *Store {
	return &Store{
		client: client,
		bucket: bucket,
	}
}

// WriteArtifact persists the artifact immutably.
// It uses If-None-Match: "*" to prevent overwrites, enforcing WORM.
func (s *Store) WriteArtifact(ctx context.Context, teamID string, art *models.CommitmentArtifact) error {
	data, err := json.Marshal(art)
	if err != nil {
		return fmt.Errorf("serialization failed: %w", err)
	}

	key := fmt.Sprintf("%s/%s.json", teamID, art.ArtifactHash)

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("application/json"),
		IfNoneMatch: aws.String("*"),
	})

	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() == "PreconditionFailed" {
				return artifact.ErrImmutableViolation
			}
		}
		return fmt.Errorf("failed to put object: %w", err)
	}

	return nil
}

// GetArtifact retrieves an artifact by hash.
func (s *Store) GetArtifact(ctx context.Context, teamID string, artifactHash string) (*models.CommitmentArtifact, error) {
	key := fmt.Sprintf("%s/%s.json", teamID, artifactHash)

	res, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		var noSuchKey *types.NoSuchKey
		if errors.As(err, &noSuchKey) {
			return nil, artifact.ErrArtifactNotFound
		}
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "NoSuchKey" {
			return nil, artifact.ErrArtifactNotFound
		}
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	var art models.CommitmentArtifact
	if err := json.Unmarshal(data, &art); err != nil {
		return nil, fmt.Errorf("deserialize failed: %w", err)
	}

	return &art, nil
}
