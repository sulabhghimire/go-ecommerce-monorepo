package helper

import (
	"ecommerce/internal/domain"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	Secret string
}

func SetUpAuth(s string) Auth {
	return Auth{Secret: s}
}

func (a Auth) CreateHashedPassword(p string) (string, error) {

	if len(p) < 6 {
		return "", errors.New("password length should be at least 6 characters long")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(p), 10)
	if err != nil {
		// log actual error
		fmt.Printf("Error encrypting password %v", err)
		return "", errors.New("password hash failed")
	}

	return string(hashPassword), nil
}

func (a Auth) GenerateToken(id uint, email string, role string) (string, error) {

	if id == 0 || email == "" || role == "" {
		return "", errors.New("inputs missing to generate token")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(a.Secret))
	if err != nil {
		fmt.Printf("Error signing token %v", err)
		return "", errors.New("unable to sign the token")
	}

	return tokenStr, nil
}

func (a Auth) VerifyPassword(pP string, hP string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hP), []byte(pP))
	if err != nil {
		return errors.New("password doesn't match")
	}

	return nil
}

func (a Auth) VerifyToken(t string) (domain.User, error) {

	// Bearer value
	if len(t) < 1 {
		return domain.User{}, errors.New("invalid token")
	}

	tokenStr := strings.TrimPrefix(t, "Bearer ")
	tokenStr = strings.TrimSpace(tokenStr)
	if len(tokenStr) < 1 {
		return domain.User{}, errors.New("invalid token")
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("signing method error")
			return nil, fmt.Errorf("unknown signing method %v", t.Header)
		}
		return []byte(a.Secret), nil
	})
	if err != nil {
		return domain.User{}, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return domain.User{}, errors.New("token is expired")
		}

		user := domain.User{}
		user.ID = uint(claims["user_id"].(float64))
		user.Email = claims["email"].(string)
		user.UserType = claims["role"].(string)

		return user, nil
	}

	return domain.User{}, errors.New("token validation failed")
}

func (a Auth) Authorize(ctx *fiber.Ctx) error {

	authHeader := ctx.Get("Authorization")
	user, err := a.VerifyToken(authHeader)
	if err == nil && user.ID > 0 {
		ctx.Locals("user", user)
		return ctx.Next()
	} else {
		fmt.Println(err)
		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "authorization failed",
			"reason":  err.Error(),
		})
	}

}

func (a Auth) AuthorizeSeller(ctx *fiber.Ctx) error {

	authHeader := ctx.Get("Authorization")
	user, err := a.VerifyToken(authHeader)
	if err != nil {
		fmt.Println(err)
		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "authorization failed",
			"reason":  err.Error(),
		})
	} else if user.ID > 0 && user.UserType == domain.SELLER {
		ctx.Locals("user", user)
		return ctx.Next()
	} else {
		return ctx.Status(http.StatusForbidden).JSON(&fiber.Map{
			"message": "authentication failed",
			"reason":  "please join in as a seller",
		})
	}

}

func (a Auth) GetCurrentUser(ctx *fiber.Ctx) domain.User {

	user := ctx.Locals("user")
	return user.(domain.User)

}

func (a Auth) GenerateVerificationCode() (int, error) {
	return RandomNumber(6)
}
