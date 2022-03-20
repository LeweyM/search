package trigram

// intersectPair

// [0, 2, 5, 7]
// [3, 4, 5, 6, 7, 8]
// => [5, 7]

// algorithm:

// two pointers, march the lowest,
// if they point to the same value, add and march both
// if you reach the end of either list, return

// intersectPair assumes that a b are both sorted and that there are no duplicates
// It must return a sorted list without duplicates.
func intersectPair(A []int, B []int) (res []int) {
	if len(A) == 0 || len(B) == 0 {
		return res
	}
	a, b := 0, 0
	for {
		if A[a] == B[b] {
			res = append(res, A[a])
		}
		if A[a] > B[b] {
			b++
		} else {
			a++
		}
		if a >= len(A) || b >= len(B) {
			return res
		}
	}
}

// unionPair

// [0, 2, 5, 7]
// [3, 4, 5, 6, 7, 8]
// => [0, 2, 3, 4, 5, 6, 7, 8]

// algorithm:

// two pointers, march the lowest,
// if they point to the same value, add one and march both
// if you reach the end of either list, add the remainder of the other list and return

// unionPair assumes that a b are both sorted and that there are no duplicates.
// It must return a sorted list without duplicates.
func unionPair(A, B []int) (res []int) {
	if len(A) == 0 {
		return B
	}
	if len(B) == 0 {
		return A
	}

	// The list will be at least the size of the largest sublist
	if len(A) > len(B) {
		res = make([]int, 0, len(A))
	} else {
		res = make([]int, 0, len(B))
	}

	a, b := 0, 0
	for {
		// if you reach the end of a list, return the results + the remainder of the other list
		if a >= len(A) {
			return append(res, B[b:]...)
		}
		if b >= len(B) {
			return append(res, A[a:]...)
		}

		// if they point to the same value, add one and march both
		if A[a] == B[b] {
			res = append(res, A[a])
			a++
			b++
			continue
		} else {
			// march the lowest value and add it to the results
			if A[a] > B[b] {
				res = append(res, B[b])
				b++
			} else {
				res = append(res, A[a])
				a++
			}
		}
	}
}
