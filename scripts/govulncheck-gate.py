#!/usr/bin/env python3
"""Fail if govulncheck reports a CALLED vulnerability that is not in the allowlist.

    govulncheck -format json ./... > vuln.json
    python3 scripts/govulncheck-gate.py vuln.json .govulncheck-allowlist

govulncheck's own exit code cannot ignore a single advisory. This gate keeps the
check strict for every other finding, and makes each exception explicit, reviewed
and visible in the log. An allowlisted advisory that is no longer reported is
printed as a notice so the entry gets removed.

Allowlist format: one OSV id per line, a '#' starts a comment.
"""
import json
import sys


def stream(path):
    """govulncheck -format json prints a stream of concatenated JSON objects."""
    dec = json.JSONDecoder()
    text = open(path, encoding="utf-8").read()
    i = 0
    while i < len(text):
        while i < len(text) and text[i].isspace():
            i += 1
        if i >= len(text):
            break
        obj, i = dec.raw_decode(text, i)
        yield obj


def load_allowlist(path):
    allowed = {}
    for line in open(path, encoding="utf-8"):
        entry = line.split("#", 1)[0].strip()
        if entry:
            allowed[entry] = line.strip()
    return allowed


def main():
    if len(sys.argv) != 3:
        print(__doc__)
        return 2
    report, allow_path = sys.argv[1], sys.argv[2]
    allowed = load_allowlist(allow_path)

    summaries, called, seen = {}, set(), set()
    for msg in stream(report):
        if "osv" in msg:
            summaries[msg["osv"]["id"]] = msg["osv"].get("summary", "")
        finding = msg.get("finding")
        if finding:
            seen.add(finding["osv"])
            # A trace whose first frame has a function means the vulnerable code
            # is actually reachable from this module, not just present in a dependency.
            if finding.get("trace") and finding["trace"][0].get("function"):
                called.add(finding["osv"])

    blocking = sorted(called - set(allowed))
    for osv in sorted(called & set(allowed)):
        print(f"::notice::{osv} is allowlisted ({allowed[osv]}): {summaries.get(osv, '')}")
    for osv in sorted(set(allowed) - called):
        print(f"::notice::{osv} is in {allow_path} but is no longer reported as called; remove it.")
    for osv in blocking:
        print(f"::error::{osv} {summaries.get(osv, '')} https://pkg.go.dev/vuln/{osv}")

    if blocking:
        print(f"\n{len(blocking)} vulnerability(ies) reachable from this module and not allowlisted.")
        return 1
    print(f"govulncheck gate passed ({len(called)} called finding(s), {len(called & set(allowed))} allowlisted, {len(seen - called)} not called).")
    return 0


if __name__ == "__main__":
    sys.exit(main())
