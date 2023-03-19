package utils

func InIntArr(in int, arr []int) bool {
	for _, v := range arr {
		if in == v {
			return true
		}
	}
	return false
}

func InInt64Arr(in int64, arr []int64) bool {
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

func DelInIntArr(del int, arr []int) []int {
	j := 0
	for _, val := range arr {
		if val != del {
			arr[j] = val
			j++
		}
	}
	return arr[:j]
}

func DelInt64Arr(del int64, arr []int64) []int64 {
	j := 0
	for _, val := range arr {
		if val != del {
			arr[j] = val
			j++
		}
	}
	return arr[:j]
}

func DelStrArr(del string, arr []string) []string {
	j := 0
	for _, val := range arr {
		if val != del {
			arr[j] = val
			j++
		}
	}
	return arr[:j]
}
