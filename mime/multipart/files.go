package multipart

import "github.com/anton2920/gofa/trace"

type Files struct {
	Keys         []string
	Names        [][]string
	ContentTypes [][]string
	Contents     [][][]byte
}

func (fs *Files) Add(key, name, contentType string, contents []byte) {
	t := trace.Begin("")

	for i := 0; i < len(fs.Keys); i++ {
		if key == fs.Keys[i] {
			fs.Names[i] = append(fs.Names[i], name)
			fs.ContentTypes[i] = append(fs.ContentTypes[i], contentType)
			fs.Contents[i] = append(fs.Contents[i], contents)

			trace.End(t)
			return
		}
	}
	fs.Keys = append(fs.Keys, key)

	if len(fs.Names) == cap(fs.Names) {
		fs.Names = append(fs.Names, []string{name})
		fs.ContentTypes = append(fs.ContentTypes, []string{contentType})
		fs.Contents = append(fs.Contents, [][]byte{contents})

		trace.End(t)
		return
	}
	n := len(fs.Names)
	fs.Names = fs.Names[:n+1]
	fs.Names[n] = fs.Names[n][:0]
	fs.Names[n] = append(fs.Names[n], name)

	fs.ContentTypes = fs.ContentTypes[:n+1]
	fs.ContentTypes[n] = fs.ContentTypes[n][:0]
	fs.ContentTypes[n] = append(fs.ContentTypes[n], contentType)

	fs.Contents = fs.Contents[:n+1]
	fs.Contents[n] = fs.Contents[n][:0]
	fs.Contents[n] = append(fs.Contents[n], contents)

	trace.End(t)
}

func (fs *Files) Get(key string) (string, string, []byte) {
	t := trace.Begin("")

	for i := 0; i < len(fs.Keys); i++ {
		if key == fs.Keys[i] {
			trace.End(t)
			return fs.Names[i][0], fs.ContentTypes[i][0], fs.Contents[i][0]
		}
	}

	trace.End(t)
	return "", "", nil
}

func (fs *Files) GetMany(key string) ([]string, []string, [][]byte) {
	t := trace.Begin("")

	for i := 0; i < len(fs.Keys); i++ {
		if key == fs.Keys[i] {
			trace.End(t)
			return fs.Names[i], fs.ContentTypes[i], fs.Contents[i]
		}
	}

	trace.End(t)
	return nil, nil, nil
}

func (fs *Files) Has(key string) bool {
	t := trace.Begin("")

	for i := 0; i < len(fs.Keys); i++ {
		if key == fs.Keys[i] {
			trace.End(t)
			return true
		}
	}

	trace.End(t)
	return false
}

func (fs *Files) Set(key, name, contentType string, contents []byte) {
	t := trace.Begin("")

	for i := 0; i < len(fs.Keys); i++ {
		if key == fs.Keys[i] {
			fs.Names[i] = fs.Names[i][:0]
			fs.Names[i] = append(fs.Names[i], name)

			fs.ContentTypes[i] = fs.ContentTypes[i][:0]
			fs.ContentTypes[i] = append(fs.ContentTypes[i], contentType)

			fs.Contents[i] = fs.Contents[i][:0]
			fs.Contents[i] = append(fs.Contents[i], contents)

			trace.End(t)
			return
		}
	}
	fs.Keys = append(fs.Keys, key)

	if len(fs.Names) == cap(fs.Names) {
		fs.Names = append(fs.Names, []string{name})
		fs.ContentTypes = append(fs.ContentTypes, []string{contentType})
		fs.Contents = append(fs.Contents, [][]byte{contents})

		trace.End(t)
		return
	}
	n := len(fs.Names)
	fs.Names = fs.Names[:n+1]
	fs.Names[n] = fs.Names[n][:0]
	fs.Names[n] = append(fs.Names[n], name)

	fs.ContentTypes = fs.ContentTypes[:n+1]
	fs.ContentTypes[n] = fs.ContentTypes[n][:0]
	fs.ContentTypes[n] = append(fs.ContentTypes[n], contentType)

	fs.Contents = fs.Contents[:n+1]
	fs.Contents[n] = fs.Contents[n][:0]
	fs.Contents[n] = append(fs.Contents[n], contents)

	trace.End(t)
}

func (fs *Files) Reset() {
	fs.Keys = fs.Keys[:0]
	fs.Names = fs.Names[:0]
	fs.ContentTypes = fs.ContentTypes[:0]
	fs.Contents = fs.Contents[:0]
}
