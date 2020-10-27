package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/serebryakov7/aviasales/anagram"
)

type test struct {
	name      string
	word      string
	regWords  []string
	wantCode  int
	wantWords []string
}

func TestAnagramsHandlerGet(t *testing.T) {
	test := []test{
		{
			name:      "Success foobar",
			word:      "foobar",
			regWords:  []string{"foobar", "aabb", "baba", "boofar", "test"},
			wantCode:  http.StatusOK,
			wantWords: []string{"foobar", "boofar"},
		},
		{
			name:      "Success raboof",
			word:      "raboof",
			regWords:  []string{"foobar", "aabb", "baba", "boofar", "test"},
			wantCode:  http.StatusOK,
			wantWords: []string{"foobar", "boofar"},
		},
		{
			name:      "Success abba",
			word:      "abba",
			regWords:  []string{"foobar", "aabb", "baba", "boofar", "test"},
			wantCode:  http.StatusOK,
			wantWords: []string{"aabb", "baba"},
		},
		{
			name:      "Success test",
			word:      "test",
			regWords:  []string{"foobar", "aabb", "baba", "boofar", "test"},
			wantCode:  http.StatusOK,
			wantWords: []string{"test"},
		},
		{
			name:     "Success qwerty",
			word:     "qwerty",
			regWords: []string{"foobar", "aabb", "baba", "boofar", "test"},
			wantCode: http.StatusOK,
		},
		{
			name:     "No query",
			word:     "",
			wantCode: http.StatusBadRequest,
		},
	}

	logger := log.New(ioutil.Discard, "", log.LstdFlags)

	t.Parallel()
	for _, tt := range test {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				a   = anagram.NewAnagram()
				q   = url.Values{}
				req = httptest.NewRequest(http.MethodGet, "/anagrams", nil)
				res = httptest.NewRecorder()
			)

			a.Insert(tt.regWords...)
			q.Set("word", tt.word)
			req.URL.RawQuery = q.Encode()

			NewHandler(logger, a).ServeHTTP(res, req)

			assertEqual(t, tt.wantCode, res.Code)

			if tt.wantWords != nil {
				ws := make([]string, len(tt.wantWords))
				if err := json.Unmarshal(res.Body.Bytes(), &ws); err != nil {
					t.Fatal(err)
				}

				assertEqualSlices(t, ws, tt.wantWords)
			}
		})
	}
}

func assertEqual(t *testing.T, a, b interface{}) {
	if a == b {
		return
	}

	t.Fatal(fmt.Sprintf("%v is not equal to %v", a, b))
}

func assertEqualSlices(t *testing.T, a, b []string) {
	if (a == nil) != (b == nil) {
		t.Fatal(fmt.Errorf("%v is not equal to %v", a, b))
	}

	if len(a) != len(b) {
		t.Fatal(fmt.Errorf("%v is not equal to %v", a, b))
	}

	for i := range a {
		if a[i] != b[i] {
			t.Fatal(fmt.Errorf("%v is not equal to %v", a, b))
		}
	}
}
