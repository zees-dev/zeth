package amm

import (
	"bytes"
	"context"
	"encoding/json"

	uuid "github.com/satori/go.uuid"
	"github.com/zees-dev/zeth/pkg/datastore"
)

var (
	ammKey = []byte("solidity/amm")
)

type ammService struct {
	store datastore.Store
}

func NewService(store datastore.Store) *ammService {
	return &ammService{
		store,
	}
}

func (svc *ammService) Create(ctx context.Context, amm AMM) (AMM, error) {
	id := uuid.NewV4()

	// set ID on existing props
	amm.ID = id

	bodyBytes := new(bytes.Buffer)

	json.NewEncoder(bodyBytes).Encode(amm)

	return amm, svc.store.Set(ammKey, id.Bytes(), bodyBytes.Bytes())
}

func (svc *ammService) Get(ctx context.Context, id uuid.UUID) (*AMM, error) {
	dbAMM, err := svc.store.Get(ammKey, id.Bytes())
	if err != nil {
		return nil, err
	}

	amm, err := unmarshal(dbAMM)
	if err != nil {
		return nil, err
	}

	return amm, nil
}

func (svc *ammService) getFromStore(ctx context.Context) ([]AMM, error) {
	results := []AMM{}

	ammMap, err := svc.store.GetAll(ammKey)
	if err != nil {
		return nil, err
	}

	for _, s := range ammMap {
		amm, err := unmarshal(s)
		if err != nil {
			return nil, err
		}
		results = append(results, *amm)
	}

	return results, nil
}

func (svc *ammService) GetAll(ctx context.Context) ([]AMM, error) {
	// all AMMs are stored in the store
	amms, err := svc.getFromStore(ctx)
	if err != nil {
		return nil, err
	}

	// all AMMs supported by default
	for _, amm := range SupportedAMMs {
		amms = append(amms, amm...)
	}

	return amms, nil
}

func (svc *ammService) GetByChainID(ctx context.Context, id int) ([]AMM, error) {
	amms, err := svc.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	results := []AMM{}
	for _, amm := range amms {
		if amm.ChainID == uint(id) {
			results = append(results, amm)
		}
	}

	return results, nil
}

func (svc *ammService) Delete(ctx context.Context, id uuid.UUID) error {
	return svc.store.RemovePrefix(ammKey, id.Bytes())
}

func unmarshal(b []byte) (*AMM, error) {
	var amm AMM
	if err := json.Unmarshal(b, &amm); err != nil {
		return nil, err
	}
	return &amm, nil
}
