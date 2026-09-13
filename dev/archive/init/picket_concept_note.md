# Picket

**Working concept note — AI Incident Response Sprint, September 2026**

## High-level statement

**Picket is a deliberately minimal, independently sourced tripwire layer for agent execution.**

The goal is **not** to detect every dangerous action, reconstruct an agent's true intent, or replace richer monitors. The goal is narrower:

> Make selected classes of worrying behavior difficult to perform silently, even when higher-level logging or monitoring is compromised.

Modern AI-control stacks increasingly rely on multiple monitors: chain-of-thought or scratchpad monitoring, tool-call transcripts, policy checks, anomaly detectors, human review, and sandboxing. But several of these layers can share the same failure mode if they consume the same compromised execution trace.

The July 2026 agent incidents sharpen this problem. If an agent can alter or replace part of the tool-execution path, a downstream logger may faithfully preserve a false account of what actually ran. A tamper-evident log can prove that a record was preserved; it cannot prove that the record corresponded to reality.

Picket proposes a complementary design principle:

> **Defense in depth should apply to the observation layer itself. At least one monitor should observe the system through a different, harder-to-corrupt window.**

The metaphor is a military picket: not the fortress, not the whole army, and not expected to defeat every attack. Its job is to make surprise harder and raise a reliable alarm when a narrow class of boundary-crossing behavior occurs.

## What Picket might look like

A useful Picket implementation should optimize for **simplicity, independence, and survivability**, not semantic richness.

Possible properties:

- **Append-only:** no update or delete API.
- **Externally sequenced:** timestamps / sequence numbers come from the recorder, not the agent.
- **Write-only from the subject domain:** the agent cannot query or rewrite its own audit stream.
- **Outside the primary trust domain:** recorder lives on the host, hypervisor, remote collector, or another independently protected component.
- **Always-on at consequential boundaries:** tool execution, process creation, network egress, filesystem mutation, privilege changes, etc.
- **Minimal schema:** fewer components for the agent to corrupt or spoof.
- **Coarse by design:** Picket does not need to understand why an action happened to notice that something worrying occurred.

An event stream might be intentionally boring:

```text
001 tool_request   shell("curl ...")
002 process_start  /usr/bin/curl
003 net_connect    203.0.113.7:443
004 bytes_out      1824
005 file_write     /tmp/bootstrap.sh
006 tool_result    "200 OK"
```

The interesting property is not that every line is semantically complete. It is that the lower-level events come from an observation path that does **not** depend on the same recorder the agent may have spoofed.

## Defense-in-depth framing

Picket is best understood as **one defender in a line of defenders**.

A safety stack might include:

| Layer | Strength | Typical failure mode |
|---|---|---|
| Chain-of-thought / scratchpad monitor | Rich semantic signal | hidden, omitted, or strategically manipulated reasoning |
| Tool-call transcript | High-level action semantics | compromised or replaced tool executor |
| Policy / anomaly monitor | Flexible interpretation | depends on upstream logs and model generalization |
| **Picket** | Independent, coarse behavioral signal | misses behavior outside observed boundaries |
| Human review | Contextual judgment | expensive, slow, dependent on available evidence |

The goal is not for every layer to be perfect. The goal is to make catastrophic failure require **several differently shaped defenses to fail at once**.

This is closer to military survivability engineering, classic cybersecurity defense-in-depth, or hardcore ARPG character design than to a single "alignment solution": stealth does not make armor redundant; armor does not make mobility redundant; a life pool does not make resistances redundant.

## Core research questions

1. **Observation independence:** At what layer can an agent first create a divergence between what actually happened and what the audit trail says happened?
2. **Minimum useful signal:** How little telemetry is needed to reliably raise a flare on a predefined class of worrying behaviors?
3. **Failure-mode diversity:** Does adding an independently sourced Picket reduce correlated monitor failure, or merely add another layer with the same blind spots?
4. **Operational cost:** What is the runtime, storage, engineering, and interpretability cost of moving observation outside the agent's primary execution domain?
5. **False alarms:** Can Picket remain cheap and broad without turning into an unusable stream of noise?

