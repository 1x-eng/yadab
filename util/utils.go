package util

import "fmt"

func PrintTableHeader(columns []ResultColumn) {
	for _, col := range columns {
		fmt.Printf("| %s ", col.Name)
	}
	fmt.Println("|")
}

func PrintTableSeparator() {
	for i := 0; i < 20; i++ {
		fmt.Printf("=")
	}
	fmt.Println()
}

func PrintTableRow(row []Cell, columns []ResultColumn) {
	fmt.Printf("|")
	for i, cell := range row {
		typ := columns[i].Type
		s := ""
		switch typ {
		case IntType:
			s = fmt.Sprintf("%d", cell.AsInt())
		case TextType:
			s = cell.AsText()
		}

		fmt.Printf(" %s | ", s)
	}
	fmt.Println()
}
