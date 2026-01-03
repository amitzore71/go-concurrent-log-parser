package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/amitzore71/go-concurrent-log-parser/parser"
)

func main() {
    filePath := flag.String("file", "", "Path to log file (required)")
    workers := flag.Int("workers", runtime.NumCPU(), "Number of worker goroutines (default: number of CPUs)")
    
    flag.Parse()
    
    if *filePath == "" {
        fmt.Println("Error: -file flag is required")
        flag.Usage()
        os.Exit(1)
    }
    
    fmt.Printf("Starting log parser with %d workers...\n", *workers)
    fmt.Printf("Processing file: %s\n", *filePath)
    
    stats, err := parser.ParseLogFile(*filePath, *workers)
    if err != nil {
        fmt.Printf("Error parsing file: %v\n", err)
        os.Exit(1)
    }
    
    stats.Print()
}
