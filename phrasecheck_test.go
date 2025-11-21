package analysistools

import (
	//"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
)

const (
	checkPatterns = `acp
affiant
affidavit*
atorn*
attorn*
attorn* w/5 client*
attorney client privilege
attorney work product
Attorney*
attrny
atty*
ATTY-CLIENT PRIVILEDGE*
ATTY-CLIENT PRIVILEGE*
awp
complainant* w/5 brief*
complainant* w/5 complaint*
complainant* w/5 memo
complainant* w/5 memoranda
complainant* w/5 memorandum
complainant* w/5 motion*
complainant* w/5 petition*
counsel*
counselor
counsil*
declarant
declaration*
defendant* w/5 brief*
defendant* w/5 complaint*
defendant* w/5 memo
defendant* w/5 memoranda
defendant* w/5 memorandum
defendant* w/5 motion*
defendant* w/5 petition*
deponent*
Deponents
depos
depose*
Deposition
Depositions
Esq
esquire
law
laws
Lawsuit
lawyer
legal
Legal Counsel
litigat*
outside counsel
plaintiff* w/5 brief*
plaintiff* w/5 complaint*
plaintiff* w/5 memo
plaintiff* w/5 memoranda
plaintiff* w/5 memorandum
plaintiff* w/5 motion*
plaintiff* w/5 petition*
prepar* w/4 litigat
Priv*
Priveledge*
privelege*
Privelidge*
Privelige*
privil*
privilage*
Priviledge*
privilege*
privileged and confidential
Privledge*
respondent* w/5 brief*
respondent* w/5 complaint*
respondent* w/5 memo
respondent* w/5 memoranda
respondent* w/5 memorandum
respondent* w/5 motion*
respondent* w/5 petition*
suit
summary judgment
testify*
testimony*
transcript*
work w/2 product
wp
`

	inputNoMatches = `The meeting was scheduled for next week.
The report was submitted on time.
The project is on track for completion.
The team reviewed the budget proposal.
The weather was sunny and warm.
The presentation was well received.
The conference room is booked for tomorrow.
The software update was installed successfully.
The network is down for maintenance.
The inventory was counted and verified.
The shipment arrived on schedule.
The training session was informative.
The survey results were compiled.
The recipe was followed exactly.
The garden is blooming beautifully.
The car was serviced yesterday.
The book was returned to the library.
The movie was highly recommended.
The restaurant was fully booked.
The flight was delayed due to weather.
`

	inputWithMatches = `The attorney filed a motion for the client.
The complainant submitted a memorandum and a brief.
This document is privileged and confidential.
The defendant's counsel prepared a memo for the litigation.
The deposition transcript was reviewed by the attorney.
The plaintiff's brief was submitted to the court.
The respondent's petition was denied.
The affiant signed the affidavit.
The legal counsel advised on the lawsuit.
The complainant's complaint was dismissed.
The attorney-client privilege was asserted.
The work product was prepared for the case.
The deponent testified under oath.
The outside counsel was consulted.
The declarant signed the declaration.
The transcript of the testimony was filed.
The summary judgment motion was granted.
The Esq. provided legal advice.
The preparation for litigation is ongoing.
`
)

func TestParsePattern(t *testing.T) {
	tests := []struct {
		input    string
		want     *Pattern
		wantErr  bool
		errMsg   string
	}{
		{
			input: "attorn*",
			want: &Pattern{
				Type:         Keyword,
				Keyword1:      "attorn*",
				OriginalText: "attorn*",
			},
			wantErr: false,
		},
		{
			input: "attorn* w/5 client*",
			want: &Pattern{
				Type:         Proximity,
				Keyword1:     "attorn*",
				Keyword2:     "client*",
				MaxDistance:  5,
				OriginalText: "attorn* w/5 client*",
			},
			wantErr: false,
		},
		{
			input:   "attorn* w/ client*",
			want:    nil,
			wantErr: true,
			errMsg:  "invalid max distance in pattern: \"attorn* w/ client*\"",
		},
		{
			input:   "w/5 client*",
			want:    nil,
			wantErr: true,
			errMsg:  "malformed proximity pattern: \"w/5 client*\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParsePattern(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePattern(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ParsePattern(%q) error = %v, wantErrMsg %v", tt.input, err, tt.errMsg)
				return
			}
			if !tt.wantErr && !patternsEqual(got, tt.want) {
				src1, _ := json.Marshal(got)
				src2, _ := json.Marshal(tt.want)
				t.Errorf("ParsePattern(%q) returned\n\t%s,\nwant\n\t%s", tt.input, src1, src2)
			}
		})
	}
}

