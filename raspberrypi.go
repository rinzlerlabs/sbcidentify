package sbcidentify

import (
	"strings"

	"log/slog"
)

const (
	BoardTypeRaspberryPi3B      BoardType = "Raspberry Pi 3B 1GB RAM"
	BoardTypeRaspberryPi3APlus  BoardType = "Raspberry Pi 3A+ 512MB RAM"
	BoardTypeRaspberryPi3BPlus  BoardType = "Raspberry Pi 3B+ 1GB RAM"
	BoardTypeRaspberryPi4B1GB   BoardType = "Raspberry Pi 4B 1GB RAM"
	BoardTypeRaspberryPi4B2GB   BoardType = "Raspberry Pi 4B 2GB RAM"
	BoardTypeRaspberryPi4B4GB   BoardType = "Raspberry Pi 4B 4GB RAM"
	BoardTypeRaspberryPi4B8GB   BoardType = "Raspberry Pi 4B 8GB RAM"
	BoardTypeRaspberryPi4B      BoardType = "Raspberry Pi 4B"
	BoardTypeRaspberryPi44004GB BoardType = "Raspberry Pi 4 400 4GB RAM"
	BoardTypeRaspberryPi5B      BoardType = "Raspberry Pi 5B"
	BoardTypeRaspberryPi5B2GB   BoardType = "Raspberry Pi 5B 2GB RAM"
	BoardTypeRaspberryPi5B4GB   BoardType = "Raspberry Pi 5B 4GB RAM"
	BoardTypeRaspberryPi5B8GB   BoardType = "Raspberry Pi 5B 8GB RAM"
)

type raspberryPi struct {
	Model    string
	Memory   int
	Type     BoardType
	Fallback BoardType
}

var raspberryPiModels = []raspberryPi{
	{"Raspberry Pi 3 Model B", 1024, BoardTypeRaspberryPi3B, BoardTypeRaspberryPi3B},
	{"Raspberry Pi 3 Model A", 512, BoardTypeRaspberryPi3APlus, BoardTypeRaspberryPi3B},
	{"Raspberry Pi 3 Model B", 1024, BoardTypeRaspberryPi3BPlus, BoardTypeRaspberryPi3B},
	{"Raspberry Pi 4 Model B", 1024, BoardTypeRaspberryPi4B1GB, BoardTypeRaspberryPi4B},
	{"Raspberry Pi 4 Model B", 2048, BoardTypeRaspberryPi4B2GB, BoardTypeRaspberryPi4B},
	{"Raspberry Pi 4 Model B", 4096, BoardTypeRaspberryPi4B4GB, BoardTypeRaspberryPi4B},
	{"Raspberry Pi 4 Model B", 8192, BoardTypeRaspberryPi4B8GB, BoardTypeRaspberryPi4B},
	{"Raspberry Pi 400", 4096, BoardTypeRaspberryPi44004GB, BoardTypeRaspberryPi44004GB},
	{"Raspberry Pi 5 Model B", 2048, BoardTypeRaspberryPi5B2GB, BoardTypeRaspberryPi5B},
	{"Raspberry Pi 5 Model B", 4096, BoardTypeRaspberryPi5B4GB, BoardTypeRaspberryPi5B},
	{"Raspberry Pi 5 Model B", 8192, BoardTypeRaspberryPi5B8GB, BoardTypeRaspberryPi5B},
}

func NewRaspberryPiIdentifier(logger *slog.Logger) boardIdentifier {
	logger.Info("initializing Raspberry Pi identifier")
	newLogger := logger.With(slog.String("source", "RaspberryPiIdentifier"))
	return &raspberryPiIdentifier{logger: newLogger}
}

type raspberryPiIdentifier struct {
	logger *slog.Logger
}

func (r raspberryPiIdentifier) Name() string {
	return "Raspberry Pi"
}

func (r raspberryPiIdentifier) GetBoardType() (BoardType, error) {
	r.logger.Debug("getting board type")
	dtbm, err := getDeviceTreeBaseModel(r.logger)
	if err == ErrCannotIdentifyBoard {
		dtbm, err = getDeviceTreeModel(r.logger)
		if err != nil {
			return BoardTypeUnknown, err
		}
	} else if err != nil {
		return BoardTypeUnknown, err
	}
	r.logger.Debug("device tree model", slog.String("model", dtbm))
	subModels := make([]raspberryPi, 0)
	for _, m := range raspberryPiModels {
		if strings.Contains(dtbm, m.Model) {
			subModels = append(subModels, m)
		}
	}
	ramMb, err := getInstalledRAM(r.logger)
	if err != nil {
		return BoardTypeUnknown, err
	}
	for _, m := range subModels {
		if m.Memory == ramMb {
			return m.Type, nil
		}
	}
	r.logger.Debug("no matching model found, using fallback", slog.String("model", dtbm), slog.Int("ram", ramMb), slog.Int("subModels", len(subModels)), slog.Any("subModels", subModels), slog.Any("fallback", subModels[0].Fallback))
	return subModels[0].Fallback, nil
}
