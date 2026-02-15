#!/bin/bash
# esecuzione aggregazione SumSalaryBuDepartment con diversi valori di epsilon

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
EXPERIMENT_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_DIR="$(dirname "$EXPERIMENT_DIR")"
GODP="$PROJECT_DIR/main/build/godp"

EPSILONS=("0.1" "0.5" "1.0" "2.0")

cd "$PROJECT_DIR"

for eps in "${EPSILONS[@]}"; do
    $GODP fromfile --file "experiment/specs/spec_eps_${eps}.yaml" 2>/dev/null
done