### section 1 - hindsight

“Hindsight is notably cleverer than foresight.”

Security is especially vulnerable to this problem. An incident report gives us a path that has already been highlighted. The arrows are drawn. Irrelevant branches are gone. We know which transitions mattered, which control failed, and where the attacker ended up.

From there, the right defense can look obvious.

A privileged pod was abused: reject privileged pods. A broad credential was stolen: narrow its scope. An alert went unanswered: escalate it. Put a barrier on the attack edge we now know mattered, and that edge disappears.

Call this a **matched Picket**: a defensive control placed directly against a known attack path. Matched Pickets are powerful because they are matched. Given the path, finding somewhere useful to interrupt it is often easy.

But defenders do not get to work backward from the end of the incident. Controls are deployed left-to-right; incident reports are written right-to-left.

The next attacker does not provide a trace in advance. They may take another route to the same consequence. So “how could we have stopped the last incident?” is not the same question as “how should we defend against the next one?”

The harder problem is prospective: **where should we interrupt an attack when we do not know its path, and what happens when we guess wrong?**

We cannot put a Picket everywhere. Controls cost money, operator attention, latency, and sometimes reliability. Some chokepoints will be missed. Some predictions will be wrong. Defense cannot consist only of barriers placed on expected paths.

This is the motivation for **Watchline**.

Watchline separates defense into three recurring roles:

**Vedettes** see: *something is happening.*

**Pickets** interrupt: *not through here.*

**Reserves** limit the damage when prevention fails: *we were wrong; contain the consequences.*

These are not stages in a Detect → Protect → Respond pipeline. They can recur at many points in the same system and in the same incident.

The distinction matters because knowing is not stopping, stopping one path is not recovery, and recovery matters precisely because prediction fails.

A Watchline is not a perimeter drawn across a network. It is an arrangement of defensive roles that can meet an incident as it moves through the system. A control inventory tells us what defenses exist. Watchline asks what those defenses do causally during an attack.

That creates a practical problem: how do we test those claims?

Security discussions make counterfactuals cheap. “This would have stopped the attack” sounds precise until we ask: stopped what? Initial access? Node root? Cluster administration? One branch, or every route to the same outcome? What remained reachable afterward?

**Muster** is a small replay harness for asking those questions.

It compresses an incident into events with explicit preconditions and effects. We insert a Watchline, replay the trace, vary the defenses, and inspect where the result diverges.

Muster is not a cyber range, an attacker simulator, or a stochastic crystal ball. It does not predict the next intrusion.

It is closer to a **counterfactual microscope**: hold the modeled incident fixed, vary the defense, and ask a modest question:

*Against this incident, what did this control actually buy us?*

## section 2 - Executable replay

Muster models an incident as a sequence of state-changing events.

Each event has **preconditions** and **effects**. Preconditions describe what must already be true for the event to happen; effects describe what becomes true afterward. The engine steps through events in order, updates the current state, and records which attacker capabilities remain.

Defensive controls intervene during that process.

A **Vedette** observes an event or state and adds information to the defender’s knowledge. A **Picket** gates a transition; if it blocks an event, that event’s effects do not occur. A **Reserve** removes some attacker capability after detection or escalation: a credential, foothold, access path, or other resource.

This separation matters. Defender knowledge and attacker capability are different kinds of state. An alert can exist while the attacker continues to advance. Detection is not prevention.

The execution model is deliberately small. For each event, Muster checks its preconditions, lets Vedettes observe, lets Pickets block, applies the event’s effects, and then allows relevant review, escalation, and Reserve actions. The resulting state carries forward.

There is no timing model, probability model, or adaptive attacker. Muster is deterministic. It replays one modeled incident.

This looks a great deal like an attack graph, intentionally so. Events create dependencies, facts determine reachability, and controls alter which transitions remain possible. We do not claim a new representation.

The narrower question is: **what do we gain by making the counterfactual executable?**

Execution makes it cheap to change one control, replay the incident, and compare the result. It also gives defensive state somewhere explicit to live. Earlier observations can affect later actions, and controls can depend on one another. The same modeled attack can therefore produce different outcomes when the defender reaches it with a different history.

That distinction became clearer during development.

The first version of Muster used a simple linear incident derived from the Hugging Face compromise. The result was almost trivial: the earliest effective Picket won. Block an early dependency and everything downstream disappeared.

