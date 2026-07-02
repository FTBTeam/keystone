package ftb_go_utils

import (
	"crypto"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
	"os"

	"github.com/aviddiviner/go-murmur"
)

func FileHash(file *os.File, hashType crypto.Hash) (string, error) {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	var h hash.Hash
	switch hashType {
	case crypto.MD5:
		h = md5.New()
	case crypto.SHA1:
		h = sha1.New()
	case crypto.SHA256:
		h = sha256.New()
	case crypto.SHA512:
		h = sha512.New()
	default:
		return "", fmt.Errorf("unsupported hash type: %s", hashType)
	}

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

func isWhitespaceCharacter(b byte) bool {
	return b == 9 || b == 10 || b == 13 || b == 32
}
