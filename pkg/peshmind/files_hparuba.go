package peshmind

const (
	HPArubaExpScript = `#!/usr/bin/expect -f 

log_user 0
set env(TERM) vt100
spawn /usr/bin/telnet {{ .IP }}
expect "Username:"
send "{{ .Username }}\r"
expect "Password:"
send "{{ .Password }}\r"
sleep 3
expect "{{ .Name }}#"
send "term leng 1000\r"
expect "{{ .Name }}#"
send "show mac-address\r"
log_user 1
expect "{{ .Name }}#"
log_user 0
send "exit\r"
expect ">"
send "exit\r"
expect "?"
send "y\r"
exit 0
`
)
