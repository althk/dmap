package dmap

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type CustomVal struct {
	v string
}

var keyPrefixes = []string{"key", "otherkey", "oldkey", "keynew", "fookey"}
var valPrefix = "val"
var keys []string

var bm DMap[string, string] // used for benchmarking Get

func prepareTestData[V any](m DMap[string, V], nkeys int, testval V) {
	keys = make([]string, nkeys)

	for i := 0; i < nkeys; i++ {
		key := fmt.Sprintf("%s_%d", keyPrefixes[rand.Intn(len(keyPrefixes))], i)
		m.Set(key, testval)
		keys[i] = key
	}
}

func TestNew(t *testing.T) {
	m := New[string, string](10)
	require.NotNil(t, m)
	require.Equal(t, 10, len(m))
}

func TestSetGetWithStrKV(t *testing.T) {
	m := New[string, string](10)
	prepareTestData(m, 10000, "some val")

	for i := 0; i < 50; i++ {
		val, e := m.Get(keys[rand.Intn(len(keys))])
		require.Equal(t, "some val", val)
		require.True(t, e)
	}
}

func TestKeys(t *testing.T) {
	m := New[string, string](10)
	prepareTestData(m, 10000, "some val")

	got := m.Keys()
	require.ElementsMatch(t, got, keys)
}

func TestHas(t *testing.T) {
	m := New[string, string](10)
	prepareTestData(m, 10000, "some val")

	for i := 0; i < 50; i++ {
		ok := m.Has(keys[rand.Intn(len(keys))])
		require.True(t, ok)
	}
	ok := m.Has("nonexistentkey")
	require.False(t, ok)
}

func TestCount(t *testing.T) {
	m := New[string, string](10)
	prepareTestData(m, 10000, "some val")

	got := m.Count()
	require.EqualValues(t, 10000, got)
}

func BenchmarkSet(b *testing.B) {
	l := len(keyPrefixes)
	for i := 0; i < b.N; i++ {
		for _, nkeys := range []int{100000, 1000000} {
			b.Run(fmt.Sprintf("%d_keys", nkeys), func(_ *testing.B) {
				m := New[string, string](10)
				for i := 0; i < nkeys; i++ {
					key := fmt.Sprintf("%s_%d", keyPrefixes[i%l], i)
					m.Set(key, "some val")
				}

			})

		}

	}
}

func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bm.Get(keys[i])
	}
}

func BenchmarkKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bm.Keys()
	}
}

func BenchmarkCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bm.Count()
	}
}

func TestMain(m *testing.M) {
	rand.Seed(42)
	bm = New[string, string](10)
	prepareTestData(bm, 100000, "some val")

	os.Exit(m.Run())
}
