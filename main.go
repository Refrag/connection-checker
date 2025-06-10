package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ncruces/zenity"
)

// Predefined list of hostnames to traceroute - DatHost server locations
// Source: https://dathost.net/reference/server-locations-mapping
var hostnames = []string{
	"beauharnois.dathost.net",   // Canada - Toronto
	"new-york-city.dathost.net", // USA - New York
	"los-angeles.dathost.net",   // USA CA - Los Angeles
	"miami.dathost.net",         // USA FL - Miami
	"chicago.dathost.net",       // USA IL - Chicago
	"portland.dathost.net",      // USA WA - Seattle
	"dallas.dathost.net",        // USA TX - Dallas
	"atlanta.dathost.net",       // USA GA - Atlanta
	"denver.dathost.net",        // USA CO - Denver
	"copenhagen.dathost.net",    // Denmark - Copenhagen
	"helsinki.dathost.net",      // Finland - Helsinki
	"strasbourg.dathost.net",    // France - Paris
	"dusseldorf.dathost.net",    // Germany - Frankfurt
	"amsterdam.dathost.net",     // Netherlands - Amsterdam
	"warsaw.dathost.net",        // Poland - Warsaw
	"barcelona.dathost.net",     // Spain - Madrid
	"stockholm.dathost.net",     // Sweden - Stockholm
	"istanbul.dathost.net",      // Turkey - Istanbul
	"bristol.dathost.net",       // United Kingdom - London
	"sydney.dathost.net",        // Australia - Sydney
	"sao-paulo.dathost.net",     // Brazil - São Paulo
	"santiago.dathost.net",      // Chile - Santiago
	"hong-kong.dathost.net",     // Hong Kong - Hong Kong
	"mumbai.dathost.net",        // India - Mumbai
	"tokyo.dathost.net",         // Japan - Tokyo
	"singapore.dathost.net",     // Singapore - Singapore
	"johannesburg.dathost.net",  // South Africa - Johannesburg
	"seoul.dathost.net",         // South Korea - Seoul
	"oslo.dathost.net",          // Norway - Oslo
	"prague.dathost.net",        // Czechia - Prague
	"milan.dathost.net",         // Italy - Milan
	"bucharest.dathost.net",     // Romania - Bucharest
	"dublin.dathost.net",        // Ireland - Dublin
	"auckland.dathost.net",      // New Zealand - Auckland
}

// TracerouteResult holds the result of a traceroute operation
type TracerouteResult struct {
	Hostname  string
	Index     int
	Result    string
	Error     error
	StartTime time.Time
	EndTime   time.Time
	TimedOut  bool
}

