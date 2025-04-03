import matplotlib.pyplot as plt
import matplotlib.ticker
import numpy as np
import seaborn as sns

from prob_functions import pdf

matplotlib.rc('text', usetex=True)

sns.set_theme()
sns.set_context("paper")
sns.set_style("darkgrid")

fig, ax1 = plt.subplots(1, 1, figsize=(5, 5))
ax2 = ax1.twinx()

ax1.set_xlim(495, 510)
ax1.set_xlabel('x')
ax1.set_ylabel('Laplace distributions')
ax2.set_ylabel('Ratio')

ax1_pdf1 = pdf(np.linspace(495, 510, 1000), 505, 1)
ax1_pdf2 = pdf(np.linspace(495, 510, 1000), 500, 1)
ratio = ax1_pdf2 / ax1_pdf1

pdf1_line = sns.lineplot(x=np.linspace(495, 510, 1000), y=ax1_pdf1, ax=ax1, linewidth=0.5)
pdf2_line = sns.lineplot(x=np.linspace(495, 510, 1000), y=ax1_pdf2, ax=ax1, linewidth=0.5)
ration_line = sns.lineplot(x=np.linspace(495, 510, 1000), y=ratio, ax=ax2, color='black', linewidth=0.5)

ax1.annotate(r'rapporto $\le e^{5\varepsilon}$', xy=(500, pdf(500, 500, 1)), xytext=(4, 0), textcoords='offset points',
             fontsize='small')

# ax1.yaxis.set_major_locator(matplotlib.ticker.MultipleLocator(base=0.1))
l1 = ax1.get_ylim()
l2 = ax2.get_ylim()
f = lambda x: l2[0] + (x - l1[0]) / (l1[1] - l1[0]) * (l2[1] - l2[0])
ax2.yaxis.set_major_locator(matplotlib.ticker.FixedLocator(f(ax1.get_yticks())))
ax2.grid(False)
pdf1_line.set_label(r'Lap(x$\vert500,1/\varepsilon$)')
pdf2_line.set_label(r'Lap(x$\vert505,1/\varepsilon$)')
ration_line.set_label('Ratio')


fig.savefig(f'plots/double_laplace_pdf.pdf', bbox_inches='tight')
plt.show()
