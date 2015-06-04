package mcp2210

import (
	"github.com/GeertJohan/go.hid"
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
	cmdTransferSPI		= 0x42
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


// ----------------------------------------------------------------------------------
//  Changing Functions
// ----------------------------------------------------------------------------------

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

func (this *MCP2210) sendCommand(opcode byte, payload ...byte) ([]byte, error) {
	// send command to mcp
	_, err := this.hidDevice.Write(append([]byte{opcode}, payload...))
	if err != nil {
		return []byte{}, err
	} 
	
	// read the response
	response := make([]byte, 64)
	_, err = this.hidDevice.Read(response)
	return response, err
}