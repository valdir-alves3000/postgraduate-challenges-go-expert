package auction

import (
	"context"
	"time"

	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/labs-auction-goexpert/configuration/logger"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/labs-auction-goexpert/internal/entity/auction_entity"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/labs-auction-goexpert/internal/internal_error"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func (ar *AuctionRepository) CloseExpiredAuctions(ctx context.Context) *internal_error.InternalError {
	now := time.Now().Unix()

	filter := bson.M{
		"expires_at": bson.M{"$lte": now},
		"status":     auction_entity.Active,
	}

	update := bson.M{
		"$set": bson.M{"status": auction_entity.Completed},
	}

	result, err := ar.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		logger.Error("Error closing expired auctions", err)
		return internal_error.NewInternalServerError("Error closing expired auctions")
	}

	logger.Info("Automatically closed auctions", zap.Int64("modified_count", result.ModifiedCount))
	return nil
}
