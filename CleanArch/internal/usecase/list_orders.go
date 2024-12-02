package usecase

import "github.com/devfullcycle/20-CleanArch/internal/entity"

//retur all orders
type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrdersUseCase(
	OrderRepository entity.OrderRepositoryInterface) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: OrderRepository,
	}
}

func (c *ListOrdersUseCase) Execute() ([]entity.Order, error) {
    orders, err := c.OrderRepository.ListAll()
    if err != nil {
        return nil, err
    }
    return orders, nil

}

