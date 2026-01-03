package parser

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func ParseLogFile(filePath string, numWorkers int) (*LogStats, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()
    
    stats := NewLogStats()
    pool := NewWorkerPool(numWorkers, stats)
    
    pool.Start()
    
    startTime := time.Now()
    
    if err := streamFile(file, pool); err != nil {
        return nil, err
    }
    
    pool.Close()
    
    pool.Wait()
    
    stats.ProcessingTime = time.Since(startTime).Seconds()
    
    return stats, nil
}
func streamFile(file *os.File, pool *WorkerPool) error {
    scanner := bufio.NewScanner(file)
    
    const maxCapacity = 1024 * 1024 // 1MB
    buf := make([]byte, maxCapacity)
    scanner.Buffer(buf, maxCapacity)
    
    lineNumber := 0
    
    for scanner.Scan() {
        lineNumber++
        
        pool.Jobs <- LogLine{
            LineNumber: lineNumber,
            Content:    scanner.Text(),
        }
    }
    
    if err := scanner.Err(); err != nil {
        return fmt.Errorf("error reading file: %w", err)
    }
    
    return nil
}
