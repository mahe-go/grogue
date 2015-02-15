package util

func Max(x int, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}

func Min(x int, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

func Abs(x int) int {
	if x < 0 {
		return 0 - x
	} else {
		return x
	}
}
