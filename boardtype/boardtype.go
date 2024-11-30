package boardtype

import (
	"strconv"
)

type BoardType struct {
	Manufacturer string
	Model        string
	SubModel     string
	RAM          int
	BaseModel    *BoardType
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

func (b BoardType) GetBaseModel() *BoardType {
	return b.BaseModel
}

func (b BoardType) IsBoardType(boardType SBC) bool {
	return isBoardType(b, boardType)
}

func isBoardType(have SBC, want SBC) bool {
	if have.GetManufacturer() == want.GetManufacturer() && have.GetModel() == want.GetModel() && have.GetSubModel() == want.GetSubModel() && have.GetRAM() == want.GetRAM() {
		return true
	}

	if have.GetBaseModel() != nil {
		return isBoardType(have.GetBaseModel(), want)
	}

	return false
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

type SBC interface {
	GetManufacturer() string
	GetModel() string
	GetSubModel() string
	GetRAM() int
	GetPrettyName() string
	GetBaseModel() *BoardType
	IsBoardType(SBC) bool
}
