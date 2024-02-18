package usertokens

import (
	"errors"
	"time"

	"github.com/nisil-thaing/golang-gingonic-practice/internal/app/users"
	"github.com/nisil-thaing/golang-gingonic-practice/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdateUserTokens(ctx mongo.SessionContext, user users.UserSchema, secretKey string) (*UserTokensPublicInfo, error) {
	var dbClient *mongo.Client = database.GetDBInstance()
	var userTokensCollection *mongo.Collection = database.OpenCollection(dbClient, "user_tokens")
	var existingUserTokensDetails UserTokenSchema

	filter := bson.M{"user_id": user.UserID}
	err := userTokensCollection.FindOne(ctx, filter).Decode(&existingUserTokensDetails)
	currentTime, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	userTokens, generatingTokensErr := GenerateTokens(user, secretKey)

	if userTokens == nil {
		return nil, errors.New("Could not generate JWT tokens due to some unexpected errors!")
	}

	if generatingTokensErr != nil {
		return nil, generatingTokensErr
	}

	if err != nil {
		var newUserTokensDetails UserTokenSchema
		id := primitive.NewObjectID()

		newUserTokensDetails = UserTokenSchema{
			ID:           id,
			UserID:       user.UserID,
			AccessToken:  userTokens.AccessToken,
			RefreshToken: userTokens.RefreshToken,
			ExpiresAt:    userTokens.ExpiresAt,
			CreatedAt:    currentTime,
			UpdatedAt:    currentTime,
		}

		_, err = userTokensCollection.InsertOne(ctx, newUserTokensDetails)
	} else {
		var updatingData primitive.D

		upsert := false
		opt := options.UpdateOptions{Upsert: &upsert}

		updatingData = append(updatingData, bson.E{Key: "access_token", Value: userTokens.AccessToken})
		updatingData = append(updatingData, bson.E{Key: "refresh_token", Value: userTokens.RefreshToken})
		updatingData = append(updatingData, bson.E{Key: "expires_at", Value: userTokens.ExpiresAt})
		updatingData = append(updatingData, bson.E{Key: "updated_at", Value: currentTime})

		_, err = userTokensCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updatingData}}, &opt)
	}

	if err != nil {
		return nil, err
	}

	return userTokens, nil
}
