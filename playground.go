package main

import "fmt"

// Define the ListNode structure
type ListNode struct {
	Val  int
	Next *ListNode
}

// Function to merge two sorted linked lists
func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode {
	// Create a dummy node to serve as the start of the merged list
	dummy := &ListNode{}
	current := dummy

	for l1 != nil && l2 != nil {
		if l1.Val <= l2.Val {
			current = l1
			l1 = l1.Next
		} else {
			current = l2
			l2 = l2.Next
		}

		current = current.Next
	}

	// If either list still has nodes, append them to the merged list
	if l1 != nil {
		current.Next = l1
	} else if l2 != nil {
		current.Next = l2
	}

	return dummy.Next
}

// Helper function to print the linked list
func printList(head *ListNode) {
	current := head
	for current != nil {
		fmt.Printf("%d -> ", current.Val)
		current = current.Next
	}
	fmt.Println("nil")
}

func main() {
	// Create first sorted linked list: 1 -> 2 -> 4
	l1 := &ListNode{Val: 1}
	l1.Next = &ListNode{Val: 2}
	l1.Next.Next = &ListNode{Val: 4}

	// Create second sorted linked list: 1 -> 3 -> 4
	l2 := &ListNode{Val: 1}
	l2.Next = &ListNode{Val: 3}
	l2.Next.Next = &ListNode{Val: 4}

	fmt.Println("List 1:")
	printList(l1)
	fmt.Println("List 2:")
	printList(l2)

	// Merge the two sorted linked lists
	mergedList := mergeTwoLists(l1, l2)

	fmt.Println("Merged List:")
	printList(mergedList)
}
