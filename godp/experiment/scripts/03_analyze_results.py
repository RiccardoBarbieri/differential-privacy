#!/usr/bin/env python3
"""
Analizza i risultati dell'esperimento DP confrontando output GoDP con diversi epsilon.
"""

import pandas as pd
import numpy as np
import json
from pathlib import Path

SCRIPT_DIR = Path(__file__).parent
DATA_DIR = SCRIPT_DIR.parent / "data"
OUTPUT_DIR = SCRIPT_DIR.parent / "output"
VICTIM_FILE = DATA_DIR / "victim_info.json"
ATTACK_FILE = OUTPUT_DIR / "attack_results.json"

EPSILONS = [0.1, 0.5, 1.0, 2.0]


class NumpyEncoder(json.JSONEncoder):
    """Custom JSON encoder per tipi numpy."""
    def default(self, obj):
        if isinstance(obj, np.bool_):
            return bool(obj)
        if isinstance(obj, np.integer):
            return int(obj)
        if isinstance(obj, np.floating):
            return float(obj)
        if isinstance(obj, np.ndarray):
            return obj.tolist()
        return super().default(obj)


def load_json(path):
    with open(path) as f:
        return json.load(f)


def load_godp_sum(eps_dir: Path) -> dict:
    """Carica l'output SUM di GoDP per un dato epsilon."""
    sum_file = eps_dir / "output_SumSalaryByDepartment.csv"
    if not sum_file.exists():
        return {}
    df = pd.read_csv(sum_file, header=None, names=["Department", "Sum"])
    return dict(zip(df["Department"], df["Sum"]))


def generate_analysis_report(victim, attack, clear_sum, dp_results, dept, true_salary) -> str:
    """Genera il report di analisi completo."""
    lines = []

    lines.append("ANALISI RISULTATI: DE-ANONIMIZZAZIONE E PROTEZIONE CON GODP")
    lines.append("-" * 70)

    lines.append(f"\nVittima Target:")
    lines.append(f"  Reparto: {dept}")
    lines.append(f"  Età: {victim['age']}")
    lines.append(f"  Genere: {victim['gender']}")
    lines.append(f"  Stipendio reale: ${true_salary:,}")

    # Attacco senza DP
    lines.append(f"\n" + "-" * 70)
    lines.append("ATTACCO SENZA DP (Difference Attack)")
    lines.append("-" * 70)
    if attack.get("attack_successful"):
        lines.append(f"Stipendio dedotto: ${attack['inferred_salary']:,}")
        lines.append(f"Stipendio reale:   ${true_salary:,}")
        lines.append(f"Stato: ATTACCO RIUSCITO (precisione 100%)")
    else:
        lines.append("Risultato non disponibile")

    # Confronto con DP
    lines.append(f"\n" + "-" * 70)
    lines.append("PROTEZIONE CON DP (GoDP)")
    lines.append("-" * 70)
    lines.append(f"{'Epsilon':<10} {'Valore Reale':>15} {'Valore DP':>15} {'Rumore':>15} {'Rumore/Stipendio':>15} {'Protezione':<20}")
    lines.append("-" * 70)

    if clear_sum:
        for eps in EPSILONS:
            if dept in dp_results[eps]:
                dp_sum = dp_results[eps][dept]
                noise = dp_sum - clear_sum
                abs_noise = abs(noise)
                ratio = abs_noise / true_salary

                # Livello di protezione basato sul rapporto rumore/stipendio
                if ratio > 1.0:
                    protection = "FORTE"
                elif ratio > 0.5:
                    protection = "MODERATA"
                elif ratio > 0.3:
                    protection = "DEBOLE"
                else:
                    protection = "INSUFFICIENTE"

                lines.append(f"ε={eps:<7} ${clear_sum:>13,.0f} ${dp_sum:>13,.0f} {noise:>+14,.0f} {ratio:>14.1%} {protection:<20}")
            else:
                lines.append(f"ε={eps:<7} ${clear_sum:>13,.0f} {'[RIMOSSA]':>15} {'∞':>15} {'∞':>15} {'MASSIMA (rimossa)':<20}")

    return "\n".join(lines)


def main():
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)

    victim = load_json(VICTIM_FILE)
    attack = load_json(ATTACK_FILE) if ATTACK_FILE.exists() else {}

    dept = victim["department"]
    true_salary = victim["true_salary"]

    # valore reale dal dataset originale
    dataset_file = DATA_DIR / "salaries.csv"
    clear_sum = None
    if dataset_file.exists():
        df = pd.read_csv(dataset_file)
        dept_sum = df[df["Department"] == dept]["Salary"].sum()
        clear_sum = float(dept_sum)
        # print(f"DEBUG: clear_sum = {clear_sum}, dataset: {len(df)} rows")
    else:
        pass
        # print(f"DEBUG: dataset not found: {dataset_file}")

    # risultati DP per ogni epsilon
    dp_results = {}
    for eps in EPSILONS:
        eps_dir = DATA_DIR / "output" / f"eps_{eps}"
        dp_results[eps] = load_godp_sum(eps_dir)

    # Genera e salva report
    report = generate_analysis_report(victim, attack, clear_sum, dp_results, dept, true_salary)
    report_file = OUTPUT_DIR / "analysis_report.txt"
    with open(report_file, "w") as f:
        f.write(report)

    # Salva anche un JSON con i dati strutturati
    analysis_data = {
        "victim": victim,
        "attack_without_dp": {
            "successful": attack.get("attack_successful", False),
            "inferred_salary": attack.get("inferred_salary"),
            "true_salary": true_salary
        },
        "dp_protection": {}
    }

    if clear_sum:
        for eps in EPSILONS:
            if dept in dp_results[eps]:
                dp_sum = dp_results[eps][dept]
                noise = abs(dp_sum - clear_sum)

                noise_ratio = noise / true_salary

                analysis_data["dp_protection"][str(eps)] = {
                    "clear_sum": clear_sum,
                    "dp_sum": dp_sum,
                    "noise_on_query": noise,
                    "victim_salary": true_salary,
                    "noise_to_salary_ratio": noise_ratio,
                    # Protetto se il rumore è almeno il 30% dello stipendio
                    # (considerando che nel difference attack il rumore si compone)
                    "protected": noise_ratio > 0.3,
                    "protection_level": (
                        "forte" if noise_ratio > 1.0 else
                        "moderata" if noise_ratio > 0.5 else
                        "debole" if noise_ratio > 0.3 else
                        "insufficiente"
                    )
                }
            else:
                analysis_data["dp_protection"][str(eps)] = {
                    "partition_removed": True,
                    "protected": True,
                    "protection_level": "massima (partizione rimossa)"
                }

    with open(OUTPUT_DIR / "analysis_data.json", "w") as f:
        json.dump(analysis_data, f, indent=2, cls=NumpyEncoder)

    protected_count = sum(1 for eps in EPSILONS if dept not in dp_results[eps] or
                          (clear_sum and abs(dp_results[eps].get(dept, 0) - clear_sum) / true_salary > 0.3))
    print(f"Analisi: {protected_count}/{len(EPSILONS)} configurazioni proteggono la vittima (soglia: rumore > 30% stipendio)")


if __name__ == "__main__":
    main()



