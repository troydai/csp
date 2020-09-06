package csp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type eightQueenConstraints struct {
	columns []VariableKey
}

// use CSP solution to solve the eight queens problem

func buildVariables() []VariableKey {
	return []VariableKey{1, 2, 3, 4, 5, 6, 7, 8}
}

func buildDomains() map[VariableKey][]DomainKey {
	retval := make(map[VariableKey][]DomainKey)
	for _, c := range buildVariables() {
		retval[c] = []DomainKey{1, 2, 3, 4, 5, 6, 7, 8}
	}
	return retval
}

func abs(n int) int {
	if n < 0 {
		return n * -1
	}
	return n
}

func (c *eightQueenConstraints) Variables() []VariableKey {
	return c.columns
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

func printSolution(assignment Assignment) string {
	cr := []rune("ABCDEFGH")
	result := ""
	for c, r := range assignment {
		result += fmt.Sprintf("%q%v, ", cr[c-1], r)
	}

	return result
}

// TestEightQueens test eight queen solution
func TestEightQueens(t *testing.T) {
	s, err := NewSolution(buildVariables(), buildDomains())
	s.AddConstraint(&eightQueenConstraints{columns: buildVariables()})
	assert.NotNil(t, s)
	assert.NoError(t, err)

	result := s.Search()
	r := printSolution(result)

	assert.NotNil(t, result)
	assert.NotEmpty(t, result)
	assert.NotEmpty(t, r)
}
