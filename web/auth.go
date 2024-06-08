package web

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"imageProcessingTestTask/db"
	"os"
	"time"
)

var secret = []byte(os.Getenv("SECRET"))

type registrationRequest struct {
	Username string  `json:"username"`
	Email    *string `json:"email"`
	Password string  `json:"password"`
}

func (h *Handlers) Register(c echo.Context) error {
	var req registrationRequest
	if err := c.Bind(&req); err != nil {
		return echo.ErrUnprocessableEntity
	}

	if exists, _ := h.db.UsernameTaken(c.Request().Context(), req.Username); exists {
		return echo.ErrConflict
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.ErrInternalServerError
	}

	user, err := h.db.CreateUser(c.Request().Context(), db.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	token, err := makeToken(user.ID)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(200, echo.Map{
		"token": token,
	})
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handlers) Login(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return echo.ErrUnprocessableEntity
	}

	user, err := h.db.GetUserByUsername(c.Request().Context(), req.Username)
	if err != nil {
		return echo.ErrInternalServerError
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return echo.ErrUnauthorized
	}

	token, err := makeToken(user.ID)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(200, echo.Map{
		"token": token,
	})
}

func (h *Handlers) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")
		if len(auth) < 7 {
			return echo.ErrUnauthorized
		}
		token, err := jwt.Parse(auth[7:], func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secret, nil
		})
		if err != nil {
			return echo.ErrUnauthorized
		}
		c.Set("userId", int32(token.Claims.(jwt.MapClaims)["uid"].(float64)))
		return next(c)
	}
}

func makeToken(uid int32) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		UserId:     uid,
		Expiration: time.Now().Add(time.Hour * 12).Unix(),
	}).SignedString(secret)
}

type claims struct {
	UserId     int32 `json:"uid"`
	Expiration int64 `json:"exp"`
}

func (c claims) Valid() error {
	if time.Now().After(time.Unix(c.Expiration, 0)) {
		return jwt.NewValidationError("token expired", jwt.ValidationErrorExpired)
	}
	if c.UserId < 0 {
		return jwt.NewValidationError("invalid user id", jwt.ValidationErrorClaimsInvalid)
	}
	return nil
}
