package search_engine

import (
	"testing"
	"testing/fstest"
)

func TestFileFinder_FindRelevantFiles(t *testing.T) {
	// Create a test filesystem
	testFS := fstest.MapFS{
		"testData/foo/api/bar.md":           &fstest.MapFile{Data: []byte("API documentation for bar")},
		"testData/foo/bad.md":               &fstest.MapFile{Data: []byte("Some bad documentation")},
		"testData/api_authentication.md":   &fstest.MapFile{Data: []byte("Authentication API guide")},
		"testData/auth/user_guide.md":      &fstest.MapFile{Data: []byte("User authentication guide")},
		"docs/api/endpoints.md":            &fstest.MapFile{Data: []byte("API endpoints documentation")},
		"docs/setup.md":                    &fstest.MapFile{Data: []byte("Setup instructions")},
		"config/api_config.json":           &fstest.MapFile{Data: []byte(`{"api": "config"}`)},
		"other/random.txt":                 &fstest.MapFile{Data: []byte("Random text file")},
	}

	finder := NewFileFinder(testFS)

	tests := []struct {
		name          string
		query         string
		maxFiles      int
		expectContains []string // Files that should be in results
		expectOrder   []string  // Expected order of top results
	}{
		{
			name:          "API query should prioritize directory structure",
			query:         "api",
			maxFiles:      5,
			expectContains: []string{"testData/foo/api/bar.md", "docs/api/endpoints.md", "testData/api_authentication.md"},
			// Both files with "api" in directory should rank higher than filename-only matches
			// (exact order between directory matches may vary, but both should rank higher than filename matches)
		},
		{
			name:          "Authentication query",
			query:         "authentication",
			maxFiles:      3,
			expectContains: []string{"testData/api_authentication.md", "testData/auth/user_guide.md"},
		},
		{
			name:     "No matches",
			query:    "nonexistent",
			maxFiles: 5,
		},
		{
			name:          "Filename exact match",
			query:         "setup",
			maxFiles:      3,
			expectContains: []string{"docs/setup.md"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := finder.FindRelevantFiles(tt.query, tt.maxFiles)
			if err != nil {
				t.Fatalf("FindRelevantFiles() error = %v", err)
			}

			// Debug: Print all results for API query test
			if tt.name == "API query should prioritize directory structure" {
				t.Logf("Results for query '%s':", tt.query)
				for i, result := range results {
					t.Logf("  %d: %s (score: %.3f, reason: %s)", i, result.Path, result.Score, result.Reason)
				}
			}

			// Check that expected files are included
			resultPaths := make(map[string]bool)
			for _, result := range results {
				resultPaths[result.Path] = true
			}

			for _, expected := range tt.expectContains {
				if !resultPaths[expected] {
					t.Errorf("Expected file %s not found in results", expected)
				}
			}

			// Check order if specified
			if len(tt.expectOrder) > 0 && len(results) >= len(tt.expectOrder) {
				for i, expectedPath := range tt.expectOrder {
					if i < len(results) && results[i].Path != expectedPath {
						t.Errorf("Expected %s at position %d, got %s", expectedPath, i, results[i].Path)
					}
				}
			}
			
			// For API query test, verify directory matches rank higher than filename matches
			if tt.name == "API query should prioritize directory structure" && len(results) >= 4 {
				// First two should be directory matches with higher scores
				dirMatches := 0
				for i := 0; i < 2 && i < len(results); i++ {
					if results[i].Score > 0.6 { // Directory matches should score > 0.6
						dirMatches++
					}
				}
				if dirMatches != 2 {
					t.Errorf("Expected first 2 results to be directory matches with high scores, got %d", dirMatches)
				}
				
				// Next results should be filename matches with lower scores
				filenameMatches := 0
				for i := 2; i < 4 && i < len(results); i++ {
					if results[i].Score < 0.5 { // Filename matches should score < 0.5
						filenameMatches++
					}
				}
				if filenameMatches != 2 {
					t.Errorf("Expected next 2 results to be filename matches with lower scores, got %d", filenameMatches)
				}
			}

			// Verify maxFiles limit
			if len(results) > tt.maxFiles {
				t.Errorf("Expected max %d files, got %d", tt.maxFiles, len(results))
			}
		})
	}
}

