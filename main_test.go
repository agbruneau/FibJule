// main_test.go

package main

import (
	"context"
	"math/big"
	"testing"
)

// TestFibFastDoublingAlgorithm verifies the correctness of the Fast Doubling algorithm
// using a table-driven approach.
func TestFibFastDoublingAlgorithm(t *testing.T) {
	// Test cases with well-known Fibonacci values.
	testCases := []struct {
		name    string
		n       int
		want    *big.Int
		wantErr bool // If an error is expected (e.g., for n < 0)
	}{
		{"n=0", 0, big.NewInt(0), false},
		{"n=1", 1, big.NewInt(1), false},
		{"n=2", 2, big.NewInt(1), false},
		{"n=7", 7, big.NewInt(13), false},
		{"n=10", 10, big.NewInt(55), false},
		{"n=20", 20, big.NewInt(6765), false},
		{"negative n", -1, nil, true}, // Test case for negative input
	}

	pool := newIntPool()
	ctx := context.Background() // Use a background context for tests
	algoName := "Fast Doubling"
	algoFunc := fibFastDoubling

	// Iterate over each test case.
	for _, tc := range testCases {
		// t.Run creates sub-tests, making debugging easier.
		t.Run(algoName+"/"+tc.name, func(t *testing.T) {
			// Execute the algorithm function.
			// The progress channel is not needed for correctness testing.
			got, err := algoFunc(ctx, nil, tc.n, pool)

			// Check if an error was expected.
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected an error for n=%d, but got none", tc.n)
				}
				return // Test is done if an error was expected and occurred.
			}

			// Check if an unexpected error occurred.
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Compare the obtained result with the expected result.
			if got == nil && tc.want == nil {
				// This case should ideally be covered by wantErr if nil result means error
			} else if got == nil && tc.want != nil {
				t.Errorf("for F(%d), expected %s, but got nil", tc.n, tc.want.String())
			} else if got != nil && tc.want == nil {
				t.Errorf("for F(%d), expected nil, but got %s", tc.n, got.String())
			} else if got.Cmp(tc.want) != 0 {
				t.Errorf("for F(%d), expected %s, but got %s", tc.n, tc.want.String(), got.String())
			}
		})
	}
}

// TestFibonacciConsistencyForLargeN is removed as there are no other algorithms to compare against.
// If needed, specific large value tests for Fast Doubling can be added to TestFibFastDoublingAlgorithm.
// The helper function min(a,b) was part of TestFibonacciConsistencyForLargeN and is now removed.

// ------------------------------------------------------------
// Benchmarks
// ------------------------------------------------------------

// Common n for all benchmarks for fair comparison.
const benchmarkN = 100000

// BenchmarkFibFastDoubling measures the performance of the Fast Doubling algorithm.
func BenchmarkFibFastDoubling(b *testing.B) {
	pool := newIntPool()
	ctx := context.Background()
	b.ReportAllocs() // Display memory allocations.
	b.ResetTimer()   // Reset timer to exclude setup time.

	for i := 0; i < b.N; i++ {
		// The result is not verified here; focus is on performance.
		_, _ = fibFastDoubling(ctx, nil, benchmarkN, pool)
	}
}

// Other benchmarks (BenchmarkFibMatrix, BenchmarkFibBinet, BenchmarkFibIterative) are removed.
