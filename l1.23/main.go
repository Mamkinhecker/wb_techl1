package main

func Remove(slice []int, i int) {
	copy(slice[i:], slice[i+1:])
	slice = slice[:len(slice)-1]
}
