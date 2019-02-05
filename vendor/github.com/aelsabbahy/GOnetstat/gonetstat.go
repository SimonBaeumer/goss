/*
   Simple Netstat implementation.
   Get data from /proc/net/tcp and /proc/net/udp and
   and parse /proc/[0-9]/fd/[0-9].

   Author: Rafael Santos <rafael@sourcecode.net.br>
*/

package GOnetstat

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/opencontainers/runc/libcontainer/user"
)

const (
	PROC_TCP  = "/proc/net/tcp"
	PROC_UDP  = "/proc/net/udp"
	PROC_TCP6 = "/proc/net/tcp6"
	PROC_UDP6 = "/proc/net/udp6"
)

var STATE = map[string]string{
	"01": "ESTABLISHED",
	"02": "SYN_SENT",
	"03": "SYN_RECV",
	"04": "FIN_WAIT1",
	"05": "FIN_WAIT2",
	"06": "TIME_WAIT",
	"07": "CLOSE",
	"08": "CLOSE_WAIT",
	"09": "LAST_ACK",
	"0A": "LISTEN",
	"0B": "CLOSING",
}

type Process struct {
	User        string
	Name        string
	Pid         string
	Exe         string
	State       string
	Ip          string
	Port        int64
	ForeignIp   string
	ForeignPort int64
}

func getData(t string) ([]string, error) {
	// Get data from tcp or udp file.

	var proc_t string

	if t == "tcp" {
		proc_t = PROC_TCP
	} else if t == "udp" {
		proc_t = PROC_UDP
	} else if t == "tcp6" {
		proc_t = PROC_TCP6
	} else if t == "udp6" {
		proc_t = PROC_UDP6
	} else {
		fmt.Printf("%s is a invalid type, tcp and udp only!\n", t)
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(proc_t)
	if err != nil {
		return []string{}, err
	}
	lines := strings.Split(string(data), "\n")

	// Return lines without Header line and blank line on the end
	return lines[1 : len(lines)-1], nil

}

func hexToDec(h string) int64 {
	// convert hexadecimal to decimal.
	d, err := strconv.ParseInt(h, 16, 32)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return d
}

// Modified from https://github.com/hverr/status-dashboard/blob/master/widgets/connections.go
func reverseIPBytes(bytes []byte) {
	j := len(bytes)
	m := j / 2
	j -= 1
	i := 0
	for i < m {
		b := bytes[i]
		bytes[i] = bytes[j]
		bytes[j] = b

		i += 1
		j -= 1
	}
}

func bytesToIP(bytes []byte) string {
	switch len(bytes) {
	case 4:
		reverseIPBytes(bytes[0:4])
	case 16:
		reverseIPBytes(bytes[0:4])
		reverseIPBytes(bytes[4:8])
		reverseIPBytes(bytes[8:12])
		reverseIPBytes(bytes[12:16])
	default:
		return ""
	}

	return strings.ToUpper(net.IP(bytes).String())
}

func convertIp(ip string) (string, error) {
	ipb, err := hex.DecodeString(ip)
	if err != nil {
		return "", err
	}
	return bytesToIP(ipb), nil
}

func findPid(inode string) string {
	// Loop through all fd dirs of process on /proc to compare the inode and
	// get the pid.

	pid := "-"

	d, err := filepath.Glob("/proc/[0-9]*/fd/[0-9]*")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	re := regexp.MustCompile(inode)
	for _, item := range d {
		path, _ := os.Readlink(item)
		out := re.FindString(path)
		if len(out) != 0 {
			pid = strings.Split(item, "/")[2]
		}
	}
	return pid
}

func getProcessExe(pid string) string {
	exe := fmt.Sprintf("/proc/%s/exe", pid)
	path, _ := os.Readlink(exe)
	return path
}

func getProcessName(exe string) string {
	n := strings.Split(exe, "/")
	name := n[len(n)-1]
	return strings.Title(name)
}

func getUser(uid int) string {
	u, _ := user.LookupUid(uid)
	return u.Name
}

func removeEmpty(array []string) []string {
	// remove empty data from line
	var new_array []string
	for _, i := range array {
		if i != "" {
			new_array = append(new_array, i)
		}
	}
	return new_array
}

func getInode2pid() map[string]string {
	// Loop through all fd dirs of process on /proc to compare the inode and
	// get the pid.

	inode2pid := make(map[string]string)
	pid := "-"

	d, err := filepath.Glob("/proc/[0-9]*/fd/[0-9]*")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var inode string
	for _, item := range d {
		path, _ := os.Readlink(item)
		if strings.Contains(path, "socket:[") {
			inode = path[8 : len(path)-1]
			pid = strings.Split(item, "/")[2]
			inode2pid[inode] = pid
		}
	}
	return inode2pid
}

func netstat(t string, lookupPids bool) ([]Process, error) {
	// Return a array of Process with Name, Ip, Port, State .. etc
	// Require Root acess to get information about some processes.

	var Processes []Process

	data, err := getData(t)
	if err != nil {
		return Processes, nil
	}
	var inode2pid map[string]string
	if lookupPids {
		inode2pid = getInode2pid()
	}

	for _, line := range data {

		// local ip and port
		line_array := removeEmpty(strings.Split(strings.TrimSpace(line), " "))
		ip_port := strings.Split(line_array[1], ":")
		ip, err := convertIp(ip_port[0])
		if err != nil {
			return Processes, err
		}
		port := hexToDec(ip_port[1])

		// foreign ip and port
		fip_port := strings.Split(line_array[2], ":")
		fip, err := convertIp(fip_port[0])
		if err != nil {
			return Processes, err
		}
		fport := hexToDec(fip_port[1])

		state := STATE[line_array[3]]
		uid, err := strconv.Atoi(line_array[7])
		if err != nil {
			return Processes, err
		}
		userName := getUser(uid)
		var pid, exe, name string
		if lookupPids {
			pid = inode2pid[line_array[9]]
			exe = getProcessExe(pid)
			name = getProcessName(exe)
		}

		p := Process{userName, name, pid, exe, state, ip, port, fip, fport}

		Processes = append(Processes, p)

	}

	return Processes, nil
}

func Tcp(lookupPids bool) ([]Process, error) {
	// Get a slice of Process type with TCP data
	return netstat("tcp", lookupPids)
}

func Udp(lookupPids bool) ([]Process, error) {
	// Get a slice of Process type with UDP data
	return netstat("udp", lookupPids)
}

func Tcp6(lookupPids bool) ([]Process, error) {
	// Get a slice of Process type with TCP6 data
	return netstat("tcp6", lookupPids)
}

func Udp6(lookupPids bool) ([]Process, error) {
	// Get a slice of Process type with UDP6 data
	return netstat("udp6", lookupPids)
}
