package service

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

// SHA1HashUrl hash given url with the secret and return the result. '?' char
// in url would be cleaned if any.
func SHA1HashUrl(sc, in string) string {
	in = strings.ReplaceAll(in, "?", "")
	in += sc

	out := sha1.Sum([]byte(in))
	return hex.EncodeToString(out[:])
}
