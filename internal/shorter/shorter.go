package shorter

import (
	"log"
	"net/url"
	"strings"
)

const alphabet = "ynAJfoSgdXHB5VasEMtcbPCr1uNZ4LG723ehWkvwYR6KpxjTm8iQUFqz9D"

var alphabetLen = uint32(len(alphabet))

func Shorten(id uint32) string {
	var (
		digits  []uint32
		num     = id
		builder strings.Builder
	)

	for num > 0 {
		digits = append(digits, num%alphabetLen)
		num /= alphabetLen
	}

	var revers = func(s []uint32) {
		for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
			s[i], s[j] = s[j], s[i]
		}
	}

	revers(digits)

	for _, digit := range digits {
		builder.WriteString(string(alphabet[digit]))
	}

	return builder.String()
}

func PrependBaseURL(baseURL, identifier string) (string, error) {
	log.Println(identifier)
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	log.Println(parsed.Path)
	parsed.Path = identifier
	log.Println(parsed.Host)

	return parsed.String(), nil
}
