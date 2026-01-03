package parser

import (
	"fmt"
	"sync"
)


type LogStats struct {
	mu             sync.Mutex
	ErrorCounts    map[string]int
	TotalLines     int64
	TotalErrors    int64
	ProcessingTime float64
}


func NewLogStats() *LogStats {
    return &LogStats{
        ErrorCounts: make(map[string]int),
    }
}

func (s *LogStats) IncrementLines() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.TotalLines++
}

func (s *LogStats) IncrementError(errorType string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TotalErrors++
	s.ErrorCounts[errorType]++
}

func (s *LogStats) Print() {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    fmt.Println("\n=== Log Processing Statistics ===")
    fmt.Printf("Total Lines Processed: %d\n", s.TotalLines)
    fmt.Printf("Total Errors Found: %d\n", s.TotalErrors)
    fmt.Printf("Processing Time: %.2f seconds\n", s.ProcessingTime)
    
    if len(s.ErrorCounts) > 0 {
        fmt.Println("\nError Breakdown:")
        for errorType, count := range s.ErrorCounts {
            fmt.Printf("  %s: %d\n", errorType, count)
        }
    }
}
