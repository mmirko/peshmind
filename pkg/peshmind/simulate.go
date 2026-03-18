package peshmind

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	graphviz "github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

type Simulation struct {
	*Config
	SimSwitches    map[string]*cgraph.Node      // Map of switch name to its graph node
	SimSwitchesMac map[string]string            // Map of switch name to its mac address
	SimSwitchPorts map[string]map[string]string // Map of switch name to port to either:
	// - connected switch name
	// - a mac address (for end devices)
	// - empty string or missing key for unconnected ports
}

func NewSimulation(config *Config) *Simulation {
	return &Simulation{
		Config: config,
	}
}

func (s *Simulation) Simulate(graphFile string) error {

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

	switches := make(map[string]*cgraph.Node)         // Map of switch name to its graph node
	switchesMac := make(map[string]string)            // Map of switch name to its mac address
	switchesSize := make(map[string]int)              // Map of switch name to number of ports
	switchesExtraPorts := make(map[string][]string)   // Map of switch name to list of extra ports
	switchPorts := make(map[string]map[string]string) // Map of switch name to port to either:
	// - connected switch name
	// - a mac address (for end devices)
	// - empty string or missing key for unconnected ports

	// Iterate over nodes to build the switches map and initialize the switchPorts map
	for n, _ := g.FirstNode(); n != nil; n, _ = g.NextNode(n) {
		sName, _ := n.Name()
		ghost := n.GetStr("ghost")
		if ghost != "true" {
			if _, ok := s.Switches[sName]; !ok {
				return errors.New("Switch " + sName + " not found in configuration")
			}
		}

		if s.Debug {
			fmt.Println("Found switch in graph:", sName)
		}
		switches[sName] = n
		if ghost != "true" {
			switchesMac[sName] = s.Switches[sName].Mac
			switchesSize[sName] = s.Switches[sName].Port
			switchesExtraPorts[sName] = s.Switches[sName].ExtraPorts
		} else {
			// For ghost switches, we need to get the mac address from the node attributes
			switchesMac[sName] = n.GetStr("mac")
			if switchesMac[sName] == "" {
				return errors.New("MAC address missing for ghost switch " + sName)
			}

			// Set the number of ports for ghost switches to 0 by default, but we can override it if the "ports"
			// attribute is present on the node in the graph
			switchesSize[sName] = 0
			ports := n.GetStr("ports")
			if ports != "" {
				// If the "ports" attribute is present, we can use it to determine the number of ports for this ghost switch
				if n, err := strconv.Atoi(ports); err == nil {
					switchesSize[sName] = n
				}
			}

			// Similarly, we can check for extra ports on ghost switches using the "extraports" attribute
			switchesExtraPorts[sName] = []string{}
			extraPorts := n.GetStr("extraports")
			if extraPorts != "" {
				switchesExtraPorts[sName] = append(switchesExtraPorts[sName], strings.Split(extraPorts, ",")...)
			}
		}
		if s.Debug {
			fmt.Println("Switch", sName, "has", switchesSize[sName], "ports and extra ports:", switchesExtraPorts[sName])
		}
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
				if s.Debug {
					fmt.Println("Edge name:", eName, "from", srcName, "to", destName, "port:", portLink)
				}

				switchPorts[srcName][portLink] = destName
			}
		}
	}

	if s.Debug {
		fmt.Println("Switch ports mapping:", switchPorts)
	}

	swId := 0
	// Populate the switchPorts mapping with fake MAC addresses for unconnected ports
	for sName, ports := range switchPorts {
		for p := 1; p <= switchesSize[sName]; p++ {
			portName := fmt.Sprintf("%d", p)
			if _, ok := ports[portName]; !ok {
				// Use sim_generate_percentage to determine whether to generate a fake MAC address
				if s.SimGeneratePercentage < 100 {
					if rand.Intn(100) >= s.SimGeneratePercentage {
						continue
					}
				}

				// Generate a fake MAC address for this unconnected port
				fakeMac := fmt.Sprintf("00000000%02x%02x", swId, p)
				switchPorts[sName][portName] = fakeMac

			}
		}
		swId++
	}

	s.SimSwitches = switches
	s.SimSwitchesMac = switchesMac
	s.SimSwitchPorts = switchPorts

	return nil

}

func (s *Simulation) getPortMacs(switchName string, portName string) []string {
	switchPorts := s.SimSwitchPorts
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
		return s.getSwitchMacs(connected, portName)
	}
}

func (s *Simulation) getSwitchMacs(switchName string, excludePort string) []string {
	switchPorts := s.SimSwitchPorts
	macs := []string{}
	if _, ok := switchPorts[switchName]; ok {
		macs = append(macs, s.SimSwitchesMac[switchName])
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
				macs = append(macs, s.getSwitchMacs(connected, portName)...)
			}
		}
	}
	return macs
}

func (s *Simulation) EmitDot() (string, error) {
	// Emit a DOT file representing the simulated network topology, best viewed with fdp
	result := `digraph G {
graph [
    bgcolor="#f8fafc",
    pad="0.3",
    nodesep="0.45",
    ranksep="0.75",
    splines=true,
    overlap=false,
    fontname="Helvetica"
  ];

node [
    shape=invtrapezium,
    style="rounded,filled",
    fillcolor="#ecfeff",
    color="#06b6d4",
    fontcolor="#0f172a",
    fontname="Helvetica",
    fontsize=11,
    penwidth=1.2
  ];
  
edge [
    color="#64748b",
    fontcolor="#334155",
    fontname="Helvetica",
    fontsize=10,
    penwidth=1.3,
    arrowsize=0.6
  ];
  `
	for sName, ports := range s.SimSwitchPorts {
		result += fmt.Sprintf("  \"%s\" [shape=box3d, style=\"rounded,filled\", fillcolor=\"#0ea5e9\", color=\"#0369a1\", fontcolor=\"#ffffff\", penwidth=1.7];\n", sName)
		for portName, connected := range ports {
			if connected == "" {
				continue
			}
			if _, ok := s.SimSwitches[connected]; !ok {
				// Connected to a MAC address (end device)
				result += fmt.Sprintf("  \"%s\" -> \"%s\" [label=\"%s\"];\n", sName, connected, portName)
			} else {
				// Connected to another switch
				result += fmt.Sprintf("  \"%s\" -> \"%s\" [label=\"%s\"];\n", sName, connected, portName)
			}
		}
	}
	result += "}\n"
	return result, nil
}

func (s *Simulation) EmitOutput() (string, error) {
	// Emit a simple text output of the simulated network topology
	var result strings.Builder
	for sName := range s.SimSwitches {
		switchMac := s.SimSwitchesMac[sName]
		isGhost := false
		if _, ok := s.Switches[sName]; !ok {
			isGhost = true
		}
		if !isGhost {
			fmt.Fprintf(&result, "switch(a%s).\n", switchMac)
			fmt.Fprintf(&result, "switchname(a%s, '%s').\n", switchMac, sName)
			for portName, connected := range s.SimSwitchPorts[sName] {
				if connected == "" {
					continue
				}
				macs := s.getPortMacs(sName, portName)
				for _, mac := range macs {
					fmt.Fprintf(&result, "seen(a%s,a%s,%s).\n", switchMac, mac, portName)
				}
			}
		}
	}
	return result.String(), nil
}
