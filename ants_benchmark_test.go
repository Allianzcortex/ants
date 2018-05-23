// MIT License

// Copyright (c) 2018 Andy Pan

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package ants_test

import (
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/panjf2000/ants"
)

const (
	_   = 1 << (10 * iota)
	KiB // 1024
	MiB // 1048576
	GiB // 1073741824
	TiB // 1099511627776             (超过了int32的范围)
	PiB // 1125899906842624
	EiB // 1152921504606846976
	ZiB // 1180591620717411303424    (超过了int64的范围)
	YiB // 1208925819614629174706176
)
const RunTimes = 10000000
const loop = 5

func demoPoolFunc(args interface{}) error {
	// m := args.(int)
	// var n int
	// for i := 0; i < m; i++ {
	// 	n += i
	// }
	// return nil
	n := args.(int)
	time.Sleep(time.Duration(n) * time.Millisecond)
	return nil
}

func BenchmarkGoroutine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for j := 0; j < RunTimes; j++ {
			wg.Add(1)
			go func() {
				demoFunc()
				wg.Done()
			}()
		}
		wg.Wait()
	}
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	b.Logf("total memory usage:%dG", mem.TotalAlloc/GiB)
}

func BenchmarkAntsPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for j := 0; j < RunTimes; j++ {
			wg.Add(1)
			ants.Push(func() {
				demoFunc()
				wg.Done()
			})
		}
		wg.Wait()
	}
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	b.Logf("total memory usage:%dG", mem.TotalAlloc/GiB)
}

func BenchmarkGoroutineWithFunc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		b.ResetTimer()
		for j := 0; j < RunTimes; j++ {
			wg.Add(1)
			go func() {
				demoPoolFunc(loop)
				wg.Done()
			}()
		}
		wg.Wait()
	}
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	b.Logf("total memory usage:%dG", mem.TotalAlloc/GiB)
}

func BenchmarkAntsPoolWithFunc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		p, _ := ants.NewPoolWithFunc(50000, func(i interface{}) error {
			demoPoolFunc(i)
			wg.Done()
			return nil
		})
		b.ResetTimer()
		for j := 0; j < RunTimes; j++ {
			wg.Add(1)
			p.Serve(loop)
		}
		wg.Wait()
	}
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	b.Logf("total memory usage:%dG", mem.TotalAlloc/GiB)
}
