package gojobqueue

import (
	"errors"
	"testing"
	"time"
)

func TestQueue_AddJob(t *testing.T) {
	result := "foo"
	res := &result
	func1 := func() error {
		time.Sleep(2 * time.Second)
		*res += "bar"
		return nil
	}
	func2 := func() error {
		*res += "baz"
		return nil
	}
	func3 := func() error {
		*res += "blub"
		return nil
	}
	queue := make(Queue, 4)
	err1 := queue.AddJob(func1, nil)
	err2 := queue.AddJob(func2, nil)
	queue.StartWorking()
	err3 := queue.AddJob(func2, nil)
	func3()
	time.Sleep(4 * time.Second)

	if err1 != nil || err2 != nil || err3 != nil {
		t.Errorf("Failed adding jobs to queue with the following errors: %+v \n %+v \n %+v \n", err1, err2, err3)
	}
	if result != "fooblubbarbazbaz" {
		t.Errorf("Failed executing jobs in correct order, expected result to be fooblubbarbaz, got %s", result)
	}
}

func TestQueue_Close(t *testing.T) {
	queue := make(Queue, 4)
	queue.StartWorking()
	queue.Close()
	err := queue.AddJob(func() error {return nil}, nil)
	if err == nil {
		t.Errorf("Adding jobs to closed Queue was still possible.")
	}
}

func TestQueue_StartWorking(t *testing.T) {
	result := "foo"
	res := &result
	queue := make(Queue, 4)
	queue.StartWorking()
	queue.AddJob(func() error {
		time.Sleep(1 * time.Second)
		*res += "bar"
		return nil
	}, nil)
	queue.AddJob(func() error {
		return errors.New("blub")
	}, func(e error) {
		*res += e.Error()
		*res += "baz"
	})
	time.Sleep(2 * time.Second)
	if result != "foobarblubbaz" {
		t.Errorf("Expected result to be 'foobarblubbaz', got %s", result)
	}
}