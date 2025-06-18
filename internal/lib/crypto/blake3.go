package crypto

import (
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"github.com/zeebo/blake3"
)

//-------------------------------------------------------------------------------------------------

func Blake3(value any) string {
	switch v := value.(type) {
	case string:
		return blake3FromString(v)
	case io.Reader:
		return blake3FromReader(v)
	default:
		panic(fmt.Sprintf("unsupported %T", value))
	}
}

func blake3FromString(value string) string {
	return blake3FromReader(strings.NewReader(value))
}

func blake3FromReader(r io.Reader) string {
	hasher := blake3.New()
	if _, err := io.Copy(hasher, r); err != nil {
		panic(err)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

//-------------------------------------------------------------------------------------------------
