// Copyright 2015 Sermo Digital, LLC. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package helpers

import (
	"errors"
	"net"
)

// Length finds the number of digits in a uint64. For example, 12 returns 2,
// 100 returns 3, and 1776 returns 4. The minimum width is 1.
func Length(x uint64) int {
	// TODO: use math/bits when it's merged in 1.9

	// Loop: for loop
	// Log10: math.Log10
	// Asm: https://graphics.stanford.edu/~seander/bithacks.html#IntegerLog10
	// 	with "IntegerLogBase2(v)" in assembly implemented similar to GCC's
	// 	__builtin_clzll.
	// Cond: the switch below
	//
	// 	Seq: sequential numbers
	// 	Rand: random numbers
	//
	// BenchmarkIntLenLoopSeq-4     	200000000         7.87 ns/op
	// BenchmarkIntLenLog10Seq-4    	30000000        43.8 ns/op
	// BenchmarkIntLenAsmSeq-4      	200000000         7.38 ns/op
	// BenchmarkIntLenCondSeq-4     	500000000         3.67 ns/op
	//
	// BenchmarkIntLenLoopRand-4    	10000000000000000000        21.8 ns/op
	// BenchmarkIntLenLog10Rand-4   	30000000        44.5 ns/op
	// BenchmarkIntLenAsmRand-4     	200000000         8.62 ns/op
	// BenchmarkIntLenCondRand-4    	200000000         8.09 ns/op

	switch {
	case x < 10:
		return 1
	case x < 100:
		return 2
	case x < 1000:
		return 3
	case x < 10000:
		return 4
	case x < 100000:
		return 5
	case x < 1000000:
		return 6
	case x < 10000000:
		return 7
	case x < 100000000:
		return 8
	case x < 1000000000:
		return 9
	case x < 10000000000:
		return 10
	case x < 100000000000:
		return 11
	case x < 1000000000000:
		return 12
	case x < 10000000000000:
		return 13
	case x < 100000000000000:
		return 14
	case x < 1000000000000000:
		return 15
	case x < 10000000000000000:
		return 16
	case x < 100000000000000000:
		return 17
	case x < 1000000000000000000:
		return 18
	case x < 10000000000000000000:
		return 19
	default:
		panic("unreachable")
	}
}

// ParseIP returns a valid IP address for the given *http.Request.RemoteAddr
// if available.
func ParseIP(remoteaddr string) (string, error) {
	host, _, err := net.SplitHostPort(remoteaddr)
	if err != nil {
		return "", err
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return "", errors.New("Invalid IP Address")
	}
	return ip.String(), nil
}