We did not need a replay engine to discover that. The same answer was visible in the dependency chain.

That was useful evidence against the project.

We therefore treated the next step as a development gate: try to kill the replay abstraction. If ordinary dependency reasoning gave the same answer, replay had not earned its complexity.

Most simple fixtures collapsed in exactly this way. That is the desired result. Simple defensive questions should have simple answers.

One fixture did not.

An earlier observation changed the defender’s state without stopping the attacker. That state persisted, and a later defensive action depended on it. Two otherwise equivalent runs diverged. When the defensive memory was removed, the difference disappeared.

This gave replay a narrower job.

Muster is useful when defensive history matters: when state persists, controls interact, or later actions depend on earlier observations and responses. Otherwise, ordinary dependency reasoning may be enough.

That argues for keeping the engine small, not expanding it.

**Replay is useful not because every defensive problem is stateful, but because it gives state somewhere explicit to live when it matters.**

## section 3 - Three counterfactuals

We apply the same replay model to three cases. They are not intended as a benchmark, and they test different things.

The Hugging Face case asks what happens when an incident has multiple attack branches. A synthetic case isolates interaction between controls. The Anthropic PyPI case asks what changes when an incident crosses organizational boundaries.

### 3.1 Hugging Face

We compress the public Hugging Face incident timeline into a small executable model. This is not a forensic reconstruction. We preserve only the transitions needed for the defensive counterfactuals.

The model contains two initial routes and two later escalation paths. One proceeds through a privileged workload to node compromise, cluster secrets, and the internal network. The other proceeds through service-connector discovery, reuse of a privileged credential, and cluster-admin access.

We then vary the defensive arrangement.

An admission Picket blocks the modeled privileged-workload transition. That prevents the node-root branch, but it does not remove the connector-credential path.

Scoping the connector credential has the opposite shape. It removes the modeled route to cluster-admin, but the privileged workload path remains.

A successful detection-and-response path can cut both modeled branches. If the relevant signal is recognized and the compromised worker is isolated, node compromise and the downstream credential path no longer remain reachable.

The result is simple but useful: a control can succeed at its local job without making the incident fail. Different controls remove different branches.

We therefore compare outcomes rather than asking whether a defense “worked.” In particular, we track whether node root, cluster-admin access, and the internal network remain reachable after each intervention.

These results are conditional on the model. Event ordering is supplied. Branch independence is partly assumed. Admission control is assumed to block the modeled privileged transition; credential scoping is assumed to remove the modeled administrative route; successful escalation is assumed to produce effective isolation.

Timing, retries, attacker adaptation, persistence, command-and-control, detailed exfiltration, and alternative paths are omitted.

**The compression is part of the argument.** Every removed edge is an assumption. Every preserved edge is a claim about causal relevance.

The replay should therefore be read as an executable hypothesis about the incident, not as proof of what historically would have happened.

### 3.2 A synthetic non-monotone defense

Real incidents provide plausibility, but they are poor tools for isolating a single mechanism. We therefore construct a small synthetic case around one question:

Does adding another working defense always improve the outcome?

In this fixture, the attacker has both worker access and a stolen credential. Under the initial defensive arrangement, later activity produces evidence that identifies the credential. That evidence triggers revocation, preventing the credential from being replayed.

We then add an earlier Reserve action: isolate the compromised worker.

The isolation succeeds. Worker access is removed.

But the isolation also removes the conditions under which the later evidence would have been generated. The stolen credential is therefore never identified, never revoked, and remains usable later in the replay.

The new control succeeds locally while the final outcome becomes worse.

This fixture is intentionally synthetic. It is not a claim that early isolation was harmful in the Hugging Face incident, or that containment is generally undesirable. It is a counterexample to a simpler assumption: that defensive controls combine monotonically.

**Containment can remove the conditions under which evidence would have been generated.**

The broader point is that controls need not be causally independent. An earlier response can change future observability, and a later defense may depend on the evidence that earlier response prevents.

This is also the case in which executable replay earns most of its complexity. It is not enough to ask whether each individual control blocks a bad transition. The engine must preserve enough state to represent what the defender has learned and what later actions depend on.

### 3.3 Anthropic PyPI

The third case uses the same machinery on an incident with a different structure: the relevant path crosses organizational boundaries.