func TestLoadPatterns(t *testing.T) {
	// Create a temporary file for testing
	tmpfile, err := os.CreateTemp("", "testpatterns")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := []byte("attorn*\nattorn* w/5 client*\n")
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if _, err := tmpfile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	got, err := LoadPatterns(tmpfile.Name())
	if err != nil {
		t.Fatalf("LoadPatterns() error = %v", err)
	}

	want := []*Pattern{
		{
			Type:         Keyword,
			Keyword1:      "attorn*",
			OriginalText: "attorn*",
		},
		{
			Type:         Proximity,
			Keyword1:     "attorn*",
			Keyword2:     "client*",
			MaxDistance:  5,
			OriginalText: "attorn* w/5 client*",
		},
	}

	if len(got) != len(want) {
		t.Fatalf("LoadPatterns() got %d patterns, want %d", len(got), len(want))
	}
	for i := range got {
		if !patternsEqual(got[i], want[i]) {
			t.Errorf("LoadPatterns()[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestCheckProximity(t *testing.T) {
	txt := "the quick brown fox jumps"
	tokens, err := Tokenizer(txt)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	tests := []struct {
		tokens      []*Token
		keyword1    string
		keyword2    string
		maxDistance int
		want        bool
	}{
		{
			tokens: tokens,
			keyword1:    "quick",
			keyword2:    "fox",
			maxDistance: 2,
			want:        true,
		},
		{
			tokens: tokens,
			keyword1:    "quick",
			keyword2:    "jumps",
			maxDistance: 2,
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(txt, func(t *testing.T) {
			_, got := CheckProximity(tt.tokens, tt.keyword1, tt.keyword2, tt.maxDistance)
			if got != tt.want {
				t.Errorf("CheckProximity(%v, %q, %q, %d) = %v, want %v", tt.tokens, tt.keyword1, tt.keyword2, tt.maxDistance, got, tt.want)
			}
		})
	}
}

func TestPhraseCheckReader(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		patterns []*Pattern
		want     string
		wantErr bool
	}{
		{
			name:  "Keyword match",
			input: "The attorn represented the client.",
			patterns: []*Pattern{
				{
					Type:         Keyword,
					Keyword1:      "attorn",
					OriginalText: "attorn*",
				},
			},
			want: "0,\"attorn*\",\"attorn\"",
			wantErr: false,
		},
		{
			name:  "keyword match",
			input: "The attorn represented the client.",
			patterns: []*Pattern{
				{
					Type:         Keyword,
					Keyword1:     "attorn",
					OriginalText: "attorn",
				},
			},
			want: "0,\"attorn\",\"attorn\"",
			wantErr: false,
		},
		{
			name:  "keyword match",
			input: "The attorn represented the client.",
			patterns: []*Pattern{
				{
					Type:         Keyword,
					Keyword1:     "*sented",
					OriginalText: "*sented",
				},
			},
			want: "0,\"*sented\",\"represented\"",
			wantErr: false,
		},
		{
			name:  "proximity match",
			input: "The attorn represented the client.",
			patterns: []*Pattern{
				{
					Type:         Proximity,
					Keyword1:     "attorn*",
					Keyword2:     "client*",
					MaxDistance:  5,
					OriginalText: "attorn* w/5 client*",
				},
			},
			want: "0,\"attorn* w/5 client*\",\"attorn\"",
			wantErr: false,
		},
		{
			name: "proximity match",
			input: "Regarding ATTY-CLIENT PRIVILEDGE materials you sent",
			patterns: []*Pattern{
				{
					Type:     Proximity,
					Keyword1: "ATTY-CLIENT",
					Keyword2: "PRIVILEDGE",
					MaxDistance:  1,
					OriginalText: "ATTY-CLIENT PRIVILEDGE",
				},
			},
			want: "0,\"ATTY-CLIENT PRIVILEDGE\",\"ATTY-CLIENT\"",
			wantErr: false,
		},
		{
			name: "proximity match",
			input: "Regarding ATTY-CLIENT PRIVILEDGE materials you sent",
			patterns: []*Pattern{
				{
					Type:     Proximity,
					Keyword1: "ATTY-CLIENT",
					Keyword2: "PRIVILEDGE*",
					MaxDistance:  1,
					OriginalText: "ATTY-CLIENT PRIVILEDGE*",
				},
			},
			want: "0,\"ATTY-CLIENT PRIVILEDGE*\",\"ATTY-CLIENT\"",
			wantErr: false,
		},
	}

	for i, tt := range tests {
		matched, err := PhraseCheck(tt.input, tt.patterns, false)
		if err == nil && tt.wantErr  {
			t.Errorf("expected error, got matched %s for input %q", MatchedStrings(matched), tt.input)
			continue
		}
		if err != nil && tt.wantErr == false {
			t.Errorf("unexpected error %q for input %q", err, tt.input)
			continue
		}
		got := MatchedStrings(matched)
		if got != tt.want {
			t.Errorf("expected (%d) %q, got %q for input %q", i, tt.want, got, tt.input)
		}
	}
}

// Helper functions for testing
func patternsEqual(a, b *Pattern) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Type == b.Type &&
		a.Keyword1 == b.Keyword1 &&
		a.Keyword2 == b.Keyword2 &&
		a.MaxDistance == b.MaxDistance &&
		a.OriginalText == b.OriginalText
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// For capturing stdout in tests
var stdout io.Writer = os.Stdout
