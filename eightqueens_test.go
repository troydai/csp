package csp

// use CSP solution to solve the eight queens problem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type eightQueenConstraints struct{}

func buildVariables() []int {
	return []int{1, 2, 3, 4, 5, 6, 7, 8}
}

func buildDomains() map[int][]int {
	retval := make(map[int][]int)
	for _, c := range buildVariables() {
		retval[c] = []int{1, 2, 3, 4, 5, 6, 7, 8}
	}
	return retval
}

func abs(n int) int {
	if n < 0 {
		return n * -1
	}
	return n
}

func (c *eightQueenConstraints) Variables() []int {
	return buildVariables()
}

func (c *eightQueenConstraints) Satisfied(assignment Assignment) bool {
	for q1c, q1r := range assignment {
		// first queen's position is q1c / q1r
		// no other queen at the same row
		for q2c := q1c + 1; int(q2c) < len(assignment)+1; q2c++ {
			// for each remaining queens in following columns
			q2r := assignment[q2c]
			if q1r == q2r {
				// same row
				return false
			}
			if abs(int(q1c-q2c)) == abs(int(q1r-q2r)) {
				// diagnoal
				return false
			}
		}

	}

	return true
}

// TestEightQueens test eight queen solution
func TestEightQueens(t *testing.T) {
	// s, err := NewSolution(buildVariables(), buildDomains())
	// s.AddConstraint(&eightQueenConstraints{columns: buildVariables()})
	result, err := BacktrackingCSP(
		buildVariables(),
		buildDomains(),
		[]Constraint{&eightQueenConstraints{}})

	assert.NotNil(t, result)
	assert.NoError(t, err)
	assert.Equal(t, 8, len(result))
}
