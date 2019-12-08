package syntax

type Location struct {
	Pos      int // position within the file
	Line     int
	Filename string
}
