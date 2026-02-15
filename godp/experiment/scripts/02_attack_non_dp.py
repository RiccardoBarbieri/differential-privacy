#!/usr/bin/env python3
"""
Dimostra l'attacco di de-anonimizzazione (difference attack) su aggregazioni non-DP.

L'attaccante conosce le caratteristiche della vittima (Department, Age, Gender)
e usa query aggregate per dedurre il suo stipendio esatto.
"""

import pandas as pd
import numpy as np
import json
from pathlib import Path

# Configurazione
SCRIPT_DIR = Path(__file__).parent
DATA_DIR = SCRIPT_DIR.parent / "data"
OUTPUT_DIR = SCRIPT_DIR.parent / "output"
DATASET_FILE = DATA_DIR / "salaries.csv"
VICTIM_FILE = DATA_DIR / "victim_info.json"
RESULTS_FILE = OUTPUT_DIR / "attack_results.json"


def load_data():
    """Carica il dataset e le informazioni sulla vittima."""
    df = pd.read_csv(DATASET_FILE)

    with open(VICTIM_FILE, "r") as f:
        victim_info = json.load(f)

    return df, victim_info


def non_dp_query_sum(df: pd.DataFrame, filters: dict = None) -> float:
    """
    Query NON differenzialmente privata: restituisce la somma esatta.

    Args:
        df: DataFrame
        filters: Dizionario di filtri {colonna: valore} o {colonna: {"op": "!=", "value": valore}}

    Returns:
        Somma degli stipendi che soddisfano i filtri
    """
    result = df.copy()

    if filters:
        for column, condition in filters.items():
            if isinstance(condition, dict):
                if condition["op"] == "!=":
                    result = result[result[column] != condition["value"]]
                elif condition["op"] == "==":
                    result = result[result[column] == condition["value"]]
            else:
                result = result[result[column] == condition]

    return result["Salary"].sum()


def non_dp_query_count(df: pd.DataFrame, filters: dict = None) -> int:
    """Query NON differenzialmente privata: restituisce il conteggio esatto.

    Args:
        df: DataFrame
        filters: Dizionario di filtri {colonna: valore} o {colonna: {"op": "!=", "value": valore}}

    Returns:
        Somma degli stipendi che soddisfano i filtri
    """
    result = df.copy()

    if filters:
        for column, condition in filters.items():
            if isinstance(condition, dict):
                if condition["op"] == "!=":
                    result = result[result[column] != condition["value"]]
                elif condition["op"] == "==":
                    result = result[result[column] == condition["value"]]
            else:
                result = result[result[column] == condition]

    return len(result)


def convert_to_native(obj):
    """Converte tipi numpy in tipi nativi per json dump."""
    if isinstance(obj, np.bool_):
        return bool(obj)
    elif isinstance(obj, (np.integer,)):
        return int(obj)
    elif isinstance(obj, (np.floating,)):
        return float(obj)
    elif isinstance(obj, np.ndarray):
        return obj.tolist()
    elif isinstance(obj, dict):
        return {k: convert_to_native(v) for k, v in obj.items()}
    elif isinstance(obj, list):
        return [convert_to_native(i) for i in obj]
    return obj


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


