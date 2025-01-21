package auction

import (
	"context"
	"os"
	"time"

	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/labs-auction-goexpert/configuration/logger"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/labs-auction-goexpert/internal/entity/auction_entity"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/labs-auction-goexpert/internal/internal_error"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
	ExpiresAt   int64                           `bson:"expires_at"`
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

	expires := getAuctionInterval()

	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
		ExpiresAt:   time.Now().Add(time.Duration(expires)).Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	return nil
}

func getAuctionInterval() time.Duration {
	defaultDuration := time.Hour * 24
	auctionInterval := os.Getenv("AUCTION_INTERVAL")

	if auctionInterval == "" {
		logger.Info("AUCTION_INTERVAL environment variable is not set, using default value")
		return defaultDuration
	}

	duration, err := time.ParseDuration(auctionInterval)
	if err != nil {
		return defaultDuration
	}

	return duration
}
