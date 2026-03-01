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

_palette = sns.color_palette("pastel")

epsilons_available = []
noise_q1 = []
noise_q2 = []
noise_q3 = []

for eps_str, res in sorted(data["dp_attack_results"].items(), key=lambda x: float(x[0])):
    if res["all_queries_available"]:
        epsilons_available.append(float(eps_str))
        noise_q1.append(abs(res["noise_Q1"]))
        noise_q2.append(abs(res["noise_Q2"]))
        noise_q3.append(abs(res["noise_Q3"]))

x = np.arange(len(epsilons_available))
width = 0.22

fig, ax = plt.subplots(figsize=(7, 4))

bars1 = ax.bar(x - width, noise_q1, width, label=r"$|$Rumore $S_1|$", color=_palette[0], edgecolor="black", lw=0.8)
bars2 = ax.bar(x, noise_q2, width, label=r"$|$Rumore $S_2|$", color=_palette[1], edgecolor="black", lw=0.8)
bars3 = ax.bar(x + width, noise_q3, width, label=r"$|$Rumore $S_3|$", color=_palette[2], edgecolor="black", lw=0.8)

ax.set_xticks(x)
ax.set_xticklabels([f"$\\varepsilon={e}$" for e in epsilons_available])
ax.set_ylabel("Rumore assoluto (\\$)")
ax.set_title("Rumore aggiunto per query al variare di $\\varepsilon$")
ax.legend(fontsize=8)
ax.yaxis.set_major_formatter(plt.FuncFormatter(lambda x, _: f"\\${x / 1e3:.0f}k"))

# linea teorica: scala di Laplace b = 170000 / (0.9 * eps)
eps_range = np.linspace(min(epsilons_available) - 0.5, max(epsilons_available) + 0.5, 100)
b_theoretical = 170000 / (0.9 * eps_range)
ax2 = ax.twinx()
ax2.plot(np.interp(eps_range, epsilons_available, x), b_theoretical,
         color="red", ls="--", lw=1.2, label=r"$b = \Delta f / \varepsilon_{\mathrm{agg}}$ (teorico)")
ax2.set_ylabel(r"Scala di Laplace $b$")
ax2.legend(fontsize=7, loc="upper right")
ax2.yaxis.set_major_formatter(plt.FuncFormatter(lambda x, _: f"\\${x / 1e3:.0f}k"))

fig.tight_layout()
fig.savefig(SCRIPT_DIR / "noise_vs_epsilon.pdf", bbox_inches="tight")
fig.savefig(SCRIPT_DIR / "noise_vs_epsilon.png", bbox_inches="tight", dpi=200)
