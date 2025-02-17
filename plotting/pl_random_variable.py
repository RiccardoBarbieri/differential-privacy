from typing import List

import matplotlib
import matplotlib.pyplot as plt
import numpy as np
import seaborn as sns
from matplotlib.axes import Axes

from prob_functions import pdf, cdf, plrv

DISTRIBUTION = 'normal'

matplotlib.rc('text', usetex=True)
sns.set_theme()
sns.set_context("paper")
sns.set_style("darkgrid")


def get_scale(epsilon: float, distribution: str = 'laplace'):
    if distribution == 'laplace':
        scale = 1 / epsilon
    elif distribution == 'normal':
        # with variance 3
        scale = np.sqrt(3)
    return scale


epsilon = np.log(3)
location1 = 1000
location2 = 1001
location_mid = (location1 + location2) / 2

values = np.linspace(990, 1050, 10000)

axes: List[Axes]
fig1, ax1 = plt.subplots(1, 1, figsize=(8, 4))
fig2, ax2 = plt.subplots(1, 1, figsize=(8, 4))
fig3, ax3 = plt.subplots(1, 1, figsize=(8, 4))
fig4, ax4 = plt.subplots(1, 1, figsize=(8, 4))

scale = get_scale(epsilon, DISTRIBUTION)

pdf_loc1_values = pdf(values, loc=location1, scale=scale, distribution=DISTRIBUTION)
pdf_loc2_values = pdf(values, loc=location2, scale=scale, distribution=DISTRIBUTION)
plrv_values = plrv(values, location1, location2, scale, distribution=DISTRIBUTION)
cdf_values = cdf(values, loc=location1, scale=scale, distribution=DISTRIBUTION)

ax1.plot(values, pdf_loc1_values, linewidth=0.8)
ax1.plot(values, pdf_loc2_values, linewidth=0.8)

ax2.plot(values, plrv_values)

if DISTRIBUTION == "normal":
    ax3.hlines(y=3, xmin=0, xmax=0.053, linewidth=0.8, linestyles="dashed")
    ax3.plot([0.053], [3], color='red', marker='o', markersize=3, zorder=9)
ax3.plot(cdf_values, np.exp(plrv_values), linewidth=0.5)

ax4.plot(cdf_values, 3 / np.exp(plrv_values), linewidth=0.5)
ax4.hlines(y=1, xmin=0, xmax=1, linewidth=0.8, linestyles="dashed")
ax4.fill_between(cdf_values[:1201], 3 / np.exp(plrv_values[:1201]), 1, alpha=0.3)

if DISTRIBUTION == "laplace":
    points = [999, 1000.5, 1003]
elif DISTRIBUTION == "normal":
    points = [995.5, 997]

for i, point in enumerate(points):
    index = int(((point - values.min()) / (values.max() - values.min())) * len(values))
    pdf_p_l1 = pdf_loc1_values[index]
    pdf_p_l2 = pdf_loc2_values[index]
    cdf_p_l1 = cdf_values[index]
    plrv_p = plrv_values[index]
    exp_plrv_p = np.exp(plrv_values[index])

    if pdf_p_l1 >= pdf_p_l2:
        ymax_ax1 = pdf_p_l1
    else:
        ymax_ax1 = pdf_p_l2
    ax1.plot([point], [ymax_ax1], color='red', marker='o', markersize=3, zorder=9)
    ax2.plot([point], [plrv_p], color='red', marker='o', markersize=3, zorder=9)
    ax3.plot([cdf_p_l1], [exp_plrv_p], color='red', marker='o', markersize=3, zorder=9)

    print(f"O{i + 1} on {location1} = ({point}, {pdf_p_l1})")
    print(f"O{i + 1} on {location2} = ({point}, {pdf_p_l2})")
    knowledge_gain = round(np.exp(plrv_p), 4)
    print(f"Knowledge gain:   L({cdf_p_l1}) = {knowledge_gain}")

    ax1.plot([point, point], [0, ymax_ax1], linestyle='dashed', color='black', linewidth=0.8)
    ax1.annotate(f'O{i + 1}', xy=(point, 0), xytext=(4, -10), textcoords='offset points', fontsize='medium')

    ax3.plot([cdf_p_l1] * 2,
             [0, exp_plrv_p], linestyle='dashed', color='black', linewidth=0.8)
    ax3.annotate(f'O{i + 1}', xy=(cdf_p_l1, 0), xytext=(4, -10), textcoords='offset points', fontsize='medium')

    print("_" * 50)

ax1.set_xlim(left=995, right=1005)
ax1.set_ylim(bottom=0)

ax2.set_xlim(left=995, right=1005)

ax3.set_xlim(left=0, right=1)
if DISTRIBUTION == "laplace":
    ax3.set_ylim(bottom=0)
elif DISTRIBUTION == "normal":
    ax3.set_ylim(bottom=0, top=6)
ax3.set_ylabel(r'$e^{\mathcal{L}}$')
ax3.set_xlabel(r'$CDF_{\mathcal{L}}$')
ax4.set_ylim(bottom=0, top=1.5)
ax4.set_xlim(left=0, right=0.15)
ax4.set_ylabel(r'$e^{\mathcal{\varepsilon}}/e^{\mathcal{L}}$')
ax4.set_xlabel(r'$CDF_{\mathcal{L}}$')



fig1.savefig(f'plots/{DISTRIBUTION}_plrv_two_dists.pdf', bbox_inches='tight')
fig2.savefig(f'plots/{DISTRIBUTION}_plrv.pdf', bbox_inches='tight')
fig3.savefig(f'plots/{DISTRIBUTION}_e^plrv.pdf', bbox_inches='tight')
fig4.savefig(f'plots/{DISTRIBUTION}_inverse_normalized_plrv.pdf', bbox_inches='tight')
plt.show()
