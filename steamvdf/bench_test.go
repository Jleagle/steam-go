package steamvdf

import (
	"os"
	"testing"
)

func benchmarkRead(b *testing.B, path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(len(data)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := ReadBytes(data); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkReadText parses the largest text document in the corpus.
func BenchmarkReadText(b *testing.B) {
	benchmarkRead(b, "testdata/app_221680.vdf")
}

// BenchmarkReadBinary parses the largest binary document in the corpus; its
// "depotids" object has over five thousand children.
func BenchmarkReadBinary(b *testing.B) {
	benchmarkRead(b, "testdata/package_17906.bin")
}
