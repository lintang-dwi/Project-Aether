package workspace

import (
	"os"
	"path/filepath"
	"strings"

	"aether/model"
	ignore "github.com/sabhiram/go-gitignore"
)

// Scanner scans project directories and discovers files while honoring ignore rules.
type Scanner struct {
	rootPath string
	matcher  *ignore.GitIgnore
}

// NewScanner creates a new Scanner for the target workspace root.
func NewScanner(rootPath string, customIgnorePatterns []string) (*Scanner, error) {
	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}

	var patterns []string
	patterns = append(patterns, ".git", "node_modules", "vendor", "bin", "dist", ".aether")
	patterns = append(patterns, customIgnorePatterns...)

	// Read .gitignore if present in root
	gitignorePath := filepath.Join(absPath, ".gitignore")
	if _, err := os.Stat(gitignorePath); err == nil {
		matcher, err := ignore.CompileIgnoreFileAndLines(gitignorePath, patterns...)
		if err == nil {
			return &Scanner{rootPath: absPath, matcher: matcher}, nil
		}
	}

	matcher := ignore.CompileIgnoreLines(patterns...)
	return &Scanner{
		rootPath: absPath,
		matcher:  matcher,
	}, nil
}

func (s *Scanner) matches(cleanRelPath string, isDir bool) bool {
	if s.matcher == nil {
		return false
	}
	if s.matcher.MatchesPath(cleanRelPath) {
		return true
	}
	if isDir && s.matcher.MatchesPath(cleanRelPath+"/") {
		return true
	}
	// Check parent directory components
	parts := strings.Split(cleanRelPath, "/")
	for i := 1; i < len(parts); i++ {
		parent := strings.Join(parts[:i], "/")
		if s.matcher.MatchesPath(parent) || s.matcher.MatchesPath(parent+"/") {
			return true
		}
	}
	return false
}

// Scan Walk through the workspace directory and return all non-ignored file metadata.
func (s *Scanner) Scan() ([]model.FileMetadata, error) {
	var files []model.FileMetadata

	err := filepath.Walk(s.rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		relPath, err := filepath.Rel(s.rootPath, path)
		if err != nil || relPath == "." {
			return nil
		}

		// Normalize slash for ignore matching
		cleanRelPath := filepath.ToSlash(relPath)
		if s.matches(cleanRelPath, info.IsDir()) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !info.IsDir() {
			files = append(files, model.FileMetadata{
				Path:      cleanRelPath,
				Size:      info.Size(),
				ModTime:   info.ModTime().UTC(),
				IsIgnored: false,
			})
		}
		return nil
	})

	return files, err
}

// IsIgnored checks whether a given relative path is ignored.
func (s *Scanner) IsIgnored(relPath string) bool {
	clean := filepath.ToSlash(relPath)
	info, err := os.Stat(filepath.Join(s.rootPath, relPath))
	isDir := err == nil && info.IsDir()
	return s.matches(clean, isDir)
}

// GetRootPath returns the absolute root path of the workspace scanner.
func (s *Scanner) GetRootPath() string {
	return s.rootPath
}
