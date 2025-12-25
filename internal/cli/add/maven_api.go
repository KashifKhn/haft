package add

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type MavenSearchResponse struct {
	Response struct {
		NumFound int             `json:"numFound"`
		Docs     []MavenArtifact `json:"docs"`
	} `json:"response"`
}

type MavenArtifact struct {
	GroupId       string `json:"g"`
	ArtifactId    string `json:"a"`
	LatestVersion string `json:"latestVersion"`
	Version       string `json:"v"`
}

type MavenClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewMavenClient() *MavenClient {
	return &MavenClient{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		baseURL: "https://search.maven.org/solrsearch/select",
	}
}

func (c *MavenClient) VerifyDependency(groupId, artifactId string) (*MavenArtifact, error) {
	query := fmt.Sprintf("g:%s AND a:%s", groupId, artifactId)
	params := url.Values{
		"q":    {query},
		"rows": {"1"},
		"wt":   {"json"},
	}

	reqURL := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Maven Central: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Maven Central returned status %d", resp.StatusCode)
	}

	var searchResp MavenSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to parse Maven Central response: %w", err)
	}

	if searchResp.Response.NumFound == 0 {
		return nil, nil
	}

	return &searchResp.Response.Docs[0], nil
}

func (c *MavenClient) SearchDependencies(query string, limit int) ([]MavenArtifact, error) {
	if limit <= 0 {
		limit = 20
	}

	params := url.Values{
		"q":    {query},
		"rows": {fmt.Sprintf("%d", limit)},
		"wt":   {"json"},
	}

	reqURL := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Maven Central: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Maven Central returned status %d", resp.StatusCode)
	}

	var searchResp MavenSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to parse Maven Central response: %w", err)
	}

	return searchResp.Response.Docs, nil
}
