package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const emptyUUIDV4 = "00000000-0000-0000-0000-000000000000"

func Test_SatisfiesSettingsInterface(t *testing.T) {
	is := assert.New(t)
	is.Implements((*Settings)(nil), NewService(nil))
}

// func Test_SettingsCrud(t *testing.T) {
// 	is := assert.New(t)

// 	store, cleanup := badgerdbtest.MustNewTestBadgerDB()
// 	defer cleanup()

// 	settingsStore := NewStore(store)

// 	t.Run("creating setting assigns UUID", func(t *testing.T) {
// 		s := Setting{}
// 		is.Equal(s.ID.String(), emptyUUIDV4)

// 		s2, err := settingsStore.Create(context.Background(), &s)
// 		is.NoError(err)

// 		is.NotEqual(s.ID.String(), emptyUUIDV4)
// 		is.NotEqual(s2.ID.String(), emptyUUIDV4)

// 		store.DropAll()
// 	})

// 	t.Run("created setting can be retrieved via UUID", func(t *testing.T) {
// 		s := Setting{}
// 		_, err := settingsStore.Create(context.Background(), &s)
// 		is.NoError(err)

// 		s2, err := settingsStore.Get(context.Background())
// 		is.NoError(err)

// 		is.Equal(s, s2)

// 		store.DropAll()
// 	})

// 	t.Run("updating setting can be retrieved via UUID", func(t *testing.T) {
// 		s := Setting{}
// 		_, err := settingsStore.Create(context.Background(), &s)
// 		is.NoError(err)

// 		// updating existing key with new setting which has empty UUID (dont do this in code)
// 		err = settingsStore.Update(context.Background(), Setting{})
// 		is.NoError(err)

// 		s2, err := settingsStore.Get(context.Background())
// 		is.NoError(err)

// 		is.NotEqual(s.ID, s2.ID)

// 		store.DropAll()
// 	})
// }
