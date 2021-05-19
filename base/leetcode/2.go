package main

import "fmt"

/*
	两数相加
*/

func main() {
	n1 := &ListNode{
		Val: 2,
		Next: &ListNode{
			Val: 4,
			Next: &ListNode{
				Val:  3,
				Next: nil,
			},
		},
	}

	n2 := &ListNode{
		Val: 5,
		Next: &ListNode{
			Val: 6,
			Next: &ListNode{
				Val: 4,
				Next: &ListNode{
					Val:  9,
					Next: nil,
				},
			},
		},
	}

	addTwoNumbers(n1, n2).Print()
}

/*Input: (2 -> 4 -> 3) + (5 -> 6 -> 4)
Output: 7 -> 0 -> 8*/

type ListNode struct {
	Val  int
	Next *ListNode
}

func (l *ListNode) Print() {

	if l != nil {
		fmt.Println(l.Val)
	}

	if l.Next != nil {
		l.Next.Print()
	}
}

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	head := &ListNode{0, nil}
	current := head
	carry := 0
	for l1 != nil || l2 != nil || carry > 0 {
		sum := carry
		if l1 != nil {
			sum += l1.Val
			l1 = l1.Next
		}
		if l2 != nil {
			sum += l2.Val
			l2 = l2.Next
		}
		carry = sum / 10
		current.Next = new(ListNode)
		current.Next.Val = sum % 10
		current = current.Next
	}
	return head.Next
}
