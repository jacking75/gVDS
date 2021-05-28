package scommon

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
)

func TestDequeAppend(t *testing.T) {
	deque := NewDeque()
	sampleSize := 100

	// Append elements in the Deque and assert it does not fail
	for i := 0; i < sampleSize; i++ {
		var value string = strconv.Itoa(i)
		_, ok := deque.Append(value)

		_assert(
			t,
			ok == true,
			"deque.Append(%d) = %t; want %t", i, ok, true,
		)
	}

	_assert(
		t,
		deque.container.Len() == sampleSize,
		"deque.container.Len() = %d; want %d", deque.container.Len(), sampleSize,
	)

	_assert(
		t,
		deque.container.Front().Value == "0",
		"deque.container.Front().Value = %s; want %s", deque.container.Front().Value, "0",
	)

	_assert(
		t,
		deque.container.Back().Value == "99",
		"deque.container.Back().Value = %s; want %s", deque.container.Back().Value, "99",
	)
}

func TestDequeAppendWithCapacity(t *testing.T) {
	dequeSize := 20
	deque := NewCappedDeque(dequeSize)

	// Append the maximum number of elements in the Deque
	// and assert it does not fail
	for i := 0; i < dequeSize; i++ {
		var value string = strconv.Itoa(i)
		_, ok := deque.Append(value)

		_assert(
			t,
			ok == true,
			"deque.Append(%d) = %t; want %t", i, ok, true,
		)
	}

	// Try to overflow the Deque size limit, and make
	// sure appending fails
	_, ok := deque.Append("should not be ok")
	_assert(
		t,
		ok == false,
		"deque.Append(%s) = %t; want %t", "should not be ok", ok, false,
	)

	_assert(
		t,
		deque.container.Len() == dequeSize,
		"deque.container.Len() = %d; want %d", deque.container.Len(), dequeSize,
	)

	_assert(
		t,
		deque.container.Front().Value == "0",
		"deque.container.Front().Value = %s; want %s", deque.container.Front().Value, "0",
	)

	_assert(
		t,
		deque.container.Back().Value == "19",
		"deque.container.Back().Value = %s; want %s", deque.container.Back().Value, "19",
	)
}

func TestDequePrepend(t *testing.T) {
	deque := NewDeque()
	sampleSize := 100

	// Prepend elements in the Deque and assert it does not fail
	for i := 0; i < sampleSize; i++ {
		var value string = strconv.Itoa(i)
		_, ok := deque.Prepend(value)

		_assert(
			t,
			ok == true,
			"deque.Prepend(%d) = %t; want %t", i, ok, true,
		)
	}

	_assert(
		t,
		deque.container.Len() == sampleSize,
		"deque.container.Len() = %d; want %d", deque.container.Len(), sampleSize,
	)

	_assert(
		t,
		deque.container.Front().Value == "99",
		"deque.container.Front().Value = %s; want %s", deque.container.Front().Value, "99",
	)

	_assert(
		t,
		deque.container.Back().Value == "0",
		"deque.container.Back().Value = %s; want %s", deque.container.Back().Value, "0",
	)
}

func TestDequePrependWithCapacity(t *testing.T) {
	dequeSize := 20
	deque := NewCappedDeque(dequeSize)

	// Prepend elements in the Deque and assert it does not fail
	for i := 0; i < dequeSize; i++ {
		var value string = strconv.Itoa(i)
		_, ok := deque.Prepend(value)

		_assert(
			t,
			ok == true,
			"deque.Prepend(%d) = %t; want %t", i, ok, true,
		)
	}

	// Try to overflow the Deque size limit, and make
	// sure appending fails
	_, ok := deque.Prepend("should not be ok")
	_assert(
		t,
		ok == false,
		"deque.Prepend(%s) = %t; want %t", "should not be ok", ok, false,
	)

	_assert(
		t,
		deque.container.Len() == dequeSize,
		"deque.container.Len() = %d; want %d", deque.container.Len(), dequeSize,
	)

	_assert(
		t,
		deque.container.Front().Value == "19",
		"deque.container.Front().Value = %s; want %s", deque.container.Front().Value, "19",
	)

	_assert(
		t,
		deque.container.Back().Value == "0",
		"deque.container.Back().Value = %s; want %s", deque.container.Back().Value, "0",
	)
}

func TestDequePop_fulfilled_container(t *testing.T) {
	deque := NewDeque()
	dequeSize := 100

	// Populate the test deque
	for i := 0; i < dequeSize; i++ {
		var value string = strconv.Itoa(i)
		deque.Append(value)
	}

	// Pop elements of the deque and assert elements come out
	// in order and container size is updated accordingly
	for i := dequeSize - 1; i >= 0; i-- {
		item := deque.Pop()

		var itemValue string = item.(string)
		var expectedValue string = strconv.Itoa(i)

		_assert(
			t,
			itemValue == expectedValue,
			"deque.Pop() = %s; want %s", itemValue, expectedValue,
		)

		_assert(
			t,
			deque.container.Len() == i,
			"deque.container.Len() = %d; want %d", deque.container.Len(), i,
		)

	}
}

