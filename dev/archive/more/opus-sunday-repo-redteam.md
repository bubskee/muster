# PROMPT

https://github.com/bubskee/muster

check it out! what's the presentation missing?

# RESPONSE

I cloned it and ran it. The engine is in good shape: `go test ./...` passes clean, and both commands in the README reproduce the HF table exactly as printed. So this is a presentation review, not a code review.

## The thing that actually matters today

Submissions close tonight, 11:59 PM AoE, which is just before 5am Pacific tomorrow. And the deliverable is not the repo: it's a research report (PDF) on the official template, max 8 pages, with a required Limitations and Dual-Use Considerations appendix. A public repo and a 3–5 minute video are optional.

Right now the report doesn't exist. `writeup/draft/pickets-are-op.md` is a chat digest, not a paper, and there's no dual-use appendix anywhere in the repo. That's the gap. Everything below is second priority to converting that file into eight pages.

Also: pick one track. Your own draft already argues the Track 2 case (judged on resolvable questions, checks somebody could run tomorrow, and causal explanations that predict something), and the non-monotonicity result is exactly a causal explanation that predicts something about real IR playbooks. Spreading across 1, 2 and 5 reads as hedging.

## What the README is missing

**A finding.** This is the big one. The README explains the mechanism thoroughly and never states a result. The HF table shows that different controls block different things, which a reader will hear as close to tautological. Your actual contribution is in the draft: adding a working containment control can make the outcome strictly worse, because early isolation removes the access needed for the discovery event that would have generated the evidence to revoke the stolen credential. That belongs above the fold, with the two commands that show it:

```bash
go run . replay --scenario examples/notebook/toy-jackpot.yaml \
  --controls examples/notebook/toy-controls.yaml \
  --only-controls worker-credential-theft,worker-k8s-discovery,review-k8s-discovery,assess-cluster-integrity,escalate-cluster-integrity,revoke-cluster-credential
# terminal state has no attacker:node-access, credential revoked

# same, plus escalate-worker-containment,isolate-worker
# terminal state: attacker:node-access, attacker:cluster-session, credential survives
```

I ran both. They do what the draft says.

**A cheap way to see anything.** The README calls the notebook the easiest way in, but `muster_lab.ipynb` has no stored outputs and ipywidgets don't render on GitHub, so a judge who clicks it sees code and nothing else. The advertised path costs Go + Python + Jupyter + clone + Run All. Paste a real CLI transcript into the README and commit the notebook with outputs, or add one screenshot. Ten minutes of work, and it's the difference between "I read about it" and "I saw it."

**A version floor that's wrong.** `go.mod` declares `go 1.26.4`, but the README tells Ubuntu users `apt install golang-go`, which gives 1.22. I tried it: Go then tries to download an 80MB toolchain, and fails outright where that's blocked. I built and ran the full test suite on 1.22 with no source changes, so the directive is doing no work. `go mod edit -go=1.22 && go mod tidy` fixes it (and drops the stray `// indirect` on yaml.v3, which is a direct dependency).

**The Anthropic PyPI example.** It's in the repo, it's in the notebook walkthrough, and the README never mentions it. It's also the most on-theme example you have for this particular sprint, since it crosses evaluation operator, public registry, and third-party vendor. One paragraph.

**The model's semantics.** State is a set of facts, events apply when present/absent preconditions hold, effects add or remove facts, controls carry a role plus an action. None of that is written down outside the Go source, so a judge can't assess validity without reading `engine/`. Half a page covering the semantics and what the model deliberately doesn't represent (probability, timing, cost, attacker adaptation, partial success) does double duty as your required Limitations appendix.

**Positioning against attack graphs.** Several of the judges are detection engineers and Kubernetes security people. They will read this as an attack graph with control annotations within thirty seconds. The honest answer is that it is one, deliberately, and the contribution is the role decomposition plus the non-monotonicity. Say that yourself, in two sentences, or it lands as a hole.

## Context / open questions

- **Internal framing leaks into the public surface.** "Suggested co-author walkthrough" in the notebook, "Next increments after co-author feedback" in UX.md, and `dev/archive/gate*/` at top level. The gate process itself is genuinely unusual for a hackathon repo (frozen incident, pre-registered prediction, red-team notes). It's an asset if you name it as method in the report. It's clutter if a reviewer stumbles into it unlabeled.
- **The outcome vocabulary lives only in Python.** `UNDETECTED & UNCONTAINED`, `RESPONDED & CONTAINED`, `ADVERSE DESPITE RESPONSE` and the rest exist in the notebook, not in the Go engine or CLI output. Your writeup leans on those labels heavily, so a CLI-only reviewer can't see the vocabulary the paper is built on. Also slightly in tension with your own "Python never evaluates semantics" rule, depending on whether you count classification as display.
- **No extensibility path.** The YAML schema isn't documented and there's no "add your own incident" section. For an artifact that judges partly on whether somebody else can use it, that's a cheap win if you have time.
- Smaller: no authorship line, no CI (a `go test ./...` action is three minutes and signals the code runs), no repo topics, and `runs/` holds two gate2 JSONs but not the HF run the README table comes from.

If it helps, I can turn `pickets-are-op.md` into a structured 6–8 page draft with the limitations and dual-use appendix. That's the highest-leverage thing left before 5am.
