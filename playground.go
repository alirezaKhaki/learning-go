package main

import "fmt"

func main() {
	// Initial slice setup
	original := []int{1, 2, 3} // Slice with length 3, capacity 6
	fmt.Printf("Original slice: %v, len: %d, cap: %d\n", original, len(original), cap(original))

	// Copying the slice
	copySlice := original // Copying the slice
	fmt.Printf("Copy slice: %v, len: %d, cap: %d\n", copySlice, len(copySlice), cap(copySlice))

	// ------------------------------
	// Modifying the contents of the slice (both see the change)
	copySlice[1] = 4 // Modify the second element
	fmt.Println("\nAfter modifying copySlice:")
	fmt.Printf("Original slice: %v\n", original) // Change visible in original slice
	fmt.Printf("Copy slice: %v\n", copySlice)    // Change visible in copy slice

	// ------------------------------
	// Appending to the copy without exceeding capacity (changes only the copy's length)
	copySlice = append(copySlice, 1, 2, 3)
	fmt.Println("\nAfter appending to copySlice within its capacity:")
	fmt.Printf("Original slice: %v, len: %d, cap: %d\n", original, len(original), cap(original))
	fmt.Printf("Copy slice: %v, len: %d, cap: %d\n", copySlice, len(copySlice), cap(copySlice))

	// Original slice cannot see the new elements added to copySlice because the length is independent

	// ------------------------------
	// Appending to the copy and exceeding capacity (triggers reallocation)
	copySlice = append(copySlice, 9)
	fmt.Println("\nAfter exceeding capacity and appending to copySlice:")
	fmt.Printf("Original slice: %v, len: %d, cap: %d\n", original, len(original), cap(original))
	fmt.Printf("Copy slice: %v, len: %d, cap: %d\n", copySlice, len(copySlice), cap(copySlice))

	// Now, copySlice has a new underlying array, so changes to copySlice won't affect the original slice.
}
