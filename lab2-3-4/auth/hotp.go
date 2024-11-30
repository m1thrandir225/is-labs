package auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
)

func GenerateHOTP(
	key []byte,
	counter uint64,
	digits int,
) (string, error) {
	if digits < 6 || digits > 8 {
		return "", errors.New("digits must be between 6 and 8")
	}
	counterBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(counterBytes, counter)

	mac := hmac.New(sha1.New, key)
	mac.Write(counterBytes)
	hmacHash := mac.Sum(nil)

	offset := hmacHash[len(hmacHash)-1] & 0xf
	truncHash := binary.BigEndian.Uint32(hmacHash[offset : offset+4])

	truncHash &= 0x7FFFFFFF

	otp := truncHash % uint32(pow(10, digits))

	return fmt.Sprintf(fmt.Sprintf("%%0%dd", digits), otp), nil
}

func VerifyHOTP(key []byte, counter uint64, providedOTP string, digits int, lookAhead int) (bool, uint64) {
	for i := 0; i <= lookAhead; i++ {
		generatedOTP, err := GenerateHOTP(key, counter, digits)
		if err != nil {
			return false, counter
		}

		if generatedOTP == providedOTP {
			return true, counter + uint64(i)
		}
	}

	return false, counter
}

func pow(base, exp int) int {
	res := 1
	for exp > 0 {
		if exp&1 == 1 {
			res *= base
		}
		base *= base
		exp >>= 1
	}
	return res
}

func GenerateHOTPKey() []byte {
	key := make([]byte, 20)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	return key
}
