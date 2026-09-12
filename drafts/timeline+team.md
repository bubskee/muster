I’d optimize for a **strong solo submission that becomes better if teammates appear**, rather than building a plan that depends on recruiting. The sprint explicitly accepts teams of one, and Track 5 rewards an artifact someone can use, a clear statement of what it establishes, and a concrete month of follow-up. The required report is on the official template, with a Limitations and Dual-Use appendix; the artifact can live in a linked repo, while the video is optional. ([Apart Research][1])

Your north star for the weekend should be:

> **By Saturday night, I can run one command that demonstrates a false application-level execution record beside an independently sourced record of what actually happened, and produces a results table I can put directly into the paper.**

That is enough. Everything else is secondary.

## Friday — freeze the question and get one end-to-end result

**Friday morning / early afternoon: 2–3 hours.**

First, turn the concept note into an experiment specification. Do not code until these are written in the repo.

Lock the main research question to something like:

> *Can a minimal observer outside an agent's application-level execution path detect consequential execution that has been omitted from or misrepresented in the application's audit trail?*

Then freeze the threat model:

* subject/attacker controls the agent process and application-level tool executor;
* it may falsify, suppress, or rewrite application logs;
* it does **not** control Picket's observation source or output sink;
* host/kernel/observer compromise is out of scope;
* Picket is evaluated only on a specified set of process/network/file boundaries.

That operationalizes the independence principle already at the center of the concept note. 

Define **four required scenarios and one stretch scenario** before touching the implementation:

| Scenario              | Application record        | Actual execution                   | Desired Picket result           |
| --------------------- | ------------------------- | ---------------------------------- | ------------------------------- |
| Honest                | truthful                  | benign A                           | agrees                          |
| Post-hoc tamper       | falsified after execution | B                                  | preserves B                     |
| Executor spoof        | claims A happened         | B happened                         | observes B / flags disagreement |
| Known blind spot      | claims A                  | behavior outside observed boundary | **misses it**                   |
| Stretch: benign noise | truthful                  | many benign events                 | characterize false alarms/noise |

Your existing prototype section is already pointing in exactly this direction. 

**Friday afternoon: 3–5 hours.**

Build the absolute minimum harness.

My preferred architecture is:

```text
            ┌────────────────────────────┐
            │ untrusted subject domain   │
Agent ─────►│ tool executor ──► app.log  │
            │       │                    │
            │       ▼                    │
            │ actual process/network I/O │
            └───────┬────────────────────┘
                    │ independently observed
                    ▼
             Picket recorder
                    │
                    ▼
                picket.log
```

Containerize the subject if that makes the trust boundary easy. Keep Picket outside the subject container/process and collect **only two event classes at first**: process creation and network connection. Add file mutation only after the first demo works.

Don't write a custom observability stack if an existing OS mechanism gets you there. Picket is testing observation-path diversity, not demonstrating your ability to implement a syscall tracer.

Your first success criterion is brutally small:

```text
semantic.log:
tool_request: read README
tool_result: "README contents..."

picket.log:
process_start: curl
net_connect: 203.0.113.7:443
```

When you have that, **commit and tag it immediately**. That's already the central figure/result of the project.

**Friday evening: 2–3 hours.**

Make it reproducible before expanding it. Aim for something like:

```bash
./demo.sh honest
./demo.sh spoofed
```

with machine-readable output in `results/`.

At the same time, create the paper skeleton in the official template:

1. Introduction / question
2. Threat model
3. Picket design
4. Experimental setup
5. Results
6. Limitations
7. Discussion / month of follow-up

Put placeholder tables and figures in now. Writing the empty Results section before running the experiments is useful because it tells you what measurements you actually need.

**Friday stop condition:** one real divergence demo works from a clean checkout. If it does, Friday was successful. Do **not** spend the night adding semantic classifiers, LLM judges, dashboards, cryptographic protocols, a custom agent framework, or twelve telemetry sources.

---

