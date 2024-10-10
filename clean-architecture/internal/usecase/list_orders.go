package usecase

import "github.com/valdir-alves3000/postgraduate-challenges-go-expert/clean-architecture/internal/entity"

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrdersUseCase(
	OrderRepository entity.OrderRepositoryInterface,
) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: OrderRepository,
	}
}

func (c *ListOrdersUseCase) Execute() ([]OrderOutputDTO, error) {

	orders, err := c.OrderRepository.FindAll()
	if err != nil {
		return nil, err
	}

	var ordersOutput []OrderOutputDTO

	for _, order := range orders {
		dto := OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		}
		ordersOutput = append(ordersOutput, dto)
	}
	return ordersOutput, nil
}
