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

const ClaimTokenJWT = "claimToken" 

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("failed to load configuration: %v", err)
	}
}

func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get("Authorization")
			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing authorization token"})
			}

			rawToken := strings.TrimPrefix(tokenString, "Bearer ")

			id, role, err := ExtractTokenFromRaw(rawToken)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
			}

			c.Set("id", id)
			c.Set("role", role)
			c.Set(ClaimTokenJWT, rawToken) 

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

func ExtractTokenFromRaw(rawToken string) (string, string, error) {
	token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return "", "", errors.New("invalid authorization token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}

	id, okID := claims["id"].(string)
	role, okRole := claims["role"].(string)
	if !okID || !okRole {
		return "", "", errors.New("invalid token claims")
	}

	return id, role, nil
}
