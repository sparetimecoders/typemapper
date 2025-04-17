package typemapper

import (
	"fmt"
	"maps"
	"reflect"
	"slices"
	"sync"
)

func New() *Mapper {
	return &Mapper{
		typeToKey: make(map[reflect.Type]string),
		keyToType: make(map[string]reflect.Type),
	}
}

func NewFromMap(org map[string]any) (*Mapper, error) {
	m := New()
	for k, v := range org {
		if err := m.Add(k, v); err != nil {
			return nil, err
		}
	}
	return m, nil
}

type Mapper struct {
	s         sync.RWMutex
	typeToKey map[reflect.Type]string
	keyToType map[string]reflect.Type
}

func (m *Mapper) Add(key string, t any) error {
	if t == nil {
		return fmt.Errorf("cannot add nil type")
	}
	typ := reflect.TypeOf(t)
	if tt, ok := t.(reflect.Type); ok {
		typ = tt
	}
	m.s.Lock()
	defer m.s.Unlock()
	if registeredType, exists := m.keyToType[key]; exists && registeredType != typ {
		return fmt.Errorf("mapping for key '%s' already registered to type '%T'", key, reflect.New(registeredType).Elem().Interface())
	}
	if registeredKey, exists := m.typeToKey[typ]; exists && registeredKey != key {
		return fmt.Errorf("mapping for type '%T' already registered to key '%s'", t, registeredKey)
	}
	m.keyToType[key] = typ
	m.typeToKey[typ] = key
	return nil
}

func (m *Mapper) DeleteKey(key string) {
	typ, ok := m.Type(key)
	if !ok {
		return
	}
	m.s.Lock()
	defer m.s.Unlock()
	delete(m.keyToType, key)
	delete(m.typeToKey, typ)
}

func (m *Mapper) DeleteType(t any) {
	typ := reflect.TypeOf(t)
	key, ok := m.Key(typ)
	if !ok {
		return
	}
	m.s.Lock()
	defer m.s.Unlock()
	delete(m.keyToType, key)
	delete(m.typeToKey, typ)
}

func (m *Mapper) Type(key string) (reflect.Type, bool) {
	m.s.RLock()
	defer m.s.RUnlock()
	t, ok := m.keyToType[key]
	return t, ok
}

func (m *Mapper) Key(t any) (string, bool) {
	typ := reflect.TypeOf(t)
	key := typ
	if typ.Kind() == reflect.Ptr {
		key = typ.Elem()
	}
	if tt, ok := t.(reflect.Type); ok {
		key = tt
	}
	m.s.RLock()
	defer m.s.RUnlock()
	k, ok := m.typeToKey[key]
	return k, ok
}

func (m *Mapper) Keys() []string {
	m.s.RLock()
	defer m.s.RUnlock()
	return slices.Collect(maps.Keys(m.keyToType))
}

func (m *Mapper) Types() []reflect.Type {
	m.s.RLock()
	defer m.s.RUnlock()
	return slices.Collect(maps.Keys(m.typeToKey))
}

func (m *Mapper) KeysToTypes() map[string]reflect.Type {
	m.s.RLock()
	defer m.s.RUnlock()
	return maps.Clone(m.keyToType)
}

func (m *Mapper) TypesToKeys() map[reflect.Type]string {
	m.s.RLock()
	defer m.s.RUnlock()
	return maps.Clone(m.typeToKey)
}