## Weekend-sized prototype

A sprint prototype should stay small.

Possible demo:

1. **Honest execution:** normal tool call; semantic log and Picket agree.
2. **Post-hoc log tampering:** semantic log is altered; Picket's append-only record survives.
3. **Recorder / executor spoofing:** higher-level transcript claims action A while lower-level behavior corresponds to action B; test whether Picket still raises a flare.
4. **Known miss:** demonstrate one class of behavior Picket does *not* catch.

The last case is important. Picket's claim should remain intentionally modest.

A useful evaluation table could report:

```text
attack / tampering mode
× observation layer
× detection rate
× false-positive rate
× runtime/storage cost
```

The repo should be a minimal executable example, not a new full-stack agent framework.

---

# Related AI-control / monitoring reading

## Sam Bowman — *Putting up Bumpers*

A strategic defense-in-depth argument for building multiple largely independent safeguards around increasingly capable AI systems.

- LessWrong: [Putting up Bumpers](https://www.lesswrong.com/posts/HXJXPjzWyS5aAoRCw)

**Why it matters for Picket:** Picket can be framed as one deliberately modest "bumper" whose value comes from having a different failure mode from neighboring safeguards.

## Buck Shlegeris — trusted + untrusted monitoring

Discussion of combining stronger but potentially compromised monitors with weaker trusted monitors because cheap, imperfect defenses can still eliminate important attack strategies.

- Alignment Forum: [How to prevent collusion when using untrusted models](https://www.alignmentforum.org/posts/GCqoks9eZDfpL8L3Q/how-to-prevent-collusion-when-using-untrusted-models-to)

**Why it matters for Picket:** a defensive layer does not need to dominate the threat model to justify its existence.

## UK AISI — *Loss of Oversight*

Work on how AI systems may become harder to audit and why independent validation channels with different failure causes matter.

- LessWrong mirror/discussion: [Loss of Oversight: How AI Systems May Become Harder to Audit](https://www.lesswrong.com/posts/JvZxp554WxcZ8BQvM/loss-of-oversight-how-ai-systems-may-become-harder-to-audit-1)

**Why it matters for Picket:** this is close to the central design principle—multiple monitors are not enough if they all fail for the same reason.

## FAR.AI / collaborators — *Layered AI Defenses Have Holes*

Empirical work showing that layered defenses should themselves be attacked as a composed system rather than assumed to multiply safety automatically.

- LessWrong: [Layered AI Defenses Have Holes](https://www.lesswrong.com/posts/Z8tAn3jp2uaQf34Mm/layered-ai-defenses-have-holes-vulnerabilities-and-key)

**Why it matters for Picket:** an important caveat. Defense in depth helps only when the layers genuinely add independent barriers or signals.

## DeepMind — system-level AI safety and security

A broader framework incorporating access control, monitoring, anomaly detection, sandboxing, logging, and auditing.

- LessWrong discussion: [DeepMind: An Approach to Technical AGI Safety and Security](https://www.lesswrong.com/posts/3ki4mt4BA6eTx56Tc/deepmind-an-approach-to-technical-agi-safety-and-security)

**Why it matters for Picket:** places the proposal inside a broader system-security tradition rather than treating monitoring as an isolated alignment technique.

## *Userland Alignment*

A short argument for treating the surrounding harness, environment, and system architecture as part of the alignment / control problem.

- LessWrong: [Userland Alignment](https://www.lesswrong.com/posts/W2ShJtS4Cvk8brZf6/userland-alignment)

**Why it matters for Picket:** philosophically adjacent to moving safety guarantees out of the model and into independently controlled infrastructure.

---

# Cross-domain / domain reading

## U.S. Army — Integrated Survivability Onion

The classic survivability model distributes protection across stages such as:

**Don't be there → Don't be seen → Don't be acquired → Don't be hit → Don't be penetrated → Don't be killed.**

- U.S. Army: [The 10X Tank Platoon](https://www.army.mil/article/287764/the_10x_tank_platoon)
- U.S. Army ARMOR bulletin (Fall 2025 PDF): [Fall 2025 ARMOR](https://www.benning.army.mil/armor/earmor/content/issues/2025/fall/Fall25%20eARMOR.pdf)

**Why it matters for Picket:** mature survivability engineering assumes different defenses operate at different stages of failure. No single layer is expected to solve survivability.

**Visual note:** the official Army "Integrated Survivability Onion" figure is a strong candidate to cite and adapt visually.

## Karen O'Brien — *Toward an Integrated T&E Framework for AI-Enabled Systems*

A 2026 paper applying survivability-style reasoning to AI-enabled systems and discussing forensic reconstruction / "black box" style recording.

- ITEA Journal: [Toward an Integrated T&E Framework for AI-Enabled Systems](https://itea.org/journals/volume-47-2/toward-an-integrated-t-and-e-framework/)
- DOI: [10.61278/itea.47.2.1009](https://doi.org/10.61278/itea.47.2.1009)

**Why it matters for Picket:** unusually direct bridge from survivability engineering to AI systems and persistent forensic evidence.

## NIST — Defense in Depth

NIST's security framing emphasizes multiple layers and heterogeneous countermeasures so that attacks missed by one control may be caught by another.

- NIST glossary: [Defense in Depth](https://csrc.nist.gov/glossary/term/defense_in_depth)
- NIST SP 800-172 overview: [Enhanced Security Requirements for Protecting CUI](https://csrc.nist.gov/news/2021/nist-publishes-sp-800-172)

**Why it matters for Picket:** authoritative source for the general strategy of layered controls with different protective roles.

## CISA — Industrial Control Systems defense in depth

Guidance for consequential systems where prevention, monitoring, segmentation, detection, and recovery must coexist.

- CISA: [ICS Recommended Practices](https://www.cisa.gov/resources-tools/resources/ics-recommended-practices)
- CISA PDF: [Improving Industrial Control System Cybersecurity with Defense-in-Depth Strategies](https://www.cisa.gov/sites/default/files/2023-01/NCCIC_ICS-CERT_Defense_in_Depth_2016_S508C.pdf)

**Why it matters for Picket:** useful analogy for systems where compromise is plausible and observability / recovery remain important after primary controls fail.

## NIST AU-9 — Protection of Audit Information

NIST audit controls include stronger variants that protect audit information using physically separate systems or different operating systems.

- NIST OSCAL control browser / AU-9 references: [NIST control material](https://pages.nist.gov/oscal-tools/demos/csx/baseline-reviewer/)

**Why it matters for Picket:** directly relevant to the provenance question: where does the recorder live, and can compromise of the audited system also compromise the audit trail?

---

# Fun analogy: hardcore ARPG survivability

## Path of Exile — *Comprehensive Guide to Defense / Defense Bible*

A community guide that divides survivability across several distinct mechanisms such as avoidance, mitigation, recovery, and total life pool.

- Path of Exile forums: [Comprehensive Guide to Defense in PoE — Defense Bible](https://www.pathofexile.com/forum/view-thread/3275670)

**Why it matters for Picket:** hardcore players optimize aggressively for speed and damage while still layering differently shaped defenses because any single mitigation mechanism has holes.

This is not scholarly evidence, but it is a memorable illustration of the same systems intuition:

> **You do not survive hardcore by maximizing one defense. You survive by making several independent things have to go wrong before the run ends.**

---

# One-sentence pitch

**Picket is defense in depth applied to the observation layer: a deliberately minimal, independently sourced tripwire that does not try to understand everything an agent does, but makes selected worrying actions harder to perform silently.**

# Working tagline

> **Not the wall. One more reliable defender in the wall.**
