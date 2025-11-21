package analysistools

import (
	//"fmt"
	//"os"
	"testing"
)

func TestWordTokenizer(t *testing.T) {
	txt := `This is a humble
multilined text. Is it but a 
  poetic chorus

Yet without a double is still presents a test.
`

	tokens, err := Tokenizer(txt)
	if err != nil {
		t.Errorf("expected err == nil, got %s", err)
	}
	expected := []*Token{
		{ Value: "This", LineNo: 0, WordNo: 0 },
		{ Value: "is", LineNo: 0, WordNo: 1 },
		{ Value: "a", LineNo: 0, WordNo: 2 },
		{ Value: "humble", LineNo: 0, WordNo: 3 },
		{ Value: "multilined", LineNo: 1, WordNo: 4 },
		{ Value: "text.", LineNo: 1, WordNo: 5 },
		{ Value: "Is", LineNo: 1, WordNo: 6 },
		{ Value: "it", LineNo: 1, WordNo: 7 },
		{ Value: "but", LineNo: 1, WordNo: 8 },
		{ Value: "a", LineNo: 1, WordNo: 9 },
		{ Value: "poetic", LineNo: 2, WordNo: 10 },
		{ Value: "chorus", LineNo: 2, WordNo: 11 },
		{ Value: "Yet", LineNo: 4, WordNo: 12 },
		{ Value: "without", LineNo: 4, WordNo: 13 },
		{ Value: "a", LineNo: 4, WordNo: 14 },
		{ Value: "double", LineNo: 4, WordNo: 15 },
		{ Value: "is", LineNo: 4, WordNo: 16 },
		{ Value: "still", LineNo: 4, WordNo: 17 },
		{ Value: "presents", LineNo: 4, WordNo: 18 },
		{ Value: "a", LineNo: 4, WordNo: 19 },
		{ Value: "test.", LineNo: 4, WordNo: 20 },
	}
	if len(tokens) != len(expected) {
		t.Errorf("expected %d tokens, got %d", len(expected), len(tokens))
	}
	if len(tokens) < len(expected) {
		for i := len(tokens); i < len(expected); i++ {
			tokens = append(tokens, &Token{})
		}
	}
	for i, tok := range tokens {
		if i >= len(expected) {
			t.Errorf("unexpected token, %d %+v", i, tok)
			continue
		}
		if tok.Value != expected[i].Value {
			t.Errorf("expectd value %q, got %q, token #%d %+v", expected[i].Value, tok.Value, i, tok)
		}
		if tok.LineNo != expected[i].LineNo {
			t.Errorf("expectd LineNo %d, got %d, token #%d %+v", expected[i].LineNo, tok.LineNo, i, tok)
		}
		if tok.WordNo != expected[i].WordNo {
			t.Errorf("expectd WordNo %d, got %d, token #%d %+v", expected[i].WordNo, tok.WordNo, i, tok)
		}
	}
}