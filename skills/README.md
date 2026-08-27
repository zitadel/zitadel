# Skills

Self-contained runbooks for manual verification of Zitadel features against a **local** test
instance. Each skill is a Markdown file describing what it verifies, what it needs and how to
run it, usually next to a script that does the work.

They are deliberately tool neutral: plain Markdown and plain shell, no assistant-specific
frontmatter or directory layout. Read one and follow it yourself, or hand it to whichever AI
assistant you use.

| Skill | Verifies |
| --- | --- |
| [test-actor-in-action-v2.md](test-actor-in-action-v2.md) | The impersonation actor is passed to Actions v2 execution targets and appears in the token claims. |
| [marvin.md](marvin.md) | Running a load-test sweep against the QA instance and publishing the results as versioned docs pages. |

## Conventions for new skills

Most skills run against a local instance. [marvin.md](marvin.md) is the exception: it targets
the shared QA benchmark instance, and states its own, stricter destructive-action rules in
place of the local-only rule below.


- **Local only.** A skill may create users, flip instance settings and delete things. Refuse to
  run against anything but `localhost`, and enforce it in the script rather than only saying it
  in the prose.
- **Ask first.** Print what will be created and require a confirmation, with a `--yes` escape
  hatch for non-interactive use.
- **Clean up.** Remove what you created, including on failure, and restore settings you changed.
  Offer a `--keep` flag for debugging.
- **Assume a clean slate.** Create the users, projects and clients you need instead of relying
  on fixtures that happen to exist on the author's machine.
- **Fail loudly and usefully.** Turn Zitadel's error envelopes into messages that say what to do
  next, and assert explicitly rather than leaving output for a human to eyeball.
