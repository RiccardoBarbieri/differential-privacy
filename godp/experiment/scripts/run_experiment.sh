#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
EXPERIMENT_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_DIR="$(dirname "$EXPERIMENT_DIR")"

echo "=== Esperimento: Difference Attack con e senza DP ==="
echo ""

echo "[1/5] Generazione dataset..."
python3 "$SCRIPT_DIR/01_generate_dataset.py"

echo "[2/5] Esecuzione difference attack (senza DP)..."
python3 "$SCRIPT_DIR/02_attack_non_dp.py"

echo "[3/5] Esecuzione GoDP (3 query x 4 epsilon = 12 esecuzioni)..."
cd "$PROJECT_DIR"
make build -s 2>/dev/null || make build
bash "$SCRIPT_DIR/run_all_epsilons.sh"

echo "[4/5] Analisi risultati..."
python3 "$SCRIPT_DIR/03_analyze_results.py"

echo ""
echo "=== Esperimento completato ==="
echo ""
echo "Report in experiment/output/:"
echo "  - attack_report.txt         Report attacco senza DP"
echo "  - analysis_report.txt       Analisi confronto con/senza DP"
echo "  - attack_results.json       Risultati JSON attacco"
echo "  - analysis_data.json        Dati analisi JSON"
echo ""


