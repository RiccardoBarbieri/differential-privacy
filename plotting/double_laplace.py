import matplotlib
import matplotlib.pyplot as plt
import numpy as np
import seaborn as sns

from prob_functions import pdf

matplotlib.rc('text', usetex=True)

sns.set_theme()
sns.set_context("paper")
sns.set_style("darkgrid")

fig, ax = plt.subplots(1, 1, figsize=(5, 5))
sns.lineplot(x=np.linspace(495, 510, 1000), y=pdf(np.linspace(495, 510, 1000), 505, 1), color='blue', ax=ax,
             linewidth=0.5)
sns.lineplot(x=np.linspace(495, 510, 1000), y=pdf(np.linspace(495, 510, 1000), 500, 1), color='orange', ax=ax,
             linewidth=0.5)
ax.annotate(r'rapporto = $e^{5\varepsilon}$', xy=(500, pdf(500, 500, 1)), xytext=(4, 0), textcoords='offset points', fontsize='small')


ax.set_xlim(495, 510)

ax.set_xlabel('x')
ax.set_ylabel(r'Lap(x$\vert500,1/\varepsilon$)' + '\n' + r'Lap(x$\vert505,1/\varepsilon$)')

fig.savefig(f'plots/double_laplace_pdf.pdf', bbox_inches='tight')
plt.show()
