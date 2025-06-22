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
			name:          "最初のターン",
			turn:          core.Turn(0),
			expectedYear:  1,
			expectedMonth: 1,
		},
		{
			name:          "1年目の最後の月",
			turn:          core.Turn(11),
			expectedYear:  1,
			expectedMonth: 12,
		},
		{
			name:          "2年目の最初の月",
			turn:          core.Turn(12),
			expectedYear:  2,
			expectedMonth: 1,
		},
		{
			name:          "2年目の途中",
			turn:          core.Turn(18),
			expectedYear:  2,
			expectedMonth: 7,
		},
		{
			name:          "3年目の最初",
			turn:          core.Turn(24),
			expectedYear:  3,
			expectedMonth: 1,
		},
		{
			name:          "10年目の5月",
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
			name:     "最初のターンから次へ",
			turn:     core.Turn(0),
			expected: core.Turn(1),
		},
		{
			name:     "通常のターンから次へ",
			turn:     core.Turn(5),
			expected: core.Turn(6),
		},
		{
			name:     "年をまたぐターン",
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
			name:     "最初のターン",
			turn:     core.Turn(0),
			expected: "1年1月",
		},
		{
			name:     "1年目12月",
			turn:     core.Turn(11),
			expected: "1年12月",
		},
		{
			name:     "2年目1月",
			turn:     core.Turn(12),
			expected: "2年1月",
		},
		{
			name:     "10年目5月",
			turn:     core.Turn(112),
			expected: "10年5月",
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
