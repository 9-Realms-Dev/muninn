package util

const divisor = 4

func SetDefaultHeight(height int) int {
	return height - (divisor * 3)
}
