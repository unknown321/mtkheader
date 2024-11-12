package main

import (
	"bytes"
	"io"
	"io/fs"
	"os"
	"testing"
	"time"
)

// FakeFile implements FileLike and also fs.FileInfo.
type FakeFile struct {
	name     string
	contents string
	mode     fs.FileMode
	offset   int
}

// Reset prepares a FakeFile for reuse.
func (f *FakeFile) Reset() *FakeFile {
	f.offset = 0
	return f
}

// FileLike methods.

func (f *FakeFile) Name() string {
	// A bit of a cheat: we only have a basename, so that's also ok for FileInfo.
	return f.name
}

func (f *FakeFile) Stat() (fs.FileInfo, error) {
	return f, nil
}

func (f *FakeFile) Read(p []byte) (int, error) {
	if f.offset >= len(f.contents) {
		return 0, io.EOF
	}
	n := copy(p, f.contents[f.offset:])
	f.offset += n
	return n, nil
}

func (f *FakeFile) Close() error {
	return nil
}

// fs.FileInfo methods.

func (f *FakeFile) Size() int64 {
	return int64(len(f.contents))
}

func (f *FakeFile) Mode() fs.FileMode {
	return f.mode
}

func (f *FakeFile) ModTime() time.Time {
	return time.Time{}
}

func (f *FakeFile) IsDir() bool {
	return false
}

func (f *FakeFile) Sys() any {
	return nil
}

func (f *FakeFile) String() string {
	return fs.FormatFileInfo(f)
}

func TestPatch(t *testing.T) {
	type args struct {
		headerFile io.Reader
		out        io.Writer
		s          os.FileInfo
	}

	data, err := os.ReadFile("test/initrd_header")
	if err != nil {
		t.Errorf("cannot read data: %s", err.Error())
	}

	r := bytes.NewReader(data)
	res := []byte{}
	out := bytes.NewBuffer(res)

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				headerFile: r,
				out:        out,
				s: &FakeFile{
					name:     "test",
					contents: "123",
					mode:     0,
					offset:   0,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Patch(tt.args.headerFile, tt.args.out, tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("Patch() error = %v, wantErr %v", err, tt.wantErr)
			}

			h := &Header{}
			if err := h.FromReader(out); err != nil {
				t.Errorf("%s\n", err.Error())
			}

			if int64(h.Length) != tt.args.s.Size() {
				t.Errorf("patch new size %d, want %d", h.Length, tt.args.s.Size())
			}
		})
	}
}
