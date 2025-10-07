package review

import "errors"

// Request validation errors
var (
	ErrInvalidSourceType  = errors.New("invalid source type")
	ErrEmptyCode          = errors.New("code cannot be empty for arbitrary reviews")
	ErrMissingRepository  = errors.New("repository path is required for git-based reviews")
	ErrMissingCommitSHA   = errors.New("commit SHA is required for commit reviews")
	ErrMissingProvider    = errors.New("provider must be specified")
	ErrInvalidReviewDepth = errors.New("review depth must be 'quick' or 'thorough'")
)

// Provider errors
var (
	ErrProviderNotAvailable = errors.New("provider is not available or not configured")
	ErrProviderTimeout      = errors.New("provider request timed out")
	ErrProviderAPIError     = errors.New("provider API error")
)
