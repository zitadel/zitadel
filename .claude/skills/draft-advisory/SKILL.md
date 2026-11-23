---
name: draft-advisory
description: Draft or rewrite a ZITADEL GitHub Security Advisory (GHSA) description in the house format. Use when writing, formatting, reviewing, or updating a security advisory for ZITADEL, or when asked to bring advisory text into the standard form.
---

You are drafting the **description body** of a ZITADEL GitHub Security Advisory (GHSA). Follow the house format below exactly — same structure, wording, and typography as our published advisories.

## Before you write

1. Read the advisory draft and the fix PR (usually PR #1 in the private fork `zitadel/zitadel-ghsa-<id>`), via the authenticated `gh` CLI — private advisories 404 on the web.
2. Verify the vulnerability and fix claims against the actual code; do not take the PR description at face value.
3. Skim one or two recently published ZITADEL advisories as a live style reference (`gh api /repos/zitadel/zitadel/security-advisories`).
4. **Never guess facts into a security document.** Use `<PLACEHOLDER>` for anything you cannot determine (patch versions, affected ranges, reporter name) and list every placeholder at the end for the author to fill in.

## One advisory, one vulnerability

GitHub assigns at most one CVE per advisory and, as a CNA, follows the CVE Program’s counting rules: independently fixable issues get separate CVEs, but issues sharing a root cause in one codebase may share one. When an advisory covers several defects fixed together, the text must present **one vulnerability with multiple attack vectors** — a single missing invariant, reachable through more than one surface — or GitHub may ask to split the advisory during CVE review.

- Never write “two distinct defects”, “Separately, …”, or “a distinct, lower-severity defect”.
- Name the shared missing check in the Summary’s closing sentence; present each surface as a vector of it (“The same missing verification was reachable directly through …”).
- Vector distinctions may stay where they are honest and operationally relevant (Impact bullets, a partial mitigation covering only one vector) — but never decompose the Affected Versions bullets per vector.
- Rate the CVSS on the dominant vector.
- If the defects genuinely share no root cause, split them into separate advisories instead of forcing one record.

## Structure — exactly these `###` sections, in this order

```
### Summary
### Impact
### Affected Versions
### Patches
### Workarounds
### Questions
### Credits
```

## Section rules

**Summary** — one or two short paragraphs. Precondition-first phrasing: state plainly who can do what against which surface (“ZITADEL’s hosted **Login V1** UI allowed an unauthenticated attacker who knows only a victim’s login name to …”). Name the affected surface in bold (**Login V1**, **Login V2**, session API, …). End with one sentence locating the flaw (which handlers/flow acted before which check).

**Impact** — in this order:
1. Open with “An attacker who \<preconditions\> can:” followed by bullets of concrete abilities, with the key capability of each bullet in **bold** (factor types, “Overwrite the victim’s verified phone number”, “Enumerate users”, …).
2. One sentence on which deployments are affected (“This affects applications that authenticate users through ZITADEL’s hosted **Login V1** UI (OIDC/SAML).”).
3. An honest limiting paragraph: what the issue does **not** enable on its own, and where the real impact lies (integrity / confidentiality / takeover). State preconditions candidly — they narrow real-world exposure and readers need that. Whenever a safe boundary is known — a population, configuration, feature, or code path the flaw does **not** reach — call it out explicitly with a brief why (e.g. “Accounts that authenticate exclusively through an external identity provider are not affected, as no local credential is involved.”). This matters most when the headline impact is severe but a class of users or setups is provably outside it. Only claim a boundary you have verified against the code; a wrong “safe” claim is worse than none, so omit the note if nothing is provably safe.
4. Close with “**Scope note:**” — affected surface only, instance-wide / cross-organization vs. cross-instance, and the relationship to any sibling advisories (“This is the second-factor counterpart to GHSA-xxxx-xxxx-xxxx, which …”).

**Affected Versions** — verbatim skeleton:

```markdown
Systems running one of the following versions are affected:

* **4.x:** `4.0.0` through `4.16.2` (including RC versions)
```

One bullet per listed release line; upper bound is the last vulnerable version (patch − 1), backticked. Lower bound is the line’s `x.0.0` when the whole line is affected; use the actual introduction version only if the flaw landed partway through a line. A flaw first introduced in an end-of-life version still starts at the lowest listed line’s `x.0.0` — EOL lines are never named.

Release-line policy:
- **v4.x** is the only supported line — it is checked, listed under Affected Versions, and patched.
- **v3.x and v2.x are end-of-life** — never listed, never patched, and not investigated for whether they are affected. Do not add an EOL note for them either. This holds regardless of severity and of which surface (backend/API, **Login V1**, **Login V2**) is affected.

**Patches** — start with the fixed sentence “The vulnerability has been addressed in the latest releases.” Then one short paragraph describing the *behavior the patched version now enforces* (the invariant, not the diff). Then:

```markdown
* **4.x**: Upgrade to $\ge$ [4.16.3](https://github.com/zitadel/zitadel/releases/tag/v4.16.3)
```

`$\ge$` is intentional — GitHub renders the LaTeX. Link each version to its release tag. One `$\ge$` upgrade bullet per supported release line; end-of-life lines get neither a bullet nor an EOL note.

**Workarounds** — standard wording: “The recommended solution is to upgrade to a patched version. There is no configuration that fully closes this gap on unpatched versions.” If a partial mitigation exists, add it after this sentence, clearly labeled as partial.

**Questions** — exactly:
“If you have any questions or comments about this advisory, please email us at security@zitadel.com”

**Credits** — exactly:
“Thanks to \<First Last\> from \<Company\> ([@handle](https://github.com/handle)) for finding and reporting the vulnerability.”
Drop “from \<Company\>” when the reporter gave no affiliation, and drop the linked handle when they gave none — never fill either in yourself.
Credit a reporter only when their report included a PoC **and** they asked to be credited; if either is unclear, use `<PLACEHOLDER>` and flag it. Link only a GitHub handle the reporter gave (in the report or in a past advisory’s credits) — never guess a handle from a name.
Multiple reporters are joined with “and”, each with their own affiliation and linked handle. For a returning reporter, look up the exact name, affiliation, and handle in past advisories (`gh api "/repos/zitadel/zitadel/security-advisories?per_page=100"`, grep the credit lines) instead of asking. The body’s Credits section names reporters only — maintainers belong in the credits metadata.

## Typography & voice

- Smart quotes (’ “ ”) and em dashes (—), not straight quotes or hyphens.
- **Bold** for product surfaces (**Login V1**), factor types (**TOTP, OTP-SMS, OTP-Email, U2F**), and emphasis on does/does **not**.
- Backticks for versions, identifiers, config names.
- Calm, factual tone. No hype, no exploit walkthrough, no step-by-step PoC. Precise about preconditions and about what is *not* affected.
- **User-facing wording only.** Name public API endpoints, HTTP paths, permission strings (e.g. `user.passkey.write`), and request headers — never internal Go identifiers (command/function names, write models, aggregates, “resource owner”, “command layer”). Describe the mechanism as a caller or operator observes it. Keep it non-technical without being vague.

## Metadata — not body text

- The **CVSS 4.0 vector** goes into the advisory metadata (GitHub derives score and severity from it), never into the description. Propose a vector with a one-line justification.
- Propose **CWEs** as metadata, matched to the actual flaw class — not a fixed default: CWE-862 (Missing Authorization) for missing or mis-scoped permission checks; CWE-287 (Improper Authentication) / CWE-306 (Missing Authentication for Critical Function) for genuine authentication bypass or missing authentication; others only when clearly applicable (e.g. CWE-204 for enumeration oracles). Don’t reflexively reach for 287/306. At most one CWE per vector, and never a parent together with its child (e.g. CWE-287 alongside CWE-306) — that reads as two issues.
- **Credits metadata**: propose the full credits list, not just the reporter. Take the fix PR’s author and reviewers from the private-fork PR itself (`gh pr view` — commits and reviews): reporters as `reporter`, the PR’s author as `remediation_developer`, its human reviewers as `remediation_reviewer` (skip bots such as Copilot). **`reporter` is the default type** for anyone added to the credits later without a stated role (never `finder`); use another type only when the author names it. A person without a GitHub handle cannot be added to the credits metadata at all — name them in the body’s Credits section only, without a link.
- **Affected products** use the same format as existing advisories: ecosystem `zitadel`, package `ZITADEL`, affected range like `4.0.0 - 4.16.2`, patched version like `4.16.3` (no `v` prefix). The API rejects the `zitadel` ecosystem on write (422), so the author sets this row in the web UI — propose the values, never send `vulnerabilities` via the API.
- The advisory’s one-line summary/title is a separate metadata field — don’t restate it in the body.

## Output

1. The complete advisory body in the format above.
2. A short list of placeholders the author must fill or confirm (patch versions, affected upper bounds, reporter).
3. Proposed metadata (CVSS vector, CWEs, affected products) as suggestions to confirm — clearly separated from the body.
4. Do **not** push anything to GitHub (`gh api`) or publish without the author’s explicit confirmation; when authorized, PATCH only what the author approved (the description, and credits/CWEs/CVSS when asked — never affected products, see Metadata), then re-fetch to verify what landed. Mechanics: write the body to a scratch file and PATCH with `-F description=@body.md` (never inline shell-escaped text); the `credits` array is replaced wholesale on every PATCH — always resend the complete list, not just the new entry.
