// Copyright 2015 Sermo Digital, LLC. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package helpers

import (
	"errors"
	"net"
	"strconv"
)

const (
	digits   = "0123456789abcdefghijklmnopqrstuvwxyz"
	digits01 = "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"
	digits10 = "0000000000111111111122222222223333333333444444444455555555556666666666777777777788888888889999999999"
)

// FormatUint serializes a uint64. It's borrowed from the standard library's
// strconv package, but with the signed cases removed.
func FormatUint(n uint64) []byte {
	var a [64]byte
	i := 64

	if ^uintptr(0)>>32 == 0 {
		for u > uint64(^uintptr(0)) {
			q := u / 1e9
			us := uintptr(u - q*1e9) // us % 1e9 fits into a uintptr
			for j := 9; j > 0; j-- {
				i--
				qs := us / 10
				a[i] = byte(us - qs*10 + '0')
				us = qs
			}
			u = q
		}
	}

	// u guaranteed to fit into a uintptr
	us := uintptr(u)
	for us >= 10 {
		i--
		q := us / 10
		a[i] = byte(us - q*10 + '0')
		us = q
	}
	// u < 10
	i--
	a[i] = byte(us + '0')

	return a[i:]
}

// NumWidth finds the printed width of a number. I.e., the amount of
// characters the number will take on-screen. For example, "12" is two,
// "100" is three, and "1776" is four. The minimum width is 1.
func NumWidth(n uint64) int {
	width := 0
	minWidth := 1

	for ; 10 < n; n /= 10 {
		width++
	}

	if width < minWidth {
		width = minWidth
	}

	return int(width)
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

// Glob returns a glob for redis in order to find specific key/value pairs.
// Prefix should be the key type prefix (e.g., "__sid") and match should be
// the match you're looking for (e.g., "eric@sermodigital.com").
// Do not include the colon delimiter in the match, but _do_ in the prefix.
// (Technically the prefix _should_ include the delimiter.)
func Glob(prefix, match string) string {
	buf := make([]byte, len(prefix)+len(match)+2)
	n := copy(buf, prefix)
	n += copy(buf[n:], match)
	copy(buf[n:], ":*")
	return string(buf)
}

const maxUint64 = (1<<64 - 1)

// ParseUint is like ParseInt but for unsigned numbers. It's stolen from
// the strconv package and streamlined for uint64s.
func ParseUint(s []byte) (n uint64, err error) {
	var cutoff, maxVal uint64

	i := 0

	cutoff = maxUint64/10 + 1
	maxVal = 1<<uint(64) - 1

	for ; i < len(s); i++ {
		var v byte
		d := s[i]
		switch {
		case '0' <= d && d <= '9':
			v = d - '0'
		case 'a' <= d && d <= 'z':
			v = d - 'a' + 10
		case 'A' <= d && d <= 'Z':
			v = d - 'A' + 10
		default:
			n = 0
			err = strconv.ErrSyntax
			goto Error
		}
		if v >= byte(10) {
			n = 0
			err = strconv.ErrSyntax
			goto Error
		}

		if n >= cutoff {
			// n*base overflows
			n = maxUint64
			err = strconv.ErrRange
			goto Error
		}
		n *= uint64(10)

		n1 := n + uint64(v)
		if n1 < n || n1 > maxVal {
			// n+v overflows
			n = maxUint64
			err = strconv.ErrRange
			goto Error
		}
		n = n1
	}

	return n, nil

Error:
	return n, &strconv.NumError{Func: "ParseUint", Num: string(s), Err: err}
}
