# Refrag ConnectionChecker

A Go application that performs traceroutes to DatHost gaming servers worldwide and saves the results for Refrag support analysis.

## Features

- **Concurrent Execution**: Runs traceroutes to all 34 DatHost gaming server locations simultaneously for maximum speed
- **Smart Timeouts**: 2-minute timeout per traceroute to prevent hanging on unresponsive routes
- **Public IP Detection**: Uses ipify.org API to get your public IP address for accurate network diagnostics
- **GUI file dialog** - Native file save dialog for easy file selection
- **Results Summary**: Clear overview of successful, timed-out, and failed traceroutes
- **User-Friendly Interface**: Window stays open for review, closes only when user is ready
- Saves results to a well-formatted text file with clear separation
- Cross-platform compatible (works on macOS, Linux, and Windows)
- Real-time progress tracking with completion times
- Branded for Refrag support with clear instructions

## Performance

🚀 **Significantly Faster**: All traceroutes run concurrently instead of one-by-one, reducing total execution time from potentially 30+ minutes to just a few minutes!

- **Concurrent (new)**: All 34 traceroutes simultaneously = 2-5 minutes total
- **Timeout Protection**: No single traceroute can hang for more than 2 minutes

## Build System

The project includes a comprehensive Makefile for cross-platform builds supporting 9 different OS/architecture combinations.

### Quick Start

```bash
# Build for your current platform
make local

# Build for all supported platforms
make build-all

# Create release packages
make release

# Show all available targets
make help
```

### Supported Platforms

- **Linux**: amd64, arm64, 386
- **macOS**: amd64 (Intel), arm64 (Apple Silicon)
- **Windows**: amd64, 386
- **FreeBSD**: amd64
- **OpenBSD**: amd64

### Build Targets

| Target | Description |
|--------|-------------|
| `make help` | Show all available targets and supported platforms |
| `make local` | Build for current platform only (fastest) |
| `make build-all` | Build for all supported platforms |
| `make release` | Create packaged releases for distribution |
| `make clean` | Remove all build artifacts |
| `make deps` | Download and verify Go dependencies |
| `make test` | Run tests |
| `make fmt` | Format Go code |
| `make vet` | Run go vet |

### Platform-Specific Builds

```bash
make darwin-amd64    # macOS Intel
make darwin-arm64    # macOS Apple Silicon
make linux-amd64     # Linux 64-bit
make windows-amd64   # Windows 64-bit
```

### Development Workflow

```bash
# Quick development cycle
make dev             # Build and run immediately

# Install locally
make install         # Install to /usr/local/bin (Unix)

# Clean up
make clean           # Remove build artifacts
```

## Default Hostnames

The application traces routes to all DatHost server locations worldwide (34 locations total), including:
- beauharnois.dathost.net (Canada - Toronto)
- new-york-city.dathost.net (USA - New York)
- los-angeles.dathost.net (USA CA - Los Angeles)
- copenhagen.dathost.net (Denmark - Copenhagen)
- tokyo.dathost.net (Japan - Tokyo)
- sydney.dathost.net (Australia - Sydney)
- And 28 more locations across the globe

Complete list sourced from: https://dathost.net/reference/server-locations-mapping

## Usage

### Method 1: Using Make (Recommended)

1. **Build the application:**
   ```bash
   make local
   ```

2. **Run the application:**
   ```bash
   ./build/connectionchecker
   ```

### Method 2: Direct Go Build

1. **Build the application:**
   ```bash
   go build -o connectionchecker
   ```

2. **Run the application:**
   ```bash
   ./connectionchecker
   ```

### Using the Application

3. **Use the GUI file dialog:**
   - A native file save dialog will automatically open
   - Choose where to save your results file
   - The dialog defaults to your home directory with a timestamped filename
   - Click "Save" to proceed

4. **Watch the concurrent execution:**
   - All 34 traceroutes start simultaneously
   - Real-time progress updates show completion status
   - Individual timing information for each traceroute
   - Timeout protection (2-minute limit per traceroute)
   - Total execution time displayed

5. **Review the results:**
   - Results summary shows successful, timed-out, and failed traceroutes
   - Window stays open for you to review the information
   - Press Enter when ready to close the application

6. **Send results to Refrag support:**
   - The application will clearly indicate when it's finished
   - Send the generated file to Refrag support as requested
   - The file contains network diagnostic information needed for troubleshooting

## GUI Features

