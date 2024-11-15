package boardtype

import "strconv"

type BoardType struct {
	Manufacturer string
	Model        string
	SubModel     string
	RAM          int
}

func (b BoardType) GetManufacturer() string {
	return b.Manufacturer
}

func (b BoardType) GetModel() string {
	return b.Model
}

func (b BoardType) GetSubModel() string {
	return b.SubModel
}

func (b BoardType) GetRAM() int {
	return b.RAM
}

func (b BoardType) GetPrettyName() string {
	if b.RAM > 0 {
		ram := b.RAM
		if ram < 1024 {
			return b.Manufacturer + " " + b.Model + " " + b.SubModel + " " + strconv.Itoa(b.RAM) + "MB"
		} else {
			return b.Manufacturer + " " + b.Model + " " + b.SubModel + " " + strconv.Itoa(b.RAM/1024) + "GB"
		}

	}
	return b.Manufacturer + " " + b.Model + " " + b.SubModel
}

var (
	BoardTypeUnknown BoardType = BoardType{"Unknown", "Unknown", "Unknown", 0}
)

type SBC interface {
	GetManufacturer() string
	GetModel() string
	GetSubModel() string
	GetRAM() int
	GetPrettyName() string
}
