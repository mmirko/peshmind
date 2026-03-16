package peshmind

import (
	"context"
	"errors"
	"fmt"
	"os"

	graphviz "github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

func (c *Config) Simulate(graphFile string) error {

	if !elementInSlice(graphFile, c.Simulations) {
		return errors.New("Graph file not in simulations list")
	}

	ctx := context.Background()
	var graph *graphviz.Graphviz

	if g, err := graphviz.New(ctx); err != nil {
		return err
	} else {
		graph = g
	}

	var data []byte

	// Load a graph from a file
	if graphFile != "" {
		var err error
		data, err = os.ReadFile(graphFile)
		if err != nil {
			return err
		}
	} else {
		return errors.New("No graph file specified")
	}

	g, err := graphviz.ParseBytes(data)
	if err != nil {
		return err
	}

	defer func() {
		if err := graph.Close(); err != nil {
			return
		}
		g.Close()
	}()

	switches := make(map[string]*cgraph.Node)
	switchPorts := make(map[string]map[string]string) // Map of switch name to port to either:
	// - connected switch name
	// - a mac address (for end devices)
	// - empty string or missing key for unconnected ports

	// Iterate over nodes to build the switches map and initialize the switchPorts map
	for n, _ := g.FirstNode(); n != nil; n, _ = g.NextNode(n) {
		sName, _ := n.Name()
		if _, ok := c.Switches[sName]; !ok {
			return errors.New("Switch " + sName + " not found in configuration")
		}

		fmt.Println("Switch name:", sName)
		switches[sName] = n
		switchPorts[sName] = make(map[string]string)
	}

	// Iterate over edges to build the switchPorts mapping
	for _, v := range switches {
		for e, _ := g.FirstEdge(v); e != nil; e, _ = g.NextEdge(e, v) {
			eName, _ := e.Name()
			if eName != "" {
				destName, _ := e.Node().Name()
				srcName, _ := v.Name()
				portLink := e.GetStr("port")
				if portLink == "" {
					return errors.New("Port attribute missing for edge " + eName)
				}
				fmt.Println("Edge name:", eName, "from", srcName, "to", destName, "port:", portLink)

				switchPorts[srcName][portLink] = destName
			}
		}
	}

	fmt.Println("Switch ports mapping:", switchPorts)

	swId := 0
	// Populate the switchPorts mapping with fake MAC addresses for unconnected ports
	for sName, ports := range switchPorts {
		for p := 1; p <= c.Switches[sName].Port; p++ {
			portName := fmt.Sprintf("%d", p)
			if _, ok := ports[portName]; !ok {
				// Generate a fake MAC address for this unconnected port
				fakeMac := fmt.Sprintf("00000000%02x%02x", swId, p)
				switchPorts[sName][portName] = fakeMac

			}
		}
		swId++
	}

	fmt.Println(c.getPortMacs(switchPorts, "privsw-0-1", "23"))

	return nil

}

func (c *Config) getPortMacs(switchPorts map[string]map[string]string, switchName string, portName string) []string {
	if _, ok := switchPorts[switchName]; !ok {
		return []string{}
	}

	if _, ok := switchPorts[switchName][portName]; !ok {
		return []string{}
	}

	connected := switchPorts[switchName][portName]
	if connected == "" {
		return []string{}
	}

	// If the connected value is a MAC address (not in switches map), return it
	if _, ok := switchPorts[connected]; !ok {
		return []string{connected}
	} else {
		// If the connected value is another switch, we need to find the corresponding port on that switch
		return c.getSwitchMacs(switchPorts, connected, portName)
	}
}

func (c *Config) getSwitchMacs(switchPorts map[string]map[string]string, switchName string, excludePort string) []string {
	macs := []string{}
	if _, ok := switchPorts[switchName]; ok {
		for portName, connected := range switchPorts[switchName] {
			if portName == excludePort {
				continue
			}
			if connected == "" {
				continue
			}
			if _, ok := switchPorts[connected]; !ok {
				macs = append(macs, connected)
			} else {
				macs = append(macs, c.getSwitchMacs(switchPorts, connected, portName)...)
			}
		}
	}
	return macs
}
