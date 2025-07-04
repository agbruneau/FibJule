package main

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"
)

// ------------------------------------------------------------
// Progress Display Management
// ------------------------------------------------------------

const progressRefreshInterval = 100 * time.Millisecond

// progressData encapsulates progress information for a task.
// This is the canonical definition.
type progressData struct {
	name string  // Name of the task
	pct  float64 // Percentage of progress
}

// progressPrinter manages consolidated progress display for all tasks.
// It refreshes the display at regular intervals or upon receiving new data.
//
// Concept:
// A dedicated goroutine continuously listens on a shared channel (progress).
// It collects percentages from the task and refreshes a single line
// on the terminal to display the overall status. The `\r` (carriage return) trick
// allows rewriting on the same line, creating a smooth progress animation.
func progressPrinter(ctx context.Context, progress <-chan progressData, taskName string) {
	var currentPct float64 = 0.0 // Initialize progress to 0%

	ticker := time.NewTicker(progressRefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case p, ok := <-progress:
			if !ok { // Channel is closed, signifies end of progress updates.
				printStatus(taskName, currentPct) // Print one last time
				fmt.Println()                     // Move to a new line after all progress is done
				return
			}
			// Ensure the name matches, though with one task it always should.
			if p.name == taskName {
				currentPct = p.pct
			}
			printStatus(taskName, currentPct) // Print current status

		case <-ticker.C:
			// Periodically refresh display to show the program is still active,
			// even if no new progress updates have been received.
			printStatus(taskName, currentPct)

		case <-ctx.Done():
			// Main context is done (e.g., timeout or cancellation), stop displaying.
			// Print one last status before exiting, then a newline.
			printStatus(taskName, currentPct)
			fmt.Println()
			return
		}
	}
}

// printStatus displays the current progress status on a single line.
func printStatus(taskName string, pct float64) {
	var b strings.Builder
	b.WriteString("\r") // Carriage return to overwrite the previous line

	// Format string for aligned display: Task Name: XX.YY%
	fmt.Fprintf(&b, "%-15s %6.2f%%", taskName+":", pct)

	// Add trailing spaces to clear any remnants of a longer previous line.
	// Adjust the number of spaces if task names or formatting changes significantly.
	b.WriteString("                    ") // Increased padding
	fmt.Print(b.String())
}

// ------------------------------------------------------------
// *big.Int Object Pool for Memory Reuse
// ------------------------------------------------------------
//
// Memory Optimization Concept (sync.Pool):
// Calculations for large Fibonacci numbers require handling integers
// that exceed the capacity of standard types (e.g., int64). Go's `math/big.Int` is used.
// The problem: Creating numerous `big.Int` objects, especially in loops for complex
// algorithms, puts significant pressure on the Garbage Collector (GC). Frequent GC cycles
// can pause the program and degrade performance.
// The solution: A `sync.Pool` provides a way to reuse objects that are otherwise
// short-lived. Instead of allocating a new `big.Int` each time one is needed,
// the program requests one from the pool. After the object is used, it's returned
// to the pool. This drastically reduces the number of allocations and, consequently,
// the GC overhead, leading to improved performance for memory-intensive operations.

// newIntPool creates a new sync.Pool specifically for *big.Int objects.
// The New function in the pool is called when Get is invoked on an empty pool.
func newIntPool() *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			// Allocate a new *big.Int instance when the pool is empty.
			return new(big.Int)
		},
	}
}
