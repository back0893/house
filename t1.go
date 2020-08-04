package main

func minimumTotal(triangle [][]int) int {
	sum := triangle[0][0]
	j := 0
	for i := 1; i < len(triangle); i++ {
		if triangle[i][j] < triangle[i][j+1] {
			sum += triangle[i][j]
		} else {
			sum += triangle[i][j+1]
			j = j + 1
		}
	}
	return sum
}
