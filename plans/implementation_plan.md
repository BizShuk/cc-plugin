# Implementation Plan - Update Go Skills Guidelines

This plan outlines the modifications to the Go guidelines (`golang-code-quality`, `golang-mvc`, `golang-naming`, and `golang-dev` skills) to align with the user's explicit preferences.

## User Review Required

> [!IMPORTANT]
> The following rules will be applied across all Go-related skill files:
>
> 1. All constants must use `SCREAMING_SNAKE_CASE`.
> 2. Configuration loading must use `config.Default()` from `github.com/bizshuk/gosdk`, falling back to `viper` only if not supported.
> 3. Global state is acceptable/good for client, handler, and configuration if they are immutable.
> 4. Use `model` (singular) as the default package for domain models, unless there are more than 30 models (then split into domain-specific packages).

## Proposed Changes

---

### golang-code-quality

#### [MODIFY] [SKILL.md](file:///Users/shuk/projects/cc-plugin/skills/golang-code-quality/SKILL.md)

- Update package singular rule to explicitly state `model` is used as default, but must be split into domain-specific packages if there are more than 30 models.
- Relax the `No global state` rule: allow global state for client, handler, and configuration if they are immutable.
- Explicitly add the rule that constants must use `SCREAMING_SNAKE_CASE`.
- Specify that configuration loading uses `config.Default()` from the Go SDK, falling back to `viper`.

---

### golang-mvc

#### [MODIFY] [SKILL.md](file:///Users/shuk/projects/cc-plugin/skills/golang-mvc/SKILL.md)

- Update Step 1 `Model first` to specify `model/` (singular) is the default package, unless there are >30 models (then split into domain-specific packages).
- Update Step 5 `Wiring` to require `config.Default()` from the SDK, falling back to `viper`.
- Update error/constant rules to enforce `SCREAMING_SNAKE_CASE` for constants.
- Allow global state for client, handler, and configuration if they are immutable.

---

### golang-naming

#### [MODIFY] [SKILL.md](file:///Users/shuk/projects/cc-plugin/skills/golang-naming/SKILL.md)

- Update Rule P5: Specify that `model` (singular) is the default package for domain types, unless there are more than 30 models (then group/split them).

---

### golang-dev

#### [MODIFY] [SKILL.md](file:///Users/shuk/projects/cc-plugin/skills/golang-dev/SKILL.md)

- Update Section 2 `Configuration (viper)` to prefer `config.Default()` from the Go SDK, falling back to raw `viper` manual setup only if the SDK is not supported/available.

---

## Verification Plan

### Manual Verification

## Verification Plan

### Manual Verification

- Review all modified markdown files to ensure the updated rules are clear, consistent, and do not contain any `**` bolding (only `backtick` highlighting).
