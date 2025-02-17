import matplotlib
import numpy as np
from scipy.stats import laplace
import matplotlib.pyplot as plt
import seaborn as sns

# This script creates a graph that represents the probability that adding noise to a category
# with count 1 brings the new value above the threshold
# It uses the Survival Function, defined as 1 - CDF (CDF(X) = P(X <= x))
# so it calculates for each threshold T the value of SF(T - 1) = P(X > T - 1)
# it calculates for (T - 1) because the category already has count of 1

matplotlib.rc('text', usetex=True)
sns.set_theme()
sns.set_context("paper")
sns.set_style("darkgrid")

epsilon = np.log(3)

location = 0
scale = 1 / epsilon

thresholds = np.linspace(1, 21, 21)

# fig, (ax1, ax2) = plt.subplots(1, 2)
fig, ax1 = plt.subplots(1, 1)

ax1.set_yscale("log")
# ax2.set_yscale("log")
sns.lineplot(x=thresholds, y=laplace.sf(thresholds - 1, loc=location, scale=scale), ax=ax1)
# sns.lineplot(x=thresholds, y=laplace.sf(thresholds, loc=location, scale=scale), ax=ax2)


ax1.set_ylim(10e-11, 1)
ax1.set_xlabel('Soglia')
ax1.set_ylabel('Probabilità di un evento distintivo')

# fig.savefig(f'plots/threshold_probability.pdf', bbox_inches='tight')
# plt.show()

print(laplace.sf(15 - 1, loc=location, scale=scale))
