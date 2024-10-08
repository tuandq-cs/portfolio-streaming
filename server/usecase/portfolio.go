package usecase

import (
	"context"
	"log"
	"portfolio/server/model"
	"portfolio/server/repo/market"
)

type StreamPortUC interface {
	GetNAV(ctx context.Context, port model.Portfolio, navCh chan int64) error
}

type streamPortUC struct {
	Market market.Market
}

func NewStreamPortUC() StreamPortUC {
	return &streamPortUC{}
}

type MarketChan map[string]chan int64

func (uc *streamPortUC) GetNAV(ctx context.Context, port model.Portfolio, navCh chan int64) error {
	portCh, err := uc.Market.Register(ctx, port)
	if err != nil {
		return err
	}
	navCh <- port.GetNAV()
	// listen change
	for stream := range portCh {
		if err := port.UpdatePrice(stream.Symbol, stream.Price); err != nil {
			// TODO: unsubscribe
			log.Print(err)
		}
		navCh <- port.GetNAV()
	}
	return nil
}
