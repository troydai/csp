package csp

import "fmt"

// Implement CSP in a standalone public function

type (
	// Constraint defines an interface use which the user can implement logic for specific problems
	Constraint interface {
		// Variables return the varibles this constrain apply
		Variables() []int

		// Satisfied returns true if the given assignment satisfy this constrain
		Satisfied(assignment Assignment) bool
	}

	// Assignment is a set of domain assigned to variables
	Assignment map[int]int

	// ConstraintsMap group constraints for variables
	ConstraintsMap map[int][]Constraint

	// DomainMap lists the allowed domain for each variable
	DomainMap map[int][]int

	// VariableList defines a list of variable
	VariableList []int

	collector struct {
		results []Assignment
	}
)

// BacktrackingCSP tries to solve a constraint-satisfaction problem using backtracking
func BacktrackingCSP(
	variables VariableList,
	domain DomainMap,
	constrains []Constraint,
	firstResult bool,
) ([]Assignment, error) {
	// init
	constraintsMap := make(ConstraintsMap)
	for _, c := range constrains {
		for _, v := range c.Variables() {
			constraintsMap[v] = append(constraintsMap[v], c)
		}
	}

	for _, v := range variables {
		if ds, ok := domain[v]; !ok || len(ds) == 0 {
			return nil, fmt.Errorf("variable %v doesn't have any domains", v)
		}
		if _, found := constraintsMap[v]; !found {
			constraintsMap[v] = make([]Constraint, 0)
		}
	}

	c := &collector{}
	if err := search(variables, domain, constraintsMap, make(Assignment), firstResult, c); err != nil {
		return nil, err
	}
	return c.Results(), nil
}

func search(
	variables []int,
	domain map[int][]int,
	constrains ConstraintsMap,
	assignment Assignment,
	firstResult bool,
	c *collector,
) error {
	// found the assignment
	if len(assignment) == len(variables) {
		c.Add(assignment)
		return nil
	}

	// retrive the first unassigned variable
	var unassigned int
	for _, v := range variables {
		if _, found := assignment[v]; !found {
			unassigned = v
			break
		}
	}

	// try each domain and backtrack if fails
	for _, d := range domain[unassigned] {
		newAssignment := make(Assignment)
		for v, d := range assignment {
			newAssignment[v] = d
		}
		newAssignment[unassigned] = d

		// test if this domain of the variable satified all constraints
		satified := true
		for _, c := range constrains[unassigned] {
			if !c.Satisfied(newAssignment) {
				satified = false
				break
			}
		}

		// found one domain, try to build next variable base on current assignment
		// if this branch doesn't workout, continue to try next domain
		if satified {
			if err := search(variables, domain, constrains, newAssignment, firstResult, c); err != nil {
				return err
			}

			if firstResult && len(c.results) > 0 {
				return nil
			}
		}
	}

	return nil
}

func (c *collector) Add(assignment Assignment) {
	c.results = append(c.results, assignment)
}

func (c *collector) Results() []Assignment {
	return c.results
}
