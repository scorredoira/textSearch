package search_engine

import (
	"io/fs"
	"regexp"
	"sort"
	"strings"
)

// ContentExtractor handles extracting relevant content from files
type ContentExtractor struct {
	fs fs.FS
}

// NewContentExtractor creates a new ContentExtractor instance
func NewContentExtractor(filesystem fs.FS) *ContentExtractor {
	return &ContentExtractor{fs: filesystem}
}

// ExtractRelevantContent extracts content relevant to the query from a file
func (ce *ContentExtractor) ExtractRelevantContent(filePath, query string, contextLines int) (string, error) {
	content, err := fs.ReadFile(ce.fs, filePath)
	if err != nil {
		return "", err
	}

	contentStr := string(content)
	queryTerms := normalizeQuery(query)

	if len(queryTerms) == 0 {
		// If no specific terms, return a reasonable sample
		return ce.getContentSample(contentStr, 1000), nil
	}

	// Find relevant sections
	relevantSections := ce.findRelevantSections(contentStr, queryTerms, contextLines)

	if len(relevantSections) == 0 {
		// No specific matches, return beginning of file
		return ce.getContentSample(contentStr, 1000), nil
	}

	// Combine and format the relevant sections
	return ce.formatRelevantSections(contentStr, relevantSections), nil
}

// findRelevantSections finds sections of the content that are relevant to the query
func (ce *ContentExtractor) findRelevantSections(content string, queryTerms []string, contextLines int) []ContentSection {
	lines := strings.Split(content, "\n")
	var sections []ContentSection

	// Score each line based on relevance
	for i, line := range lines {
		score := ce.scoreLineRelevance(line, queryTerms)
		if score > 0 {
			sections = append(sections, ContentSection{
				LineNumber: i,
				Score:      score,
				Content:    line,
			})
		}
	}

	// Sort by score (highest first)
	sort.Slice(sections, func(i, j int) bool {
		return sections[i].Score > sections[j].Score
	})

	// Expand highly relevant sections with context
	expandedSections := ce.expandSectionsWithContext(lines, sections, contextLines)

	return expandedSections
}

// scoreLineRelevance scores how relevant a line is to the query terms
func (ce *ContentExtractor) scoreLineRelevance(line string, queryTerms []string) float64 {
	if len(queryTerms) == 0 {
		return 0
	}

	lineLower := strings.ToLower(line)
	score := 0.0
	matchedTerms := 0

	for _, term := range queryTerms {
		termLower := strings.ToLower(term)

		if strings.Contains(lineLower, termLower) {
			matchedTerms++

			// Higher score for different types of matches
			if ce.isExactWordMatch(lineLower, termLower) {
				score += 1.0 // Exact word match
			} else {
				score += 0.5 // Partial match
			}

			// Bonus for lines that look like important content
			if ce.isImportantLine(line) {
				score += 0.3
			}
		}
	}

	// Bonus for matching multiple terms
	if matchedTerms > 1 {
		score += 0.2 * float64(matchedTerms-1)
	}

	return score
}

// isExactWordMatch checks if a term appears as a complete word
func (ce *ContentExtractor) isExactWordMatch(text, term string) bool {
	pattern := `\b` + regexp.QuoteMeta(term) + `\b`
	matched, _ := regexp.MatchString(pattern, text)
	return matched
}

// isImportantLine identifies lines that are likely to contain important information
func (ce *ContentExtractor) isImportantLine(line string) bool {
	line = strings.TrimSpace(line)

	// Headers
	if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "##") {
		return true
	}

	// Code blocks or endpoints
	if strings.Contains(line, "http://") || strings.Contains(line, "https://") {
		return true
	}

	// API request examples
	if strings.Contains(line, "GET ") || strings.Contains(line, "POST ") ||
		strings.Contains(line, "PUT ") || strings.Contains(line, "DELETE ") {
		return true
	}

	// curl commands
	if strings.Contains(line, "curl") {
		return true
	}

	// JSON/code examples
	if strings.HasPrefix(line, "{") || strings.HasPrefix(line, "[") {
		return true
	}

	// Code blocks
	if strings.HasPrefix(line, "```") {
		return true
	}

	// Parameter definitions
	if strings.Contains(line, ":") && (strings.Contains(line, "string") ||
		strings.Contains(line, "int") || strings.Contains(line, "bool")) {
		return true
	}

	return false
}

// expandSectionsWithContext adds context lines around relevant sections
func (ce *ContentExtractor) expandSectionsWithContext(lines []string, sections []ContentSection, contextLines int) []ContentSection {
	if len(sections) == 0 {
		return sections
	}

	var expandedSections []ContentSection
	usedLines := make(map[int]bool)

	// Process sections by score (highest first)
	for _, section := range sections {
		if len(expandedSections) >= 5 { // Limit to top 5 sections
			break
		}

		start := maxInt(0, section.LineNumber-contextLines)
		end := minInt(len(lines), section.LineNumber+contextLines+1)

		// Check if this section overlaps with already used lines
		hasOverlap := false
		for i := start; i < end; i++ {
			if usedLines[i] {
				hasOverlap = true
				break
			}
		}

		if !hasOverlap {
			// Mark lines as used
			for i := start; i < end; i++ {
				usedLines[i] = true
			}

			// Combine the lines for this section
			var contextContent []string
			for i := start; i < end; i++ {
				contextContent = append(contextContent, lines[i])
			}

			expandedSections = append(expandedSections, ContentSection{
				LineNumber: start,
				Score:      section.Score,
				Content:    strings.Join(contextContent, "\n"),
				EndLine:    end - 1,
			})
		}
	}

	return expandedSections
}

// formatRelevantSections formats the relevant sections into readable text
func (ce *ContentExtractor) formatRelevantSections(originalContent string, sections []ContentSection) string {
	if len(sections) == 0 {
		return ce.getContentSample(originalContent, 1000)
	}

	var result []string

	for i, section := range sections {
		if i >= 10 { // Show up to 10 relevant sections
			break
		}

		// Add section header
		if len(sections) > 1 {
			result = append(result, "--- Relevant Section ---")
		}

		result = append(result, section.Content)

		if i < len(sections)-1 && len(sections) > 1 {
			result = append(result, "") // Empty line between sections
		}
	}

	content := strings.Join(result, "\n")

	// Return full content without truncation for LLM usage
	return content
}

// getContentSample returns a sample of content (beginning)
func (ce *ContentExtractor) getContentSample(content string, maxLength int) string {
	if len(content) <= maxLength {
		return content
	}

	// Try to cut at a line boundary
	sample := content[:maxLength]
	if lastNewline := strings.LastIndex(sample, "\n"); lastNewline > maxLength/2 {
		sample = content[:lastNewline]
	}

	return sample
}

// ContentSection represents a section of content with relevance score
type ContentSection struct {
	LineNumber int
	EndLine    int
	Score      float64
	Content    string
}

// Helper functions

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
