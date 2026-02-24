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

    Strategia dell'attacco (3 query):
      Q1: SUM(Salary) WHERE Dept = Engineering
      Q2: SUM(Salary) WHERE Dept = Engineering AND Gender != 'F'
      Q3: SUM(Salary) WHERE Dept = Engineering AND Gender = 'F' AND Age != 35

    Calcolo:
      Diff1 = Q1 - Q2  →  somma stipendi di tutte le donne in Engineering
      Stipendio_vittima = Diff1 - Q3  →  isola la vittima (unica donna 35y in Eng)

    Restituisce i risultati.
    """
    department = victim_info["department"]
    age = victim_info["age"]
    gender = victim_info["gender"]
    true_salary = victim_info["true_salary"]

    # --- Q1: Somma stipendi nel reparto ---
    filters_q1 = {"Department": department}
    sum_q1 = non_dp_query_sum(df, filters_q1)
    count_q1 = non_dp_query_count(df, filters_q1)

    # --- Q2: Escludendo il genere della vittima ---
    filters_q2 = {
        "Department": department,
        "Gender": {"op": "!=", "value": gender}
    }
    sum_q2 = non_dp_query_sum(df, filters_q2)
    count_q2 = non_dp_query_count(df, filters_q2)

    # Calcolo intermedio: Diff1 = Q1 - Q2 = tutte le donne in Engineering
    diff1 = sum_q1 - sum_q2
    count_diff1 = count_q1 - count_q2

    # --- Q3: Stesse donne del reparto, escludendo l'eta della vittima ---
    filters_q3 = {
        "Department": department,
        "Gender": {"op": "==", "value": gender},
        "Age": {"op": "!=", "value": age}
    }
    sum_q3 = non_dp_query_sum(df, filters_q3)
    count_q3 = non_dp_query_count(df, filters_q3)

    # Calcolo finale: Stipendio vittima = Diff1 - Q3
    inferred_salary = diff1 - sum_q3
    count_victims = count_diff1 - count_q3

    # Risultato
    attack_success = (count_victims == 1 and inferred_salary == true_salary)

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
            {"description": f"Q1: Sum salaries in {department}", "result": sum_q1, "count": count_q1},
            {"description": f"Q2: Sum salaries in {department}, gender != {gender}", "result": sum_q2, "count": count_q2},
            {"description": f"Q3: Sum salaries in {department}, gender = {gender}, age != {age}", "result": sum_q3, "count": count_q3}
        ],
        "clear_queries": {
            "Q1": sum_q1,
            "Q2": sum_q2,
            "Q3": sum_q3,
            "diff1": diff1,
            "inferred": inferred_salary
        }
    }

    return result



def main():
    """Funzione principale."""
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)

    # Carica dati
    df, victim_info = load_data()

    # Esegui attacco principale
    attack_result = difference_attack(df, victim_info)

    # Salva risultati JSON
    with open(RESULTS_FILE, "w") as f:
        json.dump(attack_result, f, indent=2, cls=NumpyEncoder)

    status = "RIUSCITO" if attack_result["attack_successful"] else "FALLITO"
    print(f"Attacco: {status} | Dedotto: ${attack_result['inferred_salary']:,} | Reale: ${attack_result['true_salary']:,}")


if __name__ == "__main__":
    main()







