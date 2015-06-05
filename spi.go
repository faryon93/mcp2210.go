package mcp2210

import (
	"errors"
)


// ----------------------------------------------------------------------------------
//  Constants
// ----------------------------------------------------------------------------------

const (
	maxBytesPerWrite	= 60
	
	spiTransferFinished = 0x10
	spiTransferSuccess  = 0x00
)


// ----------------------------------------------------------------------------------
//  Informational Functions
// ----------------------------------------------------------------------------------

// Returns the currently used SPI clock in Hz.
func (this *MCP2210) GetSPIClock() uint32 {
	return uint32(this.spiSettings[3]) << 24 |
		   uint32(this.spiSettings[2]) << 16 |
		   uint32(this.spiSettings[1]) << 8 |
		   uint32(this.spiSettings[0])
}

func (this *MCP2210) GetBytesPerTransfer() int {	
	return int(this.spiSettings[15]) << 8 | 
		   int(this.spiSettings[14])
		   
}


// ----------------------------------------------------------------------------------
//  Changing Functions
// ----------------------------------------------------------------------------------

// Sends some bytes and returns the received bytes
func (this *MCP2210) Xfer(sendBuffer []byte) ([]byte, error) {
	if this.hidDevice == nil {
		return nil, errors.New("device not opened")
	}
	
	if (len(sendBuffer) > maxBytesPerWrite) {
		return nil, errors.New("cannot send more than 60 bytes at once")
	}
	
	// reset the bytes per transfer if needed
	if (this.GetBytesPerTransfer() != len(sendBuffer)) {
		this.setTransferBytes(uint16(len(sendBuffer)))
	}
	
	// send the initial transfer command with the data to send
	response, err := this.sendCommand(
		cmdTransferSPI,	// opcode			
		append([]byte{
			byte(len(sendBuffer)),	// number of bytes to send
			0x00,			// reserved
			0x00,			// reserved
		},
		sendBuffer...)...,
	)
	if err != nil {
		return nil, err
	}
	
	// data sent by the slave
	var slaveData []byte 
	
	// read the data sent by the slave
	for response[1] == spiTransferSuccess &&
		response[3] != spiTransferFinished {		
		response, err = this.sendCommand(
			cmdTransferSPI,	// opcode
			byte(0),	// number of bytes to send
		)
		
		// the received bytes start at index 4 and going until 4 + byte-count
		slaveData = append(slaveData, response[4:(4 + response[2])]...)
	} 
	
	return slaveData, nil
}

func (this *MCP2210) setTransferBytes(byteCount uint16) error {
	if this.hidDevice == nil {
		return errors.New("device not opened")
	}
	
	this.spiSettings[14] = byte(byteCount & 0xFF)
	this.spiSettings[15] = byte(byteCount >> 8)
		
	// assemble and send command
	_, err := this.sendCommand(
		cmdSetSPISettings,
		append([]byte{
			0x00,		// reserved
			0x00,		// reserved
			0x00,		// reserved
		},
		this.spiSettings...)...,
	)
	if err != nil {
		return err
	}
	
	return nil
}

func (this *MCP2210) updateSPISettings() error {
	if this.hidDevice == nil {
		return errors.New("device not opened")
	}	
	
	// get the current SPI settings
	response, err := this.sendCommand(
		cmdGetSPISettins,
	)
	if err != nil {
		return  err
	}
	
	this.spiSettings = response[4:]
	return nil
}