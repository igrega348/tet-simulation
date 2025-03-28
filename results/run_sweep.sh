#!/bin/zsh
zmodload zsh/mathfunc
# Detect platform and architecture
if [[ "$OSTYPE" == "darwin"* ]]; then
    PLATFORM="darwin"
elif [[ "$OSTYPE" == "linux"* ]]; then
    PLATFORM="linux" 
else
    echo "Unsupported platform: $OSTYPE"
    exit 1
fi

ARCH=$(uname -m)
if [[ "$ARCH" == "x86_64" ]]; then
    ARCH="amd64"
elif [[ "$ARCH" == "arm64" ]]; then
    ARCH="arm64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

EXE_PATH="../bin/tet-simulation-${PLATFORM}-${ARCH}"
mkdir sweep

for omega in 0.015 0.016 0.017 0.019 0.020 0.022 0.023 0.025 0.027 0.029 0.031 0.033 0.036 0.039 0.041 0.045 0.048 0.052 0.055 0.060 0.064 0.069 0.074 0.080 0.086 0.092 0.099 0.107 0.115 0.123 0.133 0.143 0.154 0.165 0.178 0.191 0.205 0.221 0.237 0.255 0.275 0.295 0.318 0.341 0.367 0.395 0.425 0.457 0.491 0.528 0.568 0.611 0.657 0.706 0.760 0.817 0.878 0.945 1.016 1.093 1.175 1.263 1.359 1.461 1.571 1.690 1.817 1.954 2.101 2.260 2.430 2.613 2.810 3.022 3.250 3.495 3.759 4.042 4.347 4.674 5.027 5.406 5.813 6.252 6.723 7.230 7.775 8.361 8.991 9.669 10.398 11.182 12.025 12.931 13.906 14.954 16.082 17.294 18.598 20.000; do
    # Set tmax to 10x the period of driving frequency but no less than 500
    tmax=$((10 * 6.28318530718 / $omega))
    tmax=$(( int(rint($tmax)) ))
    tmax=$(( $tmax > 500 ? $tmax : 500 ))
    echo "omega: $omega, tmax: $tmax"
    sed -i '' "s/omega:.*/omega: $omega/" params.yaml
    sed -i '' "s/tmax:.*/tmax: $tmax.0/" params.yaml
    $EXE_PATH --params ./params.yaml
    mv simulation_results.csv sweep/simulation_results_$omega.csv
    mv simulation_plot_1.png sweep/simulation_plot_1_$omega.png
    mv simulation_plot_2.png sweep/simulation_plot_2_$omega.png
done