package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/dshills/mcp-pr/internal/git"
)

// setupTestRepo creates a temporary git repository for testing
func setupTestRepo(t *testing.T) (string, func()) {
	t.Helper()

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "mcp-git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
	}

	// Initialize git repo
	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@example.com"},
		{"git", "config", "user.name", "Test User"},
	}

	for _, cmd := range cmds {
		c := exec.Command(cmd[0], cmd[1:]...)
		c.Dir = tmpDir
		if err := c.Run(); err != nil {
			cleanup()
			t.Fatalf("Failed to run %v: %v", cmd, err)
		}
	}

	return tmpDir, cleanup
}

// createAndStageFile creates a file and stages it in git
func createAndStageFile(t *testing.T, repoPath, filename, content string) {
	t.Helper()

	filePath := filepath.Join(repoPath, filename)
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file %s: %v", filename, err)
	}

	cmd := exec.Command("git", "add", filename)
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file %s: %v", filename, err)
	}
}

// modifyFile modifies an existing file without staging
func modifyFile(t *testing.T, repoPath, filename, content string) {
	t.Helper()

	filePath := filepath.Join(repoPath, filename)
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to modify file %s: %v", filename, err)
	}
}

// commitChanges commits all staged changes
func commitChanges(t *testing.T, repoPath, message string) {
	t.Helper()

	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}
}

// TestGitClientStagedDiff tests retrieval of staged changes
func TestGitClientStagedDiff(t *testing.T) {
	repoPath, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create initial commit
	createAndStageFile(t, repoPath, "README.md", "# Test Repo\n")
	commitChanges(t, repoPath, "Initial commit")

	// Create and stage a new file
	createAndStageFile(t, repoPath, "main.go", `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`)

	// Get staged diff
	client := git.NewClient(repoPath)
	diff, err := client.GetStagedDiff()
	if err != nil {
		t.Fatalf("GetStagedDiff() error = %v", err)
	}

	if diff == "" {
		t.Fatal("GetStagedDiff() returned empty diff, want non-empty")
	}

	// Verify diff contains the new file
	if !contains(diff, "main.go") {
		t.Errorf("Diff doesn't contain 'main.go': %s", diff)
	}

	if !contains(diff, "Hello, World!") {
		t.Errorf("Diff doesn't contain expected content: %s", diff)
	}

	t.Logf("Staged diff:\n%s", diff)
}

// TestGitClientUnstagedDiff tests retrieval of unstaged changes
func TestGitClientUnstagedDiff(t *testing.T) {
	repoPath, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create initial file and commit
	createAndStageFile(t, repoPath, "config.json", `{"version": "1.0.0"}`)
	commitChanges(t, repoPath, "Initial commit")

	// Modify file without staging
	modifyFile(t, repoPath, "config.json", `{"version": "1.0.1", "env": "production"}`)

	// Get unstaged diff
	client := git.NewClient(repoPath)
	diff, err := client.GetUnstagedDiff()
	if err != nil {
		t.Fatalf("GetUnstagedDiff() error = %v", err)
	}

	if diff == "" {
		t.Fatal("GetUnstagedDiff() returned empty diff, want non-empty")
	}

	// Verify diff contains the changes
	if !contains(diff, "config.json") {
		t.Errorf("Diff doesn't contain 'config.json': %s", diff)
	}

	if !contains(diff, "1.0.1") {
		t.Errorf("Diff doesn't contain version change: %s", diff)
	}

	t.Logf("Unstaged diff:\n%s", diff)
}

// TestGitClientCommitDiff tests retrieval of specific commit diff
func TestGitClientCommitDiff(t *testing.T) {
	repoPath, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create initial commit
	createAndStageFile(t, repoPath, "initial.txt", "Initial content\n")
	commitChanges(t, repoPath, "Initial commit")

	// Create second commit
	createAndStageFile(t, repoPath, "feature.go", `package main

func NewFeature() string {
	return "feature"
}
`)
	commitChanges(t, repoPath, "Add feature")

	// Get latest commit SHA
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to get commit SHA: %v", err)
	}
	commitSHA := string(output[:len(output)-1]) // Remove trailing newline

	// Get commit diff
	client := git.NewClient(repoPath)
	diff, err := client.GetCommitDiff(commitSHA)
	if err != nil {
		t.Fatalf("GetCommitDiff() error = %v", err)
	}

	if diff == "" {
		t.Fatal("GetCommitDiff() returned empty diff, want non-empty")
	}

	// Verify diff contains the new file
	if !contains(diff, "feature.go") {
		t.Errorf("Diff doesn't contain 'feature.go': %s", diff)
	}

	if !contains(diff, "NewFeature") {
		t.Errorf("Diff doesn't contain expected content: %s", diff)
	}

	t.Logf("Commit diff:\n%s", diff)
}

// TestGitClientInvalidRepo tests error handling for invalid repository
func TestGitClientInvalidRepo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mcp-notgit-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	client := git.NewClient(tmpDir)
	_, err = client.GetStagedDiff()
	if err == nil {
		t.Error("GetStagedDiff() error = nil, want error for non-git directory")
	}
}

// TestGitClientEmptyStagedArea tests behavior with no staged changes
func TestGitClientEmptyStagedArea(t *testing.T) {
	repoPath, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create initial commit
	createAndStageFile(t, repoPath, "README.md", "# Test\n")
	commitChanges(t, repoPath, "Initial commit")

	// Get staged diff (should be empty)
	client := git.NewClient(repoPath)
	diff, err := client.GetStagedDiff()
	if err != nil {
		t.Fatalf("GetStagedDiff() error = %v", err)
	}

	if diff != "" {
		t.Errorf("GetStagedDiff() = %v, want empty string for no staged changes", diff)
	}
}
