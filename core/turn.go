package core

import "fmt"

type Turn int

func (t Turn) YearMonth() (int, int) {
	return int(t/12) + 1, int(t%12) + 1
}

// Next は次のターンを返します
func (t Turn) Next() Turn {
	return t + 1
}

// String はターンを「N年M月」形式の文字列で返します
func (t Turn) String() string {
	year, month := t.YearMonth()
	return fmt.Sprintf("%d年%d月", year, month)
}
