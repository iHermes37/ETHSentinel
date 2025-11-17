package detective

import "math/big"

type Interval struct {
	StartBlock *big.Int
	EndBlock   *big.Int
}

type Detective struct {
	Interval
}
