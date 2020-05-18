package sort

import "sort"

type Int64SLice []int64

func (slice Int64SLice) Len() int {
	return len(slice)
}

func (slice Int64SLice) Less(i, j int) bool {
	return slice[i] < slice[j]
}

func (slice Int64SLice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice Int64SLice) Sort() {
	sort.Sort(slice)
}

