package memory

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"

	. "github.com/1x-eng/yadab/util"
)

var (
	ErrTableDoesNotExist  = errors.New("table does not exist")
	ErrColumnDoesNotExist = errors.New("column does not exist")
	ErrInvalidSelectItem  = errors.New("select item is not valid")
	ErrInvalidDatatype    = errors.New("invalid datatype")
	ErrMissingValues      = errors.New("missing values")
)

type MemoryBackend struct {
	Tables map[string]*Table
}

func NewMemoryBackend() *MemoryBackend {
	return &MemoryBackend{
		Tables: map[string]*Table{},
	}
}

func (mb *MemoryBackend) CreateTable(crt *CreateTableStatement) error {
	t := Table{}
	mb.Tables[crt.Name.Value] = &t
	if crt.Cols == nil {

		return nil
	}

	for _, col := range *crt.Cols {
		t.Columns = append(t.Columns, col.Name.Value)

		var dt ColumnType
		switch col.Datatype.Value {
		case "int":
			dt = IntType
		case "text":
			dt = TextType
		default:
			return ErrInvalidDatatype
		}

		t.ColumnTypes = append(t.ColumnTypes, dt)
	}

	return nil
}

func (mb *MemoryBackend) Insert(inst *InsertStatement) error {
	table, ok := mb.Tables[inst.Table.Value]
	if !ok {
		return ErrTableDoesNotExist
	}

	if inst.Values == nil {
		return nil
	}

	row := []MemoryCell{}

	if len(*inst.Values) != len(table.Columns) {
		return ErrMissingValues
	}

	for _, value := range *inst.Values {
		if value.Kind != LiteralKind {
			fmt.Println("Skipping non-literal.")
			continue
		}

		row = append(row, mb.tokenToCell(value.Literal))
	}

	table.Rows = append(table.Rows, row)
	return nil
}

func (mb *MemoryBackend) tokenToCell(t *Token) MemoryCell {
	if t.Kind == NumericKind {
		buf := new(bytes.Buffer)
		i, err := strconv.Atoi(t.Value)
		if err != nil {
			panic(err)
		}

		err = binary.Write(buf, binary.BigEndian, int32(i))
		if err != nil {
			panic(err)
		}
		return MemoryCell(buf.Bytes())
	}

	if t.Kind == StringKind {
		return MemoryCell(t.Value)
	}

	return nil
}

func (mb *MemoryBackend) Select(slct *SelectStatement) (*Results, error) {
	table, ok := mb.Tables[slct.From.Value]
	if !ok {
		return nil, ErrTableDoesNotExist
	}

	results := [][]Cell{}
	columns := []ResultColumn{}

	for i, row := range table.Rows {
		result := []Cell{}
		isFirstRow := i == 0

		for _, exp := range slct.Item {
			if exp.Kind != LiteralKind {
				// Unsupported, doesn't currently exist, ignore. For now...
				fmt.Println("Skipping non-literal expression.")
				continue
			}

			lit := exp.Literal
			if lit.Kind == IdentifierKind {
				found := false
				for i, tableCol := range table.Columns {
					if tableCol == lit.Value {
						if isFirstRow {
							columns = append(columns, struct {
								Type ColumnType
								Name string
							}{
								Type: table.ColumnTypes[i],
								Name: lit.Value,
							})
						}

						result = append(result, row[i])
						found = true
						break
					}
				}

				if !found {
					return nil, ErrColumnDoesNotExist
				}

				continue
			}

			return nil, ErrColumnDoesNotExist
		}

		results = append(results, result)
	}

	return &Results{
		Columns: columns,
		Rows:    results,
	}, nil
}
