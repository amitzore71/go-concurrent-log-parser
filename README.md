# Go Concurrent Log Parser

A high-performance log analysis tool that can handle massive log files (15GB+) without running out of memory.

## Why I Built This

I'm a data analyst/engineer who primarily works in Python, but I hit a wall when trying to analyze a 15GB+ log file. Python's GIL prevented me from fully utilizing multiple CPU cores, and I kept running into out-of-memory errors even with chunked reading.

### Why Not PySpark or Polars?

Obviously, I could have used PySpark or Polars for this - both are excellent tools for large-scale data processing. But:

- **PySpark** felt like overkill for what's essentially a single-node task. Setting up Spark for a simple log parsing job seemed like bringing a tank to a knife fight. The overhead and complexity weren't justified for this use case.

- **Polars** is great, but I wanted something even more lightweight with minimal dependencies. Also, this was a good excuse to finally try Go .

## What It Does

- Processes massive log files (tested with 15GB+) without memory issues
- Uses all available CPU cores through Go's goroutines
- Finds and categorizes errors (`ERROR`, `FATAL`, `PANIC`, `EXCEPTION`, `FAILED`)
- Streams the file line-by-line instead of loading everything into memory
- Single binary, no runtime dependencies - perfect for running on remote servers

## Installation

I'm using Go 1.25.0.

```bash
git clone https://github.com/amitzore71/go-concurrent-log-parser.git
cd go-concurrent-log-parser
go build -o logparser main.go
```

## Usage

Basic usage:

```bash
./logparser -file path/to/your.log
```

Control the number of workers (defaults to your CPU count):

```bash
./logparser -file large_server.log -workers 8
```

## Example Output

```text
Starting log parser with 16 workers...
Processing file: test_large.log

=== Log Processing Statistics ===
Total Lines Processed: 10000000
Total Errors Found: 1500525
Processing Time: 14.76 seconds

Error Breakdown:
  ERROR: 562881
  EXCEPTION: 188133
  FAILED: 187758
  PANIC: 187883
  FATAL: 373870
```

## How It Works

The tool uses a worker pool pattern:

1. One goroutine reads the file line-by-line using `bufio.Scanner`
2. Lines are sent through a channel to a pool of worker goroutines
3. Each worker scans its lines for error patterns using regex
4. Results are collected in a thread-safe stats counter
5. Everything is done in a streaming fashion - no full file load into memory
