package multipart

import "github.com/anton2920/gofa/trace/trace_"

type File struct {
	Name     string
	Mime     string
	Contents []byte
}

type Files struct {
	Keys  []string
	Files [][]File
}

func (fs *Files) Add(key string, file File) {
	t := trace_.Begin("")

	for i := 0; i < len(fs.Keys); i++ {
		if key == fs.Keys[i] {
			fs.Files[i] = append(fs.Files[i], file)

			trace_.End(t)
			return
		}
	}
	fs.Keys = append(fs.Keys, key)

	if len(fs.Files) == cap(fs.Files) {
		fs.Files = append(fs.Files, []File{file})

		trace_.End(t)
		return
	}
	n := len(fs.Files)
	fs.Files = fs.Files[:n+1]
	fs.Files[n] = fs.Files[n][:0]
	fs.Files[n] = append(fs.Files[n], file)

	trace_.End(t)
}

func (fs *Files) Get(key string) File {
	t := trace_.Begin("")

	for i := 0; i < len(fs.Keys); i++ {
		if key == fs.Keys[i] {
			trace_.End(t)
			return fs.Files[i][0]
		}
	}

	trace_.End(t)
	return File{}
}

func (fs *Files) GetMany(key string) []File {
	t := trace_.Begin("")

	for i := 0; i < len(fs.Keys); i++ {
		if key == fs.Keys[i] {
			trace_.End(t)
			return fs.Files[i]
		}
	}

	trace_.End(t)
	return nil
}

func (fs *Files) Has(key string) bool {
	t := trace_.Begin("")

	for i := 0; i < len(fs.Keys); i++ {
		if key == fs.Keys[i] {
			trace_.End(t)
			return true
		}
	}

	trace_.End(t)
	return false
}

func (fs *Files) Set(key string, file File) {
	t := trace_.Begin("")

	for i := 0; i < len(fs.Keys); i++ {
		if key == fs.Keys[i] {
			fs.Files[i] = fs.Files[i][:0]
			fs.Files[i] = append(fs.Files[i], file)
			trace_.End(t)
			return
		}
	}
	fs.Keys = append(fs.Keys, key)

	if len(fs.Files) == cap(fs.Files) {
		fs.Files = append(fs.Files, []File{file})
		trace_.End(t)
		return
	}
	n := len(fs.Files)
	fs.Files = fs.Files[:n+1]
	fs.Files[n] = fs.Files[n][:0]
	fs.Files[n] = append(fs.Files[n], file)

	trace_.End(t)
}

func (fs *Files) Reset() {
	fs.Keys = fs.Keys[:0]
	fs.Files = fs.Files[:0]
}
