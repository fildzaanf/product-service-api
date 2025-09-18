package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)
func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("failed to load configuration: %v", err)
	}
}

func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id, role, err := ExtractToken(c)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
			}
			c.Set("id", id)
			c.Set("role", role)
			return next(c)
		}
	}
}

func GenerateToken(id string, role string) (string, error) {
	logrus.Infof("generating token for user with ID: %s and Role: %s", id, role)
	tokenClaims := jwt.MapClaims{
		"authorized": true,
		"id":         id,
		"role":       role,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ExtractToken(c echo.Context) (string, string, error) {
	tokenString := c.Request().Header.Get("authorization")
	if tokenString == "" {
		return "", "", errors.New("missing authorization token")
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return "", "", errors.New("invalid authorization token")
	}
	claims, validClaims := token.Claims.(jwt.MapClaims)
	if !validClaims {
		return "", "", errors.New("invalid token claims")
	}
	id, validID := claims["id"].(string)
	role, validRole := claims["role"].(string)
	if !validID || !validRole {
		return "", "", errors.New("invalid token claims")
	}
	return id, role, nil
}
