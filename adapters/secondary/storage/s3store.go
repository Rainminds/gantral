package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/Rainminds/gantral/internal/artifact"
	"github.com/Rainminds/gantral/pkg/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

// S3ClientAPI defines the interface for S3 operations we use, to allow mocking.
type S3ClientAPI interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

// S3Store implements ImmutableStore using AWS S3.
type S3Store struct {
	client S3ClientAPI
	bucket string
	prefix string
}

// NewS3Store creates a new S3-backed immutable store.
func NewS3Store(client S3ClientAPI, bucket, prefix string) *S3Store {
	return &S3Store{
		client: client,
		bucket: bucket,
		prefix: prefix,
	}
}

// WriteArtifact persists the artifact to S3.
// It implements WORM guarantees by using `IfNoneMatch: "*"` to fail if the object already exists.
func (s *S3Store) WriteArtifact(ctx context.Context, teamID string, art *models.CommitmentArtifact) error {
	payload, err := json.Marshal(art)
	if err != nil {
		return fmt.Errorf("failed to marshal artifact for S3: %w", err)
	}

	key := s.buildKey(teamID, art.ArtifactHash)

	// Enforce WORM: write only if object does not exist.
	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(payload),
		ContentType: aws.String("application/json"),
		IfNoneMatch: aws.String("*"),
	})

	if err != nil {
		// Detect overwrite attempt
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "PreconditionFailed" {
			slog.Warn("WORM violation attempted on S3", "key", key, "artifact_hash", art.ArtifactHash)
			return artifact.ErrImmutableViolation
		}

		// Fail closed on any other error
		return fmt.Errorf("failed to write artifact to S3: %w", err)
	}

	return nil
}

// GetArtifact retrieves the artifact from S3 using its hash.
func (s *S3Store) GetArtifact(ctx context.Context, teamID string, artifactHash string) (*models.CommitmentArtifact, error) {
	key := s.buildKey(teamID, artifactHash)

	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		// Map S3 NoSuchKey to our ErrArtifactNotFound
		var noSuchKey *types.NoSuchKey
		if errors.As(err, &noSuchKey) {
			return nil, artifact.ErrArtifactNotFound
		}

		var apiErr *types.NotFound
		if errors.As(err, &apiErr) {
			return nil, artifact.ErrArtifactNotFound
		}

		return nil, fmt.Errorf("failed to get artifact from S3: %w", err)
	}
	defer output.Body.Close()

	body, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read S3 response body: %w", err)
	}

	var art models.CommitmentArtifact
	if err := json.Unmarshal(body, &art); err != nil {
		return nil, fmt.Errorf("failed to unmarshal artifact from S3: %w", err)
	}

	return &art, nil
}

func (s *S3Store) buildKey(teamID, artifactHash string) string {
	if s.prefix == "" {
		return fmt.Sprintf("%s/%s.json", teamID, artifactHash)
	}
	return fmt.Sprintf("%s/%s/%s.json", s.prefix, teamID, artifactHash)
}
