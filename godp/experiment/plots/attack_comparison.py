import matplotlib
import matplotlib.pyplot as plt
import seaborn as sns
import json
from pathlib import Path
import numpy as np

matplotlib.rc("text", usetex=True)
sns.set_theme()
sns.set_context("paper")
sns.set_style("darkgrid")

SCRIPT_DIR = Path(__file__).parent
DATA_FILE = SCRIPT_DIR.parent / "output" / "analysis_data.json"

with open(DATA_FILE) as f:
    data = json.load(f)

true_salary = data["victim"]["true_salary"]
clear_inferred = data["attack_without_dp"]["inferred_salary"]

epsilons = []
inferred_salaries = []
errors = []
labels_eps = []

for eps_str, res in sorted(data["dp_attack_results"].items(), key=lambda x: float(x[0])):
    eps = float(eps_str)
    epsilons.append(eps)
    labels_eps.append(f"$\\varepsilon={eps_str}$")
    if res["all_queries_available"]:
        inferred_salaries.append(res["inferred_salary"])
        errors.append(res["error"])
    else:
        inferred_salaries.append(None)
        errors.append(None)

_palette = sns.color_palette("pastel")

fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(10, 4))

# --- grafico sinistro: stipendio inferito ---
x_positions = range(len(epsilons) + 1)
all_labels = ["Senza DP"] + labels_eps
all_values = [clear_inferred] + inferred_salaries
colors = [_palette[3]] + [_palette[0]] * len(epsilons)

bar_vals = []
for v in all_values:
    bar_vals.append(v if v is not None else 0)

bars = ax1.bar(x_positions, bar_vals, color=colors, edgecolor="black", lw=0.8, width=0.6)

# annota partizione rimossa
for i, v in enumerate(all_values):
    if v is None:
        ax1.text(i, true_salary * 0.5, "Partizione\nrimossa",
                ha="center", va="center", fontsize=7, fontstyle="italic", color="gray")

ax1.axhline(y=true_salary, color="red", ls="--", lw=1.2, label=f"Stipendio reale (\\${true_salary:,})")
ax1.set_xticks(x_positions)
ax1.set_xticklabels(all_labels, fontsize=8)
ax1.set_ylabel("Stipendio inferito (\\$)")
ax1.set_title("Stipendio dedotto dall'attaccante")
ax1.legend(fontsize=7, loc="upper right")
ax1.yaxis.set_major_formatter(plt.FuncFormatter(lambda x, _: f"\\${x / 1e3:.0f}k"))

# --- grafico destro: errore relativo ---
dp_eps = []
dp_errors_ratio = []
for eps_str, res in sorted(data["dp_attack_results"].items(), key=lambda x: float(x[0])):
    if res["all_queries_available"]:
        dp_eps.append(float(eps_str))
        dp_errors_ratio.append(res["error_ratio"] * 100)

ax2.bar(range(len(dp_eps)), dp_errors_ratio, color=_palette[1], edgecolor="black", lw=0.8, width=0.6)
ax2.set_xticks(range(len(dp_eps)))
ax2.set_xticklabels([f"$\\varepsilon={e}$" for e in dp_eps], fontsize=8)
ax2.set_ylabel("Errore relativo (\\%)")
ax2.set_title("Errore dell'attacco rispetto allo stipendio reale")

for i, (e, v) in enumerate(zip(dp_eps, dp_errors_ratio)):
    ax2.text(i, v + 1.5, f"{v:.1f}\\%", ha="center", va="bottom", fontsize=8)

fig.tight_layout()
fig.savefig(SCRIPT_DIR / "attack_comparison.pdf", bbox_inches="tight")
fig.savefig(SCRIPT_DIR / "attack_comparison.png", bbox_inches="tight", dpi=200)
