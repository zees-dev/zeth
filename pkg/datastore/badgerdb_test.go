package datastore

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mustNewTestBadgerDB() (*badgerStore, func() error) {
	store, err := NewBadgerDB("")
	if err != nil {
		panic(err)
	}
	return store, store.Close
}

func Test_SatisfiesStoreInterface(t *testing.T) {
	is := assert.New(t)

	store, err := NewBadgerDB("")
	is.NoError(err)
	defer store.Close()

	is.Implements((*Store)(nil), store)
}

func Test_badgerNamespaceKey(t *testing.T) {
	is := assert.New(t)

	tt := []struct {
		namespace []byte
		key       []byte
		want      []byte
	}{
		{
			namespace: []byte("settings"),
			key:       []byte("1"),
			want:      []byte("settings/1"),
		},
		{
			namespace: []byte(""),
			key:       []byte("1"),
			want:      []byte("/1"),
		},
		{
			namespace: nil,
			key:       []byte("1"),
			want:      []byte("/1"),
		},
	}

	for _, test := range tt {
		testName := string(test.namespace) + "/" + string(test.key)
		t.Run(testName, func(t *testing.T) {
			got := badgerNamespaceKey(test.namespace, test.key)
			is.Equal(got, test.want, fmt.Sprintf("want=%s, got=%s", test.want, got))
		})
	}
}

func Test_Has(t *testing.T) {
	is := assert.New(t)

	store, cleanup := mustNewTestBadgerDB()
	defer cleanup()

	// create 3 entries - 2 of the same prefix
	err := store.Set([]byte("settings"), []byte("key"), []byte("val"))
	is.NoError(err)

	err = store.Set([]byte("settings"), []byte("key2"), []byte("val2"))
	is.NoError(err)

	err = store.Set([]byte("settings"), []byte("differentPrefix"), []byte("val2"))
	is.NoError(err)

	// remove all entries with `key` prefix
	err = store.RemovePrefix([]byte("settings"), []byte("key"))
	is.NoError(err)

	// check remaining entries
	ok, err := store.Has([]byte("settings"), []byte("key"))
	is.NoError(err)
	is.False(ok)

	ok, err = store.Has([]byte("settings"), []byte("key2"))
	is.NoError(err)
	is.False(ok)

	ok, err = store.Has([]byte("settings"), []byte("differentPrefix"))
	is.NoError(err)
	is.True(ok)
}

func Test_SetGetHas(t *testing.T) {
	is := assert.New(t)

	store, cleanup := mustNewTestBadgerDB()
	defer cleanup()

	tt := []struct {
		namespace []byte
		key       []byte
		value     []byte
	}{
		{
			namespace: []byte("settings"),
			key:       []byte("key"),
			value:     []byte("val"),
		},
		{
			namespace: []byte(""),
			key:       []byte(""),
			value:     []byte(""),
		},
		{
			namespace: nil,
			key:       []byte(""),
			value:     []byte(""),
		},
		{
			namespace: nil,
			key:       nil,
			value:     []byte(""),
		},
		// value cannot be nil
	}

	for _, test := range tt {
		store.Set(test.namespace, test.key, test.value)

		ok, err := store.Has(test.namespace, test.key)
		is.NoError(err)
		is.True(ok)

		v, err := store.Get(test.namespace, test.key)
		is.NoError(err)
		is.Equal(v, test.value)
	}
}

func Test_GetAll(t *testing.T) {
	is := assert.New(t)

	store, cleanup := mustNewTestBadgerDB()
	defer cleanup()

	type entry struct {
		namespace []byte
		key       []byte
		value     []byte
	}

	tt := []struct {
		namespace []byte
		entries   []entry
		want      map[string][]byte
	}{
		{
			namespace: []byte("settings"),
			entries: []entry{
				{
					namespace: []byte("settings"),
					key:       []byte("key"),
					value:     []byte("val"),
				},
				{
					namespace: []byte("settings"),
					key:       []byte("1"),
					value:     []byte("2"),
				},
				{
					namespace: []byte("test"),
					key:       []byte("1"),
					value:     []byte("2"),
				},
			},
			want: map[string][]byte{
				"key": []byte("val"),
				"1":   []byte("2"),
			},
		},
		{
			namespace: []byte("testNew"),
			entries: []entry{
				{
					namespace: []byte("testNew"),
					key:       []byte("key"),
					value:     []byte("val"),
				},
				{
					namespace: []byte("testNew"),
					key:       []byte(""),
					value:     []byte(""),
				},
				{ // this will override above
					namespace: []byte("testNew"),
					key:       nil,
					value:     []byte(""),
				},
				{ // this will override above
					namespace: []byte("testNew"),
					key:       nil,
					value:     nil,
				},
			},
			want: map[string][]byte{
				"key": []byte("val"),
				"":    []byte(""),
			},
		},
	}

	for _, test := range tt {
		for _, entry := range test.entries {
			err := store.Set(entry.namespace, entry.key, entry.value)
			is.NoError(err)
		}

		got, err := store.GetAll(test.namespace)
		is.NoError(err)
		is.Equal(test.want, got)
	}
}
