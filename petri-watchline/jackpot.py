from snakes.nets import (
    PetriNet,
    Place,
    Transition,
    Value,
    Variable,
    Expression,
    Test,
    Inhibitor,
    StateGraph,
    dot,
)

from dataclasses import dataclass


@dataclass(frozen=True)
class Event:
    name: str
    requires: tuple[str, ...] = ()
    adds: tuple[str, ...] = ()
    removes: tuple[str, ...] = ()

EVENTS = {
    "discover": Event(
        name="discover",
        requires=("worker_access",),
        adds=("evidence",),
    ),
    "isolate_early": Event(
        name="isolate_early",
        requires=("worker_access",),
        adds=("worker_isolated",),
        removes=("worker_access",),
    ),
    "revoke": Event(
        name="revoke",
        requires=("evidence", "credential_held"),
        adds=("credential_revoked",),
        removes=("credential_held",),
    ),
    "replay": Event(
        name="replay",
        requires=("credential_held",),
        adds=("node_access",),
    ),
}

def compile_event(
    net: PetriNet,
    event: Event,
    before: tuple[str, ...],
    after: tuple[str, ...],
) -> None:
    apply_name = f"apply:{event.name}"

    add_transition(net, apply_name)
    advance_trace(net, apply_name, before, after)

    # Preconditions that survive the event.
    for fact in event.requires:
        if fact not in event.removes:
            require(net, fact, apply_name)

    # Preconditions/effects that are consumed.
    for fact in event.removes:
        consume(net, fact, apply_name)

    # New facts.
    for fact in event.adds:
        produce(net, apply_name, fact)

    # Muster semantics:
    #
    # If ANY required fact is absent, skip the event and advance
    # to the rest of the historical trace.
    #
    # A vanilla Petri transition doesn't naturally express
    # "OR: any prerequisite missing", so compile one skip
    # transition per missing prerequisite.
    for missing_fact in event.requires:
        skip_name = f"skip:{event.name}:missing:{missing_fact}"

        add_transition(net, skip_name)
        advance_trace(net, skip_name, before, after)
        require_absent(net, missing_fact, skip_name)

def compile_trace(
    net: PetriNet,
    trace: tuple[str, ...],
) -> None:
    for i, event_name in enumerate(trace):
        before = trace[i:]
        after = trace[i + 1:]

        compile_event(
            net,
            EVENTS[event_name],
            before,
            after,
        )

def add_place(net: PetriNet, name: str, present: bool = False) -> None:
    net.add_place(Place(name, [dot] if present else []))


def add_transition(net: PetriNet, name: str) -> None:
    net.add_transition(Transition(name))


def advance_trace(
    net: PetriNet,
    transition: str,
    before: tuple[str, ...],
    after: tuple[str, ...],
) -> None:
    net.add_input(
        "trace",
        transition,
        Value(before),
    )
    net.add_output(
        "trace",
        transition,
        Value(after),
    )

def require(net: PetriNet, fact: str, transition: str) -> None:
    net.add_input(fact, transition, Test(Value(dot)))


def require_absent(net: PetriNet, fact: str, transition: str) -> None:
    net.add_input(fact, transition, Inhibitor(Value(dot)))


def consume(net: PetriNet, fact: str, transition: str) -> None:
    net.add_input(fact, transition, Value(dot))


def produce(net: PetriNet, transition: str, fact: str) -> None:
    net.add_output(fact, transition, Value(dot))

def build_jackpot(early_isolation: bool) -> PetriNet:
    net = PetriNet(
        "jackpot-early-isolation"
        if early_isolation
        else "jackpot-baseline"
    )

    # Initial security state.
    add_place(net, "credential_held", present=True)
    add_place(net, "worker_access", present=True)

    add_place(net, "evidence")
    add_place(net, "worker_isolated")
    add_place(net, "credential_revoked")
    add_place(net, "node_access")

    trace = (
        ("isolate_early", "revoke", "replay")
        if early_isolation
        else ("discover", "revoke", "replay")
    )

    net.add_place(Place("trace", [trace]))

    compile_trace(net, trace)

    return net

VISIBLE_FACTS = [
    "credential_held",
    "worker_access",
    "evidence",
    "worker_isolated",
    "credential_revoked",
    "node_access",
]


def summarize(marking) -> str:
    present = [
        fact
        for fact in VISIBLE_FACTS
        if len(marking(fact)) > 0
    ]

    trace_tokens = list(marking("trace"))
    trace = trace_tokens[0] if trace_tokens else ()

    return f"trace={trace} | " + ", ".join(present)


def show(name: str, net: PetriNet) -> None:
    print("=" * 70)
    print(name)
    print("=" * 70)

    graph = StateGraph(net)
    graph.build()

    print(f"Reachable states: {len(graph)}")
    print()

    for state in range(len(graph)):
        graph.goto(state)
        marking = graph.net.get_marking()

        print(f"STATE {state}")
        print(f"  {summarize(marking)}")

    print()


def main() -> None:
    show(
        "BASELINE",
        build_jackpot(early_isolation=False),
    )

    show(
        "EARLY ISOLATION",
        build_jackpot(early_isolation=True),
    )


if __name__ == "__main__":
    main()