func TestFileFinder_calculateFileScore(t *testing.T) {
	testFS := fstest.MapFS{
		"testData/foo/api/bar.md":     &fstest.MapFile{Data: []byte("API documentation for bar service")},
		"testData/foo/bad.md":         &fstest.MapFile{Data: []byte("Some unrelated content")},
		"docs/api_reference.md":       &fstest.MapFile{Data: []byte("Complete API reference guide")},
		"auth/authentication.md":     &fstest.MapFile{Data: []byte("Authentication methods")},
	}

	finder := NewFileFinder(testFS)

	tests := []struct {
		name       string
		filePath   string
		query      string
		expectScore bool // Whether score should be > 0
		expectHigher string // This file should score higher than expectHigher file
	}{
		{
			name:        "Directory path match should score high",
			filePath:    "testData/foo/api/bar.md",
			query:       "api",
			expectScore: true,
		},
		{
			name:        "No directory match should score lower",
			filePath:    "testData/foo/bad.md", 
			query:       "api",
			expectScore: false,
		},
		{
			name:        "Filename exact match",
			filePath:    "auth/authentication.md",
			query:       "authentication",
			expectScore: true,
		},
		{
			name:        "Filename contains term",
			filePath:    "docs/api_reference.md",
			query:       "api",
			expectScore: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryTerms := normalizeQuery(tt.query)
			score, reason := finder.calculateFileScore(tt.filePath, queryTerms)

			if tt.expectScore && score <= 0 {
				t.Errorf("Expected score > 0 for %s with query '%s', got %f", tt.filePath, tt.query, score)
			}
			if !tt.expectScore && score > 0 {
				t.Errorf("Expected score = 0 for %s with query '%s', got %f", tt.filePath, tt.query, score)
			}

			if score > 0 && reason == "" {
				t.Errorf("Expected non-empty reason when score > 0")
			}

			t.Logf("File: %s, Query: %s, Score: %f, Reason: %s", tt.filePath, tt.query, score, reason)
		})
	}

	// Test directory vs non-directory scoring
	t.Run("Directory path should score higher than filename only", func(t *testing.T) {
		queryTerms := normalizeQuery("api")
		
		// File with api in directory path
		dirScore, _ := finder.calculateFileScore("testData/foo/api/bar.md", queryTerms)
		
		// File with api only in content, not path
		nonDirScore, _ := finder.calculateFileScore("testData/foo/bad.md", queryTerms)
		
		if dirScore <= nonDirScore {
			t.Errorf("Directory path match should score higher: dirScore=%f, nonDirScore=%f", dirScore, nonDirScore)
		}
	})
}

func TestNormalizeQuery(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected []string
	}{
		{
			name:     "Simple query",
			query:    "api authentication",
			expected: []string{"api", "authentication"},
		},
		{
			name:     "Query with pipe separator",
			query:    "api|auth",
			expected: []string{"api", "auth"},
		},
		{
			name:     "Query with punctuation",
			query:    "api, authentication!",
			expected: []string{"api", "authentication"},
		},
		{
			name:     "Query with stop words",
			query:    "the api and authentication",
			expected: []string{"api", "authentication"},
		},
		{
			name:     "Empty query",
			query:    "",
			expected: []string{},
		},
		{
			name:     "Query with plurals",
			query:    "apis endpoints",
			expected: []string{"apis", "api", "endpoints", "endpoint"}, // Should include stemmed versions
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeQuery(tt.query)
			
			// Check that all expected terms are present
			resultMap := make(map[string]bool)
			for _, term := range result {
				resultMap[term] = true
			}
			
			for _, expected := range tt.expected {
				if !resultMap[expected] {
					t.Errorf("Expected term '%s' not found in result %v", expected, result)
				}
			}
		})
	}
}

func TestIsDocumentationFile(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"Markdown file", "docs/readme.md", true},
		{"HTML file", "docs/index.html", true},
		{"Text file", "notes.txt", true},
		{"JSON file", "config.json", true},
		{"YAML file", "config.yaml", true},
		{"YML file", "config.yml", true},
		{"Go file", "main.go", false},
		{"Hidden file", ".gitignore", false},
		{"Binary file", "image.png", false},
		{"No extension", "README", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isDocumentationFile(tt.path)
			if result != tt.expected {
				t.Errorf("isDocumentationFile(%s) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestIsWordInString(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		word     string
		expected bool
	}{
		{"Exact word match", "hello world", "hello", true},
		{"Word with punctuation", "hello, world!", "hello", true},
		{"Partial match should fail", "hello world", "hell", false},
		{"Case insensitive", "Hello World", "hello", true},
		{"Word boundary", "api-documentation", "api", true},
		{"Not a word boundary", "application", "api", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isWordInString(tt.text, tt.word)
			if result != tt.expected {
				t.Errorf("isWordInString(%s, %s) = %v, expected %v", tt.text, tt.word, result, tt.expected)
			}
		})
	}
}

func TestGenerateTermVariations(t *testing.T) {
	tests := []struct {
		name     string
		term     string
		expected []string
	}{
		{
			name:     "Plural term",
			term:     "apis",
			expected: []string{"api"},
		},
		{
			name:     "Term with ing",
			term:     "testing",
			expected: []string{"test"},
		},
		{
			name:     "Term with ed",
			term:     "tested",
			expected: []string{"test"},
		},
		{
			name:     "Short term",
			term:     "api",
			expected: []string{},
		},
		{
			name:     "Term ending in ss",
			term:     "process",
			expected: []string{}, // Should not remove 's' from 'ss'
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateTermVariations(tt.term)
			
			if len(result) != len(tt.expected) {
				t.Errorf("generateTermVariations(%s) returned %v, expected %v", tt.term, result, tt.expected)
				return
			}
			
			for i, expected := range tt.expected {
				if i < len(result) && result[i] != expected {
					t.Errorf("generateTermVariations(%s)[%d] = %s, expected %s", tt.term, i, result[i], expected)
				}
			}
		})
	}
}

func TestIsStopWord(t *testing.T) {
	tests := []struct {
		name     string
		word     string
		expected bool
	}{
		{"English stop word", "the", true},
		{"Spanish stop word", "el", true},
		{"Not a stop word", "api", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isStopWord(tt.word)
			if result != tt.expected {
				t.Errorf("isStopWord(%s) = %v, expected %v", tt.word, result, tt.expected)
			}
		})
	}
}