func TestDequePop_empty_container(t *testing.T) {
	deque := NewDeque()
	item := deque.Pop()

	_assert(
		t,
		item == nil,
		"item = %v; want %v", item, nil,
	)

	_assert(
		t,
		deque.container.Len() == 0,
		"deque.container.Len() = %d; want %d", deque.container.Len(), 0,
	)
}

func TestDequeShift_fulfilled_container(t *testing.T) {
	deque := NewDeque()
	dequeSize := 100

	// Populate the test deque
	for i := 0; i < dequeSize; i++ {
		var value string = strconv.Itoa(i)
		deque.Append(value)
	}

	// Pop elements of the deque and assert elements come out
	// in order and container size is updated accordingly
	for i := 0; i < dequeSize; i++ {
		item := deque.Shift()

		var itemValue string = item.(string)
		var expectedValue string = strconv.Itoa(i)

		_assert(
			t,
			itemValue == expectedValue,
			"deque.Shift() = %s; want %s", itemValue, expectedValue,
		)

		_assert(
			t,
			// Len should be equal to dequeSize - (i + 1) as i is zero indexed
			deque.container.Len() == (dequeSize-(i+1)),
			"deque.container.Len() = %d; want %d", deque.container.Len(), dequeSize-i,
		)
	}
}

func TestDequeShift_empty_container(t *testing.T) {
	deque := NewDeque()

	item := deque.Shift()
	_assert(
		t,
		item == nil,
		"deque.Shift() = %v; want %v", item, nil,
	)

	_assert(
		t,
		deque.container.Len() == 0,
		"deque.container.Len() = %d; want %d", deque.container.Len(), 0,
	)
}

// 가장 앞에 넣은 원소 보기만
func TestDequeFirst_fulfilled_container(t *testing.T) {
	deque := NewDeque()
	deque.Append("1")
	item := deque.First()

	_assert(
		t,
		item == "1",
		"deque.First() = %s; want %s", item, "1",
	)

	_assert(
		t,
		deque.container.Len() == 1,
		"deque.container.Len() = %d; want %d", deque.container.Len(), 1,
	)
}

func TestDequeFirst_empty_container(t *testing.T) {
	deque := NewDeque()
	item := deque.First()

	_assert(
		t,
		item == nil,
		"deque.First() = %v; want %v", item, nil,
	)

	_assert(
		t,
		deque.container.Len() == 0,
		"deque.container.Len() = %d; want %d", deque.container.Len(), 0,
	)
}

// 마지막에 들어가 원소 보기만
func TestDequeLast_fulfilled_container(t *testing.T) {
	deque := NewDeque()

	deque.Append("1")
	deque.Append("2")
	deque.Append("3")

	item := deque.Last()

	_assert(
		t,
		item == "3",
		"deque.Last() = %s; want %s", item, "3",
	)

	_assert(
		t,
		deque.container.Len() == 3,
		"deque.container.Len() = %d; want %d", deque.container.Len(), 3,
	)
}

func TestDequeLast_empty_container(t *testing.T) {
	deque := NewDeque()
	item := deque.Last()

	_assert(
		t,
		item == nil,
		"deque.Last() = %v; want %v", item, nil,
	)

	_assert(
		t,
		deque.container.Len() == 0,
		"deque.container.Len() = %d; want %d", deque.container.Len(), 0,
	)
}

func TestDequeEmpty_fulfilled(t *testing.T) {
	deque := NewDeque()
	deque.Append("1")

	_assert(
		t,
		deque.Empty() == false,
		"deque.Empty() = %t; want %t", deque.Empty(), false)
}

func TestDequeEmpty_empty_deque(t *testing.T) {
	deque := NewDeque()
	_assert(
		t,
		deque.Empty() == true,
		"deque.Empty() = %t; want %t", deque.Empty(), true,
	)
}

// full 테스트
func TestDequeFull_fulfilled(t *testing.T) {
	deque := NewCappedDeque(3)

	deque.Append("1")
	deque.Append("2")
	deque.Append("3")

	_assert(
		t,
		deque.Full() == true,
		"deque.Full() = %t; want %t", deque.Full(), true,
	)
}

// full 테스트
func TestDequeFull_non_full_deque(t *testing.T) {
	deque := NewCappedDeque(3)
	deque.Append("1")

	_assert(
		t,
		deque.Full() == false,
		"deque.Full() = %t; want %t", deque.Full(), false,
	)
}

func _assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}
