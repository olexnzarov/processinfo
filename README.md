# processinfo

processinfo is a package that lets you get process stats, such as overall CPU utilization and memory usage.

**Features**

- Ability to get process information every N seconds.
- Easy-to-use interface.
- Supports Linux and Windows. 

## Installation

```sh
go get github.com/olexnzarov/processinfo
```

## Usage

```go
func GetMemoryUsage(pid int) (uint64, error) {
  return processinfo.Get(pid).Memory
}

// This is a preferred way of getting accurate overall CPU usage of a process.
// This function will return an average for the given time duration.
// It will also block the execution for 'sampleTime' duration.
func GetAverageProcessorUsage(pid int, sampleTime time.Duration) float64 {
  ctx, cancel := context.WithCancel(context.Background())
  infoChannel := processinfo.NewSampler(ctx, pid, sampleTime)
  <-infoChannel
  info := <-infoChannel
  cancel()
  return info.CPU
}

// Print process information every 5 seconds until the context ends.
func WatchProcess(ctx context.Context, pid int) {
  infoChannel := processinfo.NewSampler(ctx, pid, time.Second * 5)
  for info := range infoChannel {
    fmt.Printf("Process %d:\n  CPU: %.2f%%\n  Memory: %dbytes\n", pid, info.CPU, info.Memory)
  }
}
```

## Examples

This package was created for use in [gofu](https://github.com/olexnzarov/gofu), a modern process manager. Check it out for examples on how to use this package.

## License

This code is available under the MIT license, allowing for free use, modification, and distribution.