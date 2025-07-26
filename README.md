# TextSearch - Intelligent Documentation Search Engine

A fast, LLM-optimized search engine designed for finding and extracting relevant content from documentation files. Built in Go with no external dependencies.

## Features

- **Smart relevance scoring** - Combines filename and content matching with intelligent weighting
- **Multi-term OR queries** - Search for multiple terms using pipe separator (`term1|term2|term3`)
- **Configurable context extraction** - Extract relevant content with customizable surrounding lines
- **LLM-optimized output** - Clean, structured results perfect for AI processing
- **No indexing required** - Direct file system scanning for simplicity
- **Multiple file format support** - Markdown, text, and other documentation formats

## Quick Start

### Build
```bash
go build -o search cmd/main.go
```

### Basic Usage
```bash
# Single term search
./search "authentication"

# Multi-term OR search  
./search "api|billing|bookings"

# Custom context lines
./search -context 5 "payment processing"
```

## Command Line Options

```bash
Usage: search [options] <query>

Options:
  -context int   Number of context lines to show around matches (default 10)
```

## Query Formats

### Single Term
```bash
./search "webhooks"
```
Searches for files containing "webhooks".

### Multi-term OR Query
```bash
./search "billing|payment|subscription"
```
Finds files containing ANY of the specified terms, returning the highest-scored results.

## Scoring Algorithm

The search engine uses a sophisticated scoring system that prioritizes different types of matches:

### Filename Scoring (High Priority)
- **Exact filename match**: 2.0 points
- **Filename contains term**: 1.5 points  
- **Word boundary match**: 1.0 points
- **Partial match**: 0.5 points

### Content Scoring (Lower Priority)
- **Multiple occurrences boost score**: Each occurrence adds 0.1 points
- **All query terms found**: Bonus 0.3 points
- **Partial word matches**: 0.05 points each
- **Content weight**: Final content score × 0.3

### Final Score
```
Total Score = Filename Score + (Content Score × 0.3)
```
Capped at 1.0 maximum, then results sorted by score descending.

## Output Format

Results are displayed with clear visual separation:

```
════════════════════════════════════════════════════════════════════════════════
📁 Path: api_authentication.md
📊 Score: 1.00
💡 Reason: filename contains 'authentication', content matches

📝 Relevant content:
────────────────────────────────────────────────────────────────────────────────
[Extracted content with configurable context lines]
────────────────────────────────────────────────────────────────────────────────
```

## Architecture

The search engine consists of three main components:

### SearchEngine Interface
Core interface providing search and content extraction capabilities.

### FileFinder
- Walks file system to find matching documents
- Calculates relevance scores based on filename and content analysis
- Supports both single and multi-term queries

### ContentExtractor  
- Identifies relevant sections within matched files
- Expands matches with configurable context lines
- Prioritizes important content (headers, code blocks, URLs)

## Use Cases

### For LLMs
- **No truncation**: Full content extraction for comprehensive context
- **Structured output**: Easy to parse results with scores and metadata
- **Relevance ranking**: Best matches first for efficient token usage

### For Documentation Teams
- **Quick searches**: Fast file system scanning without indexing overhead
- **Flexible queries**: Support for both specific and broad searches
- **Context aware**: See relevant content with surrounding context

### For API Documentation
- **Endpoint discovery**: Find API endpoints across multiple files
- **Cross-references**: Locate related concepts across documentation
- **Example extraction**: Pull relevant code examples and usage patterns

## Configuration

The search engine searches in the `testData` directory by default. To use with your documentation:

1. Place your documentation files in a `testData` directory
2. Or modify the `searchPath` in `cmd/main.go` to point to your docs
3. Build and run searches

## Performance Characteristics

- **Memory efficient**: Streams file content without loading all into memory
- **CPU optimized**: Simple string matching without complex NLP
- **Scalable**: Performance scales linearly with document count
- **Deterministic**: Consistent results across runs

## File Support

Currently optimized for documentation files:
- Markdown (`.md`)
- Text files (`.txt`) 
- API documentation
- Code documentation

The `isDocumentationFile()` function can be extended for additional formats.

## Examples

### Finding API Endpoints
```bash
./search "POST|GET|PUT|DELETE"
```

### Authentication Documentation  
```bash
./search "auth|login|token|key"
```

### Payment Processing
```bash
./search "payment|billing|charge|subscription"
```

### Error Handling
```bash
./search "error|exception|failure"
```

## Build Requirements

- Go 1.24+ 
- No external dependencies

## License

[Add your license here]