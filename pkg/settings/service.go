package settings

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/zees-dev/zeth/pkg/datastore"
)

var (
	settingsKey = []byte("settings")
)

type settingsService struct {
	bs datastore.Store
}

func NewService(bs datastore.Store) *settingsService {
	return &settingsService{
		bs: bs,
	}
}

func (ss *settingsService) Get(ctx context.Context) (Setting, error) {
	dbSetting, err := ss.bs.GetGlobal(settingsKey)
	if err != nil {
		return Setting{}, err
	}

	var setting Setting
	if err := json.Unmarshal(dbSetting, &setting); err != nil {
		return Setting{}, err
	}

	return setting, nil
}

func (ss *settingsService) Update(ctx context.Context, setting Setting) error {
	bodyBytes := new(bytes.Buffer)
	json.NewEncoder(bodyBytes).Encode(setting)
	return ss.bs.SetGlobal(settingsKey, bodyBytes.Bytes())
}
