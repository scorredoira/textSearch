package search_engine

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

// FileFinder handles finding relevant files based on queries
type FileFinder struct {
	fs fs.FS
}

// NewFileFinder creates a new FileFinder instance
func NewFileFinder(filesystem fs.FS) *FileFinder {
	return &FileFinder{fs: filesystem}
}

// FindRelevantFiles finds files most relevant to the query
func (ff *FileFinder) FindRelevantFiles(query string, maxFiles int) ([]FileMatch, error) {
	var allMatches []FileMatch
	
	// Normalize query for better matching
	queryTerms := normalizeQuery(query)
	
	// Walk through all files in the filesystem
	err := fs.WalkDir(ff.fs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors, don't fail entire search
		}
		
		if d.IsDir() {
			return nil
		}
		
		// Only search documentation files
		if !isDocumentationFile(path) {
			return nil
		}
		
		// Calculate relevance score for this file
		score, reason := ff.calculateFileScore(path, queryTerms)
		if score > 0 {
			allMatches = append(allMatches, FileMatch{
				Path:     path,
				Score:    score,
				Reason:   reason,
				FileName: filepath.Base(path),
			})
		}
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	// Sort by score (highest first)
	sort.Slice(allMatches, func(i, j int) bool {
		return allMatches[i].Score > allMatches[j].Score
	})
	
	// Limit results
	if maxFiles > 0 && len(allMatches) > maxFiles {
		allMatches = allMatches[:maxFiles]
	} else if maxFiles <= 0 {
		// Return empty results for zero or negative maxFiles
		allMatches = []FileMatch{}
	}
	
	return allMatches, nil
}

// calculateFileScore calculates how relevant a file is to the query
func (ff *FileFinder) calculateFileScore(filePath string, queryTerms []string) (float64, string) {
	if len(queryTerms) == 0 {
		return 0, ""
	}
	
	score := 0.0
	reasons := []string{}
	
	fileName := strings.ToLower(filepath.Base(filePath))
	fileNameNoExt := strings.ToLower(strings.TrimSuffix(fileName, filepath.Ext(fileName)))
	
	// Check directory path matches first (highest priority)
	dirPath := strings.ToLower(filepath.Dir(filePath))
	for _, term := range queryTerms {
		termLower := strings.ToLower(term)
		
		// Check if term appears as a directory component
		dirComponents := strings.Split(dirPath, string(filepath.Separator))
		for _, component := range dirComponents {
			if component == termLower {
				score += 2.5 // Very high score for exact directory match
				reasons = append(reasons, "directory exact match '"+term+"'")
				break
			} else if strings.Contains(component, termLower) {
				score += 1.8 // High score for directory contains term
				reasons = append(reasons, "directory contains '"+term+"'")
				break
			}
		}
		
		// Also check full directory path
		if strings.Contains(dirPath, termLower) && !strings.Contains(strings.Join(dirComponents, ""), termLower) {
			score += 1.2 // Medium-high score for path contains term (but not already counted above)
			reasons = append(reasons, "path contains '"+term+"'")
		}
	}

	// Check filename matches
	for _, term := range queryTerms {
		termLower := strings.ToLower(term)
		
		// Exact filename match (highest score)
		if fileNameNoExt == termLower {
			score += 2.0
			reasons = append(reasons, "exact filename match")
			continue
		}
		
		// Filename contains term (high score) 
		if strings.Contains(fileNameNoExt, termLower) {
			score += 1.5
			reasons = append(reasons, "filename contains '"+term+"'")
			continue
		}
		
		// Filename word boundary match (medium-high score)
		if isWordInString(fileNameNoExt, termLower) {
			score += 1.0
			reasons = append(reasons, "filename word match '"+term+"'")
		} else {
			// Check for partial word matches in filename
			words := strings.FieldsFunc(fileNameNoExt, func(r rune) bool {
				return !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9'))
			})
			
			for _, word := range words {
				if len(word) < 3 || len(termLower) < 3 {
					continue
				}
				// Check if term contains word or word contains term
				if (len(termLower) > 3 && strings.Contains(termLower, word)) ||
				   (len(word) > 3 && strings.Contains(word, termLower)) {
					score += 0.5 // Moderate score for partial filename matches
					reasons = append(reasons, "partial filename match '"+term+"'")
					break
				}
			}
		}
		
	}
	
	// Check file content for additional scoring
	contentScore, contentReason := ff.scoreFileContent(filePath, queryTerms)
	score += contentScore * 0.3 // Content match weighted lower than filename match
	if contentReason != "" {
		reasons = append(reasons, contentReason)
	}
	
	// Normalize score to 0-1 range, but allow higher scores for better ranking
	if score > 0 {
		// Keep relative differences but normalize to reasonable range
		// Don't cap too aggressively - let directory matches rank higher
		score = score / 4.0 // Scale down but preserve differences
		if score > 1.0 {
			score = 1.0
		}
	}
	
	reason := strings.Join(reasons, ", ")
	return score, reason
}

