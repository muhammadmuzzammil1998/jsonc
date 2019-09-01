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

func translate(s []byte) []byte {
	var (
		line  int
		j     []byte
		quote bool
	)
	comment := &commentData{}
	for _, ch := range []byte(s) {
		if ch == 34 { // 32 = quote (")
			quote = !quote
		}
		if (ch == 32 || ch == 9) && !quote { // 32 = space ( ), 9 = tab (	)
			continue
		}
		if ch == 10 { // 10 = new line
			line++
			if comment.isSingleLined {
				comment.stop()
			}
			continue
		}
		token := string(ch)
		if comment.startted {
			if token == "*" {
				comment.canEnd = true
				continue
			}
			if comment.canEnd && token == "/" {
				comment.stop()
				continue
			}
			continue
		}
		if comment.canStart && (token == "*" || token == "/") {
			comment.start(token)
			continue
		}
		if token == "/" {
			comment.canStart = true
			continue
		}
		j = append(j, ch)
	}
	return j
}

type commentData struct {
	canStart      bool
	canEnd        bool
	startted      bool
	isSingleLined bool
	endLine       int
}

func (c *commentData) stop() {
	c.startted = false
	c.canStart = false
}

func (c *commentData) start(token string) {
	c.startted = true
	c.isSingleLined = token == "/"
}
