#!/bin/bash
# esecuzione aggregazione SumSalaryBuDepartment con diversi valori di epsilon

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
EXPERIMENT_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_DIR="$(dirname "$EXPERIMENT_DIR")"
GODP="$PROJECT_DIR/main/build/godp"

EPSILONS=("0.5" "5" "10" "15")

cd "$PROJECT_DIR"

for eps in "${EPSILONS[@]}"; do
    $GODP fromfile --file "experiment/specs/spec_eps_${eps}.yaml" 2>/dev/null
done