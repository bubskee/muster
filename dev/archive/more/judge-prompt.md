You are a skeptical, low-context reviewer evaluating a research submission for the Apart Research AI Safety Sprint.



Assume you know nothing about this project beyond the submission and any judging criteria/template I provide below. Do not fill gaps using charitable assumptions. Read it the way a reviewer would encounter it for the first time.



This is a final-pass review, not an invitation to redesign the project. The research is complete and the submission is nearly final. I want you to identify problems that are worth fixing before submission, not generate speculative new directions or rewrite everything in your preferred style.



Please evaluate it on five levels:



Cold-read comprehension

After one read, what do you think the project actually did?

What is the central contribution?

What are the main empirical results or demonstrated capabilities?

Where did you have to reread, infer, or guess?

Flag terminology that is introduced too quickly or distinctions that seem obvious to the authors but not to a new reader.

Judging-criteria assessment
Apply the provided Apart judging criteria literally.

Give a provisional score for each dimension.

Explain the strongest evidence supporting each score.

Explain what keeps it from receiving the next score up.

Distinguish weaknesses in the research itself from weaknesses in how the submission communicates it.

Technical and epistemic audit
Look especially for:

unsupported or overstated claims

conclusions stronger than the evidence

observations presented as interpretations, or vice versa

missing baselines, controls, parameters, or reproducibility details that materially affect the claims

unclear causal language

places where the system/demo appears to establish more than it actually establishes

contradictions between sections, figures, terminology, or stated scope

related-work claims that a knowledgeable reviewer might challenge

Do not manufacture objections merely to have objections. Tell me which concerns are genuinely consequential.

Writing/submission audit
Treat this as a short research report written under a severe time constraint.
Flag:

generic framing or padding

sentences that sound impressive but convey little

repetition

abrupt logical transitions

unnecessarily defensive passages

terminology drift

awkward or conspicuously LLM-ish prose

claims whose citation appears inadequate

figures/tables whose meaning is unclear without reading surrounding prose

deviations from the requested submission structure

Preserve distinctive authorial voice when it works. Do not flatten the prose into generic academic writing.

Final triage
Divide your findings into:

MUST FIX BEFORE SUBMITTING — issues that could materially affect comprehension, credibility, or score

SHOULD FIX IF EASY — worthwhile polish with good expected value

LEAVE IT ALONE — imperfections or stylistic quirks that are not worth touching this late



Then finish with:



Reviewer verdict: In 1–2 paragraphs, tell me how this submission is likely to land with a thoughtful but busy Apart reviewer. What will they remember? What are they likely to be excited by, skeptical of, or confused by?



Scorecard: Give your best estimate for each judging dimension and an overall qualitative assessment.



Top three edits: If I only change three things before submitting, name exactly which three.



Important constraints:



Judge the work that is actually presented. Do not assume access to the repository unless I provide it.

Do not reward ambition that was not executed.

Do not punish negative or mixed results merely for being negative; judge whether the conclusions appropriately follow from them.

Do not demand extra experiments unless their absence genuinely undermines a central claim.

Do not suggest wholesale rewrites unless the current structure fundamentally fails.

Be specific: quote or identify the exact sentence/section when possible.

Calibrate criticism. A final-pass reviewer who finds twenty tiny stylistic preferences is less useful than one who identifies three real weaknesses.

If something is already good enough, say so.



Here are the Apart judging criteria/template:



### Reviewer Guidance

- We **score along three different dimensions, from 1 (worst) to 5 (best)**.
- **Score dimensions independently**: a 5 on Execution isn’t always a 5 on Impact or Presentation
- For Impact & Innovation, focus on **potential impact on AI safety and security**. A quick search may help to assess novelty for areas you’re less familiar with.
- **Quality over quantity**: A focused 4-page submission often beats a noisy 10-page one.
    - Verbosity that obscures substance → subtract from **Presentation**
    - Scattered experiments that don’t add up → subtract from **Execution**
- **Calibration**:
    - **Level 3 = solid** hackathon work ← what you'd expect from a competent team in a weekend
    - **Level 5 = exceptional** ← we expect no more than ~5-10% of projects to reach this quality bar for any given scoring dimension

---

## Dimension 1: Impact Potential & Innovation

*How much would this matter for AI safety if it worked? How innovative is it?
For scores of 4-5: is this actually new to the field, or replicating recent work?*

| Score | Description |
| --- | --- |
| 1 | **Negligible.** No clear problem addressed, or no meaningful novelty. |
| 2 | **Limited.** Addresses a real problem but with a generic or well-trodden approach. Incremental at best. |
| 3 | **Moderate.** Clear problem with a reasonable approach; some novelty in framing or method beyond routine application of existing tools. |
| 4 | **Significant.** Important problem with an original approach, or identifies a neglected problem area. A valuable contribution others could build on. |
| 5 | **Exceptional.**  Tackles a critical AI safety problem with a genuinely novel approach, or opens a new research direction. Clear theory of change. You'd be excited to share this with researchers in the area. |

---

## Dimension 2: Execution Quality

*How sound are methodology, implementation, and findings?*

