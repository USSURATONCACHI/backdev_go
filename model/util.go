package model

import (
	"encoding/base64"
	"unicode"

	"github.com/google/uuid"
)

func CapitalizeFirst(s string) string {
    if len(s) == 0 {
        return s
    }
    r := []rune(s) // Convert string to runes to handle multi-byte characters
    r[0] = unicode.ToUpper(r[0])
    return string(r)
}

func ParseUuidFromBase64(base64String string) (uuid.UUID, error) {
    var parsed uuid.UUID

	bytes, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil || len(bytes) != 16 {
		return parsed, err
	}
	
    copy(parsed[:], bytes)

    return parsed, nil
}