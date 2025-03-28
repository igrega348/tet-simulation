$EXE_PATH = "..\bin\tet-simulation-windows-amd64.exe"
$omegaValues = @(0.015, 0.016, 0.017, 0.019, 0.020, 0.022, 0.023, 0.025, 0.027, 0.029, 0.031, 0.033, 0.036, 0.039, 0.041, 0.045, 0.048, 0.052, 0.055, 0.060, 0.064, 0.069, 0.074, 0.080, 0.086, 0.092, 0.099, 0.107, 0.115, 0.123, 0.133, 0.143, 0.154, 0.165, 0.178, 0.191, 0.205, 0.221, 0.237, 0.255, 0.275, 0.295, 0.318, 0.341, 0.367, 0.395, 0.425, 0.457, 0.491, 0.528, 0.568, 0.611, 0.657, 0.706, 0.760, 0.817, 0.878, 0.945, 1.016, 1.093, 1.175, 1.263, 1.359, 1.461, 1.571, 1.690, 1.817, 1.954, 2.101, 2.260, 2.430, 2.613, 2.810, 3.022, 3.250, 3.495, 3.759, 4.042, 4.347, 4.674, 5.027, 5.406, 5.813, 6.252, 6.723, 7.230, 7.775, 8.361, 8.991, 9.669, 10.398, 11.182, 12.025, 12.931, 13.906, 14.954, 16.082, 17.294, 18.598, 20.000)

# Create sweep directory if it doesn't exist
New-Item -ItemType Directory -Force -Path "sweep"

foreach ($omega in $omegaValues) {
    # Calculate tmax (10x the period of driving frequency but no less than 500)
    $tmax = [math]::Round(10 * 6.28318530718 / $omega)
    $tmax = [math]::Max($tmax, 500)
    
    Write-Host "omega: $omega, tmax: $tmax"
    
    # Read the params file
    $paramsContent = Get-Content "params.yaml"
    
    # Replace omega and tmax values
    $paramsContent = $paramsContent -replace "omega:.*", "omega: $omega"
    $paramsContent = $paramsContent -replace "tmax:.*", "tmax: $tmax.0"
    
    # Write back to params file
    $paramsContent | Set-Content "params.yaml"
    
    # Run the simulation
    & $EXE_PATH --params ".\params.yaml"
    
    # Move results to sweep directory
    Move-Item "simulation_results.csv" "sweep\simulation_results_$omega.csv" -Force
    Move-Item "simulation_plot_1.png" "sweep\simulation_plot_1_$omega.png" -Force
    Move-Item "simulation_plot_2.png" "sweep\simulation_plot_2_$omega.png" -Force
} 