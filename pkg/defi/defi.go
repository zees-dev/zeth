package defi

import (
	"context"

	"github.com/zees-dev/zeth/pkg/defi/amm"
)

type AutomatedMarketMaker interface {
	GetByChainID(ctx context.Context, id int) ([]amm.AMM, error)
	GetAll(ctx context.Context) ([]amm.AMM, error)
	Create(ctx context.Context, amm amm.AMM) (amm.AMM, error)
}
