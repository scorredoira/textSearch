package search_engine

import (
	"io/fs"
)

// SearchEngine provides intelligent search capabilities for documentation files
type SearchEngine interface {
	// FindRelevantFiles finds files that are most relevant to the query
	FindRelevantFiles(query string, maxFiles int) ([]FileMatch, error)
	
	// ExtractRelevantContent extracts relevant content from a specific file for the query
	ExtractRelevantContent(filePath, query string, contextLines int) (string, error)
	
	// GetFileContent reads the complete content of a file
	GetFileContent(filePath string) (string, error)
}

// FileMatch represents a file that matches a search query
type FileMatch struct {
	Path     string  `json:"path"`     // Relative path to the file
	Score    float64 `json:"score"`    // Relevance score (0.0 to 1.0)
	Reason   string  `json:"reason"`   // Human-readable explanation of why this file matches
	FileName string  `json:"filename"` // Just the filename for quick reference
}

// ContentMatch represents relevant content within a file
type ContentMatch struct {
	Content    string `json:"content"`     // The relevant content extracted
	LineStart  int    `json:"line_start"`  // Starting line number
	LineEnd    int    `json:"line_end"`    // Ending line number
	Context    string `json:"context"`     // Additional context around the match
	Confidence float64 `json:"confidence"` // How confident we are this content is relevant
}

// SearchEngineImpl implements the SearchEngine interface
type SearchEngineImpl struct {
	fs           fs.FS
	fileFinder   *FileFinder
	extractor    *ContentExtractor
}

// NewSearchEngine creates a new SearchEngine instance
func NewSearchEngine(filesystem fs.FS) SearchEngine {
	return &SearchEngineImpl{
		fs:           filesystem,
		fileFinder:   NewFileFinder(filesystem),
		extractor:    NewContentExtractor(filesystem),
	}
}

// FindRelevantFiles implements SearchEngine.FindRelevantFiles
func (se *SearchEngineImpl) FindRelevantFiles(query string, maxFiles int) ([]FileMatch, error) {
	return se.fileFinder.FindRelevantFiles(query, maxFiles)
}

// ExtractRelevantContent implements SearchEngine.ExtractRelevantContent
func (se *SearchEngineImpl) ExtractRelevantContent(filePath, query string, contextLines int) (string, error) {
	return se.extractor.ExtractRelevantContent(filePath, query, contextLines)
}

// GetFileContent implements SearchEngine.GetFileContent
func (se *SearchEngineImpl) GetFileContent(filePath string) (string, error) {
	content, err := fs.ReadFile(se.fs, filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}