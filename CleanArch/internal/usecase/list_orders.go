package usecase

import (
	"github.com/devfullcycle/20-CleanArch/internal/entity"
	"github.com/devfullcycle/20-CleanArch/pkg/events"
)
type OrdersList struct {
	Orders []entity.Order `json:"orders"`
}

type Order struct {
    ID         string  `json:"id"`
    Price      float64 `json:"price"`
    Tax        float64 `json:"tax"`
    FinalPrice float64 `json:"final_price"`
}
type ListAllOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	OrderListed events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}
func NewListAllOrdersUseCase(
	OrderRepository entity.OrderRepositoryInterface,
	OrderListed events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *ListAllOrdersUseCase {
	return &ListAllOrdersUseCase{
		OrderRepository: OrderRepository,
		OrderListed: OrderListed,
		EventDispatcher: EventDispatcher,
	}
}


func (c *ListAllOrdersUseCase) Execute() (OrdersList, error) {
	orders, err := c.OrderRepository.FindAll()
	if err != nil {
		return OrdersList{}, err
	}

	dto := OrdersList{
		Orders: orders,
	}
	c.OrderListed.SetPayload(dto)
	c.EventDispatcher.Dispatch(c.OrderListed)
	return dto, nil

}
