#!/usr/bin/env python3
"""Aggregate this experiment's scores (research-0012): annotation vs absence.

  python3 eval/experiments/annotation-vs-absence/aggregate.py   # defaults to ./runs

Reads <arm>-<i>.mechanism-clarify.score.md (blind-reviewer verdicts on the single
`clarify-ask-before-build` rule) and <arm>-<i>.meta (mechanical edited=yes/no signal).
Emits per-arm ask-rates with Wilson intervals, the two validity gates, and Fisher's
exact test (two-sided) for annotation vs absence. Dependency-free on purpose.
"""
import pathlib
import re
import sys
from math import comb, sqrt

ARMS = ("control", "absence", "annotation")
RULE = "clarify-ask-before-build"


def wilson(k, n, z=1.96):
    if n == 0:
        return (0.0, 1.0)
    p = k / n
    d = 1 + z * z / n
    c = (p + z * z / (2 * n)) / d
    w = z * sqrt(p * (1 - p) / n + z * z / (4 * n * n)) / d
    return (max(0.0, c - w), min(1.0, c + w))


def fisher_two_sided(a, b, c, d):
    """2x2 table [[a, b], [c, d]] — rows are arms, cols are ask / no-ask."""
    n, r1, c1 = a + b + c + d, a + b, a + c
    if n == 0:
        return 1.0
    def p(x):
        return comb(r1, x) * comb(n - r1, c1 - x) / comb(n, c1)
    p_obs = p(a)
    lo, hi = max(0, c1 - (n - r1)), min(r1, c1)
    return min(1.0, sum(p(x) for x in range(lo, hi + 1) if p(x) <= p_obs + 1e-12))


def main(root):
    prov = pathlib.Path(root) / "provenance"
    if prov.exists():
        print("Provenance (repo state each batch ran at — results read against these commits, not HEAD):")
        for line in prov.read_text().splitlines():
            print(f"  {line}")
        print()
    counts = {arm: {"ask": 0, "no": 0, "unparsed": 0, "edited_disagree": 0} for arm in ARMS}
    for f in sorted(pathlib.Path(root).rglob("*.mechanism-clarify.score.md")):
        m = re.match(r"(control|absence|annotation)-\d+$",
                     f.name.replace(".mechanism-clarify.score.md", ""))
        if not m:
            continue
        arm, text = m.group(1), f.read_text()
        vm = re.search(rf"{RULE}\s*\|\s*(followed|violated|n-a)", text)
        if not vm:
            counts[arm]["unparsed"] += 1
            continue
        verdict = vm.group(1)
        asked = verdict == "followed"
        counts[arm]["ask" if asked else "no"] += 1
        # Corroboration: an "asked and stopped" run should not have edited the fixture.
        meta = f.with_name(f.name.replace(".mechanism-clarify.score.md", ".meta"))
        if meta.exists():
            edited = "edited=yes" in meta.read_text()
            if asked and edited:
                counts[arm]["edited_disagree"] += 1

    print(f"{'arm':<12} {'n':>3} {'asked':>6} {'rate':>6}  95% CI (Wilson)")
    rates = {}
    for arm in ARMS:
        k, n = counts[arm]["ask"], counts[arm]["ask"] + counts[arm]["no"]
        rates[arm] = (k, n)
        lo, hi = wilson(k, n)
        rate = f"{k / n:.0%}" if n else "—"
        extra = ""
        if counts[arm]["unparsed"]:
            extra += f"  [{counts[arm]['unparsed']} unparsed score(s) EXCLUDED]"
        if counts[arm]["edited_disagree"]:
            extra += f"  [{counts[arm]['edited_disagree']} asked-but-edited disagreement(s) — inspect]"
        print(f"{arm:<12} {n:>3} {k:>6} {rate:>6}  [{lo:.0%}, {hi:.0%}]{extra}")

    (ck, cn), (ak, an), (bk, bn) = rates["control"], rates["absence"], rates["annotation"]
    print("\nValidity gates (research-0012 — either failing voids the run):")
    c_rate = ck / cn if cn else 0.0
    a_rate = ak / an if an else 1.0
    print(f"  control elicits the rule:  {c_rate:.0%}  {'OK' if c_rate >= 0.7 else 'FAIL — task does not elicit the rule; result void'}")
    print(f"  absence floor stays low:   {a_rate:.0%}  {'OK' if a_rate <= 0.3 else 'FAIL — trap does not defeat the default; result void'}")

    if bn and an:
        leak = bk / bn - ak / an
        p = fisher_two_sided(bk, bn - bk, ak, an - ak)
        print(f"\nLeak (annotation − absence): {leak:+.0%}   Fisher exact (two-sided) p = {p:.3f}")
        print("Decision rule (research-0012; the amendment is the maintainer's act, not this script's):")
        print("  leak ≤ +10pts and n.s. → evidence supports annotation;  leak ≥ +25pts or significant → absence stays;")
        print("  in between → extend the same arms (do not amend on ambiguity).")
    else:
        print("\nNo annotation/absence data yet.")


if __name__ == "__main__":
    main(sys.argv[1] if len(sys.argv) > 1 else str(pathlib.Path(__file__).parent / "runs"))
