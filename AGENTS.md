# AGENTS.md

## Writing Style

When you write technical text (documentation, READMEs, runbooks, procedures, error messages, release notes, reports), obey these rules from ASD-STE100 Simplified Technical English:

CLASSIFY FIRST. Procedural text tells the reader what to do: imperative mood, maximum 20 words per sentence, one instruction per sentence. Descriptive text explains: simple tenses, maximum 25 words per sentence, one topic per paragraph, maximum six sentences per paragraph. Never mix the two in one passage.

VERBS. Use only: infinitive, imperative, simple present, simple past, simple future, past participle as adjective. No present perfect. No "-ing" verb forms. Active voice; passive only in descriptions when the agent is unknown. Approved modals: can, will, must. Banned: should, would, may, might, could. For "should": write "must" if required, delete if optional.

SENTENCES. Keep complete grammar: no contractions, keep articles, keep "that". Put conditions before commands, with a comma. No semicolons — write two sentences. Use a vertical list for more than two items or steps.

WORDS. One word, one meaning for the whole document. Pick one term and use it everywhere.

---

## Workflow Rules

1. Work phase by phase as defined in `docs/implementation.md`.
2. Do not start the next phase until the current phase has tests and passes verification.
3. Do not commit. The user commits manually after verification.
4. If the user raises a critic or issue, update `docs/implementation.md` and the relevant code. Do not start new work until the issue is resolved.
5. After each phase, write the test cases. The user will run manual verification on top of automated tests.
6. Write all test cases before marking a phase complete.
7. Do not remove existing comments or docstrings that are unrelated to your changes.

---

## Phase Completion Checklist

Before marking any phase complete, confirm all of the following:

- [ ] Code compiles without errors.
- [ ] All test cases for the phase pass.
- [ ] Tests cover the happy path and at least one error path.
- [ ] Manual verification steps are documented in `docs/implementation.md`.
- [ ] No unused imports or variables remain.

---

## Critic Response Rules

If the user gives a critic:

1. Read the critic.
2. Identify which phase and which file it relates to.
3. Update `docs/implementation.md` to reflect the changed approach.
4. Update the code.
5. Re-run the tests.
6. Report what changed and why.

Do not ask for approval to make small corrections. Make the fix and report.
