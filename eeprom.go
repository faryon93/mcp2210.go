package mcp2210

import (
	"errors"
)


// ----------------------------------------------------------------------------------
//  Informational Functions
// ----------------------------------------------------------------------------------

func (this *MCP2210) ReadEEPROM(address byte) (byte, error) {
	if this.hidDevice == nil {
		return 0xFF, errors.New("device not opened")
	}
	
	// send command to read EEPROM
	response, err := this.sendCommand(
		cmdEEPROMRead,
		address,
	)
	if err != nil {
		return 0xFF, err
	}
	
	// Index 3 from response represents the
	// byte-value at the desired address
	return response[3], nil
}


// ----------------------------------------------------------------------------------
//  Changing Functions
// ----------------------------------------------------------------------------------

func (this *MCP2210) WriteEEPROM(address byte, data byte) (error) {
	if this.hidDevice == nil {
		return errors.New("device not opened")
	}	
	
	// send command to mcp
	response, err := this.sendCommand(
		cmdEEPROMWrite,
		address,
		data,
	)
	if err != nil {
		return err
	}
	
	// check if write was successfull
	switch response[1] {
		case 0x00:
			return nil
		
		case 0xFA:
			return errors.New("EEPROM write failed")
			
		case 0xFB:
			return errors.New("EEPROM is protected or permanently locked")		
			
		default:
			return errors.New("Unknown error occoured")
	}
}