package confluence

import "fmt"

func contains[K comparable](s []K, item K) bool {
	for _, v := range s {
		if v == item {
			return true
		}
	}
	return false
}

func moveToFirstPositionOfSlice[K comparable](slice []K, item K) []K {
	if len(slice) == 0 || (slice)[0] == item {
		fmt.Println(slice)
		return nil
	}
	if (slice)[len(slice)-1] == item {
		slice = append([]K{item}, (slice)[:len(slice)-1]...)
		fmt.Println(slice)
		return nil
	}
	for p, x := range slice {
		if x == item {
			(slice) = append([]K{item}, append((slice)[:p], (slice)[p+1:]...)...)
			break
		}
	}
	return slice
}
