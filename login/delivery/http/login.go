package http

import (
	"kakeibo-api/domain"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	service domain.LoginService
}

func NewLoginHandler(service domain.LoginService) *LoginHandler {
	return &LoginHandler{service: service}
}

// Login godoc
// @Summary Login API
// @Description Validates user credentials and generate token
// @Tags Login
// @Produce json
// @Param req body domain.LoginRequest true "user login credentials"
// @Success 200 {object} domain.LoginResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 429 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /login [post]
func (h *LoginHandler) Login(c *gin.Context) {
	req := domain.LoginRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "empty login credentials"})
		return
	}

	resp, code, err := h.service.ProcessLogin(req.Email, req.Password)
	if err != nil {
		if code == 0 {
			code = http.StatusInternalServerError
		}
		c.JSON(code, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RefreshToken godoc
// @Summary RefreshToken API
// @Description generates new tokens
// @Tags Login
// @Produce json
// @Param Authorization header string true "refresh token - 'Bearer' followed by a space and JWT token."
// @Success 200 {object} domain.RefreshTokenResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /refreshtoken [post]
func (h *LoginHandler) RefreshToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Error: "invalid credentials"})
		return
	}

	tokenArr := strings.Split(token, " ")
	if len(tokenArr) != 2 {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Error: "invalid credentials"})
		return
	}

	refreshToken := tokenArr[1]

	resp, code, err := h.service.RefreshToken(refreshToken)
	if err != nil {
		if code == 0 {
			code = http.StatusInternalServerError
		}
		c.JSON(code, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ForgotPasswordSendOTP godoc
// @Summary SendOTP for resetting password
// @Description Sends OTP for forgot password process
// @Tags Login
// @Produce json
// @Param req body domain.VerifyEmailAndSendOTPRequest true "email"
// @Success 200 {object} domain.VerifyEmailAndSendOTPResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 429 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /forgotpassword/verifyemail [post]
func (h *LoginHandler) VerifyEmailAndSendOTP(c *gin.Context) {
	req := domain.VerifyEmailAndSendOTPRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	token, code, err := h.service.VerifyEmailAndSendOTP(req.Email)
	if err != nil {
		if code == 0 {
			code = http.StatusInternalServerError
		}
		c.JSON(code, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.VerifyEmailAndSendOTPResponse{Token: token, Message: "succes"})
}

// ForgotPasswordVerifyOTP godoc
// @Summary VerifyOTP API
// @Description Verifies OTP for forgot password process
// @Tags Login
// @Produce json
// @Param Authorization header string true "token - 'Bearer' followed by a space and JWT token."
// @Param req body domain.VerifyOTPRequest true "otp"
// @Success 200 {object} domain.VerifyOTPResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 429 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /forgotpassword/verifyotp [post]
func (h *LoginHandler) ForgotPasswordVerifyOTP(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Error: "invalid credentials"})
		return
	}

	tokenArr := strings.Split(token, " ")
	if len(tokenArr) != 2 {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Error: "invalid credentials"})
		return
	}

	otpToken := tokenArr[1]

	req := domain.VerifyOTPRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	forgotPasswordToken, code, err := h.service.VerifyOTP(otpToken, req.OTP)
	if err != nil {
		if code == 0 {
			code = http.StatusInternalServerError
		}
		c.JSON(code, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.VerifyOTPResponse{Token: forgotPasswordToken, Message: "succes"})
}

// ForgotPassword godoc
// @Summary ForgotPassword API
// @Description Verifies old password and update the new password in DB
// @Tags Login
// @Produce json
// @Param Authorization header string true "token - 'Bearer' followed by a space and JWT token."
// @Param req body domain.ForgotPasswordRequest true "new password"
// @Success 200 {object} domain.ForgotPasswordResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 429 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /forgotpassword/reset [post]
func (h *LoginHandler) ForgotPassword(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Error: "invalid credentials"})
		return
	}

	tokenArr := strings.Split(token, " ")
	if len(tokenArr) != 2 {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Error: "invalid credentials"})
		return
	}

	forgotPasswordToken := tokenArr[1]

	req := domain.ForgotPasswordRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	code, err := h.service.ForgotPassword(forgotPasswordToken, req.NewPassword, req.ConfirmPassword)
	if err != nil {
		if code == 0 {
			code = http.StatusInternalServerError
		}
		c.JSON(code, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.ForgotPasswordResponse{Message: "Password Changed Successfully. Please Login."})
}
