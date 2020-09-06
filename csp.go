package csp

import "fmt"

type (
	// Constraint defines the CSP's constrains
	Constraint interface {
		Variables() []Variable
		Satisfied(assignment map[Variable]Domain) bool
	}

	// Variable define the interface of the variables in a CSP solution
	Variable interface {
	}

	// Domain define the interface of possible domains of a variable
	Domain interface {
	}

	// Solution is a CSP solution
	Solution struct {
		Variables   []Variable
		Domains     map[Variable][]Domain
		Constraints map[Variable][]Constraint
	}
)

// NewSolution creates a new CSP solution
func NewSolution(variables []Variable, domains map[Variable][]Domain) (*Solution, error) {
	s := &Solution{
		Domains:     make(map[Variable][]Domain),
		Constraints: make(map[Variable][]Constraint),
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
func (s *Solution) Consistent(v Variable, assignment map[Variable]Domain) bool {
	for _, c := range s.Constraints[v] {
		if !c.Satisfied(assignment) {
			return false
		}
	}

	return true
}

// Search for a solution using backtracking
func (s *Solution) Search() map[Variable]Domain {
	initAssignment := make(map[Variable]Domain)
	return s.backtrackingSearch(initAssignment)
}

func (s *Solution) backtrackingSearch(assignment map[Variable]Domain) map[Variable]Domain {
	if len(assignment) == len(s.Variables) {
		return assignment
	}

	// select the first unassigned variable
	var unassigned Variable
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

func cloneAssignment(assignment map[Variable]Domain) map[Variable]Domain {
	retval := make(map[Variable]Domain)
	for k, v := range assignment {
		retval[k] = v
	}

	return retval
}
