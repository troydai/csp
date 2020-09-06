package csp

import "fmt"

type (
	// Assignment is a set of domain to variable assignemtn
	Assignment map[VariableKey]DomainKey

	// Constraint defines the CSP's constrains
	Constraint interface {
		Variables() []VariableKey
		Satisfied(assignment Assignment) bool
	}

	// Solution is a CSP solution
	Solution struct {
		Variables   []VariableKey
		Domains     map[VariableKey][]DomainKey
		Constraints map[VariableKey][]Constraint
	}

	// VariableKey identify variable
	VariableKey int

	// DomainKey identify domain
	DomainKey int
)

// NewSolution creates a new CSP solution
func NewSolution(variables []VariableKey, domains map[VariableKey][]DomainKey) (*Solution, error) {
	s := &Solution{
		Domains:     make(map[VariableKey][]DomainKey),
		Constraints: make(map[VariableKey][]Constraint),
	}

	for v, ds := range domains {
		s.Domains[v] = ds
	}

	for _, v := range variables {
		if ds, found := s.Domains[v]; !found || len(ds) == 0 {
			return nil, fmt.Errorf("variable %v does not have any domains", v)
		}

		s.Variables = append(s.Variables, v)
		s.Constraints[v] = make([]Constraint, 0)
	}

	return s, nil
}

// AddConstraint adds a constraint to a solution
func (s *Solution) AddConstraint(c Constraint) {
	for _, v := range c.Variables() {
		s.Constraints[v] = append(s.Constraints[v], c)
	}
}

// Consistent check if the value assignment is consistent by checking all constraints
// for the given variable againsts it
func (s *Solution) Consistent(v VariableKey, assignment Assignment) bool {
	for _, c := range s.Constraints[v] {
		if !c.Satisfied(assignment) {
			return false
		}
	}

	return true
}

// Search for a solution using backtracking
func (s *Solution) Search() Assignment {
	initAssignment := make(Assignment)
	return s.backtrackingSearch(initAssignment)
}

func (s *Solution) backtrackingSearch(assignment Assignment) Assignment {
	if len(assignment) == len(s.Variables) {
		return assignment
	}

	// select the first unassigned variable
	var unassigned VariableKey
	for _, v := range s.Variables {
		if _, found := assignment[v]; !found {
			unassigned = v
			break
		}
	}

	for _, d := range s.Domains[unassigned] {
		assignmentCopy := cloneAssignment(assignment)
		assignmentCopy[unassigned] = d
		if s.Consistent(unassigned, assignmentCopy) {
			result := s.backtrackingSearch(assignmentCopy)
			if result != nil {
				return result
			}
		}
	}

	return nil
}

func cloneAssignment(assignment Assignment) Assignment {
	retval := make(Assignment)
	for k, v := range assignment {
		retval[k] = v
	}

	return retval
}
