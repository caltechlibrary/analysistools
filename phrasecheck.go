package analysistools

import (
	"bufio"
	"fmt"
	"io"
//	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// PatternType represents the type of pattern.
type PatternType string

const (
	Keyword   PatternType = "keyword"
	Proximity PatternType = "proximity"
)

// Matched patterns
type Matched struct {
	Text string
	Pattern string
	Line int
}

func (m *Matched) String() string {
	return fmt.Sprintf("Line %d: Match for %q in : %q", m.Line, m.Pattern, m.Text)
}

func MatchedStrings(matches []*Matched) string {
	result := []string{}
	for _, m := range matches {
		result = append(result, m.String())
	}
	return strings.TrimSpace(strings.Join(result, "\n"))
}

// Pattern represents a parsed pattern.
type Pattern struct {
	Type         PatternType
	Keyword      string
	Keyword1     string
	Keyword2     string
	MaxDistance  int
	OriginalText string
}

// ParsePattern parses a pattern into its components.
func ParsePattern(pattern string) (*Pattern, error) {
	pattern = strings.TrimSpace(pattern)
	parts := strings.Fields(pattern)
	if len(parts) == 0 {
		return nil, fmt.Errorf("missing pattern")
	}
	if len(parts) == 2 || len(parts) > 3 {
		return nil, fmt.Errorf("malformed proximity pattern: %q", pattern)
	}
	if len(parts) == 1 {
		return &Pattern{
			Type: Keyword,
			Keyword: strings.TrimSuffix(parts[0], "*"),
			OriginalText: pattern,
		}, nil
	}
	// Proximity pattern: e.g., "attorn* w/5 client*"
	p := &Pattern{
		Type: Proximity,
		Keyword1: strings.TrimSuffix(parts[0], "*"),
		Keyword2: strings.TrimSuffix(parts[2], "*"),
		OriginalText: pattern,
	}
	if strings.HasPrefix(parts[1], "w/") {
		maxDistance, err := strconv.Atoi(parts[1][2:])
		if err != nil {
			return nil, fmt.Errorf("invalid max distance in pattern: %q", pattern)
		}
		p.MaxDistance = maxDistance
	}
	return p, nil
}

// LoadPatterns loads patterns from a file, one per line.
func LoadPatterns(patternFile string) ([]*Pattern, error) {
	file, err := os.Open(patternFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var patterns []*Pattern
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		pattern, err := ParsePattern(line)
		if err != nil {
			return nil, err
		}
		patterns = append(patterns, pattern)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return patterns, nil
}

// TokenizeLine tokenizes a line into words.
func TokenizeLine(line string) []string {
	re := regexp.MustCompile(`\w+`)
	words := re.FindAllString(line, -1)
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}
	return words
}

// TokenizeFile tokenizes the input file into words, line by line.
func TokenizeFile(inputFile string) ([][]string, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines [][]string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := TokenizeLine(line)
		lines = append(lines, tokens)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

// CheckProximity checks if keyword2 appears within maxDistance words after keyword1.
func CheckProximity(tokens []string, keyword1, keyword2 string, maxDistance int) bool {
	for i, token := range tokens {
		if strings.HasPrefix(token, keyword1) {
			end := i + maxDistance + 1
			if end > len(tokens) {
				end = len(tokens)
			}
			for j := i + 1; j < end; j++ {
				if strings.HasPrefix(tokens[j], keyword2) {
					return true
				}
			}
		}
	}
	return false
}


// PhraseCheckReader evaluates the input from an io.Reader for all patterns.
func PhraseCheckReader(reader io.Reader, patterns []*Pattern, matchOne bool) ([]*Matched, error) {
	result := []*Matched{}
	scanner := bufio.NewScanner(reader)
	for lineNumber := 0; scanner.Scan(); lineNumber++ {
		line := scanner.Text()
		tokens := TokenizeLine(line)
		for _, pattern := range patterns {
			switch pattern.Type {
			case Keyword:
				for _, token := range tokens {
					if strings.HasPrefix(token, pattern.Keyword) {
						result = append(result, &Matched{
							Text: line,
							Pattern: pattern.OriginalText,
							Line: lineNumber,
						})
					}
				}
			case Proximity:
				if CheckProximity(tokens, pattern.Keyword1, pattern.Keyword2, pattern.MaxDistance) {
					result = append(result, &Matched{
						Text: line,
						Pattern: pattern.OriginalText,
						Line: lineNumber,
					})
				}
			}
			if matchOne && len(result) > 0 {
				return result, scanner.Err()
			}
		}
	}
	return result, scanner.Err()
}

type PhraseCheckApp struct {
	appName string
}

func (app *PhraseCheckApp) Check(params []string) error {
	return fmt.Errorf("Check not implemented")
}

func (app *PhraseCheckApp) PruneMatches(params []string) error {
	return fmt.Errorf("PruneMatches not implemented")
}

func (app *PhraseCheckApp) FileTypes(params []string) error {
	return fmt.Errorf("FileTypes not implemented")
}

func (app *PhraseCheckApp) Run(appName string, action string, params []string) error {
	app.appName = appName
	switch action {
	case "filetypes":
		return app.FileTypes(params)
	case "check":
		return app.Check(params)
	case "prune":
		return app.PruneMatches(params)
	default:
		return fmt.Errorf("%q action not supported", action)
	}
	return nil
}