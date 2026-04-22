package peshmind

const (
	QuantaExpScript = `#!/usr/bin/expect -f 

log_user 0
set env(TERM) vt100
spawn /usr/bin/telnet {{ .IP }}
expect "User:"
send "{{ .Username }}\r"
expect "Password:"
send "{{ .Password }}\r"
sleep 3
expect "({{ .Name }}) >"
send "enable\r"
expect "Password:"
send "\r"
expect "({{ .Name }}) >"
send "terminal length 0\r"
expect "({{ .Name }}) >"
send "show mac-addr-table\r"
log_user 1
expect {
    -re {--More--|More:} {
        send " "
        exp_continue
    }
    -re {# $} {}
}
expect "({{ .Name }}) >"
log_user 0
send "quit\r"
exit 0
`
)
