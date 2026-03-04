package main

type Entry_type int
const (
	File_entry Entry_type = iota
	Dir_entry
	Symlink_entry
	Fifo_entry
	Socket_entry
)

type (
	Dir struct {
		content []Entry
	}

	File struct {
		content []byte
	}

	Entry struct {
		entry_type Entry_type
		dir *Dir
		file *File
	}
)
