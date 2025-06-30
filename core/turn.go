package core

import "fmt"

type Turn int

func (t Turn) YearMonth() (int, int) {
	return int(t/12) + 1, int(t%12) + 1
}

// Next returns the next turn.
func (t Turn) Next() Turn {
	return t + 1
}

// String returns the turn as a string in the format "Year N, Month M".
func (t Turn) String() string {
	year, month := t.YearMonth()
	return fmt.Sprintf("Year %d, Month %d", year, month)
}
