package main

import (
	"context"
	"log"
	"time"

	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/labs-auction-goexpert/configuration/database/mongodb"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/labs-auction-goexpert/configuration/logger"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/labs-auction-goexpert/internal/entity/auction_entity"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/labs-auction-goexpert/internal/infra/api/web/controller/auction_controller"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/labs-auction-goexpert/internal/infra/api/web/controller/bid_controller"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/labs-auction-goexpert/internal/infra/api/web/controller/user_controller"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/labs-auction-goexpert/internal/infra/database/auction"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/labs-auction-goexpert/internal/infra/database/bid"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/labs-auction-goexpert/internal/infra/database/user"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/labs-auction-goexpert/internal/usecase/auction_usecase"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/labs-auction-goexpert/internal/usecase/bid_usecase"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/labs-auction-goexpert/internal/usecase/user_usecase"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := godotenv.Load("cmd/auction/.env"); err != nil {
		log.Fatal("Error trying to load env variables")
		return
	}

	databaseConnection, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	router := gin.Default()

	userController, bidController, auctionsController, auctionRepository := initDependencies(databaseConnection)

	go startAuctionsCloser(ctx, auctionRepository, 1*time.Minute)

	router.GET("/auction", auctionsController.FindAuctions)
	router.GET("/auction/:auctionId", auctionsController.FindAuctionById)
	router.POST("/auction", auctionsController.CreateAuction)
	router.GET("/auction/winner/:auctionId", auctionsController.FindWinningBidByAuctionId)
	router.POST("/bid", bidController.CreateBid)
	router.GET("/bid/:auctionId", bidController.FindBidByAuctionId)
	router.GET("/user/:userId", userController.FindUserById)

	router.Run(":8080")
}

func initDependencies(database *mongo.Database) (
	userController *user_controller.UserController,
	bidController *bid_controller.BidController,
	auctionController *auction_controller.AuctionController,
	auctionRepository *auction.AuctionRepository) {
	auctionRepository = auction.NewAuctionRepository(database)
	bidRepository := bid.NewBidRepository(database, auctionRepository)
	userRepository := user.NewUserRepository(database)

	userController = user_controller.NewUserController(
		user_usecase.NewUserUseCase(userRepository))

	auctionController = auction_controller.NewAuctionController(
		auction_usecase.NewAuctionUseCase(auctionRepository, bidRepository))

	bidController = bid_controller.NewBidController(bid_usecase.NewBidUseCase(bidRepository))

	return
}

func startAuctionsCloser(ctx context.Context, auctionRepository auction_entity.AuctionRepositoryInterface, tickerDuration time.Duration) {
	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := auctionRepository.CloseExpiredAuctions(ctx)
			if err != nil {
				logger.Error("Error processing expired auctions: %v", err)
			} else {
				logger.Info("Expired auctions processed successfully")
			}
		case <-ctx.Done():
			logger.Info("Auction closing processor closed.")
			return
		}
	}
}
