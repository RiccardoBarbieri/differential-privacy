#!/bin/bash
# Esecuzione delle 3 query del difference attack con diversi valori di epsilon.
# Per ogni epsilon esegue Q1, Q2, Q3 attraverso GoDP con rumore DP.
#
# Query:
#   Q1: SUM(Salary) WHERE Department = 'Engineering'
#   Q2: SUM(Salary) WHERE Department = 'Engineering' AND Gender != 'F'
#   Q3: SUM(Salary) WHERE Department = 'Engineering' AND Gender = 'F' AND Age != 35

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
EXPERIMENT_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_DIR="$(dirname "$EXPERIMENT_DIR")"
GODP="$PROJECT_DIR/main/build/godp"

EPSILONS=("0.5" "5" "10" "15")
QUERIES=("q1" "q2" "q3")

cd "$PROJECT_DIR"

for eps in "${EPSILONS[@]}"; do
    for q in "${QUERIES[@]}"; do
        echo "  Running eps=$eps $q..."
        $GODP fromfile --file "experiment/specs/spec_eps_${eps}_${q}.yaml" 2>/dev/null
    done
done