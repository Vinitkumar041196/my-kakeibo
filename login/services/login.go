package services

import (
	"fmt"
	"kakeibo-api/domain"
	"kakeibo-api/utils"
	"log"
	"net/http"
	"os"
)

type loginService struct {
	loginRepo domain.LoginRepository
	userRepo  domain.UserRepository
}

func NewLoginService(loginRepo domain.LoginRepository, userRepo domain.UserRepository) domain.LoginService {
	return &loginService{
		loginRepo: loginRepo,
		userRepo:  userRepo,
	}
}

func (ls *loginService) ProcessLogin(email, password string) (*domain.LoginResponse, int, error) {
	functionName := "login"
	message := ""
	defer func() {
		data := password
		if message == "success" {
			data = "" // to not expose password
		}
		ls.AddAuditLog(functionName, email, data, message)
	}()

	code, err := ls.CheckForTooManyAttempts(email, functionName)
	if err != nil {
		message = err.Error()
		return nil, code, err
	}

	user, err := ls.userRepo.SearchUserByKey("email", email)
	if err != nil {
		message = err.Error()
		return nil, http.StatusInternalServerError, err
	}

	if user.Email != email {
		message = "user not found."
		return nil, http.StatusUnauthorized, fmt.Errorf(message)
	}

	if user.Password == "" {
		message = "cannot retrieve user password"
		return nil, http.StatusUnauthorized, fmt.Errorf(message)
	}

	err = utils.CompareStringWithHash(password, user.Password)
	if err != nil {
		message = err.Error()
		log.Println(err)
		return nil, http.StatusUnauthorized, fmt.Errorf("incorrect credentials")
	}

	token, refreshToken, err := utils.GenerateTokens(user.Email)
	if err != nil {
		message = err.Error()
		log.Println(err)
		return nil, http.StatusInternalServerError, fmt.Errorf("error while generating token")
	}

	//store the token
	tokenDets := new(domain.UserRefreshTokens)
	tokenDets.Email = user.Email
	tokenDets.RefreshToken = refreshToken //updating refresh token

	err = ls.loginRepo.StoreUserRefreshToken(*tokenDets)
	if err != nil {
		message = err.Error()
		log.Println(err)
		return nil, http.StatusInternalServerError, err
	}

	message = "success"
	resp := new(domain.LoginResponse)
	resp.Token = token
	resp.RefreshToken = refreshToken
	resp.Message = "success"
	return resp, http.StatusOK, nil
}

func (ls *loginService) RefreshToken(refreshToken string) (*domain.RefreshTokenResponse, int, error) {
	claims, err := utils.VerifyToken(refreshToken, os.Getenv("REFRESH_TOKEN_SECRET"))
	if err != nil {
		return nil, http.StatusUnauthorized, err
	}

	email := claims["sub"].(string)

	tokenDets, err := ls.loginRepo.GetUserRefreshToken(email)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if refreshToken != tokenDets.RefreshToken {
		return nil, http.StatusUnauthorized, fmt.Errorf("incorrect credentials")
	}

	token, _, err := utils.GenerateTokens(email)
	if err != nil {
		log.Println(err)
		return nil, http.StatusInternalServerError, fmt.Errorf("error while generating token")
	}

	resp := new(domain.RefreshTokenResponse)
	resp.Token = token
	// resp.RefreshToken = refreshToken
	resp.Message = "success"
	return resp, http.StatusOK, nil
}

func (ls *loginService) VerifyEmailAndSendOTP(email string) (string, int, error) {
	functionName := "forgotpassword_sendotp"
	message := "success"

	defer func() {
		ls.AddAuditLog(functionName, email, email, message)
	}()

	code, err := ls.CheckForTooManyAttempts(email, functionName)
	if err != nil {
		message = err.Error()
		return "", code, err
	}

	user, err := ls.userRepo.SearchUserByKey("email", email)
	if err != nil {
		message = err.Error()
		return "", http.StatusInternalServerError, err
	}

	if user.Email != email {
		message = "no user found"
		return "", http.StatusBadRequest, fmt.Errorf(message)
	}

	otp, expiry := utils.GenerateOTP()

	//ADD logic to send OTP
	// mail := notifications.Mail{
	// Priority:      1,
	// Subject:       "OTP for password forgot",
	// RecipientMail: []notifications.RecipientEmail{{Email: user.Email}},
	// Content:       fmt.Sprintf("Hello %s %s,<br> OTP for forgotting your password is : %s <br><br>Thanks.", user.FirstName, user.LastName, otp),
	// }

	// resp, err := notifications.SendMail(mail)
	// if err != nil {
	// message = err.Error()
	// return "", http.StatusInternalServerError, fmt.Errorf("error while generating OTP")
	// }

	// log.Println("sendMail response:", resp)

	token, err := utils.GenerateJWTToken(user.Email, os.Getenv("OTP_TOKEN_SECRET"), utils.GetOtpExpiry())
	if err != nil {
		message = err.Error()
		return "", http.StatusInternalServerError, err
	}

	err = ls.loginRepo.StoreOTPDetails(user.Email, otp, token, expiry)
	if err != nil {
		message = err.Error()
		return "", http.StatusInternalServerError, fmt.Errorf("error while generating OTP")
	}

	return token, http.StatusOK, nil
}

