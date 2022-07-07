#

## DMap

A generics-based simple, vertically distrubuted map structure.
Supports keys of types that implement a `String()` method (`fmt.Stringer` interface), and any values.

### Benchmarks

```bash
go test -bench=Benchmark -benchmem -benchtime=100x .
goos: linux
goarch: amd64
pkg: github.com/althk/dmap
cpu: Intel(R) Core(TM) i7-1065G7 CPU @ 1.30GHz
BenchmarkSet/100000_keys-8         	     100	    608383 ns/op	  200617 B/op	    5024 allocs/op
BenchmarkSet/1000000_keys-8        	     100	   6633708 ns/op	 1959989 B/op	   50347 allocs/op
BenchmarkGet-8                     	     100	      1492 ns/op	      40 B/op	       2 allocs/op
BenchmarkKeys-8                    	     100	   5445090 ns/op	 8944281 B/op	      56 allocs/op
BenchmarkCount-8                   	     100	       109.7 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/althk/dmap	3.770s


```