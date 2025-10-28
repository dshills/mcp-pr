package git

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
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
	return c.GetStagedDiffContext(context.Background())
}

// GetStagedDiffContext retrieves diff of staged changes with context
func (c *Client) GetStagedDiffContext(ctx context.Context) (string, error) {
	// Add 30s timeout for git operations
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "diff", "--staged")
	cmd.Dir = c.repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("git diff --staged timed out after 30s")
		}
		return "", fmt.Errorf("failed to get staged diff: %w (output: %s)", err, string(output))
	}

	return string(output), nil
}

// GetUnstagedDiff retrieves diff of unstaged changes using `git diff`
func (c *Client) GetUnstagedDiff() (string, error) {
	return c.GetUnstagedDiffContext(context.Background())
}

// GetUnstagedDiffContext retrieves diff of unstaged changes with context
func (c *Client) GetUnstagedDiffContext(ctx context.Context) (string, error) {
	// Add 30s timeout for git operations
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "diff")
	cmd.Dir = c.repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("git diff timed out after 30s")
		}
		return "", fmt.Errorf("failed to get unstaged diff: %w (output: %s)", err, string(output))
	}

	return string(output), nil
}

// GetCommitDiff retrieves diff for a specific commit using `git show <sha>`
func (c *Client) GetCommitDiff(commitSHA string) (string, error) {
	return c.GetCommitDiffContext(context.Background(), commitSHA)
}

// GetCommitDiffContext retrieves diff for a specific commit with context
func (c *Client) GetCommitDiffContext(ctx context.Context, commitSHA string) (string, error) {
	// Validate commit SHA first
	if err := c.ValidateCommitContext(ctx, commitSHA); err != nil {
		return "", err
	}

	// Add 30s timeout for git operations
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "show", commitSHA)
	cmd.Dir = c.repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("git show timed out after 30s")
		}
		return "", fmt.Errorf("failed to get commit diff: %w (output: %s)", err, string(output))
	}

	return string(output), nil
}

// ValidateCommit checks if a commit SHA exists using `git rev-parse --verify`
func (c *Client) ValidateCommit(commitSHA string) error {
	return c.ValidateCommitContext(context.Background(), commitSHA)
}

// ValidateCommitContext checks if a commit SHA exists with context
func (c *Client) ValidateCommitContext(ctx context.Context, commitSHA string) error {
	// Add 10s timeout for validation
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--verify", commitSHA)
	cmd.Dir = c.repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("git rev-parse timed out after 10s for commit %s", commitSHA)
		}
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
