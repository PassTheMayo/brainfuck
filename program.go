package main

import (
	"context"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

const (
	OutOfBoundsError = iota
	OutOfBoundsWrap
	OutOfBoundsExtend
	OutOfBoundsClamp
)

func NewProgram(memorySize int, outOfBoundsMethod int) *Program {
	return &Program{
		Input:             os.Stdin,
		Output:            os.Stdout,
		parsedTree:        make([]Command, 0),
		defaultMemorySize: memorySize,
		memory:            make([]uint8, memorySize),
		ptrIndex:          0,
		cmdIndex:          0,
		outOfBoundsMethod: outOfBoundsMethod,
		cycles:            0,
	}
}

type Program struct {
	Input             io.Reader
	Output            io.Writer
	parsedTree        []Command
	defaultMemorySize int
	memory            []uint8
	ptrIndex          int
	cmdIndex          int
	outOfBoundsMethod int
	cycles            int
}

func (p Program) DebugCommands() {
	depth := 0

	for i, c := range p.parsedTree {
		switch v := c.(type) {
		case *MovePointerCommand:
			{
				if v.Value > 0 {
					fmt.Printf("%5d | %s> MOVE_POINTER_RIGHT (c=%d)\n", i, strings.Repeat(" ", depth*4), v.Value)
				} else {
					fmt.Printf("%5d | %s< MOVE_POINTER_LEFT (c=%d)\n", i, strings.Repeat(" ", depth*4), -v.Value)
				}

				break
			}
		case *IncrementDecrementCommand:
			{
				if v.Value > 0 {
					fmt.Printf("%5d | %s+ INCREMENT (c=%d)\n", i, strings.Repeat(" ", depth*4), v.Value)
				} else {
					fmt.Printf("%5d | %s- DECREMENT (c=%d)\n", i, strings.Repeat(" ", depth*4), -v.Value)
				}

				break
			}
		case *OutputCommand:
			{
				fmt.Printf("%5d | %s. OUTPUT\n", i, strings.Repeat(" ", depth*4))

				break
			}
		case *InputCommand:
			{
				fmt.Printf("%5d | %s, INPUT\n", i, strings.Repeat(" ", depth*4))

				break
			}
		case *JumpCommand:
			{
				if v.IfZero {
					fmt.Printf("%5d | %s[ JUMP_FORWARD (i=%d)\n", i, strings.Repeat(" ", depth*4), v.CommandIndex)
					depth++
				} else {
					depth--
					fmt.Printf("%5d | %s] JUMP_BACKWARDS (i=%d)\n", i, strings.Repeat(" ", depth*4), v.CommandIndex)
				}

				break
			}
		default:
			fmt.Printf("%5d | %sUNKNOWN_COMMAND\n", i, strings.Repeat(" ", depth*4))
		}
	}
}

func (p Program) GetCellValue() uint8 {
	return p.memory[p.ptrIndex]
}

func (p Program) GetCellValueAt(index int) uint8 {
	return p.memory[index]
}

func (p Program) SetCellValue(value uint8) {
	p.memory[p.ptrIndex] = value
}

func (p Program) SetCellValueAt(index int, value uint8) {
	p.memory[index] = value
}

func (p Program) GetMemoryPointer() int {
	return p.ptrIndex
}

func (p *Program) SetMemoryPointer(value int) error {
	p.ptrIndex = value

	if value < 0 {
		switch p.outOfBoundsMethod {
		case OutOfBoundsError:
			return fmt.Errorf("memory pointer moved out of bounds to i=%d", value)
		case OutOfBoundsWrap:
			{
				p.ptrIndex = len(p.memory) - (int(math.Abs(float64(value))) % len(p.memory))

				break
			}
		case OutOfBoundsExtend:
			return fmt.Errorf("memory pointer moved out of bounds to i=%d, cannot extend left", value)
		case OutOfBoundsClamp:
			{
				p.ptrIndex = 0

				break
			}
		default:
			return fmt.Errorf("unknown out-of-bounds strategy of type %d", p.outOfBoundsMethod)
		}
	} else if value >= len(p.memory) {
		switch p.outOfBoundsMethod {
		case OutOfBoundsError:
			return fmt.Errorf("memory pointer moved out of bounds to i=%d", value)
		case OutOfBoundsWrap:
			{
				p.ptrIndex = value % len(p.memory)

				break
			}
		case OutOfBoundsExtend:
			{
				for value < len(p.memory) {
					p.memory = append(p.memory, 0)
				}

				break
			}
		case OutOfBoundsClamp:
			{
				p.ptrIndex = len(p.memory) - 1

				break
			}
		default:
			return fmt.Errorf("unknown out-of-bounds strategy of type %d", p.outOfBoundsMethod)
		}
	}

	return nil
}

func (p *Program) NextCommand() {
	p.cmdIndex++
}

func (p Program) GetCommandIndex() int {
	return p.cmdIndex
}

func (p *Program) SetCommandIndex(value int) {
	p.cmdIndex = value
}

func (p *Program) ResetMemory() {
	p.ptrIndex = 0
	p.cmdIndex = 0
	p.cycles = 0
	p.memory = make([]uint8, p.defaultMemorySize)
}

func (p *Program) ResetAll() {
	p.ResetMemory()

	p.parsedTree = make([]Command, 0)
}

func (p *Program) ReadFile(file string) error {
	data, err := os.ReadFile(file)

	if err != nil {
		return err
	}

	return p.ReadString(string(data))
}

func (p *Program) ReadString(data string) error {
	jumpCommandStack := make([]*JumpCommand, 0)

	for _, command := range data {
		switch command {
		case '>':
			{
				if len(p.parsedTree) > 0 {
					if v, ok := p.parsedTree[len(p.parsedTree)-1].(*MovePointerCommand); ok {
						v.Value++

						break
					}
				}

				p.parsedTree = append(p.parsedTree, &MovePointerCommand{
					Value: 1,
				})

				break
			}
		case '<':
			{
				if len(p.parsedTree) > 0 {
					if v, ok := p.parsedTree[len(p.parsedTree)-1].(*MovePointerCommand); ok {
						v.Value--

						break
					}
				}

				p.parsedTree = append(p.parsedTree, &MovePointerCommand{
					Value: -1,
				})

				break
			}
		case '+':
			{
				if len(p.parsedTree) > 0 {
					if v, ok := p.parsedTree[len(p.parsedTree)-1].(*IncrementDecrementCommand); ok {
						v.Value++

						break
					}
				}

				p.parsedTree = append(p.parsedTree, &IncrementDecrementCommand{
					Value: 1,
				})

				break
			}
		case '-':
			{
				if len(p.parsedTree) > 0 {
					if v, ok := p.parsedTree[len(p.parsedTree)-1].(*IncrementDecrementCommand); ok {
						v.Value--

						break
					}
				}

				p.parsedTree = append(p.parsedTree, &IncrementDecrementCommand{
					Value: -1,
				})

				break
			}
		case '.':
			{
				p.parsedTree = append(p.parsedTree, Output)

				break
			}
		case ',':
			{
				p.parsedTree = append(p.parsedTree, Input)

				break
			}
		case '[':
			{
				cmd := &JumpCommand{CommandIndex: len(p.parsedTree), IfZero: true}

				jumpCommandStack = append(jumpCommandStack, cmd)
				p.parsedTree = append(p.parsedTree, cmd)

				break
			}
		case ']':
			{
				if len(jumpCommandStack) < 1 {
					return fmt.Errorf("jump command stack is empty, unmatched closing jump command")
				}

				lastJumpCommand := jumpCommandStack[len(jumpCommandStack)-1]
				previousIndex := lastJumpCommand.CommandIndex
				lastJumpCommand.CommandIndex = len(p.parsedTree)
				jumpCommandStack = jumpCommandStack[:len(jumpCommandStack)-1]

				p.parsedTree = append(p.parsedTree, &JumpCommand{
					CommandIndex: previousIndex,
					IfZero:       false,
				})

				break
			}
		}
	}

	if len(jumpCommandStack) > 0 {
		return fmt.Errorf("jump command stack is not empty, unmatched jump command detected")
	}

	return nil
}

func (p *Program) Run(ctx context.Context) error {
	for p.cmdIndex < len(p.parsedTree) {
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := p.parsedTree[p.cmdIndex].Execute(p); err != nil {
			return err
		}

		p.cycles++
	}

	return nil
}

func (p Program) CycleCount() int {
	return p.cycles
}
