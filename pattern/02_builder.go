package main

import "fmt"

type Computer struct {
	CPU         string
	GPU         string
	RAM         string
	MotherBoard string
}

func (c *Computer) String() string {
	return fmt.Sprintf("CPU: %v, GPU: %v, RAM: %v, MotherBoard: %v", c.CPU, c.GPU, c.RAM, c.MotherBoard)
}

type ComputerBuilder interface {
	SetUpCPU() ComputerBuilder
	SetUpGPU() ComputerBuilder
	SetUpRAM() ComputerBuilder
	SetUpMB() ComputerBuilder
	Build() *Computer
}

type GamingComputerBuilder struct {
	Computer *Computer
}

func (b *GamingComputerBuilder) SetUpCPU() ComputerBuilder {
	b.Computer.CPU = "intel i9"
	return b
}

func (b *GamingComputerBuilder) SetUpGPU() ComputerBuilder {
	b.Computer.GPU = "Nvidia RTX 5090"
	return b
}

func (b *GamingComputerBuilder) SetUpRAM() ComputerBuilder {
	b.Computer.RAM = "Kingston 32 gb"
	return b
}

func (b *GamingComputerBuilder) SetUpMB() ComputerBuilder {
	b.Computer.MotherBoard = "gigabyte Z790"
	return b
}

func (b *GamingComputerBuilder) Build() *Computer {
	return b.Computer
}

type OfficeComputerBuilder struct {
	Computer *Computer
}

func (b *OfficeComputerBuilder) SetUpCPU() ComputerBuilder {
	b.Computer.CPU = "intel i5"
	return b
}

func (b *OfficeComputerBuilder) SetUpGPU() ComputerBuilder {
	b.Computer.GPU = "integrated"
	return b
}

func (b *OfficeComputerBuilder) SetUpRAM() ComputerBuilder {
	b.Computer.RAM = "Samsung 8gb"
	return b
}

func (b *OfficeComputerBuilder) SetUpMB() ComputerBuilder {
	b.Computer.MotherBoard = "as rock 365b"
	return b
}

func (b *OfficeComputerBuilder) Build() *Computer {
	return b.Computer
}

type Director struct {
	Builder ComputerBuilder
}

func (d *Director) SetBuilder(builder ComputerBuilder) {
	d.Builder = builder
}

func (d *Director) Construct() *Computer {
	return d.Builder.SetUpCPU().SetUpGPU().SetUpRAM().SetUpMB().Build()
}
