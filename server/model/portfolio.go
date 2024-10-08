package model

import (
	"errors"
	"sync"
)

type Portfolio interface {
	GetNAV() int64
	GetPositions() map[string]*Position
	UpdatePrice(symbol string, price int64) error
}

type Stream struct {
	Symbol string
	Price  int64
}

type portfolio struct {
	mt        *sync.Mutex
	cash      int64
	positions map[string]*Position
	nav       *int64
}

func NewPortfolio(cash int64, positions []Position) Portfolio {
	mPos := make(map[string]*Position)
	for _, pos := range positions {
		mPos[pos.Symbol] = &pos
	}
	return &portfolio{
		cash:      cash,
		positions: mPos,
		mt:        &sync.Mutex{},
	}
}

func (port *portfolio) GetPositions() map[string]*Position {
	return port.positions
}

func (port *portfolio) GetNAV() int64 {
	defer port.mt.Unlock()
	port.mt.Lock()
	if port.nav == nil {
		var mVal int64 = 0
		for _, pos := range port.positions {
			mVal += pos.GetMarketValue()
		}
		nav := port.cash + mVal
		port.nav = &nav
	}
	return *port.nav
}

func (port *portfolio) UpdatePrice(symbol string, price int64) error {
	defer port.mt.Unlock()
	port.mt.Lock()
	pos, found := port.positions[symbol]
	if !found {
		return errors.New("not found")
	}
	oldVal := pos.GetMarketValue()
	port.positions[symbol].Price = price
	var dNav = pos.GetMarketValue() - oldVal
	newNav := port.GetNAV() + dNav
	port.nav = &newNav
	return nil
}

type Position struct {
	Symbol string
	Qty    int64
	Price  int64
}

func (pos *Position) GetMarketValue() int64 {
	return pos.Qty * pos.Price
}
