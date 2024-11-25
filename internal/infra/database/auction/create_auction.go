package auction

import (
	"context"
	"errors"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	go closeAuctionRoutine(ctx, auctionEntity.Timestamp.Add(getAuctionDuration()), auctionEntity.Id, ar)

	return nil
}

func (ar *AuctionRepository) CloseAuction(
	ctx context.Context,
	auctionId string) error {
	_, err := ar.Collection.UpdateOne(ctx, bson.M{"_id": auctionId}, bson.M{"$set": bson.M{"status": auction_entity.Completed}})
	if err != nil {
		logger.Error("Error trying to close auction", err)
		return errors.New("error trying to close auction")
	}

	return nil
}

type Repository interface {
	CloseAuction(ctx context.Context, auctionId string) error
}

func closeAuctionRoutine(ctx context.Context, closeTime time.Time, auctionId string, repository Repository) {
	select {
	case <-time.After(time.Until(closeTime)):
		err := repository.CloseAuction(ctx, auctionId)
		if err != nil {
			logger.Error("Error trying to close auction", err)
			return
		}
		logger.Info("Auction closed: ", zap.String("auctionId", auctionId))
	case <-ctx.Done():
		logger.Info("Context cancelled, auction not closed", zap.String("auctionId", auctionId))
		return
	}
}

func getAuctionDuration() time.Duration {
	auctionDuration := os.Getenv("AUCTION_DURATION")
	duration, err := time.ParseDuration(auctionDuration)
	if err != nil {
		return time.Minute * 5
	}

	return duration
}