// scoreFileContent scores based on file content
func (ff *FileFinder) scoreFileContent(filePath string, queryTerms []string) (float64, string) {
	content, err := fs.ReadFile(ff.fs, filePath)
	if err != nil {
		return 0, ""
	}
	
	contentStr := strings.ToLower(string(content))
	score := 0.0
	matchedTerms := 0
	
	for _, term := range queryTerms {
		termLower := strings.ToLower(term)
		if strings.Contains(contentStr, termLower) {
			matchedTerms++
			// Higher score for terms that appear multiple times
			count := strings.Count(contentStr, termLower)
			score += float64(count) * 0.1
		} else {
			// Check for partial matches (term is substring of words in content)
			// or content words are substring of term
			words := strings.Fields(contentStr)
			for _, word := range words {
				cleanWord := strings.Trim(word, ".,!?:;()[]{}\"'")
				if len(cleanWord) < 3 {
					continue
				}
				
				// Check if term contains word or word contains term
				if (len(termLower) > 3 && strings.Contains(termLower, cleanWord)) ||
				   (len(cleanWord) > 3 && strings.Contains(cleanWord, termLower)) {
					score += 0.05 // Lower score for partial matches
					break
				}
			}
		}
	}
	
	if score == 0 {
		return 0, ""
	}
	
	// Bonus for matching all terms
	if matchedTerms == len(queryTerms) {
		score += 0.3
	}
	
	reason := "content matches"
	
	return score, reason
}

// Helper functions

func normalizeQuery(query string) []string {
	// Handle pipe-separated terms (OR logic)
	var allTerms []string
	if strings.Contains(query, "|") {
		pipeTerms := strings.Split(query, "|")
		for _, pipeTerm := range pipeTerms {
			// Split each pipe term by spaces and add to allTerms
			spaceTerms := strings.Fields(strings.TrimSpace(pipeTerm))
			allTerms = append(allTerms, spaceTerms...)
		}
	} else {
		// Split query into terms and clean them
		allTerms = strings.Fields(query)
	}
	
	var normalized []string
	
	for _, term := range allTerms {
		term = strings.ToLower(strings.TrimSpace(term))
		
		// Remove common stop words that don't add value
		if isStopWord(term) {
			continue
		}
		
		// Clean punctuation but preserve pipes for regex queries
		cleaned := strings.Trim(term, ".,!?:;()[]{}\"'")
		if len(cleaned) > 1 { // Ignore single characters
			normalized = append(normalized, cleaned)
			
			// Add common variations/stems for better matching
			normalized = append(normalized, generateTermVariations(cleaned)...)
		}
	}
	
	// Remove duplicates
	seen := make(map[string]bool)
	var unique []string
	for _, term := range normalized {
		if !seen[term] {
			seen[term] = true
			unique = append(unique, term)
		}
	}
	
	return unique
}

// generateTermVariations creates common variations of a term for better matching
func generateTermVariations(term string) []string {
	var variations []string
	
	// Simple stemming - remove common suffixes
	if len(term) > 3 {
		// Remove 's' plural (but not from words ending in 'ss')
		if strings.HasSuffix(term, "s") && !strings.HasSuffix(term, "ss") && len(term) > 3 {
			variations = append(variations, term[:len(term)-1])
		}
	}
	
	if len(term) > 5 {
		// Remove 'ing' 
		if strings.HasSuffix(term, "ing") {
			variations = append(variations, term[:len(term)-3])
		}
		// Remove 'ed'
		if strings.HasSuffix(term, "ed") {
			variations = append(variations, term[:len(term)-2])
		}
	}
	
	return variations
}

func isStopWord(word string) bool {
	stopWords := map[string]bool{
		"a": true, "an": true, "and": true, "are": true, "as": true, "at": true,
		"be": true, "by": true, "for": true, "from": true, "has": true, "he": true,
		"in": true, "is": true, "it": true, "its": true, "of": true, "on": true,
		"that": true, "the": true, "to": true, "was": true, "will": true, "with": true,
		"como": true, "con": true, "del": true, "el": true, "en": true, "es": true,
		"la": true, "para": true, "por": true, "que": true, "un": true, "una": true,
	}
	return stopWords[word]
}

func isDocumentationFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	docExtensions := map[string]bool{
		".md":   true,
		".html": true,
		".txt":  true,
		".json": true,
		".yaml": true,
		".yml":  true,
		".rst":  true,
		".adoc": true,
	}
	
	// Skip hidden files
	base := filepath.Base(path)
	if strings.HasPrefix(base, ".") {
		return false
	}
	
	return docExtensions[ext]
}

func isWordInString(text, word string) bool {
	// Simple word boundary check
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9'))
	})
	
	for _, w := range words {
		if strings.ToLower(w) == word {
			return true
		}
	}
	return false
}