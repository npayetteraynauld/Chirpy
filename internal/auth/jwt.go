package auth

import (
 "time"
 "errors"

 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
 currentTime := time.Now().UTC()
 claims := jwt.RegisteredClaims{
   Issuer: "chirpy",
   IssuedAt: jwt.NewNumericDate(currentTime),
   ExpiresAt: jwt.NewNumericDate(currentTime.Add(expiresIn)),
   Subject: userID.String(),
  }

 token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
 return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
 token, err := jwt.ParseWithClaims(
  tokenString, 
  &jwt.RegisteredClaims{},
  func(token *jwt.Token) (interface{}, error) {
   return []byte(tokenSecret), nil
  },
 )
 if err != nil {
  return uuid.UUID{}, err
 }

 if !token.Valid {
  return uuid.UUID{}, errors.New("invalid or expired token")
 }

 claims, ok := token.Claims.(*jwt.RegisteredClaims)
 if !ok {
  return uuid.UUID{}, errors.New("could not parse claims")
 }
 userID := claims.Subject

 u, err := uuid.Parse(userID)
 if err != nil {
  return uuid.UUID{}, err
 }

 return u, nil
}

