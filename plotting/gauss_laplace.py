import matplotlib
import matplotlib.pyplot as plt
import numpy as np
import seaborn as sns

from prob_functions import pdf

matplotlib.rc('text', usetex=True)



epsilon = 1
values = np.linspace(-10, 10, 1000)
gauss_pdf = pdf(values, loc=0, scale=np.exp(epsilon), distribution="normal")
laplace_pdf = pdf(values, loc=0, scale=1/epsilon)

sns.set_theme()
sns.set_context("paper")
sns.set_style("darkgrid")

fig, ax = plt.subplots(1, 1, figsize=(5, 5))

sns.lineplot(x=values, y=gauss_pdf, ax=ax, linewidth=0.8)
sns.lineplot(x=values, y=laplace_pdf, ax=ax, linewidth=0.8)

ax.set_xlim(-10, 10)

fig.savefig('plots/gauss_laplace_e_1.pdf', bbox_inches='tight')
plt.show()
