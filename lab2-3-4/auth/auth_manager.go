package auth

import "sync"

type AuthenticationManager struct {
	mu       sync.RWMutex
	counters map[string]uint64
}

func NewAuthenticationManager() *AuthenticationManager {
	return &AuthenticationManager{
		counters: make(map[string]uint64),
	}
}

func (am *AuthenticationManager) GenerateAndIncrement(
	email string,
	key []byte,
	digits int,
) (string, error) {
	am.mu.Lock()
	defer am.mu.Unlock()

	counter, exists := am.counters[email]
	if !exists {
		counter = 0
	}

	otp, err := GenerateHOTP(key, counter, digits)
	if err != nil {
		return "", err
	}

	am.counters[email] = counter + 1
	return otp, nil
}

func (am *AuthenticationManager) VerifyAndIncrementCounter(
	email string,
	key []byte,
	providedOTP string,
	digits int,
	lookAhead int,
) (bool, uint64) {
	am.mu.Lock()
	defer am.mu.Unlock()

	counter, exists := am.counters[email]
	if !exists {
		counter = 0
	}

	verified, newCounter := VerifyHOTP(key, counter, providedOTP, digits, lookAhead)
	if verified {
		am.counters[email] = counter
		return true, newCounter
	}

	return false, counter
}

func (am *AuthenticationManager) SyncCounterWithDatabase(email string, dbCounter int64) {
	am.mu.Lock()
	defer am.mu.Unlock()

	am.counters[email] = uint64(dbCounter)
}
