import seaborn as sns
from prob_functions import pdf, cdf, plrv
import numpy as np
import matplotlib


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
