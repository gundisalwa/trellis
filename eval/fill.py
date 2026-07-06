#!/usr/bin/env python3
"""Fill a prompt template's {{KEY}} placeholders with file contents.

  python3 eval/fill.py eval/prompts/worker.md TASK=eval/tasks/01-ambiguous-feature.md
  python3 eval/fill.py eval/prompts/reviewer.md TASK=... TRANSCRIPT=... RUBRIC=...
"""
import sys
import pathlib

tmpl = pathlib.Path(sys.argv[1]).read_text()
for kv in sys.argv[2:]:
    key, path = kv.split("=", 1)
    tmpl = tmpl.replace("{{" + key + "}}", pathlib.Path(path).read_text())
sys.stdout.write(tmpl)
