package auth

import (
	"context"
	"errors"
	"log"
	db "m1thrandir225/lab-2-3-4/db/sqlc"
)

type OTPService struct {
	store db.Store
}

func NewOTPService(store db.Store) *OTPService {
	return &OTPService{
		store: store,
	}
}

func (s *OTPService) GenerateAndVerifyOTP(
	ctx context.Context,
	userID int64,
	userProvidedOTP string,
	otpSecret string,
) (bool, error) {
	// Try to get current counter
	currentCounterInt, err := s.store.GetCurrentCounter(ctx, userID)
	if err != nil {
		log.Printf("Error getting current counter for userID %v: %v", userID, err)
		// If no counter exists, initialize it
		arg := db.CreateHotpCounterParams{
			UserID:  userID,
			Counter: 0,
		}

		err = s.store.CreateHotpCounter(ctx, arg)
		if err != nil {
			log.Printf("Error creating hotp counter: %v", err)
			return false, err
		}
		currentCounterInt = 0
	}

	currentCounter := uint64(currentCounterInt)
	log.Println(currentCounter)

	// Check a window of counters to allow for clock drift
	for lookAhead := uint64(0); lookAhead < 5; lookAhead++ {
		testCounter := currentCounter + lookAhead
		log.Printf("Testing hotp counter for userID %v: %v", userID, testCounter)
		isValid, err := ValidateHOTP(otpSecret, userProvidedOTP, testCounter)
		if err != nil {
			log.Printf("Error validating hotp counter: %v", err)
			return false, err
		}

		if isValid {
			_, err = s.store.IncreaseCounter(ctx, userID)
			if err != nil {
				log.Printf("Error increasing counter for userID %v: %v", userID, err)
				return false, err
			}
			return true, nil
		}
	}

	return false, errors.New("invalid OTP")
}

func (s *OTPService) CleanupExpiredCounters(ctx context.Context) error {
	return s.store.CleanupExpiredCounters(ctx)
}
