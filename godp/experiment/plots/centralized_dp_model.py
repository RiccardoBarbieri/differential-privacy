import matplotlib
import matplotlib.pyplot as plt
import seaborn as sns
from matplotlib.patches import FancyBboxPatch

matplotlib.rc("text", usetex=True)
sns.set_theme()
sns.set_context("paper")
sns.set_style("darkgrid")

_palette = sns.color_palette("pastel")

COL = {
    "subject": _palette[0],  # blue
    "db": _palette[1],  # orange
    "curator": _palette[2],  # green
    "analyst": _palette[4],  # purple
    "attacker": _palette[3],  # red
    "arrow": "#444444",
    "info": _palette[5],  # brown/warm
}


def box(ax, x, y, w, h, label, color, fontsize=10, bold=False):
    r = FancyBboxPatch(
        (x, y),
        w,
        h,
        boxstyle="round,pad=0.1",
        facecolor=color,
        edgecolor="black",
        lw=1.2,
    )
    ax.add_patch(r)
    ax.text(
        x + w / 2,
        y + h / 2,
        label,
        ha="center",
        va="center",
        fontsize=fontsize,
        weight="bold" if bold else "normal",
    )


def arrow(ax, x1, y1, x2, y2, label=None, lbl_dy=0.13):
    ax.annotate(
        "",
        xy=(x2, y2),
        xytext=(x1, y1),
        arrowprops=dict(arrowstyle="-|>", color=COL["arrow"], lw=1.4),
    )
    if label:
        mx, my = (x1 + x2) / 2, (y1 + y2) / 2
        ax.text(
            mx,
            my + lbl_dy,
            label,
            ha="center",
            va="bottom",
            fontsize=8,
            color=COL["arrow"],
        )


fig, ax = plt.subplots(1, 1, figsize=(11, 4.5))
ax.set_xlim(0, 11)
ax.set_ylim(0, 4.5)
ax.axis("off")

# _____SUBJECTS_____
sx, sw, sh = 0.15, 1.05, 0.5
box_y = [0.5, 1.4, 3.5]
dot_y = 2.45
box_labels = ["Soggetto 1", "Soggetto 2", "Soggetto N"]

for yi, lbl in zip(box_y, box_labels):
    box(ax, sx, yi, sw, sh, lbl, COL["subject"], fontsize=8)

ax.text(
    sx + sw / 2, dot_y, r"$\vdots$", ha="center", va="center", fontsize=16, color="#555"
)

# _____DATABASE_____
db_x, db_w, db_h = 1.9, 1.3, 0.9
db_cy = (box_y[0] + box_y[-1] + sh) / 2
db_y = db_cy - db_h / 2
box(
    ax, db_x, db_y, db_w, db_h, "Database\noriginale", COL["db"], fontsize=10, bold=True
)

# soggetti -> database (solo i soggetti con box)
for yi in box_y:
    arrow(ax, sx + sw, yi + sh / 2, db_x, db_cy)

# _____CURATOR_____
cur_x, cur_w, cur_h = 4.2, 1.9, 2.0
cur_y = db_cy - cur_h / 2
cur_r = cur_x + cur_w
box(
    ax,
    cur_x,
    cur_y,
    cur_w,
    cur_h,
    "Curatore\n(GoDP)",
    COL["curator"],
    fontsize=11,
    bold=True,
)

arrow(ax, db_x + db_w, db_cy, cur_x, db_cy)

# _____ANALYST_____
an_x, an_w, an_h = 7.4, 1.3, 0.6
an_y = cur_y + cur_h - an_h
box(ax, an_x, an_y, an_w, an_h, "Analista", COL["analyst"], fontsize=10, bold=True)

y_an_q = cur_y + cur_h * 0.78
y_an_r = cur_y + cur_h * 0.62

# analista -> curatore
arrow(ax, an_x, an_y + an_h * 0.80, cur_r, y_an_q, label="query", lbl_dy=0.12)
# curatore -> analista
arrow(ax, cur_r, y_an_r, an_x, an_y + an_h * 0.25, label="risultato DP", lbl_dy=0.12)

# _____ATTACKER_____
at_x, at_w, at_h = 7.4, 1.3, 0.6
at_y = cur_y
box(ax, at_x, at_y, at_w, at_h, "Attaccante", COL["attacker"], fontsize=10, bold=True)

y_at_q = cur_y + cur_h * 0.38
y_at_r = cur_y + cur_h * 0.20

# attaccante -> curatore
arrow(ax, at_x, at_y + at_h * 0.80, cur_r, y_at_q, label="query", lbl_dy=0.12)
# curatore -> attaccante
arrow(ax, cur_r, y_at_r, at_x, at_y + at_h * 0.25, label="risultato DP", lbl_dy=-0.22)

# _____AUXILIARY INFO_____
info_x, info_w, info_h = 9.1, 1.5, 0.5
info_y = at_y + at_h / 2 - info_h / 2
box(ax, info_x, info_y, info_w, info_h, "Info ausiliarie", COL["info"], fontsize=8)
arrow(ax, info_x, info_y + info_h / 2, at_x + at_w, at_y + at_h / 2)

fig.savefig("experiment/plots/centralized_dp_model.pdf", bbox_inches="tight")
fig.savefig("experiment/plots/centralized_dp_model.png", bbox_inches="tight", dpi=200)
