package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// Client provides git operations on a repository
type Client struct {
	repoPath string
}

// NewClient creates a new git client for the specified repository path
func NewClient(repoPath string) *Client {
	return &Client{
		repoPath: repoPath,
	}
}

// GetStagedDiff retrieves diff of staged changes using `git diff --staged`
func (c *Client) GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	cmd.Dir = c.repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get staged diff: %w (output: %s)", err, string(output))
	}

	return string(output), nil
}

// GetUnstagedDiff retrieves diff of unstaged changes using `git diff`
func (c *Client) GetUnstagedDiff() (string, error) {
	cmd := exec.Command("git", "diff")
	cmd.Dir = c.repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get unstaged diff: %w (output: %s)", err, string(output))
	}

	return string(output), nil
}

// GetCommitDiff retrieves diff for a specific commit using `git show <sha>`
func (c *Client) GetCommitDiff(commitSHA string) (string, error) {
	// Validate commit SHA first
	if err := c.ValidateCommit(commitSHA); err != nil {
		return "", err
	}

	cmd := exec.Command("git", "show", commitSHA)
	cmd.Dir = c.repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get commit diff: %w (output: %s)", err, string(output))
	}

	return string(output), nil
}

// ValidateCommit checks if a commit SHA exists using `git rev-parse --verify`
func (c *Client) ValidateCommit(commitSHA string) error {
	cmd := exec.Command("git", "rev-parse", "--verify", commitSHA)
	cmd.Dir = c.repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("invalid commit SHA %s: %w (output: %s)", commitSHA, err, string(output))
	}

	return nil
}

// IsGitRepository checks if the path is a valid git repository
func (c *Client) IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = c.repoPath

	err := cmd.Run()
	return err == nil
}

// GetRepositoryRoot returns the root directory of the repository
func (c *Client) GetRepositoryRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = c.repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get repository root: %w (output: %s)", err, string(output))
	}

	return strings.TrimSpace(string(output)), nil
}
