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
//  Changing Functions
// ----------------------------------------------------------------------------------

func (this *MCP2210) Xfer(sendBuffer []byte) ([]byte, error) {
	if this.hidDevice == nil {
		return nil, errors.New("device not opened")
	}
	
	if (len(sendBuffer) > maxBytesPerWrite) {
		return nil, errors.New("cannot send more than 60 bytes at once")
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