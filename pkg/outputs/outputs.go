package outputs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/puppetlabs/leg/encoding/transfer"
	"github.com/puppetlabs/relay-core/pkg/model"
)

const MetadataAPIURLEnvName = "METADATA_API_URL"

var (
	ErrOutputsClientKeyEmpty      = errors.New("key is required but was empty")
	ErrOutputsClientTaskNameEmpty = errors.New("taskName is required but was empty")
	ErrOutputsClientEnvVarMissing = errors.New(MetadataAPIURLEnvName + " was expected but was empty")
	ErrOutputsClientNotFound      = errors.New("output was not found")
)

// OutputsClient is a client for storing task outputs in
// the nebula outputs storage.
type OutputsClient interface {
	SetOutput(ctx context.Context, key string, value interface{}) error
	SetOutputMetadata(ctx context.Context, key string, metadata *model.StepOutputMetadata) error
}

// DefaultOutputsClient uses the default net/http.Client to
// store task output values.
type DefaultOutputsClient struct {
	apiURL *url.URL
}

func (c DefaultOutputsClient) SetOutput(ctx context.Context, key string, value interface{}) error {
	if key == "" {
		return ErrOutputsClientKeyEmpty
	}

	loc := *c.apiURL
	loc.Path = path.Join(loc.Path, key)

	encoded, err := json.Marshal(transfer.JSONInterface{Data: value})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", loc.String(), bytes.NewBuffer(encoded))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	return nil
}

func (c DefaultOutputsClient) SetOutputMetadata(ctx context.Context, key string, metadata *model.StepOutputMetadata) error {
	if key == "" {
		return ErrOutputsClientKeyEmpty
	}

	loc := *c.apiURL
	loc.Path = path.Join(loc.Path, key, "metadata")

	encoded, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", loc.String(), bytes.NewBuffer(encoded))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	return nil
}

func NewDefaultOutputsClient(location *url.URL) OutputsClient {
	return &DefaultOutputsClient{apiURL: location}
}

func NewDefaultOutputsClientFromNebulaEnv() (OutputsClient, error) {
	locStr := os.Getenv(MetadataAPIURLEnvName)

	if locStr == "" {
		return nil, ErrOutputsClientEnvVarMissing
	}

	loc, err := url.Parse(locStr)
	if err != nil {
		return nil, err
	}

	loc.Path = "/outputs"

	return NewDefaultOutputsClient(loc), nil
}
