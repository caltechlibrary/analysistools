package analysistools

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	//"regexp"
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
	PatternType PatternType
	LineNo int
	WordNo int
}

func (m *Matched) String() string {
	return fmt.Sprintf("%d,%q,%q", m.LineNo, m.Pattern, m.Text)
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
	if len(parts) > 3 {
		return nil, fmt.Errorf("malformed proximity pattern: %q", pattern)
	}
	// Setup the first keyword in pattern
	p := &Pattern{}
	p.OriginalText = pattern
	p.Type = Keyword
	token := parts[0]
	if strings.HasPrefix(token, "w/") {
		return nil, fmt.Errorf("malformed proximity pattern: %q", pattern)
	}
	p.Keyword1 = token
	if len(parts) == 1 {
		return p, nil
	}
	// We have a proximity pattern so set that up.
	token = parts[1]
	p.Type = Proximity
	if len(parts) == 2 {
		p.Keyword2 = token
		if strings.HasPrefix(token, "w/") {
			return nil, fmt.Errorf("malformed proximity pattern: %q", pattern)
		}
		p.MaxDistance = 1
		return p, nil
	}	
	// Proximity pattern: e.g., "attorn* w/5 client*"
	p.Keyword2 = parts[2]
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

// tokenMatches compares a token which may be prefixed or suffixed by an asterix and
// determines if there is a match with the provided string.
func tokenMatches(s string, expr string) bool {
	// FIXME: we could actually use the Go RegExp engine here. Need to check if it would be faster.
	switch {
	case strings.HasPrefix(expr, "*") && strings.HasSuffix(expr, "*"):
		// Checking for middle match
		return strings.Contains(s, strings.TrimSuffix(strings.TrimPrefix(expr, "*"), "*"))
	case strings.HasPrefix(expr, "*"):
		// Checking for Suffix match
		return strings.HasSuffix(s, strings.TrimPrefix(expr, "*"))
	case strings.HasSuffix(expr, "*"):
		// Checking for Prefix match
		return strings.HasPrefix(s, strings.TrimSuffix(expr, "*"))
	default:
		// Check for exact match
		return expr == s
	}
}

// CheckProximity checks if keyword2 appears within maxDistance words after keyword1.
func CheckProximity(tokens []*Token, keyword1 string, keyword2 string, maxDistance int) (*Token, bool) {
	result := &Token{}
	// Make a copy to iterate through
	for i, token := range tokens {
		if tokenMatches(token.Value, keyword1)  {
			result = token
			// Look ahead to find the next match based on max distance
			end := i + maxDistance + 1
			if end > len(tokens) {
				end = len(tokens)
			}
			for j := i + 1; j < end; j++ {
				if tokenMatches(tokens[j].Value, keyword2) {
					return result, true
				}
			}
		}
	}
	return nil, false
}

// PharseCheck takes a string, patterns and a matchOne boolean and returns
// any matches and errors.
func PhraseCheck(s string, patterns []*Pattern, matchOne bool) ([]*Matched, error) {
	in := strings.NewReader(s)
	return PhraseCheckReader(in, patterns, matchOne)
}

// PhraseCheckReader evaluates the input from an io.Reader for all patterns.
func PhraseCheckReader(reader io.Reader, patterns []*Pattern, matchOne bool) ([]*Matched, error) {
	result := []*Matched{}
	tokens, err := TokenReader(reader)
	if err != nil {
		return nil, err
	}
	for _, pattern := range patterns {
		switch pattern.Type {
		case Keyword:
			for _, token := range tokens {
				if tokenMatches(token.Value, pattern.Keyword1) {
					result = append(result, &Matched{
						Text: token.Value,
						Pattern: pattern.OriginalText,
						LineNo: token.LineNo,
					})
				}
			}
		case Proximity:
			if token, ok := CheckProximity(tokens, pattern.Keyword1, pattern.Keyword2, pattern.MaxDistance); ok {
				result = append(result, &Matched{
					Text: token.Value,
					Pattern: pattern.OriginalText,
					LineNo: token.LineNo,
				})
			}
		}
		if matchOne && len(result) > 0 {
			return result, nil
		}
	}
	return result, nil
}

type PhraseCheckApp struct {
	appName string
}

const phraseCheckCSVHeader = "\"filename\",\"line no\",\"pattern\",\"phrase\""

// checkFile will read a file stream and display matches to standard out and return any errors
func checkFile(fName string, patterns []*Pattern) error {
	in, err := os.Open(fName)
	if err != nil {
		return err
	}
	defer in.Close()
	matches, err := PhraseCheckReader(bufio.NewReader(in), patterns, false)
	if err != nil {
		return err
	}
	for _, match := range matches {
		fmt.Printf("%q,%s\n", fName, match.String())
	}
	return nil
}

// checkDirectory takes an initial path, a set of pattens and optional exclude list and
// walks the directory and reports matches for any text files found.
func checkDirectory(startDir string, patterns []*Pattern, excludeList []string) error {
	fmt.Println(phraseCheckCSVHeader)
	err := filepath.Walk(startDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip if it's a directory and in the exclude list
		if info.IsDir() {
			for _, exclude := range excludeList {
				if strings.Contains(path, exclude) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// Skip if the file is in the exclude list
		for _, exclude := range excludeList {
			if strings.Contains(path, exclude) {
				return nil
			}
		}
		return checkFile(path, patterns)
	})
	return err
}


func (app *PhraseCheckApp) CheckFile(params []string) error {
	if len(params) < 2 {
		return fmt.Errorf("missing pattern filename and files to process")
	}
	var fName string
	fName, params = params[0], params[1:]
	patterns, err := LoadPatterns(fName)
	if err != nil {
		return err
	}
	fmt.Println(phraseCheckCSVHeader)
	for _, checkFName := range params {
		if err := checkFile(checkFName, patterns); err != nil {
			return err
		}
	}
	return err
}

func (app *PhraseCheckApp) CheckDirectory(params []string) error {
	if len(params) < 2 {
		return fmt.Errorf("missing pattern filename and directory to process")
	}
	var (
		fName string // pattern filename
		dirName string // directory to walk
		excludeList []string // an option list of paths to exclude
		err error
	)
	if len(params) >= 2 {
		fName, dirName = params[0], params[1]
	}
	if len(params) > 2 {
		excludeList, err = parseExcludeListFile(params[2])
		if err != nil {
			return err
		}
	}
	patterns, err := LoadPatterns(fName)
	if err != nil {
		return err
	}
	if err := checkDirectory(dirName, patterns, excludeList); err != nil {
		return err
	}
	return err
}

func parseExcludeListFile(excludeListName string) ([]string, error) {
	excludeList := []string{}
	src, err := os.ReadFile(excludeListName)
	if err != nil {
		return nil, fmt.Errorf("unable to read %q, %s", excludeListName, err)
	}
	txt := fmt.Sprintf("%s", src)
	for _, line := range strings.Split(txt, "\n") {
		if strings.TrimSpace(line) != "" {
			excludeList = append(excludeList, strings.TrimSpace(line))
		}
	}
	return excludeList, err
}

func (app *PhraseCheckApp) FileTypes(params []string) error {
	if len(params) == 0 {
		return fmt.Errorf("expected a starting directory to crawl")
	}
	if len(params) > 2 {
		return fmt.Errorf("too many parameters provided")
	}
	var (
		startDir string
		excludeListName string
		excludeList []string
		err error
	)
	if len(params) == 1 {
		startDir = params[0]
	}
	if len(params) == 2 {
		startDir, excludeListName = params[0], params[1]
	}
	if excludeListName != "" {
		excludeList, err = parseExcludeListFile(excludeListName)
		if err != nil {
			return err
		}
	}
	fileTypes, err := FileTypes(startDir, excludeList)
	if err != nil {
		return err
	}
	fmt.Printf("\"file path\",\"mime type\"\n")
	for file, fileType := range fileTypes {
		fmt.Printf("%q,%q\n", file, fileType)
	}
	return nil
}

func (app *PhraseCheckApp) FileTypeCounts(params []string) error {
	if len(params) == 0 {
		return fmt.Errorf("expected a starting directory to crawl")
	}
	if len(params) > 2 {
		return fmt.Errorf("too many parameters provided")
	}
	var (
		startDir string
		excludeListName string
		excludeList []string
		err error
	)
	if len(params) == 1 {
		startDir = params[0]
	}
	if len(params) == 2 {
		startDir, excludeListName = params[0], params[1]
	}
	if excludeListName != "" {
		excludeList, err = parseExcludeListFile(excludeListName)
		if err != nil {
			return err
		}
	}
	fileTypes, err := FileTypes(startDir, excludeList)
	if err != nil {
		return err
	}
	cnts := map[string]int{}
	for file, _ := range fileTypes {
		ext := path.Ext(file)
		if cnt, ok := cnts[ext]; ok {
			cnts[ext] = cnt + 1
		} else {
			cnts[ext] = 1
		}
	}
	fmt.Printf("\"file ext\",\"mime type\",\"count\"\n")
	for k, v := range cnts {
		if mimeType, ok := extensionToMIME[k]; ok {
			fmt.Printf("%q,%q,%d\n", k, mimeType, v)
		} else {
			fmt.Printf("%q,%q,%d\n", k, "", v)
		}
	}
	return nil
}

func (app *PhraseCheckApp) PruneMatches(params []string) error {
	return fmt.Errorf("PruneMatches(%+v) not implemented", params)
}


func (app *PhraseCheckApp) Run(appName string, action string, params []string) error {
	app.appName = appName
	switch action {
	case "help":
		 fmt.Fprintf(os.Stdout, "%s\n", FmtHelp(HelpText, appName, Version, ReleaseDate, ReleaseHash))
		 return nil
	case "filetypes":
		return app.FileTypes(params)
	case "filetype-counts":
		return app.FileTypeCounts(params)
	case "check-file":
		return app.CheckFile(params)
	case "check-directory":
		return app.CheckDirectory(params)
	case "prune":
		return app.PruneMatches(params)
	default:
		return fmt.Errorf("%q action not supported", action)
	}
	return nil
}