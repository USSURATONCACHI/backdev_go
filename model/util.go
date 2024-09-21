package model

import (
	"unicode"
)

func CapitalizeFirst(s string) string {
    if len(s) == 0 {
        return s
    }
    r := []rune(s) // Convert string to runes to handle multi-byte characters
    r[0] = unicode.ToUpper(r[0])
    return string(r)
}