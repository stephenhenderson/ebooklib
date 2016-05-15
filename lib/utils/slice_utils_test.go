package utils

import "testing"

func TestStringSliceEqualsTrueWhenBothNil(t *testing.T) {
	var arr1, arr2 []string
	arr1 = nil
	arr2 = nil
	if !StringSliceEquals(arr1, arr2) {
		t.Fatal("Two nil slices should be equal")
	}
}

func TestStringSliceEqualsTrueWhenSameInstance(t *testing.T) {
	slices := [][]string{
		nil,
		[]string{"a"},
		[]string{"a", "b"},
		[]string{"a string", "and another", "and a third"},
	}

	for _, slice := range slices {
		if !StringSliceEquals(slice, slice) {
			t.Fatal("any string slice should equal itself")
		}
	}
}

func TestStringSliceEqualsTrueWhenContentsAreTheSameButDifferentInstances(t *testing.T) {
	left_slices := [][]string{
		[]string{"a"},
		[]string{"a", "b"},
		[]string{"a string", "and another", "and a third"},
	}

	right_slices := [][]string{
		[]string{"a"},
		[]string{"a", "b"},
		[]string{"a string", "and another", "and a third"},
	}

	for i, left_slice := range left_slices {
		right_slice := right_slices[i]
		if !StringSliceEquals(left_slice, right_slice) {
			t.Fatalf("left: %v should equal right: %v", left_slice, right_slice)
		}
	}
}

func TestStringSliceEqualsFalseWhenContentsAreDifferent(t *testing.T) {
	left_slices := [][]string{
		nil,
		[]string{"a"},
		[]string{"a", "b"},
		[]string{"a string", "and another", "and a third"},
	}

	right_slices := [][]string{
		[]string{"not nil"},
		[]string{"b"},
		[]string{"b", "a"},
		[]string{"and a third", "and another", "a string"},
	}

	for i, left_slice := range left_slices {
		right_slice := right_slices[i]
		if StringSliceEquals(left_slice, right_slice) {
			t.Fatalf("left: %v should not equal right: %v", left_slice, right_slice)
		}
	}
}