package imHash_test

import (
	"bufio"
	"fmt"
	"immutable/imHash"
	"io"
	"os"
	"runtime"
	"testing"
)

type siPair struct {
	string
	int
}

const GETSPERSET = 100

func readData(filename string) []string {
	words := make([]string, 0)
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file ", err)
		os.Exit(1)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		word, err := r.ReadString(10) // 0x0A separator = newline
		if err == io.EOF {
			// do something here
			break
		} else if err != nil {
			return nil // if you return error
		}
		words = append(words, word)
	}
	return words

}

func BenchmarkImHash(b *testing.B) {
	words := readData("/usr/share/dict/web2")

	sh := imHash.NewStringHash()
	runtime.GC()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w := words[i%200000]
		sh = sh.Insert(w, i)
	}

	for i := 0; i < b.N*GETSPERSET; i++ {
		w := words[i*5%200000]
		k, v := sh.Find(w)
		if v != nil && k == "ksdlsjfkldfksjfd" {
			/* */
		}
	}
}

func BenchmarkGoHash(b *testing.B) {
	words := readData("/usr/share/dict/web2")

	sh := make(map[string]int)
	runtime.GC()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w := words[i%200000]
		sh[w] = i
	}

	for i := 0; i < b.N*GETSPERSET; i++ {
		w := words[i*5%200000]
		v, ok := sh[w]
		if !ok && v < 0 {
			/* */
		}
	}

}
