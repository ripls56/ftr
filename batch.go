package ftr

type FileMeta struct {
	FileName string
}

type Batch struct {
	Id   int
	Meta FileMeta

	Content []byte
}
