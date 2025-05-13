package bla2

import "github.com/manifoldco/promptui"

type prompt interface {
	Run() (err error)
	Result() string
}

var (
	_ prompt = (*selectPrompt)(nil)
	_ prompt = (*selectWithAddPrompt)(nil)
	_ prompt = (*linePrompt)(nil)
)

type selectPrompt struct {
	promptui.Select

	selectedIndex int
	selectedValue string
}

func (p *selectPrompt) Run() (err error) {
	p.selectedIndex, p.selectedValue, err = p.Select.Run()
	return err
}

func (p *selectPrompt) Result() string {
	return p.selectedValue
}

type selectWithAddPrompt struct {
	promptui.SelectWithAdd

	selectedIndex int
	selectedValue string
}

func (p *selectWithAddPrompt) Run() (err error) {
	p.selectedIndex, p.selectedValue, err = p.SelectWithAdd.Run()
	return err
}

func (p *selectWithAddPrompt) Result() string {
	return p.selectedValue
}

type linePrompt struct {
	promptui.Prompt

	selectedValue string
}

func (p *linePrompt) Run() (err error) {
	p.selectedValue, err = p.Prompt.Run()
	return err
}

func (p *linePrompt) Result() string {
	return p.selectedValue
}
