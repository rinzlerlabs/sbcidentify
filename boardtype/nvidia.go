package boardtype

var (
	NVIDIA                        = BoardType{Manufacturer: "NVIDIA", Model: "", SubModel: "", RAM: 0}
	Jetson                        = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "", RAM: 0, BaseModel: &NVIDIA}
	JetsonOrin                    = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Orin", RAM: 0, BaseModel: &Jetson}
	JetsonXavier                  = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Xavier", RAM: 0, BaseModel: &Jetson}
	JetsonOrinNX                  = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Orin NX", RAM: 0, BaseModel: &JetsonOrin}
	JetsonOrinNX16GB              = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Orin NX", RAM: 16384, BaseModel: &JetsonOrinNX}
	JetsonOrinNX8GB               = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Orin NX", RAM: 8192, BaseModel: &JetsonOrinNX}
	JetsonOrinNano                = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Orin Nano", RAM: 0, BaseModel: &JetsonOrin}
	JetsonOrinNano8GB             = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Orin Nano", RAM: 8192, BaseModel: &JetsonOrinNano}
	JetsonOrinNano4GB             = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Orin Nano", RAM: 4096, BaseModel: &JetsonOrinNano}
	JetsonOrinNanoDeveloperKit    = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Orin Nano Developer Kit", RAM: 8192, BaseModel: &JetsonOrinNano}
	JetsonAGXOrin                 = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Orin", RAM: 0, BaseModel: &Jetson}
	JetsonAGXOrin32GB             = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Orin", RAM: 32768, BaseModel: &JetsonAGXOrin}
	JetsonAGXOrin64GB             = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Orin", RAM: 65536, BaseModel: &JetsonAGXOrin}
	JetsonXavierNX                = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Xavier NX", RAM: 0, BaseModel: &Jetson}
	JetsonXavierNXDeveloperKit    = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Xavier NX Developer Kit", RAM: 0, BaseModel: &JetsonXavierNX}
	JetsonXavierNX8GB             = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Xavier NX", RAM: 8192, BaseModel: &JetsonXavierNX}
	JetsonXavierNX16GB            = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Xavier NX", RAM: 16384, BaseModel: &JetsonXavierNX}
	JetsonAGXXavier               = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Xavier", RAM: 0, BaseModel: &Jetson}
	JetsonAGXXavier8GB            = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Xavier", RAM: 8192, BaseModel: &JetsonAGXXavier}
	JetsonAGXXavier16GB           = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Xavier", RAM: 16384, BaseModel: &JetsonAGXXavier}
	JetsonAGXXavier32GB           = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Xavier", RAM: 32768, BaseModel: &JetsonAGXXavier}
	JetsonAGXXavier64GB           = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Xavier", RAM: 65536, BaseModel: &JetsonAGXXavier}
	JetsonAGXXavierIndustrial32GB = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Xavier Industrial", RAM: 32768, BaseModel: &JetsonAGXXavier}
	JetsonNano                    = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Nano", RAM: 0, BaseModel: &Jetson}
	JetsonNanoDeveloperKit        = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Nano Developer Kit", RAM: 0, BaseModel: &JetsonNano}
	JetsonNano2GB                 = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Nano", RAM: 2048, BaseModel: &JetsonNano}
	JetsonNano16GbEMMC            = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Nano", RAM: 0, BaseModel: &JetsonNano}
	JetsonNano4GB                 = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Nano", RAM: 4096, BaseModel: &JetsonNano}
	JetsonTX2NX                   = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "TX2 NX", RAM: 0, BaseModel: &Jetson}
	JetsonTX24GB                  = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "TX2", RAM: 4096, BaseModel: &JetsonTX2}
	JetsonTX2i                    = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "TX2i", RAM: 0, BaseModel: &JetsonTX2}
	JetsonTX2                     = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "TX2", RAM: 0, BaseModel: &Jetson}
	JetsonTX1                     = BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "TX1", RAM: 0, BaseModel: &Jetson}
	ClaraAGX                      = BoardType{Manufacturer: "NVIDIA", Model: "Clara", SubModel: "AGX", RAM: 0, BaseModel: &NVIDIA}
	ShieldTV                      = BoardType{Manufacturer: "NVIDIA", Model: "Shield", SubModel: "TV", RAM: 0, BaseModel: &NVIDIA}
)
