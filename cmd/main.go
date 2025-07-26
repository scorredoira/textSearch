package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	search_engine "textSearch"
)

func main() {
	// Define command line flags
	contextLines := flag.Int("context", 10, "Number of context lines to show around matches")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: search [options] <query>")
		fmt.Println("Options:")
		fmt.Println("  -context int   Number of context lines (default 10)")
		os.Exit(1)
	}

	query := args[0]

	// Get current working directory
	searchPath, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// Check if kbase directory exists
	if _, err := os.Stat(searchPath); os.IsNotExist(err) {
		fmt.Printf("Error: kbase directory not found at %s\n", searchPath)
		os.Exit(1)
	}

	// Create filesystem for kbase directory
	kbaseFS := os.DirFS(searchPath)

	// Create search engine
	engine := search_engine.NewSearchEngine(kbaseFS)

	// Check if query contains pipe-separated terms
	var allResults []search_engine.FileMatch
	fileScores := make(map[string]search_engine.FileMatch)

	if strings.Contains(query, "|") {
		// Split by pipe and search for each term
		terms := strings.Split(query, "|")
		for _, term := range terms {
			term = strings.TrimSpace(term)
			if term == "" {
				continue
			}

			// Search for this term
			results, err := engine.FindRelevantFiles(term, 10)
			if err != nil {
				fmt.Printf("Search error for term '%s': %v\n", term, err)
				continue
			}

			// Merge results, keeping highest score for each file
			for _, result := range results {
				if existing, ok := fileScores[result.Path]; !ok || result.Score > existing.Score {
					// Update reason to include which term matched
					if existing.Reason != "" && ok {
						result.Reason = fmt.Sprintf("matched terms: %s, %s", term, existing.Reason)
					} else {
						result.Reason = fmt.Sprintf("matched term: %s (%s)", term, result.Reason)
					}
					fileScores[result.Path] = result
				}
			}
		}

		// Convert map to slice and sort by score
		for _, match := range fileScores {
			allResults = append(allResults, match)
		}

		// Sort by score descending
		sort.Slice(allResults, func(i, j int) bool {
			return allResults[i].Score > allResults[j].Score
		})

		// Limit to top 10
		if len(allResults) > 10 {
			allResults = allResults[:10]
		}
	} else {
		// Single term search
		var err error
		allResults, err = engine.FindRelevantFiles(query, 10)
		if err != nil {
			fmt.Printf("Search error: %v\n", err)
			os.Exit(1)
		}
	}

	results := allResults

	// Display results
	if len(results) == 0 {
		fmt.Printf("No results found for '%s'\n", query)
		return
	}

	fmt.Printf("Found %d results for '%s':\n\n", len(results), query)

	for _, result := range results {
		// Print separator
		fmt.Println(strings.Repeat("═", 80))

		// Print file info
		fmt.Printf("📁 Path: %s\n", result.Path)
		fmt.Printf("📊 Score: %.2f\n", result.Score)
		if result.Reason != "" {
			fmt.Printf("💡 Reason: %s\n", result.Reason)
		}

		// Extract and show relevant content
		content, err := engine.ExtractRelevantContent(result.Path, query, *contextLines)
		if err == nil && len(content) > 0 {
			fmt.Println("\n📝 Relevant content:")
			fmt.Println(strings.Repeat("─", 80))
			fmt.Println(content)
			fmt.Println(strings.Repeat("─", 80))
		}
		fmt.Println()
	}

	// Final separator
	fmt.Println(strings.Repeat("═", 80))
}
