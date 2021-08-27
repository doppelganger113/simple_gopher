package storage

import "testing"

func TestPagingToLimitOffset(test *testing.T) {
	data := []struct {
		testName       string
		page           uint
		size           uint
		expectedLimit  int
		expectedOffset int
	}{
		{testName: "Empty", page: 0, size: 0, expectedLimit: 20, expectedOffset: 0},
		{testName: "First page", page: 1, size: 10, expectedLimit: 10, expectedOffset: 0},
		{testName: "Second page", page: 2, size: 10, expectedLimit: 10, expectedOffset: 10},
		{testName: "Second page", page: 3, size: 70, expectedLimit: 50, expectedOffset: 100},
	}

	for _, d := range data {
		test.Run(d.testName, func(t *testing.T) {
			limit, offset := PagingToLimitOffset(d.page, d.size)
			if limit != d.expectedLimit || offset != d.expectedOffset {
				t.Fatalf(
					"Expected offset %d to equal %d and limit %d to equal %d",
					offset, d.expectedOffset, limit, d.expectedLimit,
				)
			}
		})
	}
}
