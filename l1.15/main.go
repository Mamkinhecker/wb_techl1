package l115

var justString string

func someFunc() {
	v := createHugeString(1 << 10) // предположительно возвращает большую строку
	justString = v[:100]           // берём первые 100 символов
}

func main() {
	someFunc()
}
