package mcp2210

import (
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


// ----------------------------------------------------------------------------------
//  Setters
// ----------------------------------------------------------------------------------

func (this *MCP2210) setCurrentPinValues(low uint16, high uint16) {
	this.currentPinValues = low | (high<< 8)
}


// ----------------------------------------------------------------------------------
//  Informational Functions
// ----------------------------------------------------------------------------------

func (this *MCP2210) updateGPIOValues() error {
	if this.hidDevice == nil {
		return errors.New("device not opened")
	}
	
	// assemble and send command
	response, err := this.sendCommand(cmdGetPinValue)
	if err != nil {
		return err
	}
	
	// everything is fine, update the GPIO values
	this.setCurrentPinValues(uint16(response[4]), uint16(response[5]))
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
	
	// send the command
	response, err := this.sendCommand(
		cmdSetPinValue,	// opcode
		0x00,			// reserved
		0x00,			// reserved
		0x00,			// reserved
		byte(this.currentPinValues),		// GP 0-7
		byte(this.currentPinValues >> 8),	// GP 8			
	)
	if err != nil {
		return err
	}
		
	// set the actual GPIO values
	this.setCurrentPinValues(uint16(response[4]), uint16(response[5]))
		
	return nil
}
