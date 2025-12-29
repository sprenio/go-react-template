package uuidstr

import (
	"github.com/google/uuid"
	"math/big"
	"regexp"
)

const (
	uniqueIdMinLength     = 6
	uniqueIdMaxConvertStr = 12
)

func getUniqHex() string {
	// Generowanie UUID i konwersja do hex
	uuidStr := uuid.New().String()
	re := regexp.MustCompile("[^0-9a-f]+")
	return re.ReplaceAllString(uuidStr, "")
}

func GetUniqBase36(size int) string {
	if size < uniqueIdMinLength {
		size = uniqueIdMinLength
	}
	chunkLength := size
	if chunkLength > uniqueIdMaxConvertStr {
		chunkLength = uniqueIdMaxConvertStr
	}
	hexStr := getUniqHex()
	result := ""
	for i := 0; len(result) < size; i += chunkLength {
		end := i + chunkLength
		for end > len(hexStr) {
			hexStr += getUniqHex()
		}
		chunk := hexStr[i:end]

		// Konwersja z base16 do *big.Int
		num := new(big.Int)
		num.SetString(chunk, 16)

		// Konwersja z base10 do base36
		result += num.Text(36)
	}
	return result[:size]
}
