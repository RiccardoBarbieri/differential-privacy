import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import * as d3 from "d3";
import {loadPyodide} from "pyodide";


@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss'
})
export class AppComponent {
  title = 'dp_helper';

  protected readonly asd = asd;
}

async function asd() {

  const width = 640;
  const height = 400;
  const marginTop = 20;
  const marginRight = 20;
  const marginBottom = 30;
  const marginLeft = 40;

  let pyodide = await loadPyodide();
  let res = pyodide.runPython(`
  import sys
  sys.version
  `)

  console.log(res)

  // const x = d3.scale
}

// import numpy as np
// import pandas as pd
// def lower_bound(initial_suspicion: np.ndarray, epsilon: float):
// return initial_suspicion / (np.exp(epsilon) + (1 - np.exp(epsilon)) * initial_suspicion)
// def upper_bound(initial_suspicion: np.ndarray, epsilon: float):
// return (np.exp(epsilon) * initial_suspicion) / (1 + (np.exp(epsilon) - 1) * initial_suspicion)
//
// epsilons = np.arange(0, 8, 0.5)
// suspicions = np.linspace(0, 1, 1000)
//
// data = {'suspicions': suspicions}
// df = pd.DataFrame(data=data)
// for i, epsilon in enumerate(epsilons):
// df[f"lower_bound{i}"] = lower_bound(suspicions, epsilon)
// df[f"upper_bound{i}"] = upper_bound(suspicions, epsilon)
//
// df