func main() {
	fmt.Println("=== Refrag ConnectionChecker ===")
	fmt.Println("This tool will run traceroutes to DatHost gaming servers worldwide.")
	fmt.Printf("Total locations to test: %d\n\n", len(hostnames))

	// Get output file path from user via GUI dialog
	outputFile, err := getOutputFilePathGUI()
	if err != nil {
		log.Fatalf("File selection cancelled or failed: %v", err)
	}

	// Create or open the output file
	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer file.Close()

	fmt.Printf("Starting traceroute tests, saving results to: %s\n", outputFile)
	fmt.Println()

	// Write header with timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05 MST")
	header := fmt.Sprintf("Refrag ConnectionChecker Results - %s\n", timestamp)
	file.WriteString(header)
	file.WriteString(strings.Repeat("=", len(header)-1) + "\n\n")

	// Get and write local IP address
	localIP, err := getLocalIP()
	if err != nil {
		log.Printf("Warning: Could not determine public IP: %v", err)
		localIP = "Unknown"
	}

	file.WriteString(fmt.Sprintf("Public IP Address: %s\n\n", localIP))
	file.WriteString("NOTE: This report contains network diagnostic information for Refrag support.\n")
	file.WriteString("Please send this file to Refrag support as requested.\n\n")
	file.WriteString("TIMEOUT SETTING: Each traceroute has a 2-minute timeout limit.\n\n")
	file.WriteString(strings.Repeat("-", 80) + "\n\n")

	// Run all traceroutes concurrently
	results := runConcurrentTraceroutes()

	// Write results to file in original order
	for _, result := range results {
		// Write hostname header
		hostnameHeader := fmt.Sprintf("TRACEROUTE TO: %s", result.Hostname)
		file.WriteString(hostnameHeader + "\n")
		file.WriteString(strings.Repeat("~", len(hostnameHeader)) + "\n")
		file.WriteString(fmt.Sprintf("Started at: %s\n", result.StartTime.Format("15:04:05")))
		file.WriteString(fmt.Sprintf("Completed at: %s\n", result.EndTime.Format("15:04:05")))
		file.WriteString(fmt.Sprintf("Duration: %v\n", result.EndTime.Sub(result.StartTime).Round(time.Second)))

		if result.TimedOut {
			file.WriteString("Status: TIMED OUT (exceeded 2-minute limit)\n\n")
		} else {
			file.WriteString("Status: COMPLETED\n\n")
		}

		if result.Error != nil {
			errorMsg := fmt.Sprintf("Error running traceroute to %s: %v\n", result.Hostname, result.Error)
			file.WriteString(errorMsg)
		} else {
			file.WriteString(result.Result)
		}

		// Add separation between results
		file.WriteString("\n" + strings.Repeat("-", 80) + "\n\n")
	}

	// Write footer
	file.WriteString(fmt.Sprintf("All traceroutes completed at: %s\n", time.Now().Format("2006-01-02 15:04:05 MST")))
	file.WriteString("\n" + strings.Repeat("=", 80) + "\n")
	file.WriteString("END OF REPORT - Please send this file to Refrag support\n")

	fmt.Printf("\n✅ Traceroute results saved to: %s\n", outputFile)
	fmt.Println("\n🔔 IMPORTANT: Please send the results file to Refrag support as requested.")
	fmt.Println("   The file contains network diagnostic information needed for troubleshooting.")

	// Keep window open for user to read results
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📋 RESULTS SUMMARY:")

	successful := 0
	timedOut := 0
	failed := 0

	for _, result := range results {
		if result.TimedOut {
			timedOut++
		} else if result.Error != nil {
			failed++
		} else {
			successful++
		}
	}

	fmt.Printf("✅ Successful: %d/%d\n", successful, len(hostnames))
	fmt.Printf("⏰ Timed out: %d/%d\n", timedOut, len(hostnames))
	fmt.Printf("❌ Failed: %d/%d\n", failed, len(hostnames))
	fmt.Println(strings.Repeat("=", 60))

	fmt.Print("\n🖱️  Press Enter to close this window...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// runConcurrentTraceroutes runs all traceroutes concurrently and returns results in original order
func runConcurrentTraceroutes() []TracerouteResult {
	var wg sync.WaitGroup
	results := make([]TracerouteResult, len(hostnames))
	completed := 0
	var completedMutex sync.Mutex

	fmt.Printf("🚀 Starting %d traceroutes...\n", len(hostnames))
	fmt.Printf("This may take several minutes depending on your network connection...\n")
	fmt.Printf("Please do not close this window until the traceroutes are complete.\n\n")
	startTime := time.Now()

	// Start all traceroutes concurrently
	for i, hostname := range hostnames {
		wg.Add(1)
		go func(index int, host string) {
			defer wg.Done()

			traceStartTime := time.Now()
			result, err, timedOut := runTracerouteWithTimeout(host, 2*time.Minute)
			traceEndTime := time.Now()

			results[index] = TracerouteResult{
				Hostname:  host,
				Index:     index,
				Result:    result,
				Error:     err,
				StartTime: traceStartTime,
				EndTime:   traceEndTime,
				TimedOut:  timedOut,
			}

			// Update progress counter thread-safely
			completedMutex.Lock()
			completed++
			currentCompleted := completed
			completedMutex.Unlock()

			if timedOut {
				fmt.Printf("⏰ [%d/%d] %s - Timed out after 2 minutes\n", currentCompleted, len(hostnames), host)
			} else if err != nil {
				fmt.Printf("❌ [%d/%d] %s - Error: %v\n", currentCompleted, len(hostnames), host, err)
			} else {
				fmt.Printf("✅ [%d/%d] %s - Completed in %v\n",
					currentCompleted, len(hostnames), host,
					traceEndTime.Sub(traceStartTime).Round(time.Second))
			}
		}(i, hostname)
	}

	// Wait for all traceroutes to complete
	wg.Wait()

	totalDuration := time.Since(startTime)
	fmt.Printf("\n🎉 All traceroutes completed in %v\n", totalDuration.Round(time.Second))

	return results
}

// getOutputFilePathGUI shows a native file save dialog
func getOutputFilePathGUI() (string, error) {
	fmt.Println("📁 Opening file save dialog...")
	fmt.Println("   Please choose where to save the traceroute results file.")

	// Get user's home directory for default location
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	// Default filename
	defaultFilename := fmt.Sprintf("refrag_traceroute_results_%s.txt",
		time.Now().Format("2006-01-02_15-04-05"))

	// Show native file save dialog using zenity
	filename, err := zenity.SelectFileSave(
		zenity.Title("Save Refrag Traceroute Results"),
		zenity.FileFilter{
			Name:     "Text files",
			Patterns: []string{"*.txt"},
		},
		zenity.Filename(filepath.Join(homeDir, defaultFilename)),
	)

	if err != nil {
		return "", err
	}

	// Ensure .txt extension if not provided
	if filepath.Ext(filename) == "" {
		filename += ".txt"
	}

	fmt.Printf("📄 Selected file: %s\n\n", filename)
	return filename, nil
}

// runTracerouteWithTimeout executes the traceroute command with a timeout
func runTracerouteWithTimeout(hostname string, timeout time.Duration) (string, error, bool) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var cmd *exec.Cmd

	// Check if we're on macOS/Linux or Windows and use appropriate command
	if _, err := exec.LookPath("traceroute"); err == nil {
		// Unix-like systems (macOS, Linux)
		cmd = exec.CommandContext(ctx, "traceroute", hostname)
	} else if _, err := exec.LookPath("tracert"); err == nil {
		// Windows systems
		cmd = exec.CommandContext(ctx, "tracert", hostname)
	} else {
		return "", fmt.Errorf("neither traceroute nor tracert command found"), false
	}

	output, err := cmd.Output()
	if err != nil {
		// Check if it was a timeout
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("traceroute timed out after %v", timeout), true
		}
		return "", fmt.Errorf("command execution failed: %v", err), false
	}

	return string(output), nil, false
}

