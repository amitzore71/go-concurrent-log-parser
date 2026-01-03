package parser

import (
	"regexp"
	"strings"
	"sync"
)

type LogLine struct {
    LineNumber int
    Content    string
}


type WorkerPool struct {
    NumWorkers int
    Jobs       chan LogLine
    Stats      *LogStats
    wg         sync.WaitGroup
    errorRegex *regexp.Regexp
}


func NewWorkerPool(numWorkers int, stats *LogStats) *WorkerPool {
    return &WorkerPool{
        NumWorkers: numWorkers,
        Jobs:       make(chan LogLine, numWorkers*100),
        Stats:      stats,
        errorRegex: regexp.MustCompile(`(?i)(error|ERROR|fatal|FATAL|exception|Exception|panic|PANIC|failed|FAILED)`),
    }
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.NumWorkers; i++ {
        wp.wg.Add(1)
        go wp.worker(i)
    }
}


func (wp *WorkerPool) worker(id int) {
    defer wp.wg.Done()
    for job := range wp.Jobs {
        wp.processLine(job)
    }
}


func (wp *WorkerPool) processLine(line LogLine) {
    wp.Stats.IncrementLines()
    
    if wp.errorRegex.MatchString(line.Content) {
        errorType := wp.extractErrorType(line.Content)
        wp.Stats.IncrementError(errorType)
    }
}
func (wp *WorkerPool) extractErrorType(line string) string {
    matches := wp.errorRegex.FindStringSubmatch(line)
    if len(matches) > 0 {
        return strings.ToUpper(matches[0])
    }
    return "UNKNOWN_ERROR"
}


func (wp *WorkerPool) Wait() {
    wp.wg.Wait()
}


func (wp *WorkerPool) Close() {
    close(wp.Jobs)
}
