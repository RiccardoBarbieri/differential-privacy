#!/usr/bin/env python3
"""
Generazione mock dataset stipendi.

Il dataset include una "vittima" target per dimostrare un "difference attack"
"""

import pandas as pd
import numpy as np
from pathlib import Path

SCRIPT_DIR = Path(__file__).parent
DATA_DIR = SCRIPT_DIR.parent / "data"
OUTPUT_FILE = DATA_DIR / "salaries.csv"

np.random.seed(42)

NUM_EMPLOYEES = 1500

# distribuzione dei reparti
DEPARTMENTS = {
    "Engineering": 0.30,
    "Marketing": 0.20,
    "Sales": 0.25,
    "HR": 0.10,
    "Finance": 0.15
}

# Range di età per reparto (min, max)
AGE_RANGES = {
    "Engineering": (24, 55),
    "Marketing": (22, 50),
    "Sales": (23, 60),
    "HR": (25, 55),
    "Finance": (26, 60)
}

# range di stipendi per reparto
SALARY_RANGES = {
    "Engineering": (55000, 150000),
    "Marketing": (40000, 100000),
    "Sales": (35000, 120000),
    "HR": (38000, 85000),
    "Finance": (50000, 130000)
}

# configurazione della vittima target
VICTIM_CONFIG = {
    "ID": "VICTIM_001",
    "Department": "Engineering",
    "Age": 35,
    "Gender": "F",
    "YearsExperience": 8,
    "Salary": 95000
}


def generate_employee(employee_id: int, departments_list: list, dept_weights: list) -> dict:
    """Genera un dipendente con dati casuali."""

    department = np.random.choice(departments_list, p=dept_weights)
    age_min, age_max = AGE_RANGES[department]
    salary_min, salary_max = SALARY_RANGES[department]

    age = np.random.randint(age_min, age_max + 1)
    gender = np.random.choice(["M", "F"], p=[0.55, 0.45])

    max_experience = min(age - 22, 40)  # assumo inizio lavoro a 22 anni
    years_experience = max(0, np.random.randint(0, max_experience + 1))

    # stipendio incrementale con l'esperienza
    base_salary = np.random.uniform(salary_min, salary_max)
    experience_bonus = years_experience * np.random.uniform(500, 1500)
    salary = int(base_salary + experience_bonus)

    return {
        "ID": f"EMP_{employee_id:04d}",
        "Department": department,
        "Age": age,
        "Gender": gender,
        "YearsExperience": years_experience,
        "Salary": salary
    }


def ensure_unique_victim(df: pd.DataFrame) -> pd.DataFrame:
    # rimozione righe nel reparto Engineering che sono donne di 35 anni
    mask = (
        (df["Department"] == VICTIM_CONFIG["Department"]) &
        (df["Age"] == VICTIM_CONFIG["Age"]) &
        (df["Gender"] == VICTIM_CONFIG["Gender"]) &
        (df["ID"] != VICTIM_CONFIG["ID"])
    )

    # Cambia le caratteristiche di questi dipendenti
    df.loc[mask, "Age"] = df.loc[mask, "Age"].apply(
        lambda x: x + np.random.choice([-2, -1, 1, 2])
    )

    return df


def generate_dataset() -> pd.DataFrame:
    departments_list = list(DEPARTMENTS.keys())
    dept_weights = list(DEPARTMENTS.values())

    # genera dipendenti casuali
    employees = []
    for i in range(NUM_EMPLOYEES - 1):  # -1 per lasciare spazio alla vittima
        employees.append(generate_employee(i, departments_list, dept_weights))

    employees.append(VICTIM_CONFIG.copy())

    df = pd.DataFrame(employees)

    df = ensure_unique_victim(df)

    df = df.sample(frac=1, random_state=42).reset_index(drop=True)

    return df

def main():
    import json

    DATA_DIR.mkdir(parents=True, exist_ok=True)
    OUTPUT_DIR = SCRIPT_DIR.parent / "output"
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)

    df = generate_dataset()

    df.to_csv(OUTPUT_FILE, index=False)

    victim_info = {
        "victim_id": VICTIM_CONFIG["ID"],
        "department": VICTIM_CONFIG["Department"],
        "age": VICTIM_CONFIG["Age"],
        "gender": VICTIM_CONFIG["Gender"],
        "true_salary": VICTIM_CONFIG["Salary"]
    }

    victim_file = DATA_DIR / "victim_info.json"
    with open(victim_file, "w") as f:
        json.dump(victim_info, f, indent=2)

    print(f"Dataset: {len(df)} record -> {OUTPUT_FILE.name}")
    print(f"Vittima: {VICTIM_CONFIG['Department']}, {VICTIM_CONFIG['Age']}y, {VICTIM_CONFIG['Gender']} -> ${VICTIM_CONFIG['Salary']:,}")


if __name__ == "__main__":
    main()

