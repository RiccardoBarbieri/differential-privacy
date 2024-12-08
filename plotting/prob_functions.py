import numpy as np
from scipy.stats import laplace, norm


def plrv(values: np.ndarray | float, location1: float, location2: float, scale: float, distribution: str = "laplace"):
    if distribution == "laplace":
        return np.log(laplace.pdf(values, loc=location1, scale=scale) / laplace.pdf(values, loc=location2, scale=scale))
    elif distribution == "normal":
        return np.log(norm.pdf(values, loc=location1, scale=scale) / norm.pdf(values, loc=location2, scale=scale))


def pdf(values: np.ndarray | float, loc: float, scale: float, distribution: str = "laplace"):
    if distribution == "laplace":
        return laplace.pdf(values, loc=loc, scale=scale)
    elif distribution == "normal":
        return norm.pdf(values, loc=loc, scale=scale)


def cdf(values: np.ndarray | float, loc: float, scale: float, distribution: str = "laplace"):
    if distribution == "laplace":
        return laplace.cdf(values, loc=loc, scale=scale)
    elif distribution == "normal":
        return norm.cdf(values, loc=loc, scale=scale)