func (ls *loginService) VerifyOTP(otpToken, otp string) (string, int, error) {
	functionName := "forgotpassword_verifyotp"
	message := "success"
	defer func() {
		ls.AddAuditLog(functionName, otpToken, otp, message)
	}()

	code, err := ls.CheckForTooManyAttempts(otpToken, functionName)
	if err != nil {
		message = err.Error()
		return "", code, err
	}

	claims, err := utils.VerifyToken(otpToken, os.Getenv("OTP_TOKEN_SECRET"))
	if err != nil {
		message = err.Error()
		return "", http.StatusUnauthorized, err
	}

	email, _ := claims["sub"].(string)

	otpDets, err := ls.loginRepo.GetOTPDetails(email)
	if err != nil {
		message = err.Error()
		return "", http.StatusInternalServerError, err
	}

	if otpToken != otpDets.Token {
		message = "token mismatch"
		return "", http.StatusUnauthorized, fmt.Errorf("incorrect credentials")
	}

	if otpDets.Email != email || otpDets.OTP != otp {
		message = "incorrect otp"
		return "", http.StatusBadRequest, fmt.Errorf(message)
	}

	passwordToken, err := utils.GenerateJWTToken(email, os.Getenv("PASSWORD_TOKEN_SECRET"), utils.GetForgotPasswordTokenExpiry())
	if err != nil {
		message = err.Error()
		return "", http.StatusInternalServerError, err
	}

	err = ls.loginRepo.StoreUserForgotPasswordToken(email, passwordToken)
	if err != nil {
		message = err.Error()
		log.Println(err)
		return "", http.StatusInternalServerError, err
	}

	err = ls.loginRepo.DeleteOTPDetails(email)
	if err != nil {
		message = err.Error()
		log.Println(err)
		return "", http.StatusInternalServerError, err
	}

	return passwordToken, http.StatusOK, nil
}

func (ls *loginService) ForgotPassword(passwordToken, newPassword, confirmPassword string) (int, error) {
	functionName := "forgotpassword"
	message := "success"
	defer func() {
		data := fmt.Sprintf("new password: %s confirm password %s", newPassword, confirmPassword)
		if message == "success" {
			data = "" // to not expose password
		}
		ls.AddAuditLog(functionName, passwordToken, data, message)
	}()

	code, err := ls.CheckForTooManyAttempts(passwordToken, functionName)
	if err != nil {
		message = err.Error()
		return code, err
	}
    
	if newPassword == ""{
		return http.StatusBadRequest, fmt.Errorf("invalid new password")
	}

	if newPassword != confirmPassword {
		message = "new and confirm password mismatch"
		return http.StatusBadRequest, fmt.Errorf("new password and confirm password are not same")
	}

	claims, err := utils.VerifyToken(passwordToken, os.Getenv("PASSWORD_TOKEN_SECRET"))
	if err != nil {
		message = err.Error()
		return http.StatusUnauthorized, err
	}

	email, _ := claims["sub"].(string)
	if email == "" {
		message = "empty email"
		return http.StatusUnauthorized, fmt.Errorf("incorrect credentials")
	}

	passDets, err := ls.loginRepo.GetUserForgotPasswordToken(email)
	if err != nil {
		message = err.Error()
		return http.StatusInternalServerError, err
	}

	if passwordToken != passDets.Token {
		message = "token mismatch"
		return http.StatusUnauthorized, fmt.Errorf("incorrect credentials")
	}

	user, err := ls.userRepo.SearchUserByKey("email", email)
	if err != nil {
		message = err.Error()
		return http.StatusInternalServerError, err
	}

	if user.Email != email {
		message = "user not found"
		return http.StatusUnauthorized, fmt.Errorf("incorrect credentials")
	}

	user.Password = newPassword

	err = ls.userRepo.UpdateUser(user)
	if err != nil {
		message = err.Error()
		return http.StatusInternalServerError, fmt.Errorf("could not set password")
	}

	err = ls.loginRepo.DeleteUserForgotPasswordToken(email)
	if err != nil {
		message = err.Error()
		log.Println(err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (ls *loginService) AddAuditLog(function, email, data, message string) error {
	return ls.loginRepo.AddAuditLog(function, email, data, message)
}

func (ls *loginService) CheckForTooManyAttempts(email, function string) (int, error) {
	alog, err := ls.loginRepo.GetLastestAttempts(email, function, utils.GetCoolOffPeriod())
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if len(alog) > utils.GetRetryLimit() {
		return http.StatusTooManyRequests, fmt.Errorf("too many attempts")
	}
	return http.StatusOK, nil
}
