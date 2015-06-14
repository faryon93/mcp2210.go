// A wrapper library for Microchip's MCP2210 USB-to-SPI Bridge.
// It is heavly based on github.com/GeertJohan/go.hid.
package mcp2210

import (
	"github.com/GeertJohan/go.hid"
	"sync"
)


// ----------------------------------------------------------------------------------
//  Constants
// ----------------------------------------------------------------------------------

// for reference see http://ww1.microchip.com/downloads/en/DeviceDoc/22288A.pdf
// page11 and following
const (
	cmdGetInterrupt		= 0x12
	cmdSetPinValue		= 0x30
	cmdGetPinValue		= 0x31
	cmdSetSPISettings	= 0x40
	cmdGetSPISettins	= 0x41
	cmdTransferSPI		= 0x42
	cmdEEPROMRead		= 0x50
	cmdEEPROMWrite		= 0x51
)


// ----------------------------------------------------------------------------------
//  Types
// ----------------------------------------------------------------------------------

type MCP2210 struct {
	hidDevice *hid.Device
	currentPinValues uint16
	
	spiSettings []byte
	
	mutex sync.Mutex
	xferMutex sync.Mutex
}


// ----------------------------------------------------------------------------------
//  Constructors
// ----------------------------------------------------------------------------------

// Opens a MCP2210 device with the given VendorId and ProductId.
func Open(vendorId uint16, productId uint16) (*MCP2210, error) {
	// open the hid device
	// TODO: check serial number feature
	device, err := hid.Open(vendorId, productId, "")
	if err != nil {
		return nil, err
	}
	
	// assemble mcp instance
	mcp := MCP2210{ hidDevice: device, mutex: sync.Mutex{}, xferMutex: sync.Mutex{} }
	
	// read back current GPIO pin values
	err = mcp.updateGPIOValues()
	if err != nil {
		return nil, err
	}	
	
	// read the spi settings
	err = mcp.updateSPISettings() 
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

// Closes the connection to the device.
func (this *MCP2210) Close() {
	this.hidDevice.Close()
}


// ----------------------------------------------------------------------------------
//  Helper Functions
// ----------------------------------------------------------------------------------

func (this *MCP2210) sendCommand(opcode byte, payload ...byte) ([]byte, error) {
	this.mutex.Lock()
	// send command to mcp
	_, err := this.hidDevice.Write(append([]byte{opcode}, payload...))
	if err != nil {
		return []byte{}, err
	} 
	
	// read the response
	response := make([]byte, 64)
	_, err = this.hidDevice.Read(response)
	this.mutex.Unlock()
	
	return response, err
}