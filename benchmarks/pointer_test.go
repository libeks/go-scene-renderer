package benchmarks

import (
	"testing"
)

var num = 10000000

type Obj struct {
	Val int
}

func withPointer(x int) *Obj {
	if x%2 == 0 {
		return nil
	}
	return &Obj{Val: x}
}

func withBool(x int) (Obj, bool) {
	if x%2 == 0 {
		return Obj{}, false
	}
	return Obj{Val: x}, true
}

// the idea here is to see whether returning a pointer vs returning a struct with a bool is more efficient.
// conclusion: there is no difference

//	go test ./benchmarks -bench BenchmarkWithPointer
//
// goos: windows
// goarch: amd64
// pkg: github.com/libeks/go-scene-renderer/benchmarks
// cpu: AMD Ryzen 7 7800X3D 8-Core Processor
// BenchmarkWithPointer-16              560           2060716 ns/op
// PASS
// ok      github.com/libeks/go-scene-renderer/benchmarks  1.476s
func BenchmarkWithPointer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sum := 0
		for i := range num {
			res := withPointer(i)
			if res != nil {
				sum += res.Val
			}
		}
	}
}

//	go test ./benchmarks -bench BenchmarkWithBool
//
// goos: windows
// goarch: amd64
// pkg: github.com/libeks/go-scene-renderer/benchmarks
// cpu: AMD Ryzen 7 7800X3D 8-Core Processor
// BenchmarkWithBool-16                 561           2038594 ns/op
// PASS
// ok      github.com/libeks/go-scene-renderer/benchmarks  1.966s
func BenchmarkWithBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sum := 0
		for i := range num {
			o, ok := withBool(i)
			if ok {
				sum += o.Val
			}
		}
	}
}
