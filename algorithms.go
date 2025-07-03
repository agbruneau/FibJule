package main

import (
	"context"
	"fmt"
	"math/big"
	"math/bits"
	"sync"
)

// fibFunc is a type for functions calculating Fibonacci numbers.
// It takes a context for cancellation, a channel for progress, the index n,
// and a pool of big.Int objects for memory reuse.
type fibFunc func(ctx context.Context, progress chan<- progressData, n int, pool *sync.Pool) (*big.Int, error)

// ------------------------------------------------------------
// Fibonacci Calculation Algorithms
// ------------------------------------------------------------

// fibFastDoubling calculates F(n) using the "Fast Doubling" algorithm.
//
// Concept:
// A very efficient algorithm based on mathematical identities that allow
// transitioning from F(k) and F(k+1) to F(2k) and F(2k+1) in a few operations:
// F(2k)   = F(k) * [2*F(k+1) – F(k)]
// F(2k+1) = F(k)² + F(k+1)²
//
// Implementation:
// The algorithm iterates through the bits of index `n` from left to right (most
// significant to least significant). At each step, it applies the "doubling" formulas.
// If the current bit of `n` is 1, it takes an additional step to advance.
//
// Strengths/Weaknesses:
// Extremely fast and efficient (O(log n) complexity). It's one of the best
// algorithms for this problem. It heavily uses the `sync.Pool` to optimize
// `big.Int` allocations.
func fibFastDoubling(ctx context.Context, progress chan<- progressData, n int, pool *sync.Pool) (*big.Int, error) {
	taskName := "Fast Doubling" // Used for progress reporting
	if n < 0 {
		return nil, fmt.Errorf("negative index n is not supported: %d", n)
	}
	if n <= 1 {
		if progress != nil {
			progress <- progressData{name: taskName, pct: 100.0}
		}
		return big.NewInt(int64(n)), nil
	}

	// Initialize F(k) and F(k+1)
	// a = F(k), b = F(k+1)
	a := pool.Get().(*big.Int).SetInt64(0)
	b := pool.Get().(*big.Int).SetInt64(1)
	defer pool.Put(a) // Ensure 'a' is returned to the pool when done
	defer pool.Put(b) // Ensure 'b' is returned to the pool when done

	// Temporary variables for calculations, taken from the pool.
	t1 := pool.Get().(*big.Int)
	t2 := pool.Get().(*big.Int)
	defer pool.Put(t1)
	defer pool.Put(t2)

	totalBits := bits.Len(uint(n)) // Number of bits in n
	// Iterate from the most significant bit of n down to the least significant bit
	for i := totalBits - 1; i >= 0; i-- {
		// Cooperative context cancellation check
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Doubling Step:
		// F(2k)   = F(k) * [2*F(k+1) – F(k)]
		// F(2k+1) = F(k)² + F(k+1)²
		//
		// Current a = F(k), b = F(k+1)
		// We calculate F(2k) and F(2k+1) and store them in a and b respectively.

		// t1 = 2*F(k+1) - F(k) = 2*b - a
		t1.Lsh(b, 1)  // t1 = 2*b
		t1.Sub(t1, a) // t1 = 2*b - a

		// t2 = F(k)^2 = a^2
		t2.Mul(a, a) // t2 = a*a

		// New a = F(2k) = F(k) * (2*F(k+1) - F(k)) = a * t1
		a.Mul(a, t1) // a = a * t1

		// t1 = F(k+1)^2 = b^2  (reusing t1)
		t1.Mul(b, b) // t1 = b*b

		// New b = F(2k+1) = F(k)^2 + F(k+1)^2 = t2 + t1
		b.Add(t2, t1) // b = t2 + t1 (which is F(k)^2 + F(k+1)^2)

		// If the i-th bit of n is 1, apply the "addition" step:
		// F(m+1) = F(m) + F(m-1)
		// Here, if current a=F(2k), b=F(2k+1), and bit is 1, we need F(2k+1), F(2k+2)
		// New a' = F(2k+1) = b
		// New b' = F(2k+2) = F(2k) + F(2k+1) = a + b (using OLD a and b from before this if block,
		// but since a and b are updated to F(2k) and F(2k+1) respectively in this iteration,
		// it means the new a' = F(2k+1) (which is current b),
		// and new b' = F(2k+2) = F(2k) + F(2k+1) (which is current a + current b).
		if (uint(n)>>i)&1 == 1 {
			// t1 = F(2k) + F(2k+1) (this is the new F(k+1), i.e., F(2k+2))
			t1.Add(a, b) // t1 = current_a (F(2k)) + current_b (F(2k+1))
			// a becomes F(2k+1)
			a.Set(b) // a = current_b (F(2k+1))
			// b becomes F(2k+2)
			b.Set(t1) // b = t1 (F(2k+2))
		}

		if progress != nil {
			progress <- progressData{name: taskName, pct: (float64(totalBits-i) / float64(totalBits)) * 100.0}
		}
	}

	if progress != nil {
		progress <- progressData{name: taskName, pct: 100.0}
	}
	// Return a new instance to avoid returning a pooled object that might be modified.
	return new(big.Int).Set(a), nil
}

// progressData is defined in utils.go
// It encapsulates progress information for a task.
// type progressData struct {
// 	name string
// 	pct  float64
// }
