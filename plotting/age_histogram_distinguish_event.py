import matplotlib
import matplotlib.pyplot as plt
import pandas as pd
import seaborn as sns

epsilon = 0.5

matplotlib.rc('text', usetex=True)
# sns.set_theme()
sns.set_context("paper")
# sns.set_style("darkgrid")

# _____PLOT 1 HISTOGRAM WITHOUT 80-89_____
fig, ax1 = plt.subplots(1, 1, figsize=(5, 5))
data1 = pd.DataFrame({
    "Bevanda": ["Cioccolata"] * 2 + ['Te'] * 5 + ["Caffe"] * 7
})
sns.histplot(data1, x='Bevanda', color='orange', edgecolor='black',
             linewidth=0.5, ax=ax1)

fig.savefig(f'plots/drink_histogram_no_matcha.pdf', bbox_inches='tight')
# _____PLOT 2 HISTOGRAM WITH 80-89_____
fig, ax2 = plt.subplots(1, 1, figsize=(5, 5))

data1 = pd.concat([data1, pd.DataFrame({'Bevanda': ["Matcha"] * 1})], ignore_index=True)
sns.histplot(data1, x='Bevanda', color='orange', edgecolor='black',
             linewidth=0.5, ax=ax2)


fig.savefig(f'plots/drink_histogram_yes_matcha.pdf', bbox_inches='tight')

# _____PLOT 3 HISTOGRAM WITH 80-89 THRESHOLDED_____
fig, ax3 = plt.subplots(1, 1, figsize=(5, 5))

sns.histplot(data1, x='Bevanda', color='orange', edgecolor='black',
             linewidth=0.5, ax=ax3)

children = ax3.get_children()
rectangles = [child for child in children if isinstance(child, matplotlib.patches.Rectangle)]
rand_values = [1.11678328, 2.29188425, -2.37836796, 0.89204274]

for i, rect in enumerate(rectangles[0:4]):
    new_height = round(rect.get_height() + rand_values[i])
    if new_height <= 2:
        rect.set_height(0)
    else:
        rect.set_height(new_height)

ax3.set_ylim()


fig.savefig(f'plots/drink_histogram_yes_matcha_threshold.pdf', bbox_inches='tight')