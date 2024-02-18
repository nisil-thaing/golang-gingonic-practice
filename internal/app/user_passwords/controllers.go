package userpasswords

import (
	"context"
	"time"

	"github.com/nisil-thaing/golang-gingonic-practice/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func UpdateUserPassword(ctx context.Context, userId string, newPassword string) error {
	var dbClient *mongo.Client = database.GetDBInstance()
	var userPasswordsCollection *mongo.Collection = database.OpenCollection(dbClient, "user_passwords")

	var existingUserPasswordStored UserPasswordSchema
	var newUserPasswordDetails UserPasswordSchema

	filter := bson.M{"user_id": userId}
	err := userPasswordsCollection.FindOne(ctx, filter).Decode(&existingUserPasswordStored)

	currentTime, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	if err != nil {
		id := primitive.NewObjectID()
		salt, err := GenerateSalt(bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		combinedPasswordAndSalt := append([]byte(newPassword), []byte(salt)...)
		hashedPassword, err := bcrypt.GenerateFromPassword(combinedPasswordAndSalt, bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		newUserPasswordDetails = UserPasswordSchema{
			ID:        id,
			UserID:    userId,
			Hash:      string(hashedPassword),
			Salt:      salt,
			Algorithm: "bcrypt",
			UpdatedAt: currentTime,
		}

		_, err = userPasswordsCollection.InsertOne(ctx, newUserPasswordDetails)

		return err
	}

	combinedPasswordAndSalt := append([]byte(newPassword), []byte(existingUserPasswordStored.Salt)...)
	hashedPassword, err := bcrypt.GenerateFromPassword(combinedPasswordAndSalt, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var updatingData primitive.D
	updatingData = append(updatingData, bson.E{Key: "hash", Value: string(hashedPassword)})
	updatingData = append(updatingData, bson.E{Key: "updated_at", Value: currentTime})

	upsert := false
	opt := options.UpdateOptions{Upsert: &upsert}
	_, err = userPasswordsCollection.UpdateOne(
		ctx,
		filter,
		bson.D{{Key: "$set", Value: updatingData}},
		&opt,
	)

	return err
}
