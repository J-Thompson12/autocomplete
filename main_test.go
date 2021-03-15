package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDistance(t *testing.T) {
	word1 := "hello"
	word2 := "hel"

	distance := distance(word2, word1)
	require.Equal(t, 2, distance)
}