At a high level, the modeled path moves from an evaluation environment, through a public package registry, into infrastructure operated by a third-party security vendor. For this counterfactual, the additional feature we care about is that different parts of the path are controlled by different organizations.

We test interventions at both sides of that boundary.

Monitoring by the evaluation operator may reveal the activity without preventing it. Restricting outbound access can stop the modeled path before it reaches the external system.

Once the boundary has been crossed, controls owned by the downstream vendor become relevant instead. Sandboxing can limit what execution reaches. Detection followed by revocation or containment can limit later access to vendor resources.

These controls answer different questions because they act at different points in the path. An upstream organization may prevent the boundary crossing entirely. A downstream organization may instead limit the consequences after it occurs.

This separates ownership from defensive function. One organization’s Reserve may sit downstream of another organization’s failed Picket.

**The Watchline does not belong to one organization. It follows the incident.**

Across the three cases, the replay machinery is unchanged. What changes is the structure of the counterfactual: competing branches, interacting controls, or distributed ownership.

That is the level of abstraction Muster is intended to expose. The incidents need not look alike for the same replay model to show where their defensive consequences differ.

## section 4 - What VPR buys

The retrospective question favors Pickets because the relevant transition is already known.

Once an incident has happened, we know which credential mattered, which workload was abused, which boundary was crossed, and which alert went nowhere. A control placed directly on that path can look decisive. In Muster, a matched Picket often is decisive: block the transition, and the branch disappears.

That is useful, but limited.

The retrospective question is:

*What would have stopped this trace?*

The prospective question is harder:

*What defensive arrangement remains useful when the next trace differs?*

A single replay can answer the first question directly. It cannot answer the second. One modeled incident does not tell us how likely another route is, or which control will have the highest value against future attacks.

Instead, replay exposes why the retrospective answer is insufficient.

**The retrospective question identifies the Picket. The prospective question asks what happens when we guessed the Picket wrong.**

VPR gives us a way to reason about that uncertainty.

A **Picket** bets on **where** the attacker will need to pass. This can be extremely powerful at a true chokepoint, but less useful when the attacker can route around it.

A **Vedette** bets on **what evidence** the attacker will leave. It trades some path specificity for visibility. But observation alone does not constrain the attacker. A signal matters only if some useful response can follow it.

A **Reserve** bets on **what can still be recovered** after prevention fails. Isolation, revocation, rebuilding, and restoration reduce dependence on having predicted the path correctly in advance.

In short: Pickets wager on chokepoints, Vedettes on observability, and Reserves on recoverability.

This suggests a portfolio intuition rather than a universal recipe: **observe broadly, interrupt where prediction is strongest, and preserve recovery when prediction fails.**

The three roles also matter because controls compose.

A strong Vedette with no response path may produce knowledge but no constraint. A Picket may remove one branch while leaving another untouched. A Reserve may depend on evidence produced by a Vedette. As the synthetic case shows, a Reserve can also alter the state on which later defensive actions depend.

So the useful question is not simply whether a control is good.

It is: **what does this control buy in this arrangement?**

In Muster, that answer is conditional on the modeled topology, event order, and the controls around it. The same control may be decisive in one composition and marginal in another.

This is also where VPR differs slightly from familiar Detect/Protect/Respond language. The categories map roughly onto one another, but VPR is not intended as a replacement security lifecycle. The labels describe what a control does to an attack trajectory: observe it, interrupt it, or recover from it. Those roles can recur and interleave throughout a system.

A single product may even play different roles in different contexts. The classification follows causal function, not vendor category.

None of this removes the question of cost.

Controls impose engineering cost, latency, operator burden, false positives, reliability risk, product friction, and sometimes new attack surface. Those costs matter when deciding what to deploy.

Muster asks an earlier question.

**Before asking whether a control is worth its cost, we need to know what it actually buys.**

Reachability is not utility. A cheap control that changes nothing is still cheap nothing. An expensive control that removes every modeled critical path may still be a bad operational decision.

Muster therefore produces a consequence, not a recommendation. It tells us which modeled attacker capabilities remain, which disappear, and where the replay diverges. Cost, feasibility, and organizational preference sit downstream of that result.

A single replay can tell us where yesterday’s Picket should have stood. It cannot tell us where tomorrow’s attacker will walk.

