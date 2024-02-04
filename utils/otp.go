package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"math"
	"math/big"
	"time"
)

const (
	DEFAULT_OTP_LENGTH = 6
	DEFAULT_OTP_EXPIRY = 300 * time.Second
	DEFAULT_COOLOFF    = 30 * time.Second
)

func GenerateOTP() (string, string) {
	maxDigits := GetOtpLength()

	num, err := rand.Int(
		rand.Reader,
		big.NewInt(int64(math.Pow(10, float64(maxDigits)))),
	)

	expiry := TimeNow().Add(GetOtpExpiry()).Format(time.RFC3339)
	if err != nil {
		log.Println(err)
		return fmt.Sprint(time.Now().Nanosecond())[:6], expiry
	}
	return fmt.Sprintf("%0*d", maxDigits, num), expiry
}

func GetOtpLength() int {
	// length := config.Config.GetInt("OTP_LENGTH")
	// if length == 0 {
	return DEFAULT_OTP_LENGTH
	// }
	// return length
}

func GetOtpExpiry() time.Duration {
	// duration := config.Config.GetInt("OTP_EXPIRY_DURATION")
	// if duration == 0 {
	return DEFAULT_OTP_EXPIRY
	// }
	// return time.Duration(duration) * time.Second
}

func GetCoolOffPeriod() time.Duration {
	// duration := config.Config.GetInt("COOLOFF_DURATION")
	// if duration == 0 {
	return DEFAULT_COOLOFF
	// }
	// return time.Duration(duration) * time.Second
}

func GetRetryLimit() int {
	// limit := config.Config.GetInt("MAX_RETRY_LIMIT")
	// if limit == 0 {
	return 5
	// }
	// return limit
}
