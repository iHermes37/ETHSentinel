package arbitrageur

//========================================================

type Arbitrageur struct {
	*Monitor
	CurStrategy string
}

func (a *Arbitrageur) init() {
	GetCommonPairs()
}

func (a *Arbitrageur) SetStrategy(strategy string) {
	a.CurStrategy = strategy
}

func (a *Arbitrageur) ExecArbitrage() {
	strategy := Strategies[a.CurStrategy]

	stop := make(chan struct{})
	stopChans := make(map[string]chan struct{})
	stopChans[a.CurStrategy] = stop

	go strategy.Run(stop)
	// 阻塞等待退出
	select {}
}
