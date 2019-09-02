// MIT License

// Copyright (c) 2019 Muhammad Muzzammil

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package jsonc

const (
	quote   = 34 // ("")
	space   = 32 // ( )
	tab     = 9  // (	)
	newLine = 10 // (\n)
	star    = 42 // (*)
	slash   = 47 // (/)
)

type state int

const (
	stopped state = iota
	canStart
	started
	canStop
)

type comment struct {
	state       state
	isJSON      bool
	isMultiLine bool
}

func (cmt *comment) setState(s state) {
	cmt.state = s
}

func (cmt *comment) checkState(s state) bool {
	return cmt.state == s
}

func translate(s []byte) []byte {
	vj := make([]byte, len(s))
	i := 0
	cmt := &comment{
		state: stopped,
	}
	for _, s := range s {
		if cmt.checkState(stopped) {
			if s == quote {
				cmt.isJSON = !cmt.isJSON
			}

			if cmt.isJSON {
				vj[i] = s
				i++
				continue
			} else {
				if s == space || s == tab || s == newLine {
					continue
				}
			}
		}

		if cmt.checkState(stopped) && s == slash {
			cmt.setState(canStart)
			continue
		}

		if cmt.checkState(canStart) && (s == slash || s == star) {
			cmt.setState(started)
			if s == star {
				cmt.isMultiLine = true
			}
			continue
		}

		if cmt.checkState(started) {
			if s == star {
				cmt.setState(canStop)
				continue
			}

			if s == newLine {
				if cmt.isMultiLine {
					continue
				}
				cmt.setState(stopped)
			}

			continue
		}

		if cmt.checkState(canStop) && (s == slash) {
			cmt.setState(stopped)
			cmt.isMultiLine = false
			continue
		}

		vj[i] = s
		i++
	}

	return vj[:i]
}
