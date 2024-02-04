package repository

import (
	"fmt"
	"kakeibo-api/db"
	"kakeibo-api/domain"
	"kakeibo-api/utils"
	"time"
)

type loginRepository struct {
	db *db.PostgresDB
}

func NewLoginRepository(db *db.PostgresDB) domain.LoginRepository {
	fmt.Println("Automigrate error:", db.GormDB.AutoMigrate(&domain.UserOTPDetails{}))
	fmt.Println("Automigrate error:", db.GormDB.AutoMigrate(&domain.UserForgotPasswordTokens{}))
	fmt.Println("Automigrate error:", db.GormDB.AutoMigrate(&domain.UserRefreshTokens{}))
	fmt.Println("Automigrate error:", db.GormDB.AutoMigrate(&domain.LoginAudit{}))
	return &loginRepository{db: db}
}

func (lr *loginRepository) GetUserRefreshToken(email string) (*domain.UserRefreshTokens, error) {
	tokenDets := new(domain.UserRefreshTokens)
	err := lr.db.GormDB.Where("email = ?", email).Find(tokenDets).Error
	if err != nil {
		return nil, err
	}
	return tokenDets, nil
}

func (lr *loginRepository) DeleteUserRefreshToken(email string) error {
	tokenDets := new(domain.UserRefreshTokens)
	tokenDets.Email = email
	return lr.db.GormDB.Delete(tokenDets).Error
}

func (lr *loginRepository) StoreUserRefreshToken(u domain.UserRefreshTokens) error {
	now := utils.TimeNow().Format(time.RFC3339)
	u.UpdatedAt = &now
	err := lr.db.GormDB.Save(&u).Error
	if err != nil {
		return err
	}

	return nil
}

func (lr *loginRepository) GetLastestAttempts(email, function string, duration time.Duration) ([]domain.LoginAudit, error) {
	now := utils.TimeNow().Add(duration * time.Second * -1).Format(time.RFC3339)
	alog := []domain.LoginAudit{}
	err := lr.db.GormDB.Where("created_at >= ?", now).Where("email = ?", email).Where("function = ?", function).Order("created_at desc").Find(&alog).Error
	if err != nil {
		return nil, err
	}

	return alog, err
}

func (lr *loginRepository) AddAuditLog(function, email, data, message string) error {
	now := utils.TimeNow().Format(time.RFC3339)
	alog := new(domain.LoginAudit)
	alog.Function = function
	alog.Email = email
	alog.Data = data
	alog.Message = message
	alog.CreatedAt = &now
	return lr.db.GormDB.Create(&alog).Error
}

func (lr *loginRepository) StoreOTPDetails(email, otp, hashedToken, expiry string) error {
	now := utils.TimeNow().Format(time.RFC3339)
	otpDets := new(domain.UserOTPDetails)
	otpDets.Email = email
	otpDets.OTP = otp
	otpDets.Expiry = expiry
	otpDets.Token = hashedToken
	otpDets.CreatedAt = &now
	return lr.db.GormDB.Save(&otpDets).Error
}

func (lr *loginRepository) GetOTPDetails(email string) (*domain.UserOTPDetails, error) {
	now := utils.TimeNow().Format(time.RFC3339)
	otpDets := new(domain.UserOTPDetails)
	err := lr.db.GormDB.Where("expiry >= ?", now).Where("email = ?", email).Find(otpDets).Error
	if err != nil {
		return nil, err
	}
	if otpDets == nil {
		return nil, fmt.Errorf("no otp found")
	}

	return otpDets, nil
}

func (lr *loginRepository) DeleteOTPDetails(email string) error {
	otpDets := new(domain.UserOTPDetails)
	otpDets.Email = email
	return lr.db.GormDB.Delete(otpDets).Error
}

func (lr *loginRepository) GetUserForgotPasswordToken(email string) (*domain.UserForgotPasswordTokens, error) {
	tokenDets := new(domain.UserForgotPasswordTokens)
	err := lr.db.GormDB.Where("email = ?", email).Find(tokenDets).Error
	if err != nil {
		return nil, err
	}
	return tokenDets, nil
}

func (lr *loginRepository) DeleteUserForgotPasswordToken(email string) error {
	tokenDets := new(domain.UserForgotPasswordTokens)
	tokenDets.Email = email
	return lr.db.GormDB.Delete(tokenDets).Error
}

func (lr *loginRepository) StoreUserForgotPasswordToken(email, hashedtoken string) error {
	now := utils.TimeNow().Format(time.RFC3339)
	u := domain.UserForgotPasswordTokens{}
	u.Email = email
	u.Token = hashedtoken
	u.CreatedAt = &now
	err := lr.db.GormDB.Save(&u).Error
	if err != nil {
		return err
	}

	return nil
}
