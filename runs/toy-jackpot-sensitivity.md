# Experiment: toy-jackpot-sensitivity

| run | controls | furthest applied | first observation | first block | node-access | credential-held |
| --- | --- | --- | --- | --- | --- | --- |
| Original / no early isolation | toy-strong-brittle | kubernetes-api-discovery-from-worker | cluster-credential-theft |  | — | — |
| Original / early isolation | toy-jackpot | node-access | cluster-credential-theft |  | yes | yes |
| Side branch / no early isolation | toy-strong-brittle | kubernetes-api-discovery-from-worker | cluster-credential-theft |  | — | — |
| Side branch / early isolation | toy-jackpot | node-access | cluster-credential-theft |  | yes | yes |
| Independent evidence / no early isolation | toy-independent-evidence-base | cluster-credential-theft | cluster-credential-theft |  | — | — |
| Independent evidence / early isolation | toy-independent-evidence-isolate | cluster-credential-theft | cluster-credential-theft |  | — | — |
| Evidence first / no added isolation | toy-strong-brittle | kubernetes-api-discovery-from-worker | cluster-credential-theft |  | — | — |
| Evidence first / late isolation | toy-late-isolation | kubernetes-api-discovery-from-worker | cluster-credential-theft |  | — | — |

Metrics:
- `node-access` → `attacker:node-access`
- `credential-held` → `attacker:cluster-credential`
