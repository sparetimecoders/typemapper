package typemapper

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_IllegalMapInput(t *testing.T) {
	_, err := NewFromMap(map[string]any{
		"one": one{},
		"two": one{},
	})
	require.Error(t, err)
}

func Test_Add(t *testing.T) {
	mapper, err := NewFromMap(map[string]any{
		"one": one{},
	})
	require.NoError(t, err)
	r, b := mapper.Type("one")
	require.IsType(t, reflect.TypeOf(one{}), r)
	require.True(t, b)

	require.NoError(t, mapper.Add("one", one{}))
	require.EqualError(t, mapper.Add("one", two{}), "mapping for key 'one' already registered to type 'typemapper.one'")

	require.NoError(t, mapper.Add("two", two{}))

	require.EqualError(t, mapper.Add("three", two{}), "mapping for type 'typemapper.two' already registered to key 'two'")
}

func Test_Delete(t *testing.T) {
	mapper, err := NewFromMap(map[string]any{
		"one": one{},
		"two": two{},
	})
	require.NoError(t, err)

	require.ElementsMatch(t, mapper.Keys(), []string{"one", "two"})
	mapper.DeleteKey("one")
	r, b := mapper.Type("one")
	require.Nil(t, r)
	require.False(t, b)
	require.ElementsMatch(t, mapper.Keys(), []string{"two"})

	mapper.DeleteKey("three")
	require.ElementsMatch(t, mapper.Keys(), []string{"two"})
}

func Test_Reflect(t *testing.T) {
	mapper, err := NewFromMap(map[string]any{
		"one": reflect.TypeOf(one{}),
		"two": two{},
	})
	require.NoError(t, err)
	key, _ := mapper.Key(two{})
	require.Equal(t, "two", key)
	key, _ = mapper.Key(one{})
	require.Equal(t, "one", key)
	key, _ = mapper.Key(reflect.TypeOf(one{}))
	require.Equal(t, "one", key)

	typ, _ := mapper.Type("one")
	require.Equal(t, reflect.TypeOf(one{}), typ)
}

type one struct{}
type two struct{}
