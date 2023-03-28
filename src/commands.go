package main

type Command interface {
	Execute(*Program) error
}

type MovePointerCommand struct {
	Value int
}

func (c MovePointerCommand) Execute(p *Program) error {
	if err := p.SetMemoryPointer(p.GetMemoryPointer() + c.Value); err != nil {
		return err
	}

	p.NextCommand()

	return nil
}

type IncrementDecrementCommand struct {
	Value int
}

func (c IncrementDecrementCommand) Execute(p *Program) error {
	p.SetCellValue(uint8(int(p.GetCellValue()) + c.Value))
	p.NextCommand()

	return nil
}

type InputCommand struct{}

func (c InputCommand) Execute(p *Program) error {
	data := make([]byte, 1)

	if _, err := p.Input.Read(data); err != nil {
		return err
	}

	p.SetCellValue(data[0])
	p.NextCommand()

	return nil
}

type OutputCommand struct{}

func (c OutputCommand) Execute(p *Program) error {
	p.Output.Write([]byte{p.GetCellValue()})

	p.NextCommand()

	return nil
}

type JumpCommand struct {
	CommandIndex int
	IfZero       bool
}

func (c JumpCommand) Execute(p *Program) error {
	v := p.GetCellValue()

	if (c.IfZero && v == 0) || (!c.IfZero && v != 0) {
		p.SetCommandIndex(c.CommandIndex)
	} else {
		p.NextCommand()
	}

	return nil
}

var (
	MovePointerLeft  Command = &MovePointerCommand{Value: -1}
	MovePointerRight Command = &MovePointerCommand{Value: 1}
	Increment        Command = &IncrementDecrementCommand{Value: 1}
	Decrement        Command = &IncrementDecrementCommand{Value: -1}
	Output           Command = &OutputCommand{}
	Input            Command = &InputCommand{}
)
