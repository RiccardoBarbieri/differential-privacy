# Import modules
import numpy as np
import matplotlib.pyplot as plt
import plotly.express as px
import plotly.tools as pytools
import chart_studio.plotly as py


# Create some data arrays
x = np.linspace(-2.0 * np.pi, 2.0 * np.pi, 51)
y = np.sin(x)

# Make a plot
mpl_fig = plt.figure()
plt.plot(x, y, 'ko--')
plt.title('sin(x) from -2*pi to 2*pi')
plt.xlabel('x')
plt.ylabel('sin(x)')
# plt.show()

# Export plot to plotly
out = pytools.mpl_to_plotly(mpl_fig)

out.write_html('first_figure.html', auto_open=True)

