package main

import (
	"fmt"
	"os"
	"strconv"
	"errors"
	"strings"
	"bufio"
	"path/filepath"
)

type args struct {
	count    int
	inName   string
	outNamer func(int) string
}

func main() {
	args, err := parseArgs()
	check(err)
	lines, err := readLines(args)
	check(err)
	chunks := chunkLines(lines, args)
	writeChunks(chunks, args)
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func parseArgs() (args args, err error) {
	argc := len(os.Args)
	if argc != 3 && argc != 4 {
		formatUsage()
	}
	if args.count, err = strconv.Atoi(os.Args[1]); err != nil {
		return
	}
	if args.count < 0 {
		err = errors.New("files_out must be positive")
		return
	}
	args.inName = os.Args[2]
	outName := os.Args[2]
	if argc == 4 {
		outName = os.Args[3]
	}
	outExt := filepath.Ext(outName)
	if outExt == "" {
		args.outNamer = func(i int) string {
			return outName + strconv.Itoa(i)
		}
	} else {
		outPre := strings.TrimSuffix(args.inName, outExt)
		args.outNamer = func(i int) string {
			return outPre + strconv.Itoa(i) + outExt
		}
	}
	return
}

func formatUsage() {
	fmt.Println("\nusage:", filepath.Base(os.Args[0]), "files_out input_file [output_file]")
	fmt.Println("  files_out: number of files to split origin file into")
	fmt.Println("  input_file: input file name")
	fmt.Println("  output_file: output file name (when specified input file name is used)")
	os.Exit(0)
}

func readLines(args args) (lines []string, err error) {
	inFile, err := os.Open(args.inName)
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func chunkLines(lines []string, args args) (chunks [][]string) {
	lenLines := len(lines)
	chunkSize := ceilDiv(lenLines, args.count)
	for i := 0; i < lenLines; i += chunkSize {
		chunks = append(chunks, lines[i:min(i+chunkSize, lenLines)])
	}
	return
}

func ceilDiv(x, y int) int {
	if x%y == 0 {
		return x / y
	} else {
		return x/y + 1
	}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func writeChunks(chunks [][]string, args args) {
	fmt.Println("SPLITTING", args.inName, "INTO", args.count, "FILES")
	for i, chunk := range chunks {
		outFile, err := os.Create(args.outNamer(i))
		if err != nil {
			check(err)
		}
		writer := bufio.NewWriter(outFile)
		for _, line := range chunk {
			fmt.Fprintln(writer, line)
		}
		err = writer.Flush()
		if err != nil {
			check(err)
		}
		fmt.Print("\rWROTE", i)
	}
	fmt.Println("\rSUCCESS")
	return
}
