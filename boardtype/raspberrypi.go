package boardtype

var (
	RaspberryPi       = BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "", RAM: 0}
	RaspberryPi3      = BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "", RAM: 0, BaseModel: &RaspberryPi}
	RaspberryPi3B     = BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "3B", RAM: 1024, BaseModel: &RaspberryPi3}
	RaspberryPi3APlus = BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "3A+", RAM: 512, BaseModel: &RaspberryPi3B}
	RaspberryPi3BPlus = BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "3B+", RAM: 1024, BaseModel: &RaspberryPi3B}
	RaspberryPi4      = BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "", RAM: 0, BaseModel: &RaspberryPi}
	RaspberryPi4B     = BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "4B", RAM: 0, BaseModel: &RaspberryPi4}
	RaspberryPi4B1GB  = BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "4B", RAM: 1024, BaseModel: &RaspberryPi4B}
	RaspberryPi4B2GB  = BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "4B", RAM: 2048, BaseModel: &RaspberryPi4B}
	RaspberryPi4B4GB  = BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "4B", RAM: 4096, BaseModel: &RaspberryPi4B}
	RaspberryPi4B8GB  = BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "4B", RAM: 8192, BaseModel: &RaspberryPi4B}
	RaspberryPi4400   = BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "4 400", RAM: 4096, BaseModel: &RaspberryPi4B}
	RaspberryPiCM41GB = BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "Compute Module 4", RAM: 1024, BaseModel: &RaspberryPi4B}
	RaspberryPiCM42GB = BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "Compute Module 4", RAM: 2048, BaseModel: &RaspberryPi4B}
	RaspberryPiCM44GB = BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "Compute Module 4", RAM: 4096, BaseModel: &RaspberryPi4B}
	RaspberryPiCM48GB = BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "Compute Module 4", RAM: 8192, BaseModel: &RaspberryPi4B}
	RaspberryPi5      = BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "", RAM: 0, BaseModel: &RaspberryPi}
	RaspberryPi5B     = BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "5B", RAM: 0, BaseModel: &RaspberryPi5}
	RaspberryPi5B2GB  = BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "5B", RAM: 2048, BaseModel: &RaspberryPi5B}
	RaspberryPi5B4GB  = BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "5B", RAM: 4096, BaseModel: &RaspberryPi5B}
	RaspberryPi5B8GB  = BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "5B", RAM: 8192, BaseModel: &RaspberryPi5B}
	RaspberryPiCM51GB = BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "Compute Module 5", RAM: 1024, BaseModel: &RaspberryPi5B}
	RaspberryPiCM52GB = BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "Compute Module 5", RAM: 2048, BaseModel: &RaspberryPi5B}
	RaspberryPiCM54GB = BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "Compute Module 5", RAM: 4096, BaseModel: &RaspberryPi5B}
	RaspberryPiCM58GB = BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "Compute Module 5", RAM: 8192, BaseModel: &RaspberryPi5B}
)
