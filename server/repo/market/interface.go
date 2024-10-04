package market

type Market interface {
	GetPrice(symbol string) (int64, error)
	Subscribe(symbol string, priceChan chan int64) error
}