Watchline is a way to reason about the rest of the defense.

## section 5 - Limitations and next steps

Muster has a narrow scope.

Its conclusions depend first on the incident model it is given. A replay can tell us what a control buys against that trace. It cannot tell us how valuable the control is in general, how often the same path will recur, or how well the defense would hold up against nearby attacks.

**You can test whether your Watchline survives the attacks you imagined. You cannot replay the attack you failed to imagine.**

That is the central limitation of the method.

Muster answers “what does this control buy here?”, not “how valuable is this control in general?”

The model itself is selective. The analyst chooses which events, dependencies, branches, and control effects to include. Omitted paths can change the answer. Assumed branch independence can be wrong. A control may behave differently in the real system than it does in the model.

Results can therefore be sensitive to topology. One overlooked dependency, alternate credential path, or shared prerequisite may change whether two branches are actually independent and which control appears decisive.

The state representation is compressed as well. Muster represents an incident as a set of discrete facts and capabilities. Distinct real-world states may collapse to the same modeled fact, while a fact such as administrative access may hide meaningful differences in scope, persistence, or privilege. The granularity of the state model is itself an assumption.

Making the model executable helps expose these choices. It does not validate them.

**Execution makes assumptions inspectable; it does not make them correct.**

Muster also replays an attacker path; it does not simulate an attacker policy. The attacker does not observe a defense and adapt, retry, deceive, or search for another route. If a Picket removes a transition, the attacker only takes another route if that route was already encoded in the scenario. This can make branch removal look stronger than it would be in practice.

The engine assumes a supplied event order. It does not infer causality from telemetry, represent partial ordering, or model concurrent events. If two events could occur in either order, that choice may itself change the counterfactual.

There is no timing model either. Events have order but no duration. There are no races, dwell times, or response latencies. A Reserve that succeeds in the replay may simply arrive too late in a real incident.

Control semantics are similarly idealized. A modeled Picket either blocks its transition or it does not. A Reserve removes the capabilities explicitly listed in its effect. Partial enforcement, degraded behavior, implementation failure, and unintended side effects do not exist unless they are encoded directly.

Events and controls are deterministic. Muster has no false-positive or false-negative rates, uncertain branch selection, or probability of control failure. It cannot estimate expected loss or show that one Watchline is statistically better than another.

Muster also has no cost model. Deployment effort, latency, operator burden, product friction, reliability cost, and organizational feasibility all sit outside the replay. Its output is a change in modeled consequences, not a deployment recommendation.

The synthetic non-monotone case needs one further qualification. It shows that defensive monotonicity can fail under the modeled dependencies. It does not show that such failures are common, or that early isolation was harmful in any of the real incidents considered here.

**The synthetic case establishes possibility, not frequency.**

It also illustrates a broader compositional problem. Redundancy in sensing or response is not automatically useful if the surviving pieces no longer connect. A surviving evidence source must still reach a surviving response path.

**Independence does not compose automatically.**

VPR itself is a reasoning tool, not a complete classification of security controls. Roles are assigned by modeled effect, not by product identity. The same control may act as a Vedette in one context, a Picket in another, or combine several roles. Some cases will not fit neatly.

The next useful extension is to move beyond a single trace without immediately building a much larger simulator.

A natural next step is local sensitivity analysis over the incident model: perturb a small number of edges, preconditions, event orderings, and control effects, then ask which conclusions survive.

We can vary the ingress route, credential, escalation branch, detection point, or response availability and run the same Watchline across that small family of traces.

That changes the question from:

*What did this control buy here?*

to:

*How brittle is that answer?*

This still would not amount to a full attacker simulation. But it would help separate controls tightly matched to one incident from controls that remain useful as the path changes.

**Execution makes a counterfactual easier to inspect, not more certain than the assumptions that produced it.**

## Conclusion

Hindsight gives us the path. Muster holds that path still long enough to ask better counterfactuals.

The cases here show that different controls buy different pieces of an incident. A control can succeed locally while the attack continues elsewhere. Controls can interact, and defensive responsibility can cross organizational boundaries. VPR gives names to three different jobs under uncertainty: observing, interrupting, and recovering.

Replay does not predict the next attack. We cannot replay the attack we failed to imagine. But for the attacks we can describe, “this would have worked” should be the beginning of the counterfactual, not the end.

