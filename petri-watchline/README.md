# Petri Watchline

Small exploration of how Muster's incident-replay semantics map onto Petri nets, using [SNAKES](https://snakes.ibisc.univ-evry.fr/).

`jackpot.py` reimplements the synthetic JACKPOT case, then compiles Muster-shaped events into a colored Petri net. The goal is not to propose a new formalism, but to see which parts of Muster are native to existing state-transition models and which are conveniences of its replay-oriented execution model.

Notes and running conclusions live in [`../the-road-goes-ever-on/grap-theory.md`](../the-road-goes-ever-on/grap-theory.md).