// getLocalIP returns the public IP address of the machine using ipify.org API
func getLocalIP() (string, error) {
	// Try to get public IP from ipify.org API first
	publicIP, err := getPublicIP()
	if err == nil && publicIP != "" {
		return publicIP, nil
	}

	// Fallback to local network interface detection if API fails
	log.Printf("Warning: Could not get public IP from ipify.org (%v), falling back to local interface detection", err)
	return getLocalNetworkIP()
}

// getPublicIP gets the public IP address using ipify.org API
func getPublicIP() (string, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make request to ipify.org API
	resp, err := client.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", fmt.Errorf("failed to contact ipify.org: %v", err)
	}
	defer resp.Body.Close()

	// Check if request was successful
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ipify.org returned status %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	// Clean up the IP address (remove any whitespace)
	ip := strings.TrimSpace(string(body))

	// Basic validation - check if it looks like an IP address
	if net.ParseIP(ip) == nil {
		return "", fmt.Errorf("invalid IP address received: %s", ip)
	}

	return ip, nil
}

// getLocalNetworkIP returns the local network interface IP address (fallback method)
func getLocalNetworkIP() (string, error) {
	// Get all network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		// Skip loopback and down interfaces
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// Get addresses for this interface
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				// Return the first non-loopback IPv4 address
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String(), nil
				}
			}
		}
	}

	return "", fmt.Errorf("no active network interface found")
}
