package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/serebryakov7/aviasales/anagram"
)

func TestAnagramsHandlerPost(t *testing.T) {
	test := []test{
		{
			name:     "Success",
			regWords: []string{"foobar", "aabb", "baba", "boofar", "test"},
			wantCode: http.StatusNoContent,
		},
	}

	logger := log.New(ioutil.Discard, "", log.LstdFlags)

	t.Parallel()
	for _, tt := range test {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			if err := json.NewEncoder(&buf).Encode(tt.regWords); err != nil {
				t.Fatal(err)
			}

			var (
				a   = anagram.NewAnagram()
				req = httptest.NewRequest(http.MethodPost, "/anagrams", &buf)
				res = httptest.NewRecorder()
			)

			NewHandler(logger, a).ServeHTTP(res, req)

			assertEqual(t, tt.wantCode, res.Code)

			for _, word := range tt.regWords {
				ws := a.Find(word)
				if len(ws) < 1 {
					t.Fatal(fmt.Errorf("word %s has not been registered", word))
				}

				for i, w := range ws {
					if w == word {
						break
					} else if i > len(ws) {
						t.Fatal(fmt.Errorf("word %s not found in matching", word))
					}
				}
			}
		})
	}
}
