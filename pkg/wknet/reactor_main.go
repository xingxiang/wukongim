package wknet

import "github.com/WuKongIM/WuKongIM/pkg/wklog"

type ReactorMain struct {
	acceptor *acceptor
	eg       *Engine
	wklog.Log
}

func NewReactorMain(eg *Engine) *ReactorMain {

	return &ReactorMain{
		acceptor: newAcceptor(eg),
		eg:       eg,
		Log:      wklog.NewWKLog("ReactorMain"),
	}
}

func (m *ReactorMain) Start() error {
	return m.acceptor.Start()
}

func (m *ReactorMain) Stop() error {
	return m.acceptor.Stop()
}
