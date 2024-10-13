package service

import (
	"context"

	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/clean-architecture/internal/infra/grpc/pb"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/clean-architecture/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrdersUseCase  usecase.ListOrdersUseCase
}

func NewOrderService(createOrderUsecase usecase.CreateOrderUseCase, listOrdersUseCase usecase.ListOrdersUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUsecase,
		ListOrdersUseCase:  listOrdersUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}

	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}

	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, in *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := s.ListOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}

	var outputOrders []*pb.CreateOrderResponse
	for _, output := range orders {
		outputOrders = append(outputOrders, &pb.CreateOrderResponse{
			Id:         output.ID,
			Price:      float32(output.Price),
			Tax:        float32(output.Tax),
			FinalPrice: float32(output.FinalPrice),
		})
	}

	return &pb.ListOrdersResponse{
		Orders: outputOrders,
	}, nil

}
