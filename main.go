package main

import (
	"fmt"
	"github.com/ulricqin/goutils/filetool"
	"github.com/ulricqin/randgen/g"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func main() {

	var err error

	argc := len(os.Args)
	if argc >= 2 {
		// randgen 20-49 [9999]
		scope := os.Args[1]
		scopeArr := strings.Split(scope, "-")
		if len(scopeArr) != 2 {
			usage()
			return
		}

		g.NumMin, err = strconv.Atoi(scopeArr[0])
		if err != nil {
			fmt.Printf("parse input fail. error: %s\n", err)
			return
		}

		g.NumMax, err = strconv.Atoi(scopeArr[1])
		if err != nil {
			fmt.Printf("parse input fail. error: %s\n", err)
			return
		}

	}

	g.NumDiff = g.NumMax - g.NumMin

	if argc >= 3 {
		// randgen 20-49 9999
		g.DiCnt, err = strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("parse di cnt fail. error: %s\n", err)
			return
		}
	}

	if argc >= 4 {
		// handle result dir
		g.Dir = os.Args[3]
	}

	fmt.Printf("argument>>> di-cnt: %d, num-scope: [%d-%d], result-dir: %s\n", g.DiCnt, g.NumMin, g.NumMax, g.Dir)
	generate()
}

func usage() {
	fmt.Println(`
USAGE: e.g. randgen 20-30 9999 ./result`)
}

func generate() {
	// 在result-dir下生成一万个txt文件
	fmt.Printf("generating...")

	err := filetool.InsureDir(g.Dir)
	if err != nil {
		fmt.Printf("create dir: %s fail. error: %s\n", g.Dir, err)
		return
	}

	for i := 0; i < g.DiCnt; i++ {
		filename := fmt.Sprintf("%04d.txt", i)
		if e := genOneFile(filename, int64(i)); e != nil {
			fmt.Printf("genOneFile fail. error: %s", e)
			break
		}
	}

	fmt.Println("done")

}

func genOneFile(filename string, seed int64) (err error) {
	r := rand.New(rand.NewSource(seed))
	numCnt := r.Intn(g.NumDiff) + g.NumMin

	// 生成numCnt个数字，写入filename，数字注意格式化
	m := make(map[int]bool)
	i := 0
	for {
		a := r.Intn(1000)
		if _, ok := m[a]; !ok {
			m[a] = true
			i++
			if i == numCnt {
				break
			}
		}
	}

	size := len(m)
	j := 0
	keys := make([]int, size)
	for k, _ := range m {
		keys[j] = k
		j++
	}

	sort.Ints(keys)

	nums := make([]string, size)
	for k := 0; k < size; k++ {
		nums[k] = fmt.Sprintf("%03d\r\n", keys[k])
	}

	f := filepath.Join(g.Dir, filename)
	_, err = filetool.WriteStringToFile(f, strings.Join(nums, ""))

	return
}
