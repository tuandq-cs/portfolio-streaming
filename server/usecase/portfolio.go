package usecase

import (
	"log"
	"portfolio/server/model"
	"portfolio/server/repo/market"
)

type StreamPortUC interface {
	GetNAV(model.Portfolio, chan int64) error
}

type streamPortUC struct {
	Market market.Market
}

func NewStreamPortUC() StreamPortUC {
	return &streamPortUC{}
}

type MarketChan map[string]chan int64

func (uc *streamPortUC) GetNAV(port model.Portfolio, navCh chan int64) error {
	marketCh, err := uc.SubscribeMarketChan(port)
	if err != nil {
		return err
	}
	// close channel
	defer func() {
		for _, ch := range marketCh {
			close(ch)
		}
	}()
	navCh <- port.GetNAV()
	// listen change
	for sym, priceCh := range marketCh {
		go func(sym string, priceCh chan int64) {
			for curPrice := range priceCh {
				if err := port.UpdatePrice(sym, curPrice); err != nil {
					// TODO: unsubscribe
					log.Print(err)
				}
				navCh <- port.GetNAV()
			}
		}(sym, priceCh)
	}
	return nil
}

func (uc *streamPortUC) SubscribeMarketChan(port model.Portfolio) (MarketChan, error) {
	marketCh := make(MarketChan)
	for _, pos := range port.GetPositions() {
		priceCh := make(chan int64)
		if err := uc.Market.Subscribe(pos.Symbol, priceCh); err != nil {
			log.Printf("Failed to subscribe: symbol:%v, %v", pos.Symbol, err)
			return nil, err
		}
		marketCh[pos.Symbol] = priceCh
	}
	return marketCh, nil
}
