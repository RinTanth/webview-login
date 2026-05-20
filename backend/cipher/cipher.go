package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/argon2"
)

// argon2id parameters (OWASP recommended minimums)
const (
	argonTime    = 3
	argonMemory  = 64 * 1024 // 64 MB
	argonThreads = 4
	argonKeyLen  = 32
	argonSaltLen = 16
)

// HashPassword hashes a password with argon2id and returns an encoded string
// that contains all parameters needed to verify it later.
// Format: $argon2id$v=19$m=65536,t=3,p=4$<salt_b64>$<hash_b64>
func HashPassword(password string) (string, error) {
	salt := make([]byte, argonSaltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		argonMemory, argonTime, argonThreads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
	return encoded, nil
}

// VerifyPassword checks a plaintext password against an encoded argon2id hash.
func VerifyPassword(password, encoded string) (bool, error) {
	var version int
	var memory, time uint32
	var threads uint8
	var saltB64, hashB64 string

	_, err := fmt.Sscanf(
		encoded,
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s",
		&version, &memory, &time, &threads, &saltB64,
	)
	if err != nil {
		return false, errors.New("invalid hash format")
	}
	// saltB64 will include the trailing $hash — split manually
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}
	saltB64 = parts[4]
	hashB64 = parts[5]

	salt, err := base64.RawStdEncoding.DecodeString(saltB64)
	if err != nil {
		return false, err
	}
	expectedHash, err := base64.RawStdEncoding.DecodeString(hashB64)
	if err != nil {
		return false, err
	}

	keyLen := uint32(len(expectedHash))
	actualHash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)
	return subtle.ConstantTimeCompare(actualHash, expectedHash) == 1, nil
}

// HashEmail produces a fixed-length HMAC-SHA256 hex string.
// Uses a pepper so the hash can't be brute-forced without the key.
func HashEmail(email, pepper string) string {
	mac := hmac.New(sha256.New, []byte(pepper))
	mac.Write([]byte(strings.ToLower(strings.TrimSpace(email))))
	return hex.EncodeToString(mac.Sum(nil))
}

// EncryptEmail encrypts plaintext with AES-256-GCM.
// keyHex must be a 64-char hex string (32 bytes).
// Output is base64(nonce + ciphertext) — nonce is random per call.
func EncryptEmail(email, keyHex string) (string, error) {
	key, err := hex.DecodeString(keyHex)
	if err != nil || len(key) != 32 {
		return "", errors.New("AES_KEY must be a 64-character hex string (32 bytes)")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(email), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptEmail reverses EncryptEmail.
func DecryptEmail(encoded, keyHex string) (string, error) {
	key, err := hex.DecodeString(keyHex)
	if err != nil || len(key) != 32 {
		return "", errors.New("AES_KEY must be a 64-character hex string (32 bytes)")
	}

	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	plaintext, err := gcm.Open(nil, data[:nonceSize], data[nonceSize:], nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
