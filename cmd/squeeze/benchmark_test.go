package squeeze

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockDirEntry struct {
	name string
	dir  bool
}

func (m mockDirEntry) Name() string {
	return m.name
}

func (m mockDirEntry) IsDir() bool {
	return m.dir
}

func (m mockDirEntry) Type() fs.FileMode {
	panic("not implemented")
}

func (m mockDirEntry) Info() (fs.FileInfo, error) {
	panic("not implemented")
}

func Test_filterFile(t *testing.T) {
	type test struct {
		name     string
		path     string
		entry    mockDirEntry
		expected bool
	}

	tests := []test{
		{
			name: "valid pdf",
			path: "invoice.pdf",
			entry: mockDirEntry{
				name: "invoice.pdf",
			},
			expected: true,
		},
		{
			name: "valid xml",
			path: "invoice.xml",
			entry: mockDirEntry{
				name: "invoice.xml",
			},
			expected: true,
		},
		{
			name: "hidden path",
			path: ".env",
			entry: mockDirEntry{
				name: ".env",
			},
			expected: false,
		},
		{
			name: "directory",
			path: "my-dir",
			entry: mockDirEntry{
				name: "my-dir",
				dir:  true,
			},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expected, filterFile(test.path, test.entry))
		})
	}
}
