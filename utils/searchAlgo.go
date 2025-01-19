package utils

import "strings"

func LevenshteinDistance(s1, s2 string) int {
	// Convert strings to lowercase for case-insensitive comparison
	s1 = strings.ToLower(s1)
	s2 = strings.ToLower(s2)

	// Create matrix of size (m+1)x(n+1)
	m, n := len(s1), len(s2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	// Initialize first row and column
	for i := 0; i <= m; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j
	}

	// Fill the matrix
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s1[i-1] == s2[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = minNum(
					dp[i-1][j]+1,   // deletion
					dp[i][j-1]+1,   // insertion
					dp[i-1][j-1]+1, // substitution
				)
			}
		}
	}
	return dp[m][n]
}

func minNum(numbers ...int) int {
	result := numbers[0]
	for _, num := range numbers[1:] {
		if num < result {
			result = num
		}
	}
	return result
}
