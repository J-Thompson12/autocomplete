package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {

	c := testing.Coverage()
	require.GreaterOrEqual(t, .8, c, fmt.Sprintf("Coverage failed at: %v", c))
}

func TestDistance(t *testing.T) {
	word1 := "hello"
	word2 := "hel"

	distance := distance(word2, word1)
	require.Equal(t, 2, distance)
}
