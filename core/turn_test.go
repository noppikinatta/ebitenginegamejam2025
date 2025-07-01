package core_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

func TestTurn_YearMonth(t *testing.T) {
	tests := []struct {
		name          string
		turn          core.Turn
		expectedYear  int
		expectedMonth int
	}{
		{
			name:          "First turn",
			turn:          core.Turn(0),
			expectedYear:  1,
			expectedMonth: 1,
		},
		{
			name:          "Last month of year 1",
			turn:          core.Turn(11),
			expectedYear:  1,
			expectedMonth: 12,
		},
		{
			name:          "First month of year 2",
			turn:          core.Turn(12),
			expectedYear:  2,
			expectedMonth: 1,
		},
		{
			name:          "Middle of year 2",
			turn:          core.Turn(18),
			expectedYear:  2,
			expectedMonth: 7,
		},
		{
			name:          "Beginning of year 3",
			turn:          core.Turn(24),
			expectedYear:  3,
			expectedMonth: 1,
		},
		{
			name:          "May of year 10",
			turn:          core.Turn(112), // (10-1)*12 + (5-1) = 108 + 4 = 112
			expectedYear:  10,
			expectedMonth: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			year, month := tt.turn.YearMonth()
			if year != tt.expectedYear {
				t.Errorf("Year = %d, want %d", year, tt.expectedYear)
			}
			if month != tt.expectedMonth {
				t.Errorf("Month = %d, want %d", month, tt.expectedMonth)
			}
		})
	}
}

func TestTurn_Next(t *testing.T) {
	tests := []struct {
		name     string
		turn     core.Turn
		expected core.Turn
	}{
		{
			name:     "From first turn to next",
			turn:     core.Turn(0),
			expected: core.Turn(1),
		},
		{
			name:     "From regular turn to next",
			turn:     core.Turn(5),
			expected: core.Turn(6),
		},
		{
			name:     "Turn crossing years",
			turn:     core.Turn(11),
			expected: core.Turn(12),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.turn.Next()
			if result != tt.expected {
				t.Errorf("Next() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestTurn_String(t *testing.T) {
	tests := []struct {
		name     string
		turn     core.Turn
		expected string
	}{
		{
			name:     "First turn",
			turn:     core.Turn(0),
			expected: "Year 1, Month 1",
		},
		{
			name:     "December of year 1",
			turn:     core.Turn(11),
			expected: "Year 1, Month 12",
		},
		{
			name:     "January of year 2",
			turn:     core.Turn(12),
			expected: "Year 2, Month 1",
		},
		{
			name:     "May of year 10",
			turn:     core.Turn(112),
			expected: "Year 10, Month 5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.turn.String()
			if result != tt.expected {
				t.Errorf("String() = %q, want %q", result, tt.expected)
			}
		})
	}
}
