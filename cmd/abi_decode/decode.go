package main

import (
	"encoding/hex"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type DecodedLog struct {
	Name    string
	Params  []Param
	Address common.Address // contract address
}

type Param struct {
	Name  string
	Value interface{}
	Type  string
}
type MethodData struct {
	Name   string
	Params []Param
}

// ABIDecoder ethereum transaction data decoder
type ABIDecoder struct {
	myabi abi.ABI
}

func NewABIDecoder() *ABIDecoder {
	return &ABIDecoder{}
}

func (d *ABIDecoder) SetABI(contractAbi string) {
	myabi, err := abi.JSON(strings.NewReader(contractAbi))
	if err != nil {
		log.Fatal(err)
	}
	d.myabi = myabi
}

func (d *ABIDecoder) DecodeMethod(txData string) (MethodData, error) {
	if strings.HasPrefix(txData, "0x") {
		txData = txData[2:]
	}

	decodedSig, err := hex.DecodeString(txData[:8])
	if err != nil {
		return MethodData{}, err
	}

	method, err := d.myabi.MethodById(decodedSig)
	if err != nil {
		return MethodData{}, err
	}

	decodedData, err := hex.DecodeString(txData[8:])
	if err != nil {
		return MethodData{}, err
	}

	inputs, err := method.Inputs.Unpack(decodedData)
	if err != nil {
		return MethodData{}, err
	}

	nonIndexedArgs := method.Inputs.NonIndexed()

	retData := MethodData{}
	retData.Name = method.Name
	for i, input := range inputs {
		arg := nonIndexedArgs[i]
		param := Param{
			Name:  arg.Name,
			Value: input,
			Type:  arg.Type.String(),
		}
		retData.Params = append(retData.Params, param)
	}

	return retData, nil
}

func (d *ABIDecoder) ABI() abi.ABI {
	return d.myabi
}
