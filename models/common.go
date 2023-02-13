package models

var count int

func fibonacci() int {
	count++
	switch count {
	case 1:
		return 0
	case 2:
		return 1
	case 3:
		return 1
	case 4:
		return 2
	case 5:
		return 3
	case 6:
		return 5
	case 7:
		return 8
	default:
		return fibonacci() + fibonacci()
	}
}

// func fibonacci() func() int {
// 	x, y := 0, 1
// 	return func() int {
// 		x, y = y, x+y
// 		if x > 8 {
// 			return x
// 		}
// 		return x
// 	}
// }
