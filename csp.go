package csp

import "fmt"

// Implement CSP in a standalone public function

// Constraint defines an interface use which the user can implement logic for specific problems
type Constraint interface {
	// Variables return the varibles this constrain apply
	Variables() []int

	// Satisfied returns true if the given assignment satisfy this constrain
	Satisfied(assignment Assignment) bool
}

// Assignment is a set of domain assigned to variables
type Assignment map[int]int

// ConstraintsMap group constraints for variables
type ConstraintsMap map[int][]Constraint

// DomainMap lists the allowed domain for each variable
type DomainMap map[int][]int

// VariableList defines a list of variable
type VariableList []int

// BacktrackingCSP tries to solve a constraint-satisfaction problem using backtracking
func BacktrackingCSP(
	variables VariableList,
	domain DomainMap,
	constrains []Constraint,
) (Assignment, error) {
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

	return search(variables, domain, constraintsMap, make(Assignment))
}

func search(
	variables []int,
	domain map[int][]int,
	constrains ConstraintsMap,
	assignment Assignment,
) (Assignment, error) {
	// found the assignment
	if len(assignment) == len(variables) {
		return assignment, nil
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
			result, err := search(variables, domain, constrains, newAssignment)
			if err != nil {
				return nil, err
			}
			if result != nil {
				return result, nil
			}
		}
	}

	return nil, nil
}
