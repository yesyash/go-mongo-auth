package gomongoauth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type JwtSignedDetails struct {
	ID    string `json:"_id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.StandardClaims
}

/**
main functions
*/
func Signup(client *mongo.Client, user User) AuthResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	database := client.Database("pointform")
	collection := database.Collection("users")

	// Validate user input with the struct
	if validateErr := validateInput(user); validateErr != nil {
		return AuthResponse{
			Status: 400,
			Msg:    "Invalid Input",
			Err:    validateErr,
		}
	}

	// check if user already exists
	if data, _ := getUserByEmail(*collection, ctx, user.Email); data.Email == user.Email {
		return AuthResponse{
			Status: 409,
			Msg:    "User already exists",
			Err:    errors.New("User Err : user already exists"),
		}
	}

	// if user does not exist
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 13)

	if err != nil {
		fmt.Println(err)
	}

	newUser := User{
		Name:     user.Name,
		Email:    user.Email,
		Password: string(password),
	}

	if _, err := collection.InsertOne(ctx, newUser); err != nil {
		return AuthResponse{
			Status: 500,
			Msg:    "Error creating new user",
			Err:    err,
		}
	}

	res := AuthResponse{
		Status: 200,
		Msg:    "Successfully created new user",
		Err:    nil,
	}

	return res
}

func Login(client *mongo.Client, user User) AuthResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	database := client.Database("pointform")
	collection := database.Collection("users")

	// Validate user input with the struct
	if validateErr := validateInput(user); validateErr != nil {
		return AuthResponse{
			Status: 400,
			Msg:    "Invalid inpput",
			Err:    validateErr,
		}
	}

	data, err := getUserByEmail(*collection, ctx, user.Email)

	if err != nil && err == mongo.ErrNoDocuments {
		return AuthResponse{
			Status: 404,
			Msg:    "No user found",
			Err:    err,
		}
	} else if err != nil && err != mongo.ErrNoDocuments {
		return AuthResponse{
			Status: 500,
			Msg:    "Internal Server Error",
			Err:    err,
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(user.Password)); err != nil {
		return AuthResponse{
			Status: 400,
			Msg:    "Invalid email/password",
			Err:    err,
		}
	}

	claims := JwtSignedDetails{
		ID:    data.ID.Hex(),
		Name:  data.Name,
		Email: data.Email,
	}

	secretKey := os.Getenv("SECRET_KEY")
	secretKeyErr := godotenv.Load()

	if secretKeyErr != nil || len(secretKey) == 0 {
		resErr := secretKeyErr

		if len(secretKey) == 0 {
			resErr = errors.New("Env Error : Unable to get SECRET_KEY")
		}

		return AuthResponse{
			Status: 500,
			Msg:    "Unable to get SECRET_KEY from environment",
			Err:    resErr,
		}
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))

	if err != nil {
		return AuthResponse{
			Status: 500,
			Msg:    "Internal server error",
			Err:    err,
		}
	}

	return AuthResponse{
		Status: 200,
		Msg:    "User found",
		Err:    nil,
		Data:   token,
	}
}
