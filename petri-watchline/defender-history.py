from dataclasses import dataclass

from snakes.nets import (
    PetriNet,
    Place,
    Transition,
    Value,
    Test,
    Inhibitor,
    StateGraph,
    dot,
)


@dataclass(frozen=True)
class Event:
    name: str
    requires: tuple[str, ...] = ()
    adds: tuple[str, ...] = ()
    removes: tuple[str, ...] = ()


def add_place(net: PetriNet, name: str, present: bool = False) -> None:
    net.add_place(Place(name, [dot] if present else []))


def add_transition(net: PetriNet, name: str) -> None:
    net.add_transition(Transition(name))


def require(net: PetriNet, fact: str, transition: str) -> None:
    net.add_input(fact, transition, Test(Value(dot)))


def require_absent(net: PetriNet, fact: str, transition: str) -> None:
    net.add_input(fact, transition, Inhibitor(Value(dot)))


def consume(net: PetriNet, fact: str, transition: str) -> None:
    net.add_input(fact, transition, Value(dot))


def produce(net: PetriNet, transition: str, fact: str) -> None:
    net.add_output(fact, transition, Value(dot))


def advance_trace(
    net: PetriNet,
    transition: str,
    before: tuple[str, ...],
    after: tuple[str, ...],
) -> None:
    net.add_input("trace", transition, Value(before))
    net.add_output("trace", transition, Value(after))


def compile_event(
    net: PetriNet,
    event: Event,
    before: tuple[str, ...],
    after: tuple[str, ...],
) -> None:
    apply_name = f"apply:{event.name}"

    add_transition(net, apply_name)
    advance_trace(net, apply_name, before, after)

    for fact in event.requires:
        if fact not in event.removes:
            require(net, fact, apply_name)

    for fact in event.removes:
        consume(net, fact, apply_name)

    for fact in event.adds:
        produce(net, apply_name, fact)

    for missing_fact in event.requires:
        skip_name = f"skip:{event.name}:missing:{missing_fact}"

        add_transition(net, skip_name)
        advance_trace(net, skip_name, before, after)
        require_absent(net, missing_fact, skip_name)


def compile_trace(
    net: PetriNet,
    trace: tuple[str, ...],
    events: dict[str, Event],
) -> None:
    for i, event_name in enumerate(trace):
        compile_event(
            net,
            events[event_name],
            trace[i:],
            trace[i + 1:],
        )

ATTACKER_FACTS = (
    "credential_held",
    "node_access",
)

DEFENDER_FACTS = (
    "alert_seen",
    "credential_revoked",
)


def build_history_case(persistent_history: bool) -> PetriNet:
    net = PetriNet(
        "persistent-history"
        if persistent_history
        else "current-state-only"
    )

    add_place(net, "credential_held", present=True)
    add_place(net, "node_access")

    add_place(net, "alert_seen")
    add_place(net, "credential_revoked")

    if persistent_history:
        observe = Event(
            name="observe",
            adds=("alert_seen",),
        )
    else:
        # Same historical observation occurs, but nothing persists
        # into defender state.
        observe = Event(
            name="observe",
        )

    events = {
        "observe": observe,

        "review": Event(
            name="review",
            requires=("alert_seen", "credential_held"),
            removes=("credential_held",),
            adds=("credential_revoked",),
        ),

        "replay": Event(
            name="replay",
            requires=("credential_held",),
            adds=("node_access",),
        ),
    }

    trace = ("observe", "review", "replay")
    net.add_place(Place("trace", [trace]))

    compile_trace(net, trace, events)

    return net


def facts_present(marking, names: tuple[str, ...]) -> list[str]:
    return [
        name
        for name in names
        if len(marking(name)) > 0
    ]


def summarize(marking) -> str:
    trace = list(marking("trace"))[0]

    attacker = facts_present(marking, ATTACKER_FACTS)
    defender = facts_present(marking, DEFENDER_FACTS)

    return (
        f"trace={trace}\n"
        f"    attacker: {', '.join(attacker) or '-'}\n"
        f"    defender: {', '.join(defender) or '-'}"
    )


def show(name: str, net: PetriNet) -> None:
    print("=" * 70)
    print(name)
    print("=" * 70)

    graph = StateGraph(net)
    graph.build()

    print(f"Reachable states: {len(graph)}\n")

    for state in range(len(graph)):
        graph.goto(state)
        marking = graph.net.get_marking()

        print(f"STATE {state}")
        print(f"  {summarize(marking)}")

    print()


def main() -> None:
    show(
        "PERSISTENT DEFENDER HISTORY",
        build_history_case(persistent_history=True),
    )

    show(
        "CURRENT-STATE-ONLY DEFENDER",
        build_history_case(persistent_history=False),
    )


if __name__ == "__main__":
    main()
