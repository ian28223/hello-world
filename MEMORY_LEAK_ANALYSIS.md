# Go pprof Memory Leak Analysis Guide

## Overview
This guide helps you identify and analyze memory leaks in Go applications using pprof.

## Quick Start

### 1. Capture Memory Profiles

If you have a running Go application with pprof enabled:

```bash
# Capture heap profile
curl http://localhost:6060/debug/pprof/heap > heap.pprof

# Capture multiple profiles over time to compare
curl http://localhost:6060/debug/pprof/heap > heap1.pprof
# Wait some time...
curl http://localhost:6060/debug/pprof/heap > heap2.pprof
```

### 2. Analyze Profiles

```bash
# Interactive analysis
go tool pprof heap.pprof

# Compare two profiles to see growth
go tool pprof -base heap1.pprof heap2.pprof

# Generate visual reports
go tool pprof -http=:8080 heap.pprof
```

## Common Memory Leak Patterns

### 1. **Goroutine Leaks**
- Goroutines that never exit
- Accumulating goroutines over time
- Check with: `curl http://localhost:6060/debug/pprof/goroutine`

### 2. **Unclosed Resources**
- HTTP response bodies not closed
- Database connections not returned to pool
- File handles not closed

### 3. **Growing Slices/Maps**
- Unbounded caches
- Slices that grow but never shrink
- Maps without cleanup logic

### 4. **Reference Retention**
- Objects held in closures
- Event listeners not removed
- Circular references

## Analysis Commands

### Top Memory Consumers
```bash
# In pprof interactive mode
(pprof) top
(pprof) top -cum
```

### List Function Code
```bash
(pprof) list <function_name>
```

### Show Call Graph
```bash
(pprof) web
(pprof) png > graph.png
```

### Compare Profiles
```bash
go tool pprof -base baseline.pprof current.pprof
(pprof) top
```

## Key Metrics to Monitor

1. **inuse_space** - Currently allocated memory
2. **inuse_objects** - Number of allocated objects
3. **alloc_space** - Total allocated memory (cumulative)
4. **alloc_objects** - Total allocated objects (cumulative)

## Analysis Workflow

1. **Establish Baseline**: Capture initial heap profile
2. **Stress Test**: Run application under load
3. **Capture Again**: Take multiple profiles over time
4. **Compare**: Look for growing allocations
5. **Investigate**: Focus on functions with most growth
6. **Fix**: Address the root cause
7. **Verify**: Confirm fix with new profiles

## Tools in This Repository

- `sample_leak.go` - Example program with memory leak
- `analyze_pprof.sh` - Script to analyze pprof files
- `capture_profiles.sh` - Script to capture multiple profiles over time
