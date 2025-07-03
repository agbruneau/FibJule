// main.go
//
// This program calculates the n-th Fibonacci number using distinct algorithms:
// 1. Fast Doubling algorithm.
//
// It executes this algorithm, displays its real-time progress,
// and its execution time and result.
// A sync.Pool is used to reduce memory allocations for big.Int objects.
//
// Usage:
//   go run . -n <index> -timeout <duration>
// Example:
//   go run . -n 100000 -timeout 1m

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"
)

// ------------------------------------------------------------
// Types and Structures
// ------------------------------------------------------------

// task represents a Fibonacci calculation task to be executed.
type task struct {
	name string  // Name of the algorithm
	fn   fibFunc // Algorithm function
}

// result stores the outcome of a calculation task.
type result struct {
	name     string        // Name of the algorithm
	value    *big.Int      // Calculated Fibonacci value
	duration time.Duration // Duration of the calculation
	err      error         // Potential error
}

// ------------------------------------------------------------
// Main Function: The Orchestrator
// ------------------------------------------------------------
//
// The `main` function orchestrates the entire process:
// 1. It reads command-line parameters (`-n`, `-timeout`).
// 2. It defines the task to execute (Fast Doubling).
//  3. It creates a `context` with a global timeout to ensure the program
//     doesn't run indefinitely. This context is passed to the calculation goroutine
//     to allow for cooperative cancellation.
//  4. It launches the `progressPrinter` goroutine for real-time display.
//  5. It launches a goroutine for each calculation task. Using goroutines
//     allows all selected algorithms to run concurrently.
//  6. It waits for all tasks to complete using a `sync.WaitGroup`.
//  7. It closes communication channels to signal recipient goroutines
//     (like `progressPrinter`) that there will be no more data.
//  8. Finally, it calls `collectAndDisplayResults` to analyze and present the results.
func main() {
	// 1. Read command-line parameters
	nFlag := flag.Int("n", 100000, "Index n of the Fibonacci term (non-negative integer)")
	timeoutFlag := flag.Duration("timeout", 1*time.Minute, "Global maximum execution time")
	flag.Parse()

	n := *nFlag
	timeout := *timeoutFlag

	if n < 0 {
		log.Fatalf("Index n must be greater than or equal to 0. Received: %d", n)
	}

	// 2. Define the task to run
	taskToRun := task{
		name: "Fast Doubling",
		fn:   fibFastDoubling,
	}
	selectedTaskNames := []string{taskToRun.name} // For progress printer

	log.Printf("Calculating F(%d) using %s with a timeout of %v...", n, taskToRun.name, timeout)

	// 3. Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // Important to release resources associated with the context

	intPool := newIntPool()

	// Channels for communication between goroutines
	progressAggregatorCh := make(chan progressData, 2) // Buffer for progress data
	resultsCh := make(chan result, 1)                  // Buffer for the single result

	// 4. Launch progress display
	var wgDisplay sync.WaitGroup
	wgDisplay.Add(1)
	go func() {
		defer wgDisplay.Done()
		progressPrinter(ctx, progressAggregatorCh, selectedTaskNames)
	}()

	// 5. Launch calculation
	var wg sync.WaitGroup
	wg.Add(1)
	log.Println("Launching calculation...")
	go func(currentTask task) {
		defer wg.Done()
		start := time.Now()
		v, err := currentTask.fn(ctx, progressAggregatorCh, n, intPool)
		duration := time.Since(start)
		resultsCh <- result{currentTask.name, v, duration, err}
	}(taskToRun)

	// 6. Wait for the calculation to finish
	wg.Wait()
	log.Println("Calculation finished.")

	// 7. Close channels to signal end of transmissions
	close(progressAggregatorCh)
	close(resultsCh)

	// Wait for the display goroutine to finish
	wgDisplay.Wait()

	// 8. Collect and display results
	collectAndDisplayResults(ctx, resultsCh, n)

	log.Println("Program finished.")
}

// collectAndDisplayResults retrieves, sorts, and displays calculation results.
//
// This function is responsible for the final presentation:
//  1. It collects all results from the `resultsCh` channel until it's closed.
//  2. It displays a clear summary.
//  3. It displays details about the calculated number.
func collectAndDisplayResults(ctx context.Context, resultsCh <-chan result, n int) {
	// Since there's only one result, we read it directly.
	r := <-resultsCh // This will block until the result is sent.

	fmt.Println("\n--------------------------- RESULT ---------------------------")

	if r.err != nil {
		// Distinguish a timeout from other errors for a clearer message.
		if err := ctx.Err(); err == context.DeadlineExceeded && r.err == context.DeadlineExceeded {
			log.Printf("⚠️ Task '%s' was interrupted by the global timeout after %v", r.name, r.duration.Round(time.Microsecond))
		} else if r.err == context.DeadlineExceeded {
			log.Printf("⚠️ Task '%s' self-terminated due to context cancellation (possibly timeout) after %v", r.name, r.duration.Round(time.Microsecond))
		} else {
			log.Printf("❌ Error for task '%s': %v (duration: %v)", r.name, r.err, r.duration.Round(time.Microsecond))
		}
		fmt.Println("------------------------------------------------------------------------")
		fmt.Println("\nThe calculation could not complete successfully.")
		return
	}

	// Display the result
	status := "OK"
	valStr := "N/A"
	if r.value != nil {
		if len(r.value.String()) > 15 {
			valStr = r.value.String()[:5] + "..." + r.value.String()[len(r.value.String())-5:]
		} else {
			valStr = r.value.String()
		}
	}
	fmt.Printf("%-16s : %-12v [%-14s] Result: %s\n", r.name, r.duration.Round(time.Microsecond), status, valStr)
	fmt.Println("------------------------------------------------------------------------")

	if r.value != nil {
		fmt.Printf("\n📊 Algorithm: %s (%v)\n", r.name, r.duration.Round(time.Microsecond))
		printFibResultDetails(r.value, n)
	} else {
		// This case should ideally be covered by r.err != nil
		fmt.Println("\nNo result value was produced, despite no explicit error.")
	}
}

// printFibResultDetails displays detailed information about the calculated Fibonacci number.
// This function remains unchanged as its logic is independent of the number of algorithms.
func printFibResultDetails(value *big.Int, n int) {
	if value == nil {
		return
	}

	digits := len(value.Text(10))
	fmt.Printf("Number of digits in F(%d): %d\n", n, digits)

	// Use scientific notation for numbers too large to display.
	if digits > 20 {
		floatVal := new(big.Float).SetPrec(uint(digits + 10)).SetInt(value)
		sci := floatVal.Text('e', 8) // 8 digits of precision for scientific notation
		fmt.Printf("Value (scientific notation) ≈ %s\n", sci)
	} else {
		fmt.Printf("Value = %s\n", value.Text(10))
	}
}
