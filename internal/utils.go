package internal

const (
	DefaultAperture = 3
)

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
