package main

import "testing"

func Benchmark_isGopher(b *testing.B) {
	tt := []struct {
		name string
		email string
	}{
		{
			name:  "benchmark",
			email: "example@example.com",
		},
		{
			name:  "benchmark",
			email: "example@golang.org",
		},
	}

	for _, tc := range tt {
		tc := tc
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				isGopher(tc.email)
			}
		})
	}
}
