package ezquake_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/streambot/ezquake"
)

var proc = ezquake.Process{Path: "/home/vikpe/code/ezquake-api/quake2/ezquake-linux-x86_64"}

func TestProcess_GetProcessID(t *testing.T) {
	assert.NotNil(t, proc.GetProcessID())
}

func TestProcess_IsStarted(t *testing.T) {
	assert.True(t, proc.IsStarted())
}

func TestHEHE(t *testing.T) {
	fmt.Println("GetProcessID", proc.GetProcessID())
	fmt.Println("IsStarted", proc.IsStarted())
	fmt.Println("IsConnected", proc.IsConnected())
	fmt.Println("GetSocketAddress", proc.GetSocketAddress())
	fmt.Println("GetUdpSocketAddress", proc.GetUdpSocketAddress())
	fmt.Println("GetTcpSocketAddress", proc.GetTcpSocketAddress())

	t.Fatal()
}
