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

clear = data["attack_without_dp"]["clear_queries"]

labels = [
    r"$S_1$ (Eng.)",
    r"$-S_2$ ($\neq$F)",
    r"$S_1 - S_2$",
    r"$-S_3$ (F, $\neq$35)",
    "Stipendio\ninferito",
]

values = [
    clear["Q1"],
    -clear["Q2"],
    clear["diff1"],
    -clear["Q3"],
    clear["inferred"],
]

_palette = sns.color_palette("pastel")
colors = [_palette[0], _palette[3], _palette[2], _palette[3], _palette[1]]

fig, ax = plt.subplots(figsize=(8, 4))

bottoms = [0, 0, 0, 0, 0]
bottoms[1] = values[0]
bottoms[3] = values[2]

bar_vals = [
    values[0],
    values[1],
    values[2],
    values[3],
    values[4],
]

bars = ax.bar(labels, bar_vals, bottom=bottoms, color=colors, edgecolor="black", lw=0.8, width=0.6)

for i, (b, v, bot) in enumerate(zip(bars, bar_vals, bottoms)):
    y_pos = bot + v / 2
    if i == 4:
        label = f"\\${v:,.0f}"
    elif abs(v) > 1e6:
        label = f"\\${v / 1e6:,.1f}M"
    else:
        label = f"\\${v:,.0f}"
    ax.text(b.get_x() + b.get_width() / 2, y_pos, label,
            ha="center", va="center", fontsize=8, fontweight="bold")

# connettori tra barre di sottrazione
for src, dst in [(0, 1), (2, 3)]:
    y_line = bottoms[src] + bar_vals[src] if bar_vals[src] > 0 else bottoms[src]
    ax.plot([src + 0.35, dst - 0.35], [y_line, y_line],
            color="gray", ls="--", lw=0.8)

ax.set_ylim(0, 19e6)
ax.set_ylabel("Somma stipendi (\\$)")
ax.set_title("Difference attack: isolamento dello stipendio della vittima")
ax.yaxis.set_major_formatter(plt.FuncFormatter(lambda x, _: f"\\${x / 1e6:.0f}M"))

fig.tight_layout()
fig.savefig(SCRIPT_DIR / "difference_attack_waterfall.pdf", bbox_inches="tight")
fig.savefig(SCRIPT_DIR / "difference_attack_waterfall.png", bbox_inches="tight", dpi=200)
