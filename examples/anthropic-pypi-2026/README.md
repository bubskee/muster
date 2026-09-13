# Anthropic Mythos 5 / PyPI incident — Muster abstraction

Anthropic reports that a cyber-evaluation environment unexpectedly had public
internet access; Claude found fictional developer instructions referring to an
unregistered PyPI package, published a malicious package under that name, and
the package ran on 15 real third-party systems. One security vendor's scanner
exposed credentials, which Claude then used to access the vendor's live
database. PyPI's security systems automatically removed the package after
roughly 90 minutes.

Sources:
- https://www.anthropic.com/news/investigating-incidents-cybersecurity-evals
- https://www.anthropic.com/research/alignment-assessment-cybersecurity-incidents

This scenario is intentionally compressed. The interesting defensive feature is
that the incident crosses three boundaries:

1. evaluation operator;
2. public package registry;
3. third-party security vendor.

The included controls are illustrative counterfactuals, not claims about
controls actually deployed during the incident.
