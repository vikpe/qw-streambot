package ezquake

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Process struct {
	clientPort int
	Path       string
}

func (p Process) GetProcessID() int {
	id, err := strconv.Atoi(shellCommand("pgrep", "-fox", p.Path))

	if err != nil {
		return 0
	}

	return id
}

func (p Process) Stop() {
	p.stopWithSignal(os.Interrupt)
}

func (p Process) Kill() {
	p.stopWithSignal(os.Kill)
}

func (p Process) stopWithSignal(signal os.Signal) {
	pid := p.GetProcessID()

	if 0 == pid {
		return
	}

	exec.Command("kill", "-s", signal.String(), string(rune(pid)))
}

func (p Process) IsStarted() bool {
	return 0 != p.GetProcessID()
}

func (p Process) IsConnected() bool {
	return p.IsQtvConnected() && p.IsServerConnected()
}

func (p Process) IsServerConnected() bool {
	return "" != p.GetUdpSocketAddress()
}

func (p Process) IsQtvConnected() bool {
	return "" != p.GetTcpSocketAddress()
}

func (p Process) GetSocketAddress() string {
	tcpHostname := p.GetTcpSocketAddress()

	if "" == tcpHostname {
		return p.GetUdpSocketAddress()
	}

	return tcpHostname
}

func (p Process) GetTcpSocketAddress() string {
	pid := p.GetProcessID()

	if 0 == pid {
		return ""
	}

	//identifier := fmt.Sprintf("%d/ezquake", pid)

	asdasd := shellCommand("netstat", "-put", "2>/dev/null")

	// args := []string{fmt.Sprintf("-nput 2>/dev/null | grep %s | awk '{{print $5}}'", identifier)}
	return asdasd
}

func (p Process) GetUdpSocketAddress() string {
	if !p.IsStarted() {
		return ""
	}

	ipAddress := GetLocalIP()
	args := fmt.Sprintf("dst %s and dst port %d -n 1 -q | grep '#1'", ipAddress, p.clientPort)
	output := bashCommandWithTimeout(500*time.Millisecond, "sudo ngrep", args)

	if len(output) > 0 {
		return output[2:strings.Index(output, " -> ")]
	}

	return ""
}

func shellCommand(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("SHELL ERROR", err, cmd.String())
		return ""
	}

	return strings.TrimSpace(string(out))
}

func bashCommand(name string, args ...string) string {
	bashArgs := append([]string{"-c", fmt.Sprintf(`"%s %s"`, name, strings.Join(args, " "))})
	cmd := exec.Command("bash", bashArgs...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("SHELL ERROR", err, cmd.String())
		return ""
	}

	return strings.TrimSpace(string(out))
}

func bashCommandWithTimeout(timeout time.Duration, name string, args ...string) string {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	out, err := cmd.Output()

	if ctx.Err() == context.DeadlineExceeded || err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
