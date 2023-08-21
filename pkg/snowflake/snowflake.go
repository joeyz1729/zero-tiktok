package snowflake

import (
	"errors"
	sf "github.com/sony/sonyflake"
	"time"
)

var (
	sonyflake     *sf.Sonyflake
	sonyStartTime string
	sonyMachineId uint16
)

func getMachineID() (uint16, error) {
	return sonyMachineId, nil
}

func Init(startTime string, machineId uint16) (err error) {
	sonyStartTime = startTime
	sonyMachineId = machineId
	t, _ := time.Parse(sonyStartTime, "2023-01-01")
	settings := sf.Settings{
		StartTime: t,
		MachineID: getMachineID,
	}
	sonyflake = sf.NewSonyflake(settings)
	return
}

func GenID() (id uint64, err error) {
	if sonyflake == nil {
		err = errors.New("sonyflake not initialized")
		return
	}
	id, err = sonyflake.NextID()
	return
}
