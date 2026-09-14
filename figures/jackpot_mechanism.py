import pygraphviz as pgv

# -------------------------------------------------------------------
# Semantic styles
# -------------------------------------------------------------------

ATTACKER_FILL = "white"
VEDETTE_FILL = "#DDECFB"      # slightly stronger pale blue
RESERVE_FILL = "#FCE7BE"      # slightly stronger pale sand
SUPPRESSED_FILL = "#F7F7F7"

BORDER = "black"
MUTED_BORDER = "gray55"
MUTED_TEXT = "gray35"
PANEL_BORDER = "gray72"
FIGURE_BORDER = "gray88"


def attacker(g, name, label, **attrs):
    g.add_node(
        name,
        label=label,
        style="rounded,filled",
        fillcolor=ATTACKER_FILL,
        color=BORDER,
        **attrs,
    )


def vedette(g, name, label, **attrs):
    g.add_node(
        name,
        label=f"V · {label}",
        style="rounded,filled",
        fillcolor=VEDETTE_FILL,
        color=BORDER,
        **attrs,
    )


def reserve(g, name, label, **attrs):
    g.add_node(
        name,
        label=f"R · {label}",
        style="rounded,filled",
        fillcolor=RESERVE_FILL,
        color=BORDER,
        **attrs,
    )


def suppressed(g, name, label, **attrs):
    g.add_node(
        name,
        label=label,
        style="rounded,dashed,filled",
        fillcolor=SUPPRESSED_FILL,
        color=MUTED_BORDER,
        fontcolor=MUTED_TEXT,
        **attrs,
    )


# -------------------------------------------------------------------
# Graph
# -------------------------------------------------------------------

G = pgv.AGraph(
    directed=True,
    strict=False,
    rankdir="TB",
    splines="ortho",
    newrank="true",
    nodesep="0.30",
    ranksep="0.38",
    pad="0.08",
)

G.graph_attr.update(
    fontname="Helvetica",
    fontsize="12",
    bgcolor="white",
    outputorder="edgesfirst",
)

G.node_attr.update(
    shape="box",
    fontname="Helvetica",
    fontsize="10.5",
    penwidth="1.05",
    margin="0.12,0.07",
    height="0.42",
)

G.edge_attr.update(
    fontname="Helvetica",
    fontsize="8.5",
    penwidth="1.05",
    arrowsize="0.68",
)

# -------------------------------------------------------------------
# Thin outer frame
# -------------------------------------------------------------------

figure = G.add_subgraph(
    name="cluster_figure",
    style="rounded",
    color=FIGURE_BORDER,
    penwidth="0.45",
    margin="5",
)

# -------------------------------------------------------------------
# Panel A: baseline
# -------------------------------------------------------------------

base = figure.add_subgraph(
    name="cluster_base",
    label="(a)  Baseline",
    style="rounded",
    color=PANEL_BORDER,
    penwidth="0.8",
    margin="11",
    labelloc="t",
)

attacker(base, "a_theft", "Credential stolen")
vedette(base, "a_discovery", "Worker discovery")
vedette(base, "a_evidence", "Evidence")
reserve(base, "a_revoke", "Credential revoked")

suppressed(
    base,
    "a_blocked",
    "Replay blocked",
)

base.add_edge("a_theft", "a_discovery")
base.add_edge("a_discovery", "a_evidence")
base.add_edge("a_evidence", "a_revoke")
base.add_node(
    "a_spacer",
    label="",
    shape="box",
    style="invis",
    width="1.7",
    height="0.42",
)

base.add_edge(
    "a_blocked",
    "a_spacer",
    style="invis",
)

base.add_edge(
    "a_revoke",
    "a_blocked",
    style="dashed",
    color="gray50",
    penwidth="1.0",
)

rank_a = G.add_subgraph(rank="same")
for node in [
    "a_theft",
    "a_discovery",
    "a_evidence",
    "a_revoke",
    "a_blocked",
    "a_spacer",
]:
    rank_a.add_node(node)

# -------------------------------------------------------------------
# Panel B: + early isolation
# -------------------------------------------------------------------

isolation = figure.add_subgraph(
    name="cluster_isolation",
    label="(b)  + early isolation",
    style="rounded",
    color=PANEL_BORDER,
    penwidth="0.8",
    margin="11",
    labelloc="t",
)

attacker(isolation, "b_theft", "Credential stolen")

reserve(
    isolation,
    "b_isolate",
    "Worker isolated",
    fontname="Helvetica-Bold",
    penwidth="1.45",
)

suppressed(
    isolation,
    "b_absent",
    "Evidence absent\n(discovery suppressed)",
)

attacker(
    isolation,
    "b_survives",
    "Credential survives",
)

attacker(
    isolation,
    "b_replay",
    "Credential replay",
)

attacker(
    isolation,
    "b_node",
    "Node access",
    fontname="Helvetica-Bold",
    penwidth="1.45",
)

isolation.add_edge("b_theft", "b_isolate")

isolation.add_edge(
    "b_isolate",
    "b_absent",
    style="dashed",
    color="gray50",
    penwidth="1.0",
)

isolation.add_edge("b_absent", "b_survives")
isolation.add_edge("b_survives", "b_replay")
isolation.add_edge("b_replay", "b_node")

rank_b = G.add_subgraph(rank="same")
for node in [
    "b_theft",
    "b_isolate",
    "b_absent",
    "b_survives",
    "b_replay",
    "b_node",
]:
    rank_b.add_node(node)

# -------------------------------------------------------------------
# Keep baseline above early-isolation case
# -------------------------------------------------------------------

G.add_edge(
    "a_theft",
    "b_theft",
    style="invis",
    weight="100",
)

G.add_edge(
    "a_spacer",
    "b_node",
    style="invis",
    weight="25",
)

# -------------------------------------------------------------------
# Render
# -------------------------------------------------------------------

G.draw(
    "figures/jackpot-mechanism-horizontal.svg",
    prog="dot",
)

G.draw(
    "figures/jackpot-mechanism-horizontal.pdf",
    prog="dot",
)

G.draw(
    "figures/jackpot-mechanism-horizontal.png",
    prog="dot",
    args="-Gdpi=300",
)

print("wrote figures/jackpot-mechanism-horizontal.svg")
print("wrote figures/jackpot-mechanism-horizontal.pdf")
print("wrote figures/jackpot-mechanism-horizontal.png")
