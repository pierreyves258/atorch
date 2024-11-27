package atorch

import (
	"encoding/binary"
	"fmt"
	"math"
	"sync"
	"time"

	"go.bug.st/serial"
)

type PX100 struct {
	sync.Mutex
	port serial.Port
}

const (
	GetIsOn byte = iota + 0x10
	GetVoltage
	GetCurrent
	GetTime
	GetCharge
	GetEnergy
	GetTemperature
	GetCurrentLimit
	GetVoltageLimit
	GetTimer
)

const (
	SetOutput byte = iota + 0x01
	SetCurrent
	SetCutoff
	SetMaxTime
	Reset
)

var setToGet = map[byte]byte{
	SetOutput:  GetIsOn,
	SetCurrent: GetCurrentLimit,
	SetCutoff:  GetVoltageLimit,
	SetMaxTime: GetTime,
	Reset:      GetCharge,
}

func NewPX100(tty string) (*PX100, error) {
	mode := &serial.Mode{
		BaudRate: 9600,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open(tty, mode)
	if err != nil {
		return nil, err
	}
	port.SetReadTimeout(1 * time.Second) // Broadcast packet every 1s

	return &PX100{
		port: port,
	}, nil
}

func (psu *PX100) Destroy() {
	psu.port.Close()
	psu = nil
}

func (px *PX100) GetData(command byte) (interface{}, error) {
	if px == nil {
		return nil, fmt.Errorf("not initialised")
	}

	err := px.sendCommand(command, []byte{0x00, 0x00})
	if err != nil {
		return 0, err
	}

	raw, err := px.readReply()
	if err != nil {
		return 0, err
	}

	switch command {
	case GetVoltage, GetCurrent, GetEnergy, GetCharge:
		return float64(parseReply(raw)) / 1000, nil
	case GetIsOn:
		return raw[4] == 1, nil
	case GetTemperature:
		return float64(parseReply(raw)), nil
	case GetCurrentLimit, GetVoltageLimit:
		return float64(parseReply(raw)) / 100, nil
	case GetTime, GetTimer:
		return time.Duration(raw[2])*time.Hour + time.Duration(raw[3])*time.Minute + time.Duration(raw[4])*time.Second, nil
		// default:
		// 	return fmt.Sprintf("%x\n", raw), nil
	}

	return nil, fmt.Errorf("unknown command %x", command)
}

func (px *PX100) readReply() ([]byte, error) {
	rawData := []byte{}
	buff := make([]byte, 1)

	timeout := time.After(time.Second)

	px.Lock()
	defer px.Unlock()

	for {
		n, err := px.port.Read(buff)
		if err != nil {
			return nil, err
		}
		if n == 0 {
			return nil, fmt.Errorf("EOF")
		}

		rawData = append(rawData, buff...)
		// fmt.Printf("%d %x , %d\n", n, rawData, len(rawData))

		if buff[0] == 0xca {
			rawData = []byte{0xca} // reset rawData
		}

		if rawData[0] == 0xca && len(rawData) == 7 {
			if rawData[5] != 0xce || rawData[6] != 0xcf {
				return nil, fmt.Errorf("invalid response")
			}
			return rawData, nil
		}

		select {
		case <-timeout:
			return nil, fmt.Errorf("no response")
		default:
		}
	}
}

func (px *PX100) sendCommand(command byte, payload []byte) error {
	if len(payload) != 2 {
		return fmt.Errorf("invalid payload len")
	}

	px.Lock()
	defer px.Unlock()

	px.port.ResetInputBuffer() // Flush read buffer
	n, err := px.port.Write([]byte{0xb1, 0xb2, command, payload[0], payload[1], 0xb6})
	if err != nil {
		return err
	}
	if n != 6 {
		return fmt.Errorf("%d byte send instead of 6", n)
	}
	time.Sleep(100 * time.Millisecond)

	return nil
}

func (px *PX100) SetData(command byte, value interface{}, ensure bool) error {
	payload := make([]byte, 2)

	switch tVal := value.(type) {
	case float64:
		i, f := math.Modf(tVal)
		payload[0] = byte(i)
		payload[1] = byte(float32(f) * 100)
	case bool:
		if tVal {
			payload[0] = 1
		}
	case time.Duration:
		sec := uint16(tVal.Seconds())
		binary.BigEndian.PutUint16(payload, sec)
	}
	// default [0, 0]

	err := px.sendCommand(command, payload)
	if err != nil {
		return err
	}

	if ensure {
		res, err := px.GetData(setToGet[command])
		if err != nil {
			return err
		}
		if command == Reset {
			value = .0 // After a reset capacity must be 0
		}
		// fmt.Printf("%d %v | %d %v\n", command, value, setToGet[command], res)
		if res != value {
			fmt.Printf("%v != %v, retry\n", value, res)
			px.SetData(command, value, ensure) // retry
		}
	}

	return nil
}

func parseReply(raw []byte) uint32 {
	return binary.BigEndian.Uint32([]byte{0x00, raw[2], raw[3], raw[4]})
}
