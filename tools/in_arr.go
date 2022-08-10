package tools

func InIntArr(in int, arr []int) bool {
	for _, v := range arr {
		if in == v {
			return true
		}
	}
	return false
}

func InStrArr(in string, arr []string) bool {
	for _, v := range arr {
		if in == v {
			return true
		}
	}
	return false
}
