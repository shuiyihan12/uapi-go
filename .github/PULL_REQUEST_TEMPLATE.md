## Purpose

<!-- What does this PR change and why? Link the related issue/task if any. -->

## Key changes

<!-- Bullet list of the substantive changes. For behavior changes, include sample request/response snippets or logs when relevant. -->

## Checklist

- [ ] `./scripts/build.sh all` fully green (CI must pass)
- [ ] Behavior changes come with test evidence (commands and output) or sample requests/responses
- [ ] New endpoints updated `docs/routing.md` (and the README capability table)
- [ ] `pkg/generated/` contains only changes produced by `./scripts/build.sh wsdl`
- [ ] Comments and docs follow the English-default convention (`*.zh-CN.md` counterparts updated when applicable)

## Domain ownership

Direct commits are reserved for **hotel / tools / system**. Changes to **air / rail / vehicle / universal / sharedBooking / uprofile / sharedUprofile / gdsQueue / terminal** or other domains must go through this PR for review by the domain owner.
