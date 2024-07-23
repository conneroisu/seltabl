package safe

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewSafeMap tests the NewSafeMap function.
func TestNewSafeMap(t *testing.T) {
	sm := NewSafeMap[string, int]()
	assert.NotNil(t, sm)
}

// TestSetAndGet tests the SafeMap's set and get methods.
func TestSetAndGet(t *testing.T) {
	sm := NewSafeMap[string, int]()
	sm.Set("key1", 10)
	val, ok := sm.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, 10, *val)
}

// TestGetNonExistentKey tests the SafeMap's get method for a non-existent key.
func TestGetNonExistentKey(t *testing.T) {
	sm := NewSafeMap[string, int]()
	val, ok := sm.Get("nonexistent")
	assert.False(t, ok)
	assert.Equal(t, 0, *val)
}

// TestOverwriteValue tests the SafeMap's overwrite behavior.
func TestOverwriteValue(t *testing.T) {
	sm := NewSafeMap[string, int]()
	sm.Set("key1", 10)
	sm.Set("key1", 20)
	val, ok := sm.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, 20, *val)
}

// TestDeleteKey tests the SafeMap's delete method.
func TestDeleteKey(t *testing.T) {
	sm := NewSafeMap[string, int]()
	sm.Set("key1", 10)
	sm.Delete("key1")
	val, ok := sm.Get("key1")
	assert.False(t, ok)
	assert.Equal(t, 0, *val)
}

// TestLen tests the SafeMap's length method.
func TestLen(t *testing.T) {
	sm := NewSafeMap[string, int]()
	sm.Set("key1", 10)
	sm.Set("key2", 20)
	assert.Equal(t, 2, sm.Len())
	sm.Delete("key1")
	assert.Equal(t, 1, sm.Len())
}

// TestConcurrentSetAndGet tests the SafeMap's concurrent set and get methods.
func TestConcurrentSetAndGet(t *testing.T) {
	sm := NewSafeMap[int, int]()
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sm.Set(i, i)
		}(i)
	}
	wg.Wait()
	for i := 0; i < 1000; i++ {
		val, ok := sm.Get(i)
		assert.True(t, ok)
		assert.Equal(t, i, *val)
	}
}

// TestSetDifferentTypes tests the SafeMap's set method with different types.
func TestSetDifferentTypes(t *testing.T) {
	smInt := NewSafeMap[string, int]()
	smInt.Set("int", 123)
	valInt, okInt := smInt.Get("int")
	assert.True(t, okInt)
	assert.Equal(t, 123, *valInt)
	smStr := NewSafeMap[string, string]()
	smStr.Set("str", "value")
	valStr, okStr := smStr.Get("str")
	assert.True(t, okStr)
	assert.Equal(t, "value", *valStr)
}

// TestSetNilValue tests the SafeMap's set method with a nil value.
func TestSetNilValue(t *testing.T) {
	sm := NewSafeMap[string, *int]()
	sm.Set("key", nil)
	val, ok := sm.Get("key")
	assert.True(t, ok)
	assert.Nil(t, *val)
}
func TestOverwriteWithNilValue(t *testing.T) {
	sm := NewSafeMap[string, *int]()
	value := 10
	sm.Set("key", &value)
	sm.Set("key", nil)
	val, ok := sm.Get("key")
	assert.True(t, ok)
	assert.Nil(t, *val)
}

// TestLenAfterDelete tests the SafeMap's length method after a delete.
func TestLenAfterDelete(t *testing.T) {
	sm := NewSafeMap[string, int]()
	sm.Set("key1", 10)
	sm.Set("key2", 20)
	sm.Delete("key1")
	assert.Equal(t, 1, sm.Len())
	sm.Delete("key2")
	assert.Equal(t, 0, sm.Len())
}

// TestOverwriteAfterDelete tests the SafeMap's overwrite behavior after a
// delete.
func TestOverwriteAfterDelete(t *testing.T) {
	sm := NewSafeMap[string, int]()
	sm.Set("key1", 10)
	sm.Delete("key1")
	sm.Set("key1", 20)
	val, ok := sm.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, 20, *val)
}

