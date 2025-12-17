package storage_test

import (
	"bedis/internal/storage"
	"bytes"
	"fmt"
	"sync"
	"testing"
)

func TestSetAndGet(t *testing.T) {
	s := storage.New(nil)
	err := s.Set("test", []byte("test"))
	if err != nil {
		t.Fatalf("unexpeted error in Set: %s", err)
	}

	val, err := s.Get("test")
	if err != nil {
		t.Fatalf("unexpeted error in Get: %s", err)
	}
	if string(val) != "test" {
		t.Fatalf("unexpeted value in Get: %s", val)
	}
}

func TestGetUnknown(t *testing.T) {
	s := storage.New(nil)
	val, err := s.Get("=-=-=-")
	if err == nil {
		t.Fatalf("unexpeted error in Get: %s", val)
	}
	if val != nil {
		t.Fatalf("unexpeted value in Get: %s", val)
	}
}

func TestGetReturnsCopy(t *testing.T) {
	s := storage.New(nil)
	original := []byte("original")
	err := s.Set("original", original)
	if err != nil {
		t.Fatalf("unexpeted error in Set: %s", err)
	}
	original[0] = 'X'

	val, err := s.Get("original")
	if err != nil {
		t.Fatalf("unexpeted error in Get: %s", err)
	}
	if bytes.Equal(original, val) {
		t.Fatalf("Get returned original exepted copy")
	}

}

func TestConcurrency(t *testing.T) {
	s := storage.New(nil)

	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", i)
			err := s.Set(key, []byte(key))
			if err != nil {
				fmt.Printf("unexpeted error in Ser: %s", err)
			}
			_, err = s.Get(key)
			if err != nil {
				fmt.Printf("unexpeted error in Get: %s", err)
			}
		}(i)
	}
	wg.Wait()
}
