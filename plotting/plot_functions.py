import seaborn as sns
from prob_functions import pdf, cdf, plrv
import numpy as np
import matplotlib




# def update(val):
#     pdf_loc1_values = pdf(values, loc=location1, scale=get_scale(val))
#     pdf_loc2_values = pdf(values, loc=location2, scale=get_scale(val))
#     plrv_values = plrv(values, location1, location2, get_scale(val))
#     cdf_values = cdf(values, loc=location1, scale=get_scale(val))
#     # lines
#     axes[0].get_lines()[0].set_ydata(pdf(values, loc=location1, scale=get_scale(val)))
#     axes[0].get_lines()[1].set_ydata(pdf(values, loc=location2, scale=get_scale(val)))
#     axes[1].get_lines()[0].set_ydata(plrv(values, location1, location2, get_scale(val)))
#     axes[2].get_lines()[0].set_ydata(np.exp(plrv(values, location1, location2, get_scale(val))))
#     # points
#     points = [axes[0].get_lines()[i].get_xdata()[0] for i in range(2, len(axes[0].get_lines()))]
#
#     for i, point in enumerate(points):
#         index = int(((point - values.min()) / (values.max() - values.min())) * len(values))
#         pdf_p_l1 = pdf_loc1_values[index]
#         pdf_p_l2 = pdf_loc2_values[index]
#         plrv_p = plrv_values[index]
#         exp_plrv_p = np.exp(plrv_values[index])
#
#         if pdf_p_l1 >= pdf_p_l2:
#             ymax_ax1 = pdf_p_l1
#         else:
#             ymax_ax1 = pdf_p_l2
#         # offset in index accounts for lines plotted first
#         axes[0].get_lines()[i + 2].set_ydata([ymax_ax1])
#         axes[1].get_lines()[i + 1].set_ydata([plrv_p])
#         # offset index accounts for first line and alternation point-vlines in third ax
#         axes[2].get_lines()[2 * i + 1].set_ydata([exp_plrv_p])
#         axes[2].get_lines()[2 * i + 2].set_ydata([0, exp_plrv_p])
#
#         patch_point_ax1 = (point, ymax_ax1)
#         patch_point_ax2 = (point, plrv_p)
#         axes[1].patches[i].xy1 = patch_point_ax1
#         axes[1].patches[i].xy2 = patch_point_ax2
#
#     # Limits at 5%
#     axes[0].set_ylim(bottom=0, top=pdf_loc1_values.max() + pdf_loc1_values.max() / 20)
#     axes[1].set_ylim(bottom=plrv_values.min() + plrv_values.min() / 20, top=plrv_values.max() + plrv_values.max() / 20)
#     axes[2].set_ylim(bottom=0, top=np.exp(plrv_values).max() + np.exp(plrv_values).max() / 20)
#
#     fig.canvas.draw_idle()

# ax_epsilon_slider: Axes = fig.add_axes((0.3, 0.05, 0.4, 0.03))
# epsilon_slider = Slider(
#     ax=ax_epsilon_slider,
#     label='$\epsilon$',
#     valmin=0.001,
#     valmax=14,
#     valinit=epsilon,
# )
# epsilon_slider.on_changed(update)

def plot_plrv(fig, ax1, ax2, ax3, values: np.ndarray, location1, location2, scale):
    location_mid = (location1 + location2) / 2
    pdf_loc1_values = pdf(values, loc=location1, scale=scale)
    pdf_loc2_values = pdf(values, loc=location2, scale=scale)
    plrv_values = plrv(values, location1, location2, scale)
    cdf_values = cdf(values, loc=location1, scale=scale)

    ax1.plot(values, pdf_loc1_values, linewidth=0.8)
    ax1.plot(values, pdf_loc2_values, linewidth=0.8)

    ax2.plot(values, plrv_values)

    ax3.plot(cdf_values, np.exp(plrv_values))

    points = [999, 1000.5, 1003]

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

        # if point < location_mid:
        #     text_voffset_border = -10
        # elif point > location_mid:
        #     text_voffset_border = 6
        # else:
        #     text_voffset_border = 0 < 1
        ax1.annotate(f'O{i + 1}', xy=(point, 0), xytext=(4, -10), textcoords='offset points', fontsize='x-small')
        # ax2.annotate(f'$\epsilon$ = {round(plrv_p, 2)}', xy=(point, plrv_p), xytext=(5, text_voffset_border),
        #              textcoords='offset points', fontsize='x-small')
        # ax3.annotate(f'$\mathcal{{L}}({round(cdf_p_l1, 2)}) = {round(exp_plrv_p, 2)}$', xy=(cdf_p_l1, exp_plrv_p),
        #              xytext=(5, text_voffset_border),
        #              textcoords='offset points', fontsize='x-small')

        print(f"O{i + 1} on {location1} = ({point}, {pdf_p_l1})")
        print(f"O{i + 1} on {location2} = ({point}, {pdf_p_l2})")
        knowledge_gain = round(np.exp(plrv_p), 4)
        print(f"Knowledge gain:   L({cdf_p_l1}) = {knowledge_gain}")

        patch_point_ax1 = (point, ymax_ax1)
        patch_point_ax2 = (point, plrv_p)
        patch = matplotlib.patches.ConnectionPatch(xyA=patch_point_ax1, xyB=patch_point_ax2, axesA=ax1, axesB=ax2,
                                                   coordsA='data', coordsB='data', linestyle=':', color='black',
                                                   linewidth=0.5)
        ax2.add_artist(patch)
        ax3.plot([cdf_p_l1] * 2,
                 [0, exp_plrv_p], linestyle='dashed', color='black', linewidth=0.8)

        print("_" * 50)

    ax1.set_xlim(left=995, right=1005)
    ax2.set_xlim(left=995, right=1005)
    ax3.set_xlim(left=0, right=1)
    ax1.set_ylim(bottom=0)
    ax3.set_ylim(bottom=0)
    ax1.sharex(ax2)
