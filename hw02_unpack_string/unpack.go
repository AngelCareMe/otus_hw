package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	runes := []rune(s)
	var sb strings.Builder

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if unicode.IsDigit(r) {
			if i == 0 || unicode.IsDigit(runes[i-1]) {
				return "", ErrInvalidString
			}
			count, _ := strconv.Atoi(string(r))
			if count == 0 {
				tmp := []rune(sb.String())
				sb.Reset()
				sb.WriteString(string(tmp[:len(tmp)-1]))
			} else {
				sb.WriteString(strings.Repeat(string(runes[i-1]), count-1))
			}
		} else {
			sb.WriteRune(r)
		}
	}

	return sb.String(), nil
}
