// Package api handles interactions with the GitHub API.
package api

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
)

// contentsResponse matches the GitHub Contents API shape.
type contentsResponse struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
	Message  string `json:"message"` // populated on error
}

// GetReadme fetches and decodes the README content for owner/repo.
// If ref is non-empty the file is fetched at that ref.
func GetReadme(owner, repo, ref string) (string, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return "", fmt.Errorf("creating API client: %w", err)
	}

	path := fmt.Sprintf("repos/%s/%s/readme", owner, repo)
	if ref != "" {
		path += "?ref=" + ref
	}

	var resp contentsResponse
	if err := client.Get(path, &resp); err != nil {
		return "", fmt.Errorf("fetching README for %s/%s: %w", owner, repo, err)
	}

	if resp.Encoding != "base64" {
		return "", fmt.Errorf("unexpected encoding %q", resp.Encoding)
	}

	// strip newlines in the base64 payload
	raw := strings.ReplaceAll(resp.Content, "\n", "")
	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return "", fmt.Errorf("decoding README content: %w", err)
	}

	return string(decoded), nil
}
