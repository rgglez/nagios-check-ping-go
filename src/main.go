/*
    check-ping-go
    Copyright (C) 2025 Rodolfo González González.

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.

*/

package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	probing "github.com/prometheus-community/pro-bing"
	"github.com/spf13/pflag"
	"github.com/xorpaul/go-nagios"
)

//-----------------------------------------------------------------------------

var (
	host     string
	warn     string
	crit     string
	count    int
	warnTime float64
	warnLoss float64
	critTime float64
	critLoss float64
)

//-----------------------------------------------------------------------------

func init() {
	pflag.StringVarP(&host, "host", "H", "", "Hostname or IP address to ping")
	pflag.StringVarP(&warn, "warn", "w", "150,5%", "Warning threshold in the form 'time,packets%'")
	pflag.StringVarP(&crit, "crit", "c", "200,10%", "Critical threshold in the form 'time,packets%'")
	pflag.IntVarP(&count, "count", "n", 5, "Number of packets to be send")
	pflag.Parse()

	// Parse warning and critical thresholds
	var err error
	warnTime, warnLoss, err = parseThreshold(warn)
	if err != nil {
		fmt.Printf("Error parsing warning threshold: %v\n", err)
		os.Exit(1)
	}
	critTime, critLoss, err = parseThreshold(crit)
	if err != nil {
		fmt.Printf("Error parsing critical threshold: %v\n", err)
		os.Exit(1)
	}
}

//-----------------------------------------------------------------------------

func parseThreshold(threshold string) (time float64, loss float64, err error) {
	parts := strings.Split(threshold, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid threshold format")
	}
	time = parseFloat(parts[0])
	loss = parseFloat(strings.TrimSuffix(parts[1], "%"))
	return time, loss, nil
}

//-----------------------------------------------------------------------------

func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

//-----------------------------------------------------------------------------

func main() {
	if host == "" {
		fmt.Println("Host argument is required")
		os.Exit(1)
	}

	// Create a new pinger
	pinger, err := probing.NewPinger(host)
	if err != nil {
		fmt.Printf("Error creating pinger: %v\n", err)
		os.Exit(1)
	}
	pinger.SetPrivileged(false) // Run as non-root user
	pinger.Size = 56            // Default ICMP payload size
	pinger.Count = count        // Number of pings to send
	pinger.Timeout = time.Second * 5

	// Run the ping
	err = pinger.Run()
	if err != nil {
		fmt.Printf("Error pinging host: %v\n", err)
		os.Exit(1)
	}

	stats := pinger.Statistics()

	// Check thresholds
	var exitCode int
	var statusText string

	if stats.PacketLoss >= critLoss || stats.AvgRtt >= time.Duration(critTime)*time.Millisecond {
		exitCode = 2 // Critical
		statusText = fmt.Sprintf("Packet loss: %.2f%%, RTT: %v", stats.PacketLoss, stats.AvgRtt)
	} else if stats.PacketLoss >= warnLoss || stats.AvgRtt >= time.Duration(warnTime)*time.Millisecond {
		exitCode = 1 // Warning
		statusText = fmt.Sprintf("Packet loss: %.2f%%, RTT: %v", stats.PacketLoss, stats.AvgRtt)
	} else {
		exitCode = 0 // OK
		statusText = fmt.Sprintf("Packet loss: %.2f%%, RTT: %v", stats.PacketLoss, stats.AvgRtt)
	}

	// Prepare performance data
	perfdata := fmt.Sprintf("rtt=%.2fms;%.2f;%.2f loss=%.2f%%;%.2f;%.2f",
		stats.AvgRtt.Seconds()*1000, warnTime, critTime, stats.PacketLoss, warnLoss, critLoss)

	// Return Nagios result
	nr := nagios.NagiosResult{
		ExitCode: exitCode,
		Text:     statusText,
		Perfdata: perfdata,
	}
	nagios.NagiosExit(nr)
}
