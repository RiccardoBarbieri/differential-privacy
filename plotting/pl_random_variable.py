from typing import List

import matplotlib.pyplot as plt
import numpy as np
from matplotlib.axes import Axes
from matplotlib.widgets import Slider

from plot_functions import plot_plrv
from prob_functions import pdf, cdf, plrv

DISTRIBUTION = 'laplace'


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

values = np.linspace(990, 1050, 5000)

axes: List[Axes]
fig, axes = plt.subplots(3, 1, figsize=(10, 15))


def update(val):
    pdf_loc1_values = pdf(values, loc=location1, scale=get_scale(val))
    pdf_loc2_values = pdf(values, loc=location2, scale=get_scale(val))
    plrv_values = plrv(values, location1, location2, get_scale(val))
    cdf_values = cdf(values, loc=location1, scale=get_scale(val))
    # lines
    axes[0].get_lines()[0].set_ydata(pdf(values, loc=location1, scale=get_scale(val)))
    axes[0].get_lines()[1].set_ydata(pdf(values, loc=location2, scale=get_scale(val)))
    axes[1].get_lines()[0].set_ydata(plrv(values, location1, location2, get_scale(val)))
    axes[2].get_lines()[0].set_ydata(np.exp(plrv(values, location1, location2, get_scale(val))))
    # points
    points = [axes[0].get_lines()[i].get_xdata()[0] for i in range(2, len(axes[0].get_lines()))]

    for i, point in enumerate(points):
        index = int(((point - values.min()) / (values.max() - values.min())) * len(values))
        pdf_p_l1 = pdf_loc1_values[index]
        pdf_p_l2 = pdf_loc2_values[index]
        plrv_p = plrv_values[index]
        exp_plrv_p = np.exp(plrv_values[index])

        if pdf_p_l1 >= pdf_p_l2:
            ymax_ax1 = pdf_p_l1
        else:
            ymax_ax1 = pdf_p_l2
        # offset in index accounts for lines plotted first
        axes[0].get_lines()[i + 2].set_ydata([ymax_ax1])
        axes[1].get_lines()[i + 1].set_ydata([plrv_p])
        # offset index accounts for first line and alternation point-vlines in third ax
        axes[2].get_lines()[2 * i + 1].set_ydata([exp_plrv_p])
        axes[2].get_lines()[2 * i + 2].set_ydata([0, exp_plrv_p])

        patch_point_ax1 = (point, ymax_ax1)
        patch_point_ax2 = (point, plrv_p)
        axes[1].patches[i].xy1 = patch_point_ax1
        axes[1].patches[i].xy2 = patch_point_ax2

    # Limits at 5%
    axes[0].set_ylim(bottom=0, top=pdf_loc1_values.max() + pdf_loc1_values.max() / 20)
    axes[1].set_ylim(bottom=plrv_values.min() + plrv_values.min() / 20, top=plrv_values.max() + plrv_values.max() / 20)
    axes[2].set_ylim(bottom=0, top=np.exp(plrv_values).max() + np.exp(plrv_values).max() / 20)

    fig.canvas.draw_idle()


ax_epsilon_slider: Axes = fig.add_axes((0.3, 0.05, 0.4, 0.03))
epsilon_slider = Slider(
    ax=ax_epsilon_slider,
    label='$\epsilon$',
    valmin=0.001,
    valmax=14,
    valinit=epsilon,
)
epsilon_slider.on_changed(update)

plot_plrv(fig, axes[0], axes[1], axes[2], values, location1, location2, get_scale(epsilon))

plt.savefig(f'{DISTRIBUTION}_plrv.svg', bbox_inches='tight')
plt.show()
