package mcp2210

import (
	"github.com/GeertJohan/go.hid"
	
	"errors"
)


// ----------------------------------------------------------------------------------
//  Constants
// ----------------------------------------------------------------------------------

const (
	DirectionIn			= 0x00
	DirectionOut		= 0x01
	
	ValueInactive		= 0x00
	ValueActive			= 0x01
	
	FunctionGPIO		= 0x00
	FunctionChipSelect	= 0x01
	FunctionAlternative	= 0x02
)

// for reference see http://ww1.microchip.com/downloads/en/DeviceDoc/22288A.pdf p11 and following
const (
	cmdSetPinValue		= 0x30
	cmdGetPinValue		= 0x33
)


// ----------------------------------------------------------------------------------
//  Types
// ----------------------------------------------------------------------------------

type MCP2210 struct {
	hidDevice *hid.Device
	currentPinValues uint16
}


// ----------------------------------------------------------------------------------
//  Constructors
// ----------------------------------------------------------------------------------

func Open(vendorId uint16, productId uint16) (*MCP2210, error) {
	// open the hid device
	// TODO: check serial number feature
	device, err := hid.Open(vendorId, productId, "")
	if err != nil {
		return nil, err
	}
	
	// assemble mcp instance
	mcp := MCP2210{ hidDevice: device }
	
	// read back current GPIO pin values
	err = mcp.updateGPIOValues()
	if err != nil {
		return nil, err
	}	
	
	return &mcp, nil
}


// ----------------------------------------------------------------------------------
//  Informational Functions
// ----------------------------------------------------------------------------------

func (this *MCP2210) updateGPIOValues() error {
	if this.hidDevice == nil {
		return errors.New("device not opened")
	}
	
	// assemble and send command
	command := []byte{
		cmdGetPinValue,
	}
	this.hidDevice.Write(command)

	// read the answer (6 bytes long)
	response, err := this.readResponse(6)
	if err != nil {
		return err
	}
	
	// everything is fine, update the GPIO values
	this.currentPinValues = uint16(response[4]) | (uint16(response[5]) << 8)
	return nil
}


// ----------------------------------------------------------------------------------
//  Changing Functions
// ----------------------------------------------------------------------------------

func (this *MCP2210) SetGPIOValue(pin uint16, state uint16) error {
	if this.hidDevice == nil {
		return errors.New("device not opened")
	}
	
	// set the pin state
	if state == ValueActive {
		this.currentPinValues |= 1 << pin	
	} else {
		this.currentPinValues &= ^(1 << pin)
	}
		
	// assemble and send command
	command := []byte{
		cmdSetPinValue,	// opcode
		0x00,			// reserved
		0x00,			// reserved
		0x00,			// reserved
		byte(this.currentPinValues),		// GP 0-7
		byte(this.currentPinValues >> 8),	// GP 8
	}
	this.hidDevice.Write(command)
	
	return nil
}

func (this *MCP2210) Close() {
	this.hidDevice.Close()
}


// ----------------------------------------------------------------------------------
//  Helper Functions
// ----------------------------------------------------------------------------------

func (this *MCP2210) readResponse(length int) ([]byte, error) {
	response := make([]byte, length)
	_, error := this.hidDevice.Read(response)
	return response, error	
}