# TET Simulation

A package for simulating Targeted Energy Transfer (TET) systems. This package provides numerical simulation capabilities for analyzing the dynamics of coupled torsional oscillators with targeted energy transfer characteristics.

## Project Structure

```
.
├── bin/
├── src/           # Source code and Go module files
│   ├── main.go
│   ├── io_func.go
│   ├── plotting.go
│   ├── go.mod
│   └── go.sum
└── results/       # Simulation results and configuration
    ├── params.yaml        # Simulation parameters
    ├── run_sweep.sh      # Parameter sweep script
    ├── run_sweep.ps1     # Parameter sweep script for Windows
    ├── plot_sweep.m      # MATLAB script for analyzing sweep results
    └── plot_sweep.ipynb  # Jupyter notebook for analyzing sweep results
```

## Quick Start

### Linux and macOS

1. Download the appropriate executable for your system from the releases page and save to bin directory:
   - Linux: `./bin/tet-simulation-linux-amd64`
   - macOS (Intel): `./bin/tet-simulation-darwin-amd64`
   - macOS (Apple Silicon): `./bin/tet-simulation-darwin-arm64`

2. Make the executable file executable:
```bash
chmod +x ./bin/tet-simulation*
```

3. Run the simulation. For example:
```bash
./bin/tet-simulation-darwin-arm64 --params results/params.yaml -v
```


### Windows

1. Download the Windows executable from the releases page and save to
```
/bin/tet-simulation-windows-amd64.exe
```

2. Run the simulation using PowerShell. For example
```cmd
.\bin\tet-simulation-windows-amd64.exe --params results\params.yaml -v
```

Note: If you encounter any permission issues on Windows, you may need to unblock the executable file by right-clicking it, selecting Properties, and checking "Unblock" in the Security section.

## Configuration

System parameters are configured through a YAML file (`results/params.yaml`). The configuration includes:
- System masses and moments of inertia
- Spring constants
- Damping coefficients
- Simulation time and step size
- Initial conditions

Verbose flag `-v` can be set to enable more detailed output.
Try 
```cmd
tet-simulation-windows-amd64.exe --params params.yaml -v
```

## Parameter Sweep

The package includes a parameter sweep script (`results/run_sweep.sh`) that automatically runs simulations across a range of frequencies. This is useful for:
- Finding optimal energy transfer conditions
- Analyzing system behavior across different excitation frequencies
- Generating comprehensive datasets for analysis

To run the parameter sweep:
```bash
cd results
./run_sweep.sh
```

Or on Windows:
```cmd
cd results
powershell run_sweep.ps1 
```

The script will:
1. Iterate through predefined frequencies
2. Update the parameters file
3. Run simulations
4. Save results with frequency-specific filenames

## Output

The simulation generates:
- CSV files containing simulation results
- PNG plots visualizing the system dynamics
- Log files with simulation progress and diagnostics

## Features

- Numerical simulation of coupled torsional oscillators
- Support for both 4th order Runge-Kutta and Explicit Euler integration methods
- Configurable system parameters via YAML configuration
- Visualization capabilities for simulation results
- CSV export of simulation data
- Command-line interface for easy execution
- Parameter sweep functionality for frequency analysis

## Development

If you want to build from source:

### Prerequisites

- Go 1.x or later
- Required Go packages (managed via go.mod)

### Building from Source

1. Clone the repository:
```bash
git clone https://github.com/igrega348/tet-simulation.git
cd tet-simulation
```

2. Install dependencies:
```bash
cd src
go mod download
```

3. Build the executable:
```bash
./build.sh
```

4. Run the simulation:
```bash
./build/tet-simulation simulate --params results/params.yaml
```

## Implementation Details

The simulation implements:
- Coupled differential equations for torsional oscillators
- Matrix assembly for system dynamics
- Numerical integration methods (RK4 and Explicit Euler)
- Energy transfer analysis between oscillators

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Author

Ivan Grega
