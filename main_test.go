package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	rc := m.Run()

	// rc 0 means we've passed,
	// and CoverMode will be non empty if run with -cover
	if rc == 0 && testing.CoverMode() != "" {
		c := testing.Coverage()
		if c < 0.8 {
			fmt.Println("Tests passed but coverage failed at", c)
			rc = -1
		}
	}
}

func TestDistance(t *testing.T) {
	word1 := "hello"
	word2 := "hel"

	distance := distance(word2, word1)
	require.Equal(t, 2, distance)
}