// TestDeleteNonExistentKey tests the SafeMap's delete method for a non-existent
// key.
func TestDeleteNonExistentKey(t *testing.T) {
	sm := NewSafeMap[string, int]()
	sm.Delete("nonexistent")
	val, ok := sm.Get("nonexistent")
	assert.False(t, ok)
	assert.Equal(t, 0, *val)
}

// TestConcurrentLen tests the SafeMap's concurrent length method.
func TestConcurrentLen(t *testing.T) {
	sm := NewSafeMap[int, int]()
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sm.Set(i, i)
		}(i)
	}
	wg.Wait()
	assert.Equal(t, 1000, sm.Len())
}

// TestNewSafeMap2 tests the NewSafeMap function with a non-pointer type.
func TestNewSafeMap2(t *testing.T) {
	sm := NewSafeMap[string, int]()
	if sm == nil {
		t.Fatal("NewSafeMap returned nil")
	}
	if sm.m == nil {
		t.Fatal("map is nil")
	}
}

// TestSafeMap_GetSet tests the SafeMap's set and get methods.
func TestSafeMap_GetSet(t *testing.T) {
	sm := NewSafeMap[string, int]()
	// Test Set and Get
	sm.Set("key1", 10)
	val, ok := sm.Get("key1")
	if !ok || *val != 10 {
		t.Errorf("Expected (10, true), got (%v, %v)", val, ok)
	}
	// Test Get for non-existent key
	val, ok = sm.Get("key2")
	if ok || *val != 0 {
		t.Errorf("Expected (0, false), got (%v, %v)", val, ok)
	}
}

// TestSafeMap_Delete tests the SafeMap's delete method.
func TestSafeMap_Delete(t *testing.T) {
	sm := NewSafeMap[string, int]()
	sm.Set("key1", 10)
	sm.Delete("key1")
	_, ok := sm.Get("key1")
	if ok {
		t.Error("Key should have been deleted")
	}
}

// TestSafeMap_Len tests the SafeMap's length method.
func TestSafeMap_Len(t *testing.T) {
	sm := NewSafeMap[string, int]()
	if sm.Len() != 0 {
		t.Errorf("Expected length 0, got %d", sm.Len())
	}
	sm.Set("key1", 10)
	sm.Set("key2", 20)
	if sm.Len() != 2 {
		t.Errorf("Expected length 2, got %d", sm.Len())
	}
}

// TestSafeMap_Clear tests the SafeMap's clear method.
func TestSafeMap_Clear(t *testing.T) {
	sm := NewSafeMap[string, int]()
	sm.Set("key1", 10)
	sm.Set("key2", 20)
	sm.Clear()
	if sm.Len() != 0 {
		t.Errorf("Expected length 0 after Clear, got %d", sm.Len())
	}
}

// TestSafeMap_String tests the SafeMap's string representation
func TestSafeMap_String(t *testing.T) {
	sm := NewSafeMap[string, int]()
	sm.Set("key1", 10)
	sm.Set("key2", 20)
	str := sm.String()
	expected := "key1: 10\nkey2: 20\n"
	if str != expected && str != "key2: 20\nkey1: 10\n" {
		t.Errorf("Unexpected string representation:\n%s", str)
	}
}

// TestSafeMap_Concurrent tests the SafeMap's concurrent behavior.
func TestSafeMap_Concurrent(t *testing.T) {
	sm := NewSafeMap[int, int]()
	var wg sync.WaitGroup
	numGoroutines := 100
	numOperations := 1000
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := (id*numOperations + j) % (numGoroutines * numOperations)
				sm.Set(key, id)
				sm.Get(key)
				if j%2 == 0 {
					sm.Delete(key)
				}
			}
		}(i)
	}
	wg.Wait()
	// Check if the final length is correct
	expectedLen := numGoroutines * numOperations / 2
	if sm.Len() > expectedLen {
		t.Errorf("Expected length <= %d, got %d", expectedLen, sm.Len())
	}
}