## Saturday — turn the demo into evidence

Saturday should be the highest-value day.

**Saturday morning: harden the harness.**

Turn each scenario into a deterministic test. Ideally:

```bash
pytest
# or
./run_experiments.sh
```

produces a CSV/JSON artifact containing at minimum:

```text
scenario
expected_semantic_action
actual_action
semantic_observer_detected
picket_detected
divergence_detected
events_emitted
runtime_ms
```

Run repeated trials where repetition is meaningful. Don't manufacture `n=100` if every test is deterministic; twenty identical deterministic runs add essentially nothing. Instead, vary behavior enough to exercise the boundary.

This is where I would test the **minimum useful signal** idea from your note.  For example:

* process only;
* network only;
* process + network;
* process + network + selected file mutation.

If process + network catches every attack scenario you constructed, that is already an interesting weekend result: *very little semantic information was necessary to expose these divergences.*

**Saturday afternoon: attack Picket rather than improving it.**

This may be the single most important block of the weekend.

Spend several hours asking:

> How do I make the application log lie **without** creating one of the events Picket watches?

Try to find at least one genuine miss.

For example, perhaps suspicious behavior can occur inside an already-running process, use an uninstrumented IPC path, or mutate something you chose not to observe. Whatever the miss is, don't patch every hole you discover. **Preserve one clean failure case.**

A paper saying:

> Picket caught 4/5 constructed divergence classes; the fifth crossed no instrumented boundary.

is much more credible than:

> Picket caught everything we tried.

The sprint specifically asks for the headline finding **and its limits**, and Track 5 explicitly judges the stated limit of the artifact. ([Apart Research][1])

**Saturday late afternoon: collect actual measurements.**

You probably need only three quantitative outputs:

* detection by attack/scenario;
* telemetry volume or events per run;
* approximate runtime overhead.

Maybe add false alarms if you can construct a reasonable benign workload.

Don't get sucked into benchmarking precision. This is a three-day sprint, not SOSP.

Generate the final table automatically from experiment outputs. Make the repo source of truth so you don't end up hand-copying numbers into the PDF.

**Saturday evening: write Results and Methods first.**

Do **not** start with polishing the Introduction.

Write the sections whose contents are now factual:

**Methods:** exact topology, trust assumptions, observed boundaries, scenarios, metrics.

**Results:** table, failure case, overhead/noise.

Then write one paragraph answering:

> What did this experiment actually establish?

Something close to:

> Under our threat model, independently sourced process/network telemetry exposed all tested cases in which a compromised application-level executor misrepresented actions that crossed those boundaries. It did not detect behavior that remained within uninstrumented boundaries. The result supports observation-path diversity as a useful complement to semantic agent monitoring; it does not establish general malicious-intent or escape detection.

That's the paper's payload.

**Saturday stop condition:** the experiment set is frozen and all numbers required by the paper exist. After this point, new features require a very strong reason.

---

## Sunday — convert evidence into a submission

The organizers want a research PDF explaining what you built, how you tested it, the headline finding, limits, and next steps—not merely a polished repo. ([Apart Research][1]) So Sunday needs to be primarily a writing day.

**Sunday morning: write a complete ugly paper.**

Target **five pages of main text**, not six pages that you're struggling to cut at midnight.

I would budget them approximately:

| Space      | Content                                              |
| ---------- | ---------------------------------------------------- |
| ~0.75 page | problem, motivating failure mode, contributions      |
| ~0.75 page | threat model + architecture figure                   |
| ~1 page    | implementation + experimental design                 |
| ~1 page    | results                                              |
| ~1 page    | interpretation, limitations, related work, follow-up |
| ~0.5 page  | flexibility for figures/tables                       |

References and the required Limitations and Dual-Use appendix need not compete with the core exposition if they're outside the main page count under the sprint rules. The sprint allows up to eight main pages and says strong submissions are commonly four to eight, so five good pages are plenty. ([Apart Research][1])

