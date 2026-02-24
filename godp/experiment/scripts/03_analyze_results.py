#!/usr/bin/env python3
"""
Analizza i risultati dell'esperimento DP: confronta il difference attack
eseguito senza DP (query esatte) con lo stesso attacco eseguito attraverso
GoDP (query con rumore Laplace) a diversi livelli di epsilon.

Per ogni epsilon, GoDP esegue le 3 query del difference attack:
  Q1: SUM(Salary) WHERE Department = 'Engineering'
  Q2: SUM(Salary) WHERE Department = 'Engineering' AND Gender != 'F'
  Q3: SUM(Salary) WHERE Department = 'Engineering' AND Gender = 'F' AND Age != 35

L'attaccante calcola: stipendio_dedotto = (Q1 - Q2) - Q3
Con DP, il rumore si compone attraverso le sottrazioni.
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

EPSILONS = [0.5, 5, 10, 15]
QUERY_NAMES = ["Q1", "Q2", "Q3"]

ERROR_RATIO_LOW = 0.1
ERROR_RATIO_MODERATE = 0.3
ERROR_RATIO_HIGH = 1.0


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


def load_dp_query(eps_dir: Path, query_name: str, department: str):
    """
    Carica il risultato di una singola query DP da GoDP.

    Returns:
        Il valore della somma per il reparto specificato, o None se la
        partizione e' stata rimossa dal meccanismo di partition selection.
    """
    query_file = eps_dir / f"output_{query_name}.csv"
    if not query_file.exists():
        return None
    df = pd.read_csv(query_file, header=None, names=["Department", "Sum"])
    row = df[df["Department"] == department]
    if row.empty:
        return None  # Partizione rimossa dalla partition selection
    return float(row["Sum"].iloc[0])


def compute_dp_attack(q1, q2, q3):
    """
    Esegue il difference attack usando i risultati DP.

    Returns:
        Dizionario con i risultati intermedi e finali dell'attacco.
    """
    result = {
        "Q1": q1,
        "Q2": q2,
        "Q3": q3,
        "all_queries_available": all(v is not None for v in [q1, q2, q3]),
    }

    if result["all_queries_available"]:
        diff1 = q1 - q2  # Isola il gruppo eta
        inferred = diff1 - q3  # Isola la vittima
        result["diff1"] = diff1
        result["inferred_salary"] = inferred
    else:
        result["diff1"] = None
        result["inferred_salary"] = None
        # Annota quali query sono state rimosse
        result["removed_queries"] = [
            name for name, val in zip(QUERY_NAMES, [q1, q2, q3]) if val is None
        ]

    return result


def main():
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)

    victim = load_json(VICTIM_FILE)
    attack = load_json(ATTACK_FILE) if ATTACK_FILE.exists() else {}

    dept = victim["department"]
    true_salary = victim["true_salary"]

    # Query clear (senza DP) dal risultato dell'attacco
    clear_queries = attack.get("clear_queries", {})

    # Se non ci sono le clear_queries (vecchio formato), calcolale dal dataset
    if not clear_queries:
        df = pd.read_csv(DATA_DIR / "salaries.csv")
        eng = df[df["Department"] == dept]
        q1 = float(eng["Salary"].sum())
        q2 = float(eng[eng["Gender"] != victim["gender"]]["Salary"].sum())
        q3 = float(eng[(eng["Gender"] == victim["gender"]) & (eng["Age"] != victim["age"])]["Salary"].sum())
        diff1 = q1 - q2
        clear_queries = {
            "Q1": q1,
            "Q2": q2,
            "Q3": q3,
            "diff1": diff1,
            "inferred": diff1 - q3
        }

    # Carica risultati DP per ogni epsilon
    dp_results = {}
    for eps in EPSILONS:
        eps_dir = DATA_DIR / "output" / f"eps_{eps}"

        q1 = load_dp_query(eps_dir, "Q1", dept)
        q2 = load_dp_query(eps_dir, "Q2", dept)
        q3 = load_dp_query(eps_dir, "Q3", dept)

        dp_results[eps] = compute_dp_attack(q1, q2, q3)

    # Salva dati strutturati JSON
    analysis_data = {
        "victim": victim,
        "attack_without_dp": {
            "successful": attack.get("attack_successful", False),
            "inferred_salary": attack.get("inferred_salary"),
            "true_salary": true_salary,
            "clear_queries": clear_queries,
        },
        "dp_attack_results": {},
    }

    for eps in EPSILONS:
        r = dp_results[eps]
        entry = {
            "all_queries_available": r["all_queries_available"],
            "Q1_dp": r["Q1"],
            "Q2_dp": r["Q2"],
            "Q3_dp": r["Q3"],
        }

        if r["all_queries_available"]:
            inferred = r["inferred_salary"]
            error = abs(inferred - true_salary)
            error_ratio = error / true_salary

            entry.update({
                "diff1_dp": r["diff1"],
                "inferred_salary": inferred,
                "error": error,
                "error_ratio": error_ratio,
                "protected": error_ratio > ERROR_RATIO_LOW,
                "protection_level": (
                    "forte" if error_ratio > ERROR_RATIO_HIGH else
                    "moderata" if error_ratio > ERROR_RATIO_MODERATE else
                    "sufficiente"
                ),
                "noise_Q1": r["Q1"] - clear_queries["Q1"],
                "noise_Q2": r["Q2"] - clear_queries["Q2"],
                "noise_Q3": r["Q3"] - clear_queries["Q3"],
            })
        else:
            entry.update({
                "partition_removed": True,
                "protected": True,
                "protection_level": "massima (partizione rimossa)",
                "removed_queries": r.get("removed_queries", []),
            })

        analysis_data["dp_attack_results"][str(eps)] = entry

    with open(OUTPUT_DIR / "analysis_data.json", "w") as f:
        json.dump(analysis_data, f, indent=2, cls=NumpyEncoder)

    # Stampa riepilogo
    protected_count = sum(
        1 for eps in EPSILONS
        if not dp_results[eps]["all_queries_available"]
        or abs(dp_results[eps]["inferred_salary"] - true_salary) / true_salary > 0.3
    )
    print(f"Analisi: {protected_count}/{len(EPSILONS)} configurazioni proteggono la vittima (soglia errore > 30% stipendio)")


if __name__ == "__main__":
    main()



