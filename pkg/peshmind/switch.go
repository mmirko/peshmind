package peshmind

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"
)

type Switch struct {
	Name        string `json:"name"`
	Mac         string `json:"mac"`
	Description string `json:"description"`
	IP          string `json:"ip"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Model       string `json:"model"`
	Data        string `json:"data"`
}

type switchTemplate struct {
	Name     string
	IP       string
	Username string
	Password string
	funcMap  template.FuncMap
}

type SwitchPort struct {
	SwitchID string
	Name     string
	EndPoint *SwitchPort
}

type SwitchDOT struct {
	ID    string
	Ports map[string]*SwitchPort
}

func (s *Switch) modelTemplate() (string, error) {
	switch s.Model {
	case "HP-Aruba":
		return HPArubaExpScript, nil
	case "HP-Aruba-Press":
		return HPArubaExpScriptPress, nil
	default:
		return "", fmt.Errorf("unsupported switch model: %s", s.Model)
	}
}

func (s *Switch) createTemplate() *switchTemplate {
	return &switchTemplate{
		Name:     s.Name,
		IP:       s.IP,
		Username: s.Username,
		Password: s.Password,
		funcMap:  template.FuncMap{},
	}
}

func (s *Switch) ApplyTemplate() error {
	// Get the template for the switch model
	modelTemplate, err := s.modelTemplate()
	if err != nil {
		return err
	}
	// Create a new template and parse the model template
	tmpl := s.createTemplate()
	t, err := template.New("switch").Funcs(tmpl.funcMap).Parse(modelTemplate)
	if err != nil {
		return err
	}
	// Execute the template with the switch data
	var result bytes.Buffer
	err = t.Execute(&result, tmpl)
	if err != nil {
		return err
	}
	// Store the result in the switch data field
	s.Data = result.String()
	return nil
}

func (s *Switch) FetchData() error {
	// Choose a temporary directory for the switch data
	tempDir, err := os.MkdirTemp("", "switch-data")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir) // Clean up the temp directory after use

	// Create a file to store the switch fectch script
	scriptFileName := fmt.Sprintf("%s/%s", tempDir, s.Name)
	scriptFile, err := os.Create(scriptFileName)
	if err != nil {
		return fmt.Errorf("failed to create script file: %w", err)
	}
	// Write the switch data to the script file
	_, err = scriptFile.WriteString(s.Data)
	if err != nil {
		return fmt.Errorf("failed to write script file: %w", err)
	}
	scriptFile.Close()

	// Set the execution permissions for the script file
	err = os.Chmod(scriptFileName, 0700)
	if err != nil {
		return fmt.Errorf("failed to set script file permissions: %w", err)
	}

	// Execute the script file
	outBuffer := new(bytes.Buffer)
	cmd := exec.Command(scriptFileName)
	cmd.Dir = tempDir
	cmd.Stdout = outBuffer
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute script file: %w", err)
	}

	// Read the output from the command and store it in the switch data field
	s.Data = outBuffer.String()

	return nil
}

func (s *Switch) SaveKB(kbpool string) error {
	// Save the switch data to a file
	fileName := fmt.Sprintf("%s/%s.pl", kbpool, s.Name)
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create KB file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("switch(a%s).\n", s.Mac))
	_, err = file.WriteString(fmt.Sprintf("switchname(a%s, '%s').\n", s.Mac, s.Name))

	for _, line := range strings.Split(s.Data, "\n") {
		if len(line) > 0 {
			re := regexp.MustCompile("\\s+(?P<mac>[a-z0-9]+-[a-z0-9]+)\\s+(?P<port>[a-zA-Z0-9]+)\\s+(?P<vlan>[0-9]+)\\s+")
			if re.MatchString(line) {
				mac := re.ReplaceAllString(string(line), "${mac}")
				mac = strings.ReplaceAll(mac, "-", "")
				port := re.ReplaceAllString(string(line), "${port}")
				port = strings.ToLower(port)
				// vlan := re.ReplaceAllString(string(line), "${vlan}")
				// Write the data to the file
				_, err := file.WriteString(fmt.Sprintf("seen(a%s,a%s,%s).\n", s.Mac, mac, port))
				if err != nil {
					return fmt.Errorf("failed to write to KB file: %w", err)
				}
				continue
			}

			re = regexp.MustCompile("\\s+(?P<mac>[a-z0-9]+-[a-z0-9]+)\\s+(?P<port>[a-zA-Z0-9]+)\\s+")
			if re.MatchString(line) {
				mac := re.ReplaceAllString(string(line), "${mac}")
				mac = strings.ReplaceAll(mac, "-", "")
				port := re.ReplaceAllString(string(line), "${port}")
				port = strings.ToLower(port)
				// vlan := re.ReplaceAllString(string(line), "${vlan}")
				// Write the data to the file
				_, err := file.WriteString(fmt.Sprintf("seen(a%s,a%s,%s).\n", s.Mac, mac, port))
				if err != nil {
					return fmt.Errorf("failed to write to KB file: %w", err)
				}
			}
		}
	}

	return nil
}
