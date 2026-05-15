package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	repogithub "repoready/internal/github"
)

func friendlyError(err error) error {
	switch {
	case errors.Is(err, repogithub.ErrInvalidURL):
		return fmt.Errorf("invalid GitHub URL. Use a URL like https://github.com/user/repo")
	case errors.Is(err, repogithub.ErrRepositoryNotFound):
		return fmt.Errorf("GitHub repository not found")
	case errors.Is(err, repogithub.ErrPrivateOrUnavailable):
		return fmt.Errorf("repository not found or private. Set GITHUB_TOKEN for private repositories")
	case errors.Is(err, repogithub.ErrFolderExists):
		return fmt.Errorf("target folder already exists. Choose another clone directory or remove the existing folder")
	case errors.Is(err, os.ErrPermission):
		return fmt.Errorf("permission denied while accessing the requested path")
	case strings.Contains(strings.ToLower(err.Error()), "i/o timeout"),
		strings.Contains(strings.ToLower(err.Error()), "network"),
		strings.Contains(strings.ToLower(err.Error()), "lookup github.com"):
		return fmt.Errorf("network error while contacting GitHub. Check your connection and try again")
	default:
		return err
	}
}
