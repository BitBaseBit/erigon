package types

import (
	"fmt"
	"math/big"
	"math/rand"
	"sort"
	"testing"
)

func TestSortHeadersAsc(t *testing.T) {
	indexes := []int{0, 1, 2, 3, 4, 5}
	perms := permutations(indexes)

	for i, idxs := range perms {
		i := i
		idxs := idxs

		t.Run(fmt.Sprintf("permutation %v(%d)", idxs, i), func(t *testing.T) {
			hs := make([]*Header, len(indexes))

			was := make([]uint64, len(indexes))

			for i, idx := range idxs {
				hs[i] = &Header{
					Number:   big.NewInt(int64(idx)),
					GasLimit: uint64(idx),
					GasUsed:  uint64(idx),
					Time:     uint64(idx),
				}

				was[i] = hs[i].Number.Uint64()
			}

			SortHeadersAsc(hs)

			sorted := sort.SliceIsSorted(hs, func(i, j int) bool {
				return hs[i].Number.Cmp(hs[j].Number) < 0
			})

			got := make([]uint64, len(hs))
			for i, h := range hs {
				got[i] = h.Number.Uint64()
			}

			if !sorted {
				t.Errorf("not sorted: was %v, got %v", was, got)
			}

			for i, h := range hs {
				if h.GasLimit != uint64(i) {
					t.Error("GasLimit has been changed")
				}
				if h.GasUsed != uint64(i) {
					t.Error("GasUsed has been changed")
				}
				if h.Time != uint64(i) {
					t.Error("Time has been changed")
				}
			}
		})
	}
}

func BenchmarkSortHeadersAsc(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()

	for i := 0; i < b.N; i++ {
		const n = 16386 //128
		idxs := rand.Perm(n)
		hs := make([]*Header, n)

		for i, idx := range idxs {
			hs[i] = &Header{
				Number:   big.NewInt(int64(idx)),
				GasLimit: uint64(idx),
				GasUsed:  uint64(idx),
				Time:     uint64(idx),
			}
		}

		b.StartTimer()
		SortHeadersAsc(hs)
		b.StopTimer()
	}
}
func BenchmarkSortHeadersAscStd(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		const n = 16386 //128
		idxs := rand.Perm(n)
		hs := make([]*Header, n)

		for i, idx := range idxs {
			hs[i] = &Header{
				Number:   big.NewInt(int64(idx)),
				GasLimit: uint64(idx),
				GasUsed:  uint64(idx),
				Time:     uint64(idx),
			}
		}

		b.StartTimer()
		sort.Slice(hs, func(i, j int) bool {
			return hs[i].Number.Cmp(hs[j].Number) < 0
		})
		b.StopTimer()
	}
}

func permutations(arr []int) [][]int {
	var helper func([]int, int)
	res := [][]int{}

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}
