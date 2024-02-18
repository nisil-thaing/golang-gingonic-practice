package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	userpasswords "github.com/nisil-thaing/golang-gingonic-practice/internal/app/user_passwords"
	usertokens "github.com/nisil-thaing/golang-gingonic-practice/internal/app/user_tokens"
	"github.com/nisil-thaing/golang-gingonic-practice/internal/app/users"
	"github.com/nisil-thaing/golang-gingonic-practice/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var validate = validator.New()

func HandleRegistration(c *gin.Context) {
	secretKey := os.Getenv("JWT_SECRET_KEY")

	if secretKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate JWT tokens!"})
		return
	}

	var registeringUser RegisteringUserSchema
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	if err := c.ShouldBindJSON(&registeringUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		defer cancel()
		return
	}

	if err := validate.Struct(registeringUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		defer cancel()
		return
	}

	var dbClient *mongo.Client = database.GetDBInstance()
	var usersCollection *mongo.Collection = database.OpenCollection(dbClient, "users")
	numOfExistingUsers, err := usersCollection.CountDocuments(ctx, bson.M{"email": registeringUser.Email})

	defer cancel()

	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred when we tried to validate your information!"})
		return
	}

	if numOfExistingUsers > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "This email has been used before!"})
		return
	}

	session, err := dbClient.StartSession()

	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open a new transaction!"})
		return
	}

	defer session.EndSession(context.Background())

	// TODO: store user data to the database
	userID := primitive.NewObjectID()
	userType := "USER"

	newUser := users.UserSchema{
		ID:        userID,
		UserID:    userID.Hex(),
		Type:      userType,
		Email:     registeringUser.Email,
		FirstName: registeringUser.FirstName,
		LastName:  registeringUser.LastName,
	}

	var accessToken string
	var refreshToken string

	txnFunc := func(sessionCtx mongo.SessionContext) (interface{}, error) {
		result, err := usersCollection.InsertOne(sessionCtx, newUser)

		if err != nil {
			return nil, err
		}

		err = userpasswords.UpdateUserPassword(sessionCtx, newUser.UserID, registeringUser.Password)

		if err != nil {
			return nil, err
		}

		userTokens, err := usertokens.UpdateUserTokens(sessionCtx, newUser, secretKey)

		if err != nil {
			return nil, errors.New("Could not generate token due to some unexpected issues!")
		}

		accessToken = userTokens.AccessToken
		refreshToken = userTokens.RefreshToken

		return result, nil
	}

	txnOptions := options.Transaction().SetReadPreference(readpref.Primary())
	_, err = session.WithTransaction(context.Background(), txnFunc, txnOptions)

	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accessToken": accessToken, "refreshToken": refreshToken})
	return
}
