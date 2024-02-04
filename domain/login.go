package domain

import "time"

type UserOTPDetails struct {
	Email     string  `gorm:"column:email;type:varchar(300);primaryKey"`
	OTP       string  `gorm:"column:otp;type:varchar(10)"`
	Token     string  `gorm:"column:token;type:text"`
	Expiry    string  `gorm:"column:expiry;type:timestamp"`
	CreatedAt *string `gorm:"column:created_at;type:timestamp"`
}

func (UserOTPDetails) TableName() string {
	return "user_otp"
}

type UserForgotPasswordTokens struct {
	Email     string  `gorm:"column:email;type:varchar(300);primaryKey"`
	Token     string  `gorm:"column:token;type:text"`
	CreatedAt *string `gorm:"column:created_at;type:timestamp"`
}

func (UserForgotPasswordTokens) TableName() string {
	return "user_password_reset_tokens"
}

type UserRefreshTokens struct {
	Email        string  `gorm:"column:email;type:varchar(300);primaryKey"`
	RefreshToken string  `gorm:"column:refresh_token;type:text"`
	UpdatedAt    *string `gorm:"column:updated_at;type:timestamp"`
}

func (UserRefreshTokens) TableName() string {
	return "user_refresh_token"
}

type LoginAudit struct {
	ID        int64   `gorm:"<-:false;column:id;type:serial;primaryKey"`
	Email     string  `gorm:"column:email;type:varchar(300)"`
	Message   string  `gorm:"column:message;type:text"`
	Data      string  `gorm:"column:data;type:text"`
	Function  string  `gorm:"column:function;type:varchar(200)"`
	CreatedAt *string `gorm:"column:created_at;type:timestamp"`
}

func (LoginAudit) TableName() string {
	return "login_audit"
}

type LoginService interface {
	AddAuditLog(function, email, password, message string) error
	ProcessLogin(email, password string) (*LoginResponse, int, error)
	RefreshToken(refreshtoken string) (*RefreshTokenResponse, int, error)
	VerifyEmailAndSendOTP(email string) (string, int, error)
	VerifyOTP(token, otp string) (string, int, error)
	ForgotPassword(token, newPassword, confirmPassword string) (int, error)
}

type LoginRepository interface {
	AddAuditLog(function, email, password, message string) error
	GetLastestAttempts(email, function string, duration time.Duration) ([]LoginAudit, error)

	StoreUserRefreshToken(u UserRefreshTokens) error
	GetUserRefreshToken(email string) (*UserRefreshTokens, error)
	DeleteUserRefreshToken(email string) error

	StoreOTPDetails(email, otp, hashedToken, expiry string) error
	GetOTPDetails(email string) (*UserOTPDetails, error)
	DeleteOTPDetails(email string) error

	StoreUserForgotPasswordToken(email, hashedtoken string) error
	GetUserForgotPasswordToken(email string) (*UserForgotPasswordTokens, error)
	DeleteUserForgotPasswordToken(email string) error
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message      string `json:"message"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type VerifyEmailAndSendOTPRequest struct {
	Email string `json:"email"`
}

type VerifyEmailAndSendOTPResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

type VerifyOTPRequest struct {
	OTP string `json:"otp"`
}

type VerifyOTPResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

type ForgotPasswordRequest struct {
	NewPassword     string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type ForgotPasswordResponse struct {
	Message string `json:"message"`
}
