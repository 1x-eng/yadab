package lexer

import (
	"strings"

	. "github.com/1x-eng/yadab/util"
)

// longestMatch iterates through a source string starting at the given
// cursor to find the longest matching substring among the provided options
func longestMatch(source string, ic Cursor, options []string) string {
	var value []byte
	var skipList []int
	var match string

	cur := ic

	for cur.Pointer < uint(len(source)) {

		value = append(value, strings.ToLower(string(source[cur.Pointer]))...)
		cur.Pointer++

	match:
		for i, option := range options {
			for _, skip := range skipList {
				if i == skip {
					continue match
				}
			}

			// Deal with cases like INT vs INTO
			if option == string(value) {
				skipList = append(skipList, i)
				if len(option) > len(match) {
					match = option
				}

				continue
			}

			sharesPrefix := string(value) == option[:cur.Pointer-ic.Pointer]
			tooLong := len(value) > len(option)
			if tooLong || !sharesPrefix {
				skipList = append(skipList, i)
			}
		}

		if len(skipList) == len(options) {
			break
		}
	}

	return match
}