def difference_attack(df: pd.DataFrame, victim_info: dict) -> tuple:
    """
    Esegue l'attacco di differenza per dedurre lo stipendio della vittima.
    Restituisce i risultati e un report testuale.
    """
    department = victim_info["department"]
    age = victim_info["age"]
    gender = victim_info["gender"]
    true_salary = victim_info["true_salary"]

    report_lines = []
    report_lines.append("DIFFERENCE ATTACK - Senza Differential Privacy")
    report_lines.append("-" * 70)

    report_lines.append(f"\nInformazioni ausiliarie note all'attaccante:")
    report_lines.append(f"  - Reparto: {department}")
    report_lines.append(f"  - Età: {age}")
    report_lines.append(f"  - Genere: {gender}")

    # Query 1: Somma stipendi nel reparto
    filters_dept = {"Department": department}
    sum_dept = non_dp_query_sum(df, filters_dept)
    count_dept = non_dp_query_count(df, filters_dept)

    report_lines.append(f"\n--- Step 1: Query sul reparto ---")
    report_lines.append(f"Query: SUM(Salary) WHERE Department = '{department}'")
    report_lines.append(f"Risultato: ${sum_dept:,} ({count_dept} dipendenti)")

    # Query 2: Escludendo età
    filters_exclude_age = {
        "Department": department,
        "Age": {"op": "!=", "value": age}
    }
    sum_exclude_age = non_dp_query_sum(df, filters_exclude_age)
    count_exclude_age = non_dp_query_count(df, filters_exclude_age)

    report_lines.append(f"\n--- Step 2: Escludendo età vittima ---")
    report_lines.append(f"Query: SUM(Salary) WHERE Department = '{department}' AND Age != {age}")
    report_lines.append(f"Risultato: ${sum_exclude_age:,} ({count_exclude_age} dipendenti)")

    # Calcolo intermedio
    sum_age_only = sum_dept - sum_exclude_age
    count_age_only = count_dept - count_exclude_age

    report_lines.append(f"\n--- Step 3: Calcolo intermedio ---")
    report_lines.append(f"SUM({department}, Age={age}) = ${sum_dept:,} - ${sum_exclude_age:,} = ${sum_age_only:,}")

    # Query 3: Escludendo genere
    df_dept = df[df["Department"] == department]
    df_dept_age = df_dept[df_dept["Age"] == age]
    df_dept_age_not_gender = df_dept_age[df_dept_age["Gender"] != gender]

    sum_exclude_gender = df_dept_age_not_gender["Salary"].sum()
    count_exclude_gender = len(df_dept_age_not_gender)

    report_lines.append(f"\n--- Step 4: Escludendo genere vittima ---")
    report_lines.append(f"Query: SUM(Salary) WHERE Dept='{department}' AND Age={age} AND Gender!='{gender}'")
    report_lines.append(f"Risultato: ${sum_exclude_gender:,} ({count_exclude_gender} dipendenti)")

    # Calcolo finale
    inferred_salary = sum_age_only - sum_exclude_gender
    count_victims = count_age_only - count_exclude_gender

    report_lines.append(f"\n--- Step 5: Deduzione finale ---")
    report_lines.append(f"Stipendio vittima = ${sum_age_only:,} - ${sum_exclude_gender:,} = ${inferred_salary:,}")
    report_lines.append(f"Individui con queste caratteristiche: {count_victims}")

    # Risultato
    report_lines.append(f"\n" + "-" * 70)
    report_lines.append("RISULTATO")
    report_lines.append("-" * 70)
    report_lines.append(f"Stipendio dedotto: ${inferred_salary:,}")
    report_lines.append(f"Stipendio reale:   ${true_salary:,}")

    attack_success = (count_victims == 1 and inferred_salary == true_salary)
    if attack_success:
        report_lines.append(f"Stato: ATTACCO RIUSCITO (precisione 100%)")
    else:
        report_lines.append(f"Stato: ATTACCO FALLITO")

    result = {
        "attack_type": "difference_attack",
        "dp_protection": False,
        "victim_department": department,
        "victim_age": age,
        "victim_gender": gender,
        "true_salary": true_salary,
        "inferred_salary": inferred_salary,
        "attack_successful": attack_success,
        "num_matching_individuals": count_victims,
        "queries_used": [
            {"description": f"Sum salaries in {department}", "result": sum_dept},
            {"description": f"Sum salaries in {department}, age != {age}", "result": sum_exclude_age},
            {"description": f"Sum salaries in {department}, age = {age}, gender != {gender}", "result": sum_exclude_gender}
        ]
    }

    return result, "\n".join(report_lines)



def main():
    """Funzione principale."""
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)

    # Carica dati
    df, victim_info = load_data()

    # Esegui attacco principale
    attack_result, attack_report = difference_attack(df, victim_info)

    # Salva report
    report_file = OUTPUT_DIR / "attack_report.txt"
    with open(report_file, "w") as f:
        f.write(attack_report)

    # Salva risultati JSON
    with open(RESULTS_FILE, "w") as f:
        json.dump(attack_result, f, indent=2, cls=NumpyEncoder)

    status = "RIUSCITO" if attack_result["attack_successful"] else "FALLITO"
    print(f"Attacco: {status} | Dedotto: ${attack_result['inferred_salary']:,} | Reale: ${attack_result['true_salary']:,}")


if __name__ == "__main__":
    main()