- **Native File Dialog**: Uses your operating system's native file save dialog
- **Smart Defaults**: Automatically suggests a timestamped filename in your home directory
- **File Extension Handling**: Automatically adds `.txt` extension if not provided
- **Filter Support**: Shows only text files in the dialog for easy selection
- **Cross-Platform**: Works on macOS, Linux, and Windows with native look and feel

## Concurrent Execution Features

- **Parallel Processing**: All traceroutes run simultaneously using goroutines
- **Thread-Safe Progress**: Real-time progress counter with mutex protection
- **Timeout Protection**: 2-minute maximum per traceroute prevents hanging
- **Ordered Output**: Results are written to file in the original hostname order
- **Individual Timing**: Each traceroute shows start time, end time, and duration
- **Status Tracking**: Clear indication of completed, timed-out, and failed traceroutes
- **Error Handling**: Failed traceroutes don't block others from completing
- **Performance Metrics**: Total execution time and efficiency information

## Requirements

- Go 1.21 or higher (for building)
- `traceroute` command available on Unix-like systems (macOS, Linux)
- `tracert` command available on Windows systems
- GUI environment (for the file save dialog)
- `make` (for using the build system)

## Output Format

The application generates a text file with:
- Refrag ConnectionChecker branded header with timestamp
- Local machine IP address
- Note about 2-minute timeout setting
- Note about sending to Refrag support
- Traceroute results for each hostname with:
  - Start time, completion time, and duration
  - Status (COMPLETED or TIMED OUT)
  - Full traceroute output
  - Clear separation between results
- Overall completion timestamp and final instructions

## Dependencies

The application uses the following Go packages:
- `github.com/ncruces/zenity` - For cross-platform native file dialogs
- `context` - For timeout management
- `sync` - For concurrent execution coordination
- Standard library packages for networking and system operations

## Release Management

The Makefile includes comprehensive release management:

```bash
# Create release packages for all platforms
make release
```

This creates:
- **Linux/macOS/FreeBSD/OpenBSD**: `.tar.gz` archives
- **Windows**: `.zip` archives
- Each package includes the binary and README

Output structure:
```
build/
├── connectionchecker-dev-darwin-amd64.tar.gz
├── connectionchecker-dev-darwin-arm64.tar.gz
├── connectionchecker-dev-linux-amd64.tar.gz
├── connectionchecker-dev-windows-amd64.zip
└── ...
```

## Example Workflow

```
=== Refrag ConnectionChecker ===
This tool will run traceroutes to DatHost gaming servers worldwide.
Total locations to test: 34

📁 Opening file save dialog...
   Please choose where to save the traceroute results file.
📄 Selected file: /Users/username/Desktop/refrag_traceroute_results_2024-01-15_14-30-25.txt

Starting traceroute tests, saving results to: /Users/username/Desktop/refrag_traceroute_results_2024-01-15_14-30-25.txt

🚀 Starting 34 traceroutes...
This may take several minutes depending on your network connection...
Please do not close this window until the traceroutes are complete.

✅ [3/34] copenhagen.dathost.net - Completed in 2s
✅ [7/34] stockholm.dathost.net - Completed in 3s
⏰ [12/34] some-host.dathost.net - Timed out after 2 minutes
✅ [15/34] tokyo.dathost.net - Completed in 4s
...

🎉 All traceroutes completed in 3m 15s
💡 Running concurrently saved significant time compared to sequential execution!

✅ Traceroute results saved to: /Users/username/Desktop/refrag_traceroute_results_2024-01-15_14-30-25.txt

🔔 IMPORTANT: Please send the results file to Refrag support as requested.
   The file contains network diagnostic information needed for troubleshooting.

============================================================
📋 RESULTS SUMMARY:
✅ Successful: 31/34
⏰ Timed out: 2/34
❌ Failed: 1/34
============================================================

🖱️  Press Enter to close this window...
```

## Important Notes

- **Always send the results file to Refrag support** as instructed by the application
- The application requires appropriate permissions to run traceroute commands
- On some systems, you may need to run with elevated privileges (sudo)
- **Much faster execution**: Concurrent operation reduces total time significantly
- **Timeout protection**: Each traceroute is limited to 2 minutes maximum
- A GUI environment is required for the file save dialog
- Progress updates appear in real-time as traceroutes complete
- **Window stays open**: You can review results before closing

## Customization

To modify the list of hostnames, edit the `hostnames` variable in `main.go`. The current list includes all DatHost gaming server locations for comprehensive network analysis.

To change the timeout limit, modify the `2*time.Minute` value in the `runTracerouteWithTimeout` function calls.
