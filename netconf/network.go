package netconf

import (
	"encoding/json"
	"io/ioutil"
)

// Network defines a Scrit network.
type Network struct {
	NetworkEpochs []NetworkEpoch // global list of signing epochs
}

// Load a network configuration from filename and return the Network struct.
func Load(filename string) (*Network, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var net Network
	if err := json.Unmarshal(data, &net); err != nil {
		return nil, err
	}
	return &net, err
}

// Validate the net configuration.
func (net *Network) Validate() error {
	// validate each network epoch
	for _, e := range net.NetworkEpochs {
		if err := e.Validate(); err != nil {
			return err
		}
	}
	// validate network epoch transitions
	for i := 1; i < len(net.NetworkEpochs); i++ {
		// sign end n-1 == sign start n
		if net.NetworkEpochs[i-1].SignEnd != net.NetworkEpochs[i].SignStart {
			return ErrSignEpochWrongBoundaries
		}
	}

	// validation end n <= sign end n+1
	// TODO

	return nil
}

// Marshal net as string.
func (net *Network) Marshal() string {
	jsn, err := json.MarshalIndent(net, "", "  ")
	if err != nil {
		panic(err) // should never happen
	}
	return string(jsn)
}
