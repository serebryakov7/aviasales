package anagram

import (
	"sync"
	"unicode"
)

type Anagram struct {
	mu sync.RWMutex

	runes    map[rune]int64
	anagrams map[int64][]string

	seqCh <-chan int64
}

func NewAnagram() *Anagram {
	seqCh := make(chan int64, 1)

	go generate(seqCh)

	return &Anagram{
		runes:    make(map[rune]int64),
		anagrams: make(map[int64][]string),
		seqCh:    seqCh,
	}
}

func (m *Anagram) Find(w string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	prime := m.primify(w)

	return m.anagrams[prime]
}

func (m *Anagram) InsertWords(words ...string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, w := range words {
		prime := m.parseWordRunes(w)
		m.anagrams[prime] = append(m.anagrams[prime], w)
	}
}

func (m *Anagram) primify(w string) int64 {
	var out int64 = 1
	for _, r := range w {
		r = unicode.ToLower(r)

		n, ok := m.runes[r]
		if !ok {
			return 1
		}

		out *= n
	}

	return out
}

func (m *Anagram) parseWordRunes(w string) int64 {
	var out int64 = 1
	for _, r := range w {
		r = unicode.ToLower(r)

		next, ok := m.runes[r]
		if !ok {
			next := <-m.seqCh
			m.runes[r] = next
		}

		out *= next
	}

	return out
}
