package github

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

var (
	ErrInvalidURL           = errors.New("invalid GitHub repository URL")
	ErrRepositoryNotFound   = errors.New("GitHub repository not found")
	ErrPrivateOrUnavailable = errors.New("repository is private or unavailable")
	ErrFolderExists         = errors.New("target folder already exists")
)

type Repository struct {
	URL   string
	Owner string
	Name  string
}

type Client struct {
	httpClient *http.Client
	token      string
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 15 * time.Second},
		token:      strings.TrimSpace(os.Getenv("GITHUB_TOKEN")),
	}
}

func ParseRepositoryURL(raw string) (Repository, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return Repository{}, ErrInvalidURL
	}

	if strings.HasPrefix(raw, "git@github.com:") {
		re := regexp.MustCompile(`^git@github\.com:([^/]+)/([^/]+?)(?:\.git)?$`)
		matches := re.FindStringSubmatch(raw)
		if len(matches) != 3 {
			return Repository{}, ErrInvalidURL
		}
		return Repository{URL: raw, Owner: matches[1], Name: matches[2]}, nil
	}

	parsed, err := url.Parse(raw)
	if err != nil || parsed.Hostname() != "github.com" {
		return Repository{}, ErrInvalidURL
	}
	if parsed.Scheme != "https" && parsed.Scheme != "http" {
		return Repository{}, ErrInvalidURL
	}

	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return Repository{}, ErrInvalidURL
	}

	name := strings.TrimSuffix(parts[1], ".git")
	if name == "" {
		return Repository{}, ErrInvalidURL
	}

	return Repository{URL: raw, Owner: parts[0], Name: name}, nil
}

func SafeClonePath(root, owner, repo string) (string, error) {
	root = filepath.Clean(root)
	target := filepath.Join(root, owner, repo)
	rel, err := filepath.Rel(root, target)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("unsafe clone path")
	}
	return target, nil
}

func (c *Client) Validate(ctx context.Context, repo Repository) error {
	requestURL := fmt.Sprintf("https://api.github.com/repos/%s/%s", repo.Owner, repo.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("network error while contacting GitHub: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return ErrPrivateOrUnavailable
	case http.StatusNotFound:
		if c.token == "" {
			return ErrPrivateOrUnavailable
		}
		return ErrRepositoryNotFound
	default:
		return fmt.Errorf("GitHub returned unexpected status %s", resp.Status)
	}
}

func (c *Client) Clone(ctx context.Context, repo Repository, target string) error {
	if stat, err := os.Stat(target); err == nil && stat.IsDir() {
		entries, readErr := os.ReadDir(target)
		if readErr != nil {
			return readErr
		}
		if len(entries) > 0 {
			return ErrFolderExists
		}
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}

	opts := &git.CloneOptions{
		URL:      repo.URL,
		Progress: io.Discard,
	}
	if c.token != "" {
		opts.Auth = &githttp.BasicAuth{
			Username: "x-access-token",
			Password: c.token,
		}
	}

	_, err := git.PlainCloneContext(ctx, target, false, opts)
	return err
}