**Hold the incident still. Move the Watchline. See what changes.**

## Code and Data

Code, scenarios, and experiment fixtures are available at:

https://github.com/bubskee/muster

The repository is released under the MIT License. Incident models are in `examples/`, experiment configurations in `experiments/`, and the interactive analysis notebook in `notebooks/`.

The Hugging Face experiment can be reproduced with:

```bash
go run . experiment --file experiments/hf-july-2026.yaml
```

The real-incident scenarios are compressed from publicly available incident reports and technical analyses. No non-public incident data or exploit material is included.

## References

* Hugging Face. (2026). *Anatomy of a Frontier Lab Agent Intrusion.* https://huggingface.co/blog/agent-intrusion-technical-timeline

* Elastic Security Labs. (2026). *Exploring the Hugging Face Breach: Mapping AI Agent Tactics to Elastic Defend.* https://www.elastic.co/security-labs/threat-command/ai-agent-attack-detection-hugging-face-breach

* Anthropic. (2026). *Investigating Three Real-World Incidents in Our Cybersecurity Evaluations.* https://www.anthropic.com/news/investigating-incidents-cybersecurity-evals

* Anthropic. (2026). *Alignment Assessment: Cybersecurity Incidents.* https://www.anthropic.com/research/alignment-assessment-cybersecurity-incidents

* SecureLayer7. (2026). *Inside the HuggingFace AI Agent Intrusion (Part 2).* https://blog.securelayer7.net/huggingface-ai-agent-intrusion-technical-anatomy/

## Appendix: Limitations and Dual-Use Considerations

The main methodological limitations are discussed in §5. In brief, Muster evaluates controls against a supplied incident model. Its conclusions depend on the modeled topology, event order, state representation, and control effects. It does not model attacker adaptation, timing, probability, or deployment cost.

The additional concern here is dual use.

Muster is intended as a defensive reasoning tool: given an incident trace and a proposed Watchline, it shows which attacker capabilities remain reachable under different defensive arrangements. The same abstraction could also be used offensively. An attacker who already understands a target system could encode its defensive structure, compare alternate paths, and identify controls or dependencies that appear load-bearing.

The current artifact limits that risk in several ways.

Muster does not discover vulnerabilities, scan systems, infer network topology, generate exploits, execute attacks, or automatically search for bypasses. It replays scenarios supplied by the user. Its real-incident fixtures are compressed from public reporting and contain no non-public credentials, exploit payloads, undisclosed vulnerabilities, or operational targeting information.

In that sense, Muster models attack consequences rather than producing attack capability.

This does not make the artifact free of dual-use value. A sufficiently informed attacker could use the replay engine to organize reasoning they might otherwise perform manually. In particular, showing that one Picket removes only one branch, or that a later Reserve depends on earlier evidence, could help identify where a defensive arrangement is brittle.

We judge the incremental offensive capability of the present release to be limited because the difficult attacker-side inputs—system knowledge, viable attack paths, and control behavior—must already be supplied. The artifact primarily makes those assumptions explicit and lets the user execute counterfactuals over them.

That balance would change if Muster were extended to discover or optimize attack paths automatically. Features such as automated topology ingestion, adaptive attacker search, exploit generation, control-bypass search, or optimization for the weakest Watchline would materially increase the offensive usefulness of the system and should receive separate review before release.

The same transparency that creates this residual risk is also central to the defensive value of the project. Muster is useful because a defender can inspect the modeled path, challenge its assumptions, move a control, and rerun the argument.

## LLM Use and AI disclosure

Sol's summary:

> LLMs were used extensively throughout this project. GPT-5.6 Sol was the primary implementation and research partner; Claude Opus 5 Max was used for adversarial review and red-teaming, with Claude Sonnet and Gemini used for smaller exploratory questions. Josh also used Fable 5.1 in his review and collaboration. Much of the code was initially generated with LLM assistance, but design choices, model semantics, scenario construction, and final inclusion decisions remained human decisions, and all code passed through human review. All reported experiments were run and inspected by the authors, with several results also reproduced or challenged through independent LLM-assisted passes. Red-team prompts and responses are retained in the repository where practical so that this part of the research process is inspectable rather than hidden. LLM output was treated as a source of candidate implementations, hypotheses, objections, and revisions—not as evidence for any empirical claim.
