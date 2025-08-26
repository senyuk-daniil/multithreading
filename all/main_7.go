//**Условие**:
//
//Напиши две функции:
//
//1. `AddJWTToContext(ctx context.Context, userID int) (context.Context, error)`
//2. `ExtractUserIDFromContext(ctx context.Context) (int, error)`
//
//`AddJWTToContext` должен:
//
//- Создавать JWT-токен (используй `github.com/golang-jwt/jwt/v5`).
//- Зашифровывать в него `userID`.
//- Возвращать новый `context.Context`, в который записан JWT-токен.
//
//`ExtractUserIDFromContext` должен:
//
//- Извлекать JWT-токен из контекста.
//- Расшифровывать `userID`.
//- Вывести его на экран
//
//**Дополнительное условие**:
//
//Создай горутину, в которой будет использоваться `ExtractUserIDFromContext`. Покажи, что передача контекста работает и данные можно безопасно извлекать между горутинами.

package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

type contextKey string

const jwtContextKey contextKey = "jwtToken"

var privateKey *ecdsa.PrivateKey
var publicKey *ecdsa.PublicKey

func init() {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	privateKey = key
	publicKey = &key.PublicKey
}

func AddJWTToContext(ctx context.Context, userID int) (context.Context, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return ctx, fmt.Errorf("error signing token: %w", err)
	}

	return context.WithValue(ctx, jwtContextKey, tokenString), nil
}

func ExtractUserIDFromContext(ctx context.Context) (int, error) {
	tokenStr, ok := ctx.Value(jwtContextKey).(string)
	if !ok {
		return 0, fmt.Errorf("error extracting user id from context")
	}

	parsedToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return 0, fmt.Errorf("error parsing token: %w", err)
	}

	//проверка подписи
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		if sub, ok := claims["sub"].(float64); ok {
			return int(sub), nil
		}
		return 0, fmt.Errorf("error extracting user id from token")
	}

	return 0, fmt.Errorf("error extracting user id from token")
}

func main() {
	ctx := context.Background()

	ctx, err := AddJWTToContext(ctx, 42)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})

	go func(ctx context.Context) {
		defer close(done)
		userID, err := ExtractUserIDFromContext(ctx)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("User ID:", userID)
	}(ctx)

	<-done
}
