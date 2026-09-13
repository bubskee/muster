# FABLE

Notes for Bobby: where the remaining leverage is

Who's reading it. The judges most likely to draw this are detection and Kubernetes security engineers who have built the real versions of the controls in Table 1. That's an audience that will immediately understand the two-path result and the paged-isolation row, and will also notice anything that reads as jargon or filler. Write for them: plain verbs first (observe / block / respond), the role names second.

The report. The structure is already right: one question, three tables, one mechanism, honest limits. What matters now is that every sentence is yours. The guidelines are explicit that generated-reading reports don't get scored, and the strongest thing this project has going for it is a clear-eyed author who kept the tool tiny on purpose. Let that voice through. Under three pages is fine.

The repo's first screen. A judge will open the README before anything else. Make the top of it say, in your words: what Muster does, one command, Table 1. Everything metaphorical can move below the fold or go.

dev/archive/. It's visible and a judge will click it. A one-line README inside saying these are LLM-assisted working notes, plus matching honesty in the LLM statement, turns it from a liability into evidence of process. The gate story only helps if it's told as what actually happened; if the "preregistered" framing was added after the fact, describe the tests instead.

Cost. Leaving cost out of the engine is the right call, and the reason is good: it varies by organization and cost-benefit is a separate calculation. Say that in one sentence in the report, because several of the track's example projects ask for a cost column and a silent gap reads differently from a scoping decision. Framing: Muster answers "what does this control leave reachable"; CBA takes that table as an input.

The tautology objection. Someone will think "you built the model, so of course the table says what it says." The honest answer is already in the draft (compression choices are load-bearing; the timeline is the source). Keep that visible near Table 1 rather than only in Limitations.

Two small corrections to carry over from the audit: the ablation wording (isolate-worker alone reaches node access; the interaction needs the full early-containment path on top of revocation), and the PyPI metric (the package fact is removed by the final event, so use execution and vendor DB access).

If there's time left over. One command per table. The experiment files in the zip do that; without them, Tables 2 and 3 require hand-typed control lists a judge won't reproduce.
