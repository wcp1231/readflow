package shaarli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/integration/webhook"
	"github.com/ncarlier/readflow/pkg/model"
)

// shaarliEntry is the structure definition of a Shaarli article
type shaarliEntry struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	URL         *string `json:"url,omitempty"`
	Private     bool    `json:"private"`
	Created     string  `json:"created"`
	Updated     string  `json:"updated"`
}

// shaarliProviderConfig is the structure definition of a Shaarli API configuration
type shaarliProviderConfig struct {
	Endpoint string `json:"endpoint"`
	Secret   string `json:"secret"`
	Private  bool   `json:"private"`
}

// shaarliProvider is the structure definition of a Shaarli webhook provider
type shaarliProvider struct {
	config   shaarliProviderConfig
	endpoint *url.URL
}

func newShaarliProvider(srv model.OutgoingWebhook, conf config.Config) (webhook.Provider, error) {
	config := shaarliProviderConfig{}
	if err := json.Unmarshal([]byte(srv.Config), &config); err != nil {
		return nil, err
	}

	// Validate endpoint URL
	endpoint, err := url.ParseRequestURI(config.Endpoint)
	if err != nil {
		return nil, err
	}

	// Validate config
	if config.Secret == "" {
		return nil, fmt.Errorf("shaarli: missing secret")
	}

	provider := &shaarliProvider{
		config:   config,
		endpoint: endpoint,
	}

	return provider, nil
}

// Send article to Shaarli endpoint.
func (p *shaarliProvider) Send(ctx context.Context, article model.Article) error {
	token, err := p.getAccessToken()
	if err != nil {
		return err
	}

	entry := shaarliEntry{
		Title:       article.Title,
		Description: article.Text,
		URL:         article.URL,
		Private:     p.config.Private,
		Created:     article.CreatedAt.Format(time.RFC3339),
		Updated:     article.UpdatedAt.Format(time.RFC3339),
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(entry)

	req, err := http.NewRequest("POST", p.getAPIEndpoint("/api/v1/links"), b)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", constant.UserAgent)
	req.Header.Set("Content-Type", constant.ContentTypeJSON)
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 201 {
		if err == nil {
			err = fmt.Errorf("bad status code: %d", resp.StatusCode)
		}
		return err
	}

	return nil
}

func (p *shaarliProvider) getAPIEndpoint(path string) string {
	baseURL := *p.endpoint
	baseURL.Path = path
	return baseURL.String()
}

func (p *shaarliProvider) getAccessToken() (string, error) {
	claims := new(jwt.StandardClaims)
	claims.IssuedAt = time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString([]byte(p.config.Secret))
}

func init() {
	webhook.Register("shaarli", &webhook.Def{
		Name:   "Shaarli",
		Desc:   "Send article(s) to Shaarli instance.",
		Create: newShaarliProvider,
	})
}
