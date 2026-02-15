#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
EXPERIMENT_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_DIR="$(dirname "$EXPERIMENT_DIR")"

echo "Generazione dataset..."
python3 "$SCRIPT_DIR/01_generate_dataset.py"

echo "Esecuzione attacco (senza DP)..."
python3 "$SCRIPT_DIR/02_attack_non_dp.py"

echo "Esecuzione GoDP (epsilon: 0.1, 0.5, 1.0, 2.0)..."
cd "$PROJECT_DIR"
make build -s 2>/dev/null || make build
bash "$SCRIPT_DIR/run_all_epsilons.sh" 2>/dev/null

echo "Analisi risultati..."
python3 "$SCRIPT_DIR/03_analyze_results.py"

echo "Output in experiment/output/:"
echo "  - attack_report.txt         Report attacco"
echo "  - analysis_report.txt       Analisi confronto DP"
echo "  - attack_results.json       Risultati JSON"
echo "  - analysis_data.json        Dati analisi JSON"
echo ""
