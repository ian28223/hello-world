package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

// MEMORY LEAK EXAMPLE 1: Goroutine leak
// These goroutines never exit and accumulate over time
func goroutineLeakExample() {
	for i := 0; i < 10; i++ {
		go func() {
			// This goroutine blocks forever - LEAK!
			select {} // blocks forever
		}()
	}
}

// MEMORY LEAK EXAMPLE 2: Growing slice without bounds
var globalCache [][]byte

func sliceLeakExample() {
	// Continuously append to slice without cleanup - LEAK!
	data := make([]byte, 1024*1024) // 1MB
	globalCache = append(globalCache, data)
}

// MEMORY LEAK EXAMPLE 3: Map without cleanup
var userSessions = make(map[string]*UserSession)

type UserSession struct {
	UserID    string
	Data      []byte // Large data
	CreatedAt time.Time
}

func mapLeakExample(userID string) {
	// Add sessions but never remove old ones - LEAK!
	session := &UserSession{
		UserID:    userID,
		Data:      make([]byte, 1024*100), // 100KB per session
		CreatedAt: time.Now(),
	}
	userSessions[userID] = session
}

// MEMORY LEAK EXAMPLE 4: Substring/Subslice retention
func substringLeakExample() []string {
	results := make([]string, 0)
	
	for i := 0; i < 100; i++ {
		// Create large string
		largeString := string(make([]byte, 1024*1024)) // 1MB
		
		// Taking substring keeps entire backing array in memory - LEAK!
		small := largeString[:10]
		results = append(results, small)
		
		// FIX: Use strings.Clone() or copy to new string
		// small := strings.Clone(largeString[:10])
	}
	
	return results
}

// HEALTHY EXAMPLE: Proper resource management
func healthyExample() {
	resp, err := http.Get("http://example.com")
	if err != nil {
		return
	}
	defer resp.Body.Close() // Properly close resources
	
	// Process response...
}

func main() {
	// Start pprof server
	go func() {
		log.Println("Starting pprof server on :6060")
		log.Println("Access at http://localhost:6060/debug/pprof/")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	fmt.Println("Memory Leak Example Server Starting...")
	fmt.Println("This program intentionally leaks memory for demonstration")
	fmt.Println("\nAccess pprof at: http://localhost:6060/debug/pprof/")
	fmt.Println("\nCapture heap profile with:")
	fmt.Println("  curl http://localhost:6060/debug/pprof/heap > heap.pprof")
	fmt.Println("\nAnalyze with:")
	fmt.Println("  go tool pprof heap.pprof")
	fmt.Println("\nPress Ctrl+C to stop")
	fmt.Println("\n" + "=".repeat(60))

	ticker := time.NewTicker(1 * time.Second)
	counter := 0

	for range ticker.C {
		counter++
		
		// Trigger different leak examples
		goroutineLeakExample()
		sliceLeakExample()
		mapLeakExample(fmt.Sprintf("user-%d", counter))
		
		if counter%10 == 0 {
			substringLeakExample()
		}

		// Print stats
		if counter%5 == 0 {
			fmt.Printf("[%d] Goroutine leaks: %d, Slice items: %d, Map sessions: %d\n",
				counter,
				counter*10, // 10 goroutines per iteration
				len(globalCache),
				len(userSessions))
		}
	}
}
