package market

import (
	"context"
	"portfolio/server/model"
	"sync"

	"github.com/google/uuid"
)

type instrument struct {
	Symbol    string
	Price     *int64
	listeners map[string]*listener
}

type listener struct {
	id   string
	ctx  context.Context
	port model.Portfolio
	ch   chan model.Stream
}

func newListener(ctx context.Context, port model.Portfolio, ch chan model.Stream) *listener {
	return &listener{
		id:   uuid.New().String(),
		ctx:  ctx,
		port: port,
		ch:   ch,
	}
}

type random struct {
	mt     *sync.Mutex
	ins    map[string]*instrument
	stream chan model.Stream
}

func NewRandom() Market {
	new := &random{
		mt:     &sync.Mutex{},
		ins:    make(map[string]*instrument),
		stream: make(chan model.Stream),
	}
	go new.generate()
	go new.broadcast()
	return new
}

func (r *random) Register(ctx context.Context, port model.Portfolio) (chan model.Stream, error) {
	newCh := make(chan model.Stream)
	// new listener
	lis := newListener(ctx, port, newCh)
	for symbol := range port.GetPositions() {
		r.mt.Lock()
		ins := r.getInstrument(symbol)
		// add new listener
		ins.listeners[lis.id] = lis
		r.mt.Unlock()
	}
	return newCh, nil
}

func (r *random) generate() {
	// TODO: implement
}

func (r *random) broadcast() {
	for stream := range r.stream {
		// think about optimize
		r.mt.Lock()
		ins := r.getInstrument(stream.Symbol)
		ins.Price = &stream.Price
		for _, lis := range ins.listeners {
			select {
			case lis.ch <- stream:
			case <-lis.ctx.Done(): // listeners are done
				delete(ins.listeners, lis.id)
			}
		}
		r.mt.Unlock()
	}
}

func (r *random) getInstrument(symbol string) *instrument {
	ins, found := r.ins[symbol]
	if !found {
		ins = &instrument{
			Symbol:    symbol,
			listeners: make(map[string]*listener),
		}
		r.ins[symbol] = ins
	}
	return ins
}
