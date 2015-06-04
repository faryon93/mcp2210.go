package mcp2210

import (
	"errors"
)


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