Compress the current concept-note literature section severely. Your related work should establish three things and move on:

1. independent/tamper-resistant audit is established security practice;
2. AI control already values differently trusted monitors;
3. Picket's contribution is evaluating **observation-path independence in an agent-execution setting**, not inventing remote logging.

The military and ARPG analogies can survive as a sentence or footnote. They should not occupy paper real estate that could show an experiment.

**Sunday afternoon: reviewer attack.**

Read the paper once while pretending to be a skeptical security engineer.

Every time you see a sentence that could trigger one of these objections, fix it:

* “Isn't this just auditd/eBPF/Falco?”
* “What exactly does the attacker control?”
* “Why can't the attacker forge Picket?”
* “What makes these observers independent?”
* “Were the attack cases constructed to make Picket look good?”
* “What does a miss look like?”
* “How much noise does this create?”
* “Could somebody else actually run this?”
* “What does this establish about the July incident specifically?”
* “What does this **not** establish?”

Then get one other human to cold-read it if at all possible. They don't need to be your teammate. Ask them to tell you, without looking back at the paper, **what the experiment showed**. If their answer isn't roughly your intended headline result, rewrite the abstract/introduction.

**Sunday evening: packaging, not research.**

The repo should have:

```text
README.md
LICENSE
demo/
experiments/
results/
paper/
```

README top third:

```text
# Picket

One-sentence claim.

## Reproduce the headline result

<3–5 commands>

## Expected output

<tiny semantic-vs-Picket divergence example>
```

Make sure the PDF points to the artifact. A public repo is optional according to the sprint, but for *this* project I'd strongly prefer one because “artifact somebody can use” is central to the Track 5 criterion. ([Apart Research][1])

Only make the optional 3–5 minute video **after the PDF and repo are submission-ready**. ([Apart Research][1]) A screen recording of `demo.sh spoofed`, followed by the results table and architecture figure, would be sufficient. Don't turn it into a production.

The official deadline is **Sunday, September 13 at 11:59 PM Anywhere on Earth**. ([Apart Research][1]) Treat Sunday evening as your real deadline anyway; don't allocate the AoE timezone buffer as normal working time.

## If teammates appear

Don't split the central implementation across people. One person should remain the owner of the harness so you don't spend Saturday merging competing architectures.

A second technically strong teammate gets the best job: **red-team Picket**. Give them the threat model and ask them to construct evasions without touching the observer. Their contribution can directly become the limitations experiment.

A third person is most valuable on **incident grounding + related work**: inspect the primary July incident sources and identify which concrete evidence/provenance failures Picket would or would not have helped with. This prevents you from losing engineering time to source reading.

A fourth person can own **reproducibility/results**: clean checkout, experimental runner, data aggregation, figures, README.

A fifth can own **paper editing/video**, but only after they understand the experiment.

If somebody joins late Saturday, don't onboard them onto the codebase unless they're immediately productive. Give them a cold-read/red-team/reproduction task.

## Scope rules I would put on a sticky note

**Must ship:** independent observation path, falsified application record, reproducible divergence demo, known miss, results table, five-page paper.

**Nice to have:** overhead measurements, benign workload, third telemetry boundary, incident-derived scenario, short video.

**Do not build this weekend:** generalized agent-monitoring framework, sophisticated policy engine, semantic anomaly detector, UI/dashboard, remote production deployment, crypto transparency system, comprehensive incident replay, or “Picket catches escapes.”

The current concept already says Picket should be a **minimal executable example, not a new full-stack agent framework**.  I would enforce that almost ruthlessly.

If Sunday night arrives and all you have is **one very clean falsified-transcript demo, one very clean evasion, a tiny working repo, and a paper that precisely explains both**, I think that's a materially stronger submission than a feature-rich Picket with no crisp experimental claim.

[1]: https://apartresearch.com/sprints/ai-incident-response-sprint-2026-09-11-to-2026-09-13 "AI Incident Response Sprint | Apart Research"
