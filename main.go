package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/1x-eng/yadab/memory"
	"github.com/1x-eng/yadab/parser"
	"github.com/1x-eng/yadab/util"
)

func main() {
	mb := memory.NewMemoryBackend()

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("  Y   Y      A       DDDD      A       BBBBB  ")
	fmt.Println("   Y Y      A A      D   D    A A      B    B ")
	fmt.Println("    Y      AAAAAA    D   D   AAAAAA    BBBB   ")
	fmt.Println("    Y     A      A   D   D  A      A   B    B ")
	fmt.Println("    Y    A        A  DDDD  A        A  BBBBB  ")
	fmt.Println("")
	fmt.Println("    Welcome! This is `yadab` - yet another (relational) database ;)")

	for {
		fmt.Print("# ")
		text, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}
		text = strings.Replace(text, "\n", "", -1)

		ast, err := parser.Parse(text)

		if err != nil {
			log.Fatal(err)
		}

		for _, stmt := range ast.Statements {
			switch stmt.Kind {
			case util.CreateTableKind:
				err = mb.CreateTable(stmt.CreateTableStatement)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("ok")
			case util.InsertKind:
				err = mb.Insert(stmt.InsertStatement)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("ok")
			case util.SelectKind:
				results, err := mb.Select(stmt.SelectStatement)
				if err != nil {
					log.Fatal(err)
				}

				util.PrintTableHeader(results.Columns)
				util.PrintTableSeparator()

				for _, result := range results.Rows {
					util.PrintTableRow(result, results.Columns)
				}

				fmt.Println("ok")
			}
		}
	}
}