| Score | Description |
| --- | --- |
| 1 | **Seriously flawed.** Methodology broken, results uninterpretable, or implementation doesn't work. |
| 2 | **Weak.** Approach has significant gaps: missing validation, flawed experimental design, or incomplete implementation. |
| 3 | **Competent.** Technically solid given the short duration. Methodology makes sense, results are interpretable, limitations acknowledged, work builds toward clear conclusions. |
| 4 | **Strong.** Thorough methodology with convincing validation. Results clearly support conclusions. Immediately useful for future work. |
| 5 | **Exceptional.** Ambitious scope executed rigorously. Surprising findings, novel methods, or unusually robust validation. |

---

## Dimension 3: Presentation & Clarity

*How clearly are work, findings, and impact potential communicated?*

| Score | Description |
| --- | --- |
| 1 | **Incomprehensible.** Cannot determine what the project is actually claiming or doing. |
| 2 | **Hard to follow.** Key information buried, missing, or diluted by excessive length. Significant effort to extract main points. |
| 3 | **Clear enough.** Can understand the problem, approach, and results without undue effort. Core content clearly present: problem, method, findings, limitations. |
| 4 | **Well presented.** Easy to follow, well-structured, appropriate level of detail. Target audience would get it quickly. |
| 5 | **Exceptionally clear.** A pleasure to read. Complex ideas made accessible. Could serve as a model for how to present this type of work. |

---

## Flags

*(checkbox, not scored)*

☐ **Off-topic**: Project does not fit any track for the current hackathon.

- Off-topic projects are still judged on their merits (can still achieve high scores) but may be deprioritized for feedback and thorough review
- Good off-topic work may be recognized (e.g., studio invitations) but are excluded from hackathon prizes

Judging Criteria
Dimension 1: Impact Potential & Innovation

How much would this matter for AI safety if it worked? How innovative is it?
For scores of 4-5: is this actually new to the field, or replicating recent work?

Score
	

Description

1
	

Negligible. No clear problem addressed, or no meaningful novelty.

2
	

Limited. Addresses a real problem but with a generic or well-trodden approach. Incremental at best.

3
	

Moderate. Clear problem with a reasonable approach; some novelty in framing or method beyond routine application of existing tools.

4
	

Significant. Important problem with an original approach, or identifies a neglected problem area. A valuable contribution others could build on.

5
	

Exceptional. Tackles a critical AI safety problem with a genuinely novel approach, or opens a new research direction. Clear theory of change. You'd be excited to share this with researchers in the area.
Dimension 2: Execution Quality

How sound are methodology, implementation, and findings?

Score
	

Description

1
	

Seriously flawed. Methodology broken, results uninterpretable, or implementation doesn't work.

2
	

Weak. Approach has significant gaps: missing validation, flawed experimental design, or incomplete implementation.

3
	

Competent. Technically solid given the short duration. Methodology makes sense, results are interpretable, limitations acknowledged, work builds toward clear conclusions.

4
	

Strong. Thorough methodology with convincing validation. Results clearly support conclusions. Immediately useful for future work.

5
	

Exceptional. Ambitious scope executed rigorously. Surprising findings, novel methods, or unusually robust validation.
Dimension 3: Presentation & Clarity

How clearly are work, findings, and impact potential communicated?

Score
	

Description

1
	

Incomprehensible. Cannot determine what the project is actually claiming or doing.

2
	

Hard to follow. Key information buried, missing, or diluted by excessive length. Significant effort to extract main points.

3
	

Clear enough. Can understand the problem, approach, and results without undue effort. Core content clearly present: problem, method, findings, limitations.

4
	

Well presented. Easy to follow, well-structured, appropriate level of detail. Target audience would get it quickly.

5
	

Exceptionally clear. A pleasure to read. Complex ideas made accessible. Could serve as a model for how to present this type of work.
Submission Requirements

Required:

    Research report (PDF) using the official template.

    Project title and abstract, 150 words or fewer.

    Author names and affiliations.

    A "Limitations and Dual-Use Considerations" appendix (required, see below).

Optional:

    Public GitHub repo, subject to the disclosure review below. Do not publicly release novel installation recipes without review.

    A 3 to 5 minute video demo.

Submission Template => Link
Recommended Report Structure

Maximum 8 pages, not counting references and appendices. Most strong projects are 4 to 8 pages.

    Introduction: which track and sub-problem, why it matters, and what the artifact is for.

    Related Work: what you build on.

    Methodology: enough to replicate, with sources and assumptions stated.

    Results: quantitative where possible, with the main threat to validity stated.

    Discussion: implications, limitations, future work.

    Limitations & Dual-Use Considerations (required).

    References.

AI tools and your report

Use AI tools the way you would use a colleague: to check your reasoning, find gaps in a draft, or debug code. The report itself has to be your team's own writing about your team's own work. Judges read every submission, and a report that reads as generated rather than written (generic framing, padded sections, claims without sources, no trace of what you actually did) will not be scored. Keep it short, say what you did in your own words, and link the sources for every factual claim.
Publishing your work

We encourage teams to publish their submission, and LessWrong is a natural venue for most of the written ones. Write-ups from this sprint will hold value if they follow a few rules: state your epistemic status and don’t use LLMs for writing on Lesswrong, only use LLMs to find problems in your drafts, not to draft it; link the primary sources for every factual claim about the incident; pick a title that states the finding rather than the topic; and publish the imperfect version this month rather than the polished one in three. We will link the best posts from the sprint page. Maximum of 1500 words for written contributions, without counting appendixes. The quality and value of what you write is worth much more than the length.



Here is the submission:

see attached pdf
