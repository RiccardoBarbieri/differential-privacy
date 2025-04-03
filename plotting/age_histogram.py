import matplotlib
import matplotlib.pyplot as plt
import pandas as pd
import seaborn as sns

epsilon = 0.5

matplotlib.rc('text', usetex=True)
# sns.set_theme()
sns.set_context("paper")
# sns.set_style("darkgrid")

# _____PLOT 1 NORMAL HISTOGRAM_____
fig, ax1 = plt.subplots(1, 1, figsize=(5, 5))
data1 = pd.DataFrame({
    'Età': [20, 21, 21, 22, 22, 23, 23, 24, 24, 25, 30, 31, 31, 32, 32, 33, 34, 34, 35, 35, 36, 36, 37, 37, 38, 41, 42,
            43, 43, 43, 44, 45, 51, 52, 52, 52, 53, 54, 61, 62, 62, 63, 71, 72]
})
sns.histplot(data1, x='Età', bins=[10, 20, 30, 40, 50, 60, 70, 80, 90], color='orange', edgecolor='black',
             linewidth=0.5, ax=ax1)
fig.savefig(f'plots/age_histogram.pdf', bbox_inches='tight')

# _____PLOT 2 NEGATIVE HISTOGRAM_____
fig, ax2 = plt.subplots(1, 1, figsize=(5, 5))


sns.histplot(data1, x='Età', bins=[10, 20, 30, 40, 50, 60, 70, 80, 90], color='orange', edgecolor='black',
             linewidth=0.5,
             weights=[1 if x < 80 else -1 for x in data1['Età']], ax=ax2)
children = ax2.get_children()
rectangles = [child for child in children if isinstance(child, matplotlib.patches.Rectangle)]
rand_values = [1.11678328, 2.29188425, -0.65355953, 3.09204274, 1.40710785, -0.19950759,
               -1.37836796, -1.67339696]
for i, rect in enumerate(rectangles[0:8]):
    rect.set_height(rect.get_height() + rand_values[i])

ax2.set_ylim(bottom=-2)
fig.savefig(f'plots/age_histogram_negative.pdf', bbox_inches='tight')

# _____PLOT 3 NEGATIVE HISTOGRAM POST PROCESSED_____
fig, ax3 = plt.subplots(1, 1, figsize=(5, 5))
sns.histplot(data1, x='Età', bins=[10, 20, 30, 40, 50, 60, 70, 80, 90], color='orange', edgecolor='black',
             linewidth=0.5,
             weights=[1 if x < 80 else -1 for x in data1['Età']], ax=ax3)
children = ax3.get_children()
rectangles = [child for child in children if isinstance(child, matplotlib.patches.Rectangle)]
for i, rect in enumerate(rectangles[0:8]):
    new_height = rect.get_height() + rand_values[i]
    if new_height < 0:
        rect.set_height(0)
    else:
        rect.set_height(round(new_height, 0))

fig.savefig(f'plots/age_histogram_negative_post_processed.pdf', bbox_inches='tight')
plt.show()
