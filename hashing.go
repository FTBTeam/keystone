package keystone

import (
	"crypto"
	"fmt"
	"io"
	"os"

	"github.com/aviddiviner/go-murmur"
)

// CryptoFromString Helper function to convert a string representation of a hash type to the corresponding crypto.Hash value.
// Deprecated: This function is deprecated and may be removed in future versions. Use the crypto package directly for hash type
func CryptoFromString(hashType string) (crypto.Hash, error) {
	switch hashType {

	case "md5":
		return crypto.MD5, nil
	case "sha1":
		return crypto.SHA1, nil
	case "sha256":
		return crypto.SHA256, nil
	case "sha512":
		return crypto.SHA512, nil

	default:
		return crypto.SHA1, fmt.Errorf("unknown hash type: %s", hashType)
	}
}

func FileHash[T string | *os.File](inputType T, hashType crypto.Hash) (string, error) {
	var file *os.File
	var err error

	switch input := any(inputType).(type) {
	case string:
		file, err = os.Open(input)
		if err != nil {
			return "", err
		}
		defer file.Close()
	case *os.File:
		file = input
	default:
		return "", fmt.Errorf("unsupported type: %T", inputType)
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	h := hashType.New()

	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func CFMurmurHash(file *os.File) (uint32, error) {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return 0, err
	}

	info, err := file.Stat()
	if err != nil {
		return 0, err
	}
	data := make([]byte, info.Size())
	_, err = io.ReadFull(file, data)

	if err != nil {
		return 0, fmt.Errorf("reading file %s: %w", file.Name(), err)
	}

	filtered := data[:0]
	for _, b := range data {
		if !isWhitespaceCharacter(b) {
			filtered = append(filtered, b)
		}
	}

	murmurHash := murmur.MurmurHash2(filtered, 1)
	return murmurHash, nil
}

// helper bits

func isWhitespaceCharacter(b byte) bool {
	return b == 9 || b == 10 || b == 13 || b == 32
}
