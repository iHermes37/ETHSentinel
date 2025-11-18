package mempool

type Monitor interface {
	CollectPendingTx()
	MonitorWhaleRefTx()
}

type Radar struct {
}

func (r *Radar) CollectPendingTx() {

}

func (r *Radar) MonitorWhaleRefTx() {

}
