package mcp2210

import (
	"errors"
)


// ----------------------------------------------------------------------------------
//  Constants
// ----------------------------------------------------------------------------------

// GPIO constants
const (
	// Pin direction: input
	DirectionIn			= 0x00
	// Pin direction: output
	DirectionOut		= 0x01
	
	// Pin value: inactive/low
	ValueInactive		= 0x00
	// Pin value: active/high
	ValueActive			= 0x01
	
	// Pin function: GPIO
	FunctionGPIO		= 0x00
	// Pin function: SPI Chip-Select
	FunctionChipSelect	= 0x01
	// Pin function: Alternative
	FunctionAlternative	= 0x02
)

const (
	counterReset	= 0x00
	counterKeep		= 0x26
)


// ----------------------------------------------------------------------------------
//  Setters
// ----------------------------------------------------------------------------------

func (this *MCP2210) setCurrentPinValues(low uint16, high uint16) {
	this.currentPinValues = low | (high << 8)
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

// Get the current value of a GPIO pin.
// Returns ValueInactive or ValueActive. 
func (this *MCP2210) GetGPIOValue(pin uint16) (uint8, error) {
	if this.hidDevice == nil {
		return 0xFF, errors.New("device not opened")
	}
	
	// update the GPIO values to get the most recent state
	err := this.updateGPIOValues()
	if err != nil {
		return 0xFF, err
	}
	
	// 1 = active/high, 0 = inactive/low
	return uint8((this.currentPinValues & (1 << pin)) >> pin), nil
}

// Gets the external interrupt counter value and resets it to zero.
// Interrupt signaling is available on GP6, if it is configured
// with alternative pin function.
func (this *MCP2210) GetInterruptCount() (uint16, error) {
	if this.hidDevice == nil {
		return 0xFF, errors.New("device not opened")
	}

	// submit command
	response, err := this.sendCommand(
		cmdGetInterrupt,
		counterReset,	// reset the counter
	)
	if err != nil {
		return 0xFF, err
	}

	return uint16(response[4]) | (uint16(response[5]) << 8), nil
}


// ----------------------------------------------------------------------------------
//  Changing Functions
// ----------------------------------------------------------------------------------

// Sets the value of a GPIO pin.
// pin Number of the pin to set.
// state New pinstate, use ValueInactive or ValueActive
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
