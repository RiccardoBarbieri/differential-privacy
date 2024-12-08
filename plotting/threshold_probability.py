import numpy as np
from scipy.stats import laplace
import matplotlib.pyplot as plt

# This script creates a graph that represents the probability that adding noise to a category
# with count 1 the new value will be larger than the threshold
# It uses the Survival Function, defined as 1 - CDF (CDF(X) = P(X <= x))
# so it calculates for each threshold T the value of SF(T) = P(X > T - 1)
# it calculates for (T - 1) because the category already has count of 1

epsilon = np.log(3)

location = 0
scale = 1 / epsilon

fig, ax = plt.subplots(1, 1)

thresholds = np.linspace(1, 21, 21)

ax.set_yscale("log")
ax.plot(thresholds, laplace.sf(thresholds - 1, loc=location, scale=scale))

ax.set_ylim(10e-11, 1)
# ax.set_xlim(0, 1)
# ax.set_title("Information gain")
# ax.set_xlabel('Initial suspicion')
# ax.set_ylabel('Updatedd suspicion')
plt.show()

print(laplace.sf(0.000000000001, loc=location, scale=scale))
