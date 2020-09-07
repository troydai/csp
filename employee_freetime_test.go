package csp

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

// https://leetcode.com/problems/employee-free-time/description/

type (
	Interval struct {
		Start int
		End   int
	}

	employeeFreetimeTestCase struct {
		Schedule [][]*Interval
		Expected []*Interval
	}

	employeeFreetimeConstraint struct {
		UserIDs []int
	}
)

func (c *employeeFreetimeConstraint) Variables() []int {
	return c.UserIDs
}

func (c *employeeFreetimeConstraint) Satisfied(assignment Assignment) bool {
	if len(assignment) == 1 {
		return true
	}

	pickedTime := -1
	for _, t := range assignment {
		if pickedTime == -1 {
			pickedTime = t
			continue
		}

		if t != pickedTime {
			return false
		}
	}

	return true
}

func employeeFreeTime(schedule [][]*Interval) []*Interval {
	// find domain range
	low := -1
	high := -1
	for _, slots := range schedule {
		for _, slot := range slots {
			if low == -1 {
				low = slot.Start
			} else if slot.Start < low {
				low = slot.Start
			}

			if high == -1 {
				high = slot.End
			} else if slot.End > high {
				high = slot.End
			}
		}
	}

	userIDs := make([]int, 0, len(schedule))
	userTimeSlots := make(map[int][]int)
	for i, slots := range schedule {
		userIDs = append(userIDs, i)
		userTimeSlots[i] = revertSchedule(low, high, slots)
	}

	assignemtns, err := BacktrackingCSP(
		userIDs,
		userTimeSlots,
		[]Constraint{&employeeFreetimeConstraint{UserIDs: userIDs}},
		false,
	)
	if err != nil {
		return nil
	}

	var slots []int
	for _, v := range assignemtns {
		slots = append(slots, v[0])
	}
	sort.Ints(slots)

	// convert it into intervals
	retval := []*Interval{{Start: slots[0], End: slots[0] + 1}}
	for i := 1; i < len(slots); i++ {
		last := retval[len(retval)-1]
		if slots[i] == last.End {
			last.End = last.End + 1
		} else {
			retval = append(retval, &Interval{Start: slots[i], End: slots[i] + 1})
		}
	}

	return retval
}

func revertSchedule(low, high int, intervals []*Interval) []int {
	occupied := make(map[int]bool)
	for _, interval := range intervals {
		for i := interval.Start; i < interval.End; i++ {
			occupied[i] = true
		}
	}
	var retval []int
	for i := low; i < high; i++ {
		if _, found := occupied[i]; !found {
			retval = append(retval, i)
		}
	}

	return retval
}

func TestEmployeeFreeTimeSolution(t *testing.T) {
	testcases := []*employeeFreetimeTestCase{
		{
			Schedule: [][]*Interval{
				[]*Interval{{Start: 1, End: 2}, {Start: 5, End: 6}},
				[]*Interval{{Start: 1, End: 3}},
				[]*Interval{{Start: 4, End: 10}},
			},
			Expected: []*Interval{
				{Start: 3, End: 4},
			},
		},
		{
			Schedule: [][]*Interval{
				[]*Interval{{Start: 1, End: 3}, {Start: 6, End: 7}},
				[]*Interval{{Start: 2, End: 4}},
				[]*Interval{{Start: 2, End: 5}, {Start: 9, End: 12}},
			},
			Expected: []*Interval{
				{Start: 5, End: 6},
				{Start: 7, End: 9},
			},
		},
	}

	for i, tc := range testcases {
		t.Run(fmt.Sprintf("test case %v", i), func(t *testing.T) {
			retval := employeeFreeTime(tc.Schedule)
			assert.NotEmpty(t, retval)
			assert.Equal(t, len(tc.Expected), len(retval))
			for i, expected := range tc.Expected {
				assert.EqualValues(t, expected, retval[i])
			}
		})
	}
}
