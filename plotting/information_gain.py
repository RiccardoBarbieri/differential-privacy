import matplotlib
import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
import seaborn as sns
from matplotlib.colors import ListedColormap

matplotlib.rc('text', usetex=True)
sns.set_theme()
sns.set_context("paper")
sns.set_style("darkgrid")


# This script creates a plot that shows the Information Gain, given an Initial Suspicion that a sample
# belongs to a database, given the result of the mechanism on said database with varying values of epsilon

epsilons = np.arange(0, 8, 0.5)
suspicions = np.linspace(0, 1, 1000)


def lower_bound(initial_suspicion: np.ndarray, epsilon: float):
    return initial_suspicion / (np.exp(epsilon) + (1 - np.exp(epsilon)) * initial_suspicion)


def upper_bound(initial_suspicion: np.ndarray, epsilon: float):
    return (np.exp(epsilon) * initial_suspicion) / (1 + (np.exp(epsilon) - 1) * initial_suspicion)


fig, ax = plt.subplots(1, 1, figsize=(5.9, 5.9))

palette = sns.color_palette("rocket_r", n_colors=len(epsilons)).as_hex()
palette = sns.color_palette("rainbow", n_colors=len(epsilons)).as_hex()
cmap = ListedColormap(palette)

data = {'suspicions': suspicions}
df = pd.DataFrame(data=data)
for i, epsilon in enumerate(epsilons):
    df[f"lower_bound{i}"] = lower_bound(suspicions, epsilon)
    df[f"upper_bound{i}"] = upper_bound(suspicions, epsilon)

sns.lineplot(x=np.linspace(0, 1, 10), y=np.linspace(0, 1, 10), color='black', linewidth=0.5, ax=ax)
for i, _ in reversed(list(enumerate(epsilons))):
    sns.lineplot(data=df, x="suspicions", y=f"lower_bound{i}", color='black', linewidth=0.5, ax=ax)
    sns.lineplot(data=df, x="suspicions", y=f"upper_bound{i}", color='black', linewidth=0.5, ax=ax)
    if i == 0:
        continue
    color = palette[i - 1]
    ax.fill_between(df['suspicions'], y1=df[f"upper_bound{i}"], y2=df[f"lower_bound{i}"], color=color)

# Create a discrete colorbar
norm = plt.Normalize(vmin=0, vmax=len(epsilons))
sm = plt.cm.ScalarMappable(cmap=cmap, norm=norm)
sm.set_array([])

# Add the colorbar to the figure
cb = fig.colorbar(sm, ticks=np.append(epsilons, 8), ax=ax, shrink=0.7)
cb.set_label('Epsilon')

# Customize the colorbar labels
# cb.set_ticklabels(epsilons)

cb.mappable.set_clim(0, 8)

ax.set_ylim(0, 1)
ax.set_xlim(0, 1)
ax.set_aspect('equal', adjustable='box')
ax.set_title('Guadagno di informazione')
ax.set_xlabel('Sospetto iniziale')
ax.set_ylabel('Sospetto finale')

plt.savefig("plots/information_gain.pdf", bbox_inches='tight')
plt.show()
