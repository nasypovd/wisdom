package pow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPoW(t *testing.T) {
	pow := NewPoW(2)
	assert.Equal(t, 2, pow.Difficulty, "Expected Difficulty to be 2")
}

func TestGenerate(t *testing.T) {
	pow := NewPoW(2)
	challenge := pow.Generate()
	assert.Equal(t, 2, challenge.Difficulty, "Expected Difficulty to be 2")
}

func TestSolveAndVerify(t *testing.T) {
	pow := NewPoW(2)
	solver := NewSolver()
	challenge := pow.Generate()
	solution := solver.Solve(challenge)

	assert.True(t, pow.Verify(challenge, solution), "Expected solution to be correct for the generated challenge")

	// Negative test
	assert.False(t, pow.Verify(challenge, "wrong solution"), "Expected 'wrong solution' to be incorrect for the generated challenge")
}
