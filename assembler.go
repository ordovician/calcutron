package calcutron

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

// The operands to an instruction
type Operands struct {
	regs     []uint8
	constant int8
}

// A table containing the memory address of labels in the code
func readSymTable(reader io.Reader) map[string]uint8 {
	scanner := bufio.NewScanner(reader)
	labels := make(map[string]uint8)
	for address := 0; scanner.Scan(); address++ {
		line := strings.Trim(scanner.Text(), " \t")
		n := len(line)

		if n == 0 {
			continue
		}

		if i := strings.IndexRune(line, ':'); i >= 0 {
			labels[line[0:i]] = uint8(address)

			// is there anything beyond the label?
			if n == i+1 {
				continue
			}
		}

	}
	return labels
}

func assembleLine(labels map[string]uint8, line string) (int16, error) {
	code := strings.Trim(line, " \t")
	i := len(code)
	if j := strings.Index(code, "//"); j >= 0 {
		i = j
	}
	n := len(code)

	if n == 0 || code[n-1] == ':' {
		return -1, nil
	}

	code = code[0:i]
	if i = strings.IndexRune(code, ' '); i < 0 {
		i = n
	}
	mnemonic := code[0:i]
	var operands []string = strings.SplitN(code[i:], ",", 3)
	opcode := ParseOpcode(mnemonic)
	var instruction uint16 = uint16(opcode)

	registers, err := parseOperands(labels, operands)
	if err != nil {
		return 0, err
	}
	if strings.EqualFold(mnemonic, "DAT") {
		return int16(registers[0]), nil
	}

	switch opcode {
	case ADD, SUB:
		switch len(operands) {
		case 3:
		}
		return fmt.Sprintf("%-4v x%d, x%d, x%d", inst.opcode, inst.dst, inst.src, inst.offset)
	case SUBI, LSH, RSH:
		return fmt.Sprintf("%-4v x%d, x%d, %d", inst.opcode, inst.dst, inst.src, inst.offset)
	case LD, ST, BRZ, BGT:
		return fmt.Sprintf("%-4v x%d, %d", inst.opcode, inst.dst, inst.addr)
	case HLT:
		return fmt.Sprintf("%-4v", inst.opcode)
	default:
		break
	}

	return -1, nil
}

func parseOperands(labels map[string]uint8, operands []string) ([]int8, error) {
	registers := make([]int8, 0)

	for _, operand := range operands {
		if addr, ok := labels[operand]; ok {
			registers = append(registers, int8(addr))
		} else if len(operand) == 2 && unicode.IsDigit(rune(operand[1])) {
			registers = append(registers, int8(operand[1]-'0'))
		} else {
			offset, err := strconv.Atoi(operand)
			if err != nil {
				return registers, fmt.Errorf("unable to parse offset/address %s because %w", operand, err)
			}

			registers = append(registers, int8(offset))
		}
	}

	return registers, nil
}
