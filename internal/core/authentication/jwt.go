package authentication

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

type JwtTokenProvider struct {
	SecretKey              string
	ExpirationAccessToken  time.Duration
	ExpirationRefreshToken time.Duration
}

// NewJwtTokenProvider creates a new JwtTokenProvider instance
func NewJwtTokenProvider(secretKey string, expirationAccessToken, expirationRefreshToken time.Duration) *JwtTokenProvider {
	return &JwtTokenProvider{
		SecretKey:              secretKey,
		ExpirationAccessToken:  expirationAccessToken,
		ExpirationRefreshToken: expirationRefreshToken,
	}
}

func (j *JwtTokenProvider) GenerateToken(value string) (string, error) {
	log.Println("Generating accessToken")

	claims := jwt.MapClaims{
		"sub": value,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(j.ExpirationAccessToken).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		log.Println("Error signing token:", err)
		return "", err
	}

	return signedToken, nil
}

func (j *JwtTokenProvider) GenerateRefreshToken(value string) (string, error) {
	log.Println("Generating refreshToken")

	claims := jwt.MapClaims{
		"sub": value,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(j.ExpirationRefreshToken).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		log.Println("Error signing token:", err)
		return "", err
	}

	return signedToken, nil
}

func (j *JwtTokenProvider) DecodeTokenClaims(tokenString string) (jwt.MapClaims, error) {
	log.Printf("Decoding token: %s\n", tokenString)

	// ตรวจสอบว่ามีการส่ง token มาหรือไม่
	if tokenString == "" {
		log.Println("Token string is empty")
		return nil, fmt.Errorf("token is empty")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// ตรวจสอบ signing method ว่าตรงกับที่คาดหวังหรือไม่
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Invalid signing method: %v", token.Method)
			return nil, fmt.Errorf("invalid signing method")
		}

		// คืนค่า SecretKey สำหรับตรวจสอบลายเซ็นต์
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		log.Println("Error parsing token:", err)
		return nil, err
	}

	// ตรวจสอบว่า token เป็น valid และ claims สามารถถูกดึงได้หรือไม่
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//log.Println("Token is valid. Claims:", claims)
		return claims, nil
	} else {
		log.Println("Invalid token or failed to parse claims")
		return nil, fmt.Errorf("invalid token")
	}
}

func (j *JwtTokenProvider) ValidateToken(tokenString string) bool {
	log.Printf("Validating token: %s\n", tokenString)
	_, err := j.DecodeTokenClaims(tokenString)
	if err != nil {
		log.Printf("Invalid token: %s", err)
		return false
	}
	return true
}
