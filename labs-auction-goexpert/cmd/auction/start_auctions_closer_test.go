package main

import (
	"context"
	"testing"
	"time"

	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/labs-auction-goexpert/internal/internal_error"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/labs-auction-goexpert/tests/mocks"

	"github.com/golang/mock/gomock"
)

func TestStartAuctionsCloser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuctionRepositoryInterface(ctrl)

	tickerDuration := 50 * time.Millisecond

	t.Run("should process expired auctions successfully", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockRepo.EXPECT().CloseExpiredAuctions(gomock.Any()).Return(nil).Times(1)

		go startAuctionsCloser(ctx, mockRepo, tickerDuration)

		time.Sleep(tickerDuration + 10*time.Millisecond)
	})

	t.Run("should return error when processing expired auctions", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		expectedError := internal_error.NewInternalServerError("error processing expired auctions")
		mockRepo.EXPECT().CloseExpiredAuctions(gomock.Any()).Return(expectedError).Times(1)

		go startAuctionsCloser(ctx, mockRepo, tickerDuration)

		time.Sleep(10*time.Millisecond + tickerDuration)
	})

	t.Run("should stop processing when context is cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan bool)
		go func() {
			startAuctionsCloser(ctx, mockRepo, tickerDuration)
			done <- true
		}()
		cancel()

		select {
		case <-done:
			// Success - processor stopped
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Processor did not stop after context cancellation")
		}
	})

	t.Run("should maintain periodic execution", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockRepo.EXPECT().
			CloseExpiredAuctions(gomock.Any()).
			Return(nil).
			Times(2)

		go startAuctionsCloser(ctx, mockRepo, tickerDuration)

		time.Sleep(2*tickerDuration + 10*time.Millisecond)
	})
}
