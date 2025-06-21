package core

type Turn int

func (t Turn) YearMonth() (int, int) {
	return int(t/12) + 1, int(t%12) + 1
}
