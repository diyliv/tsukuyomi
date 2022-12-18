package reading

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"sync"
)

func Process(f *os.File) error {
	linesPool := sync.Pool{New: func() interface{} {
		lines := make([]byte, 250*1024)
		return lines
	}}

	stringPool := sync.Pool{New: func() interface{} {
		lines := ""
		return lines
	}}

	r := bufio.NewReader(f)

	var wg sync.WaitGroup

	for {
		buf := linesPool.Get().([]byte)

		n, err := r.Read(buf)
		buf = buf[:n]

		if n == 0 {
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println(err)
				break
			}
			return err
		}

		nextUntillNewline, err := r.ReadBytes('\n')

		if err != io.EOF {
			buf = append(buf, nextUntillNewline...)
		}

		wg.Add(1)
		go func() {
			ProcessChunk(buf, &linesPool, &stringPool)
			wg.Done()
		}()

	}

	wg.Wait()
	return nil
}

func ProcessChunk(chunk []byte, linesPool *sync.Pool, stringPool *sync.Pool) {

	var wg2 sync.WaitGroup

	data := stringPool.Get().(string)
	data = string(chunk)

	linesPool.Put(chunk)

	dataSlice := strings.Split(data, "\n")

	stringPool.Put(data)

	chunkSize := 300
	n := len(dataSlice)
	noOfThread := n / chunkSize

	if n%chunkSize != 0 {
		noOfThread++
	}

	for i := 0; i < (noOfThread); i++ {

		wg2.Add(1)
		go func(s int, e int) {
			defer wg2.Done()
			for i := s; i < e; i++ {
				text := dataSlice[i]
				if len(text) == 0 {
					continue
				}
				fmt.Println(text)
			}

		}(i*chunkSize, int(math.Min(float64((i+1)*chunkSize), float64(len(dataSlice)))))
	}

	wg2.Wait()
	dataSlice = nil
}
