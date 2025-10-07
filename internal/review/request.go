package review

// Request represents a code review request
type Request struct {
	SourceType     string   // "arbitrary", "staged", "unstaged", "commit"
	Code           string   // Raw code text (for arbitrary) or diff content
	Provider       string   // "anthropic", "openai", "google"
	Language       string   // Programming language hint (optional)
	ReviewDepth    string   // "quick" or "thorough"
	FocusAreas     []string // Filter to specific categories (empty = all)
	RepositoryPath string   // Path to git repository (for git-based reviews)
	CommitSHA      string   // Git commit SHA (for commit reviews)
}

// Validate checks if the request is valid
func (r *Request) Validate() error {
	if r.SourceType == "" {
		return ErrInvalidSourceType
	}

	if r.SourceType == "arbitrary" && r.Code == "" {
		return ErrEmptyCode
	}

	if r.SourceType != "arbitrary" && r.RepositoryPath == "" {
		return ErrMissingRepository
	}

	if r.SourceType == "commit" && r.CommitSHA == "" {
		return ErrMissingCommitSHA
	}

	if r.Provider == "" {
		return ErrMissingProvider
	}

	if r.ReviewDepth != "" && r.ReviewDepth != "quick" && r.ReviewDepth != "thorough" {
		return ErrInvalidReviewDepth
	}

	return nil
}
