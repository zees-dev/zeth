package defi

import (
	"context"

	uuid "github.com/satori/go.uuid"
	"github.com/zees-dev/zeth/pkg/defi/amm"
)

type AutomatedMarketMaker interface {
	Get(context.Context, uuid.UUID) (*amm.AMM, error)
	GetByChainID(ctx context.Context, id int) ([]amm.AMM, error)
	GetAll(ctx context.Context) ([]amm.AMM, error)
	Create(ctx context.Context, amm amm.AMM) (amm.AMM, error)
	Delete(context.Context, uuid.UUID) error
}
