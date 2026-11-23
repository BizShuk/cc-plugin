# content-summarizer — Steps 4-6: Strategies, Business Value, Output

## Step 4 — Summarize

### Index strategy (recursive summary)

The point is NOT to relist the page. Dig into each item.

1. Extract the items (title + link).
2. Default to the `top 5-8` by prominence/relevance. If more exist, summarize
   these, state how many remain, and offer to continue. Confirm before fetching
   a large batch.
3. For each item: fetch it (re-run the Step 1 check per item) and write
   `2-4 sentences` — what it is and why it matters. If the item is technical,
   also flag any standout configuration option, setting, or operational caveat.
   Skip items that can't be fetched and note which were skipped.
4. Go `one level deep only` (index → item). Do not recurse into items-of-items
   unless asked.
5. Open with a short `roll-up`: the 2-3 themes or patterns across the items.

### Article strategy (concise)

Optimize for speed and signal. One fetch, no deep recursion.

- `TL;DR`: 1-2 sentences.
- `Key points`: 3-6 bullets, each a complete thought.
- `Pros / Cons`: include ONLY when the piece is evaluative (product, proposal,
  argument, recommendation). If it is purely informational, skip pros/cons —
  do not force it.
- `Configuration & fine print`: when the content includes technical or
  operational detail, list the easy-to-miss specifics — special configuration
  options, key settings or flags, defaults, prerequisites, limits, and version-
  or platform-specific caveats, plus any non-obvious steps. Capture exact names
  and values (a parameter, a flag, a limit), not vague paraphrases. Include
  only when such details exist.
- Do not pad the factual summary.

## Step 5 — Business value (1-3 ideas)

After the summary, add `1-3` concrete business-value or opportunity ideas drawn
from the content. This is the high-leverage part: move beyond "what it says" to
"so what — how could this matter".

- Produce `1-3` items, calibrated to how much signal the content offers. One
  sharp idea beats three generic ones.
- Make each idea `specific and tied to this content`. Bad: "could improve
  efficiency". Good: name the angle, who it helps, and the concrete move.
- Brainstorm across angles, picking what fits: a product/feature idea, a market
  or customer opportunity, a competitive/strategic implication, a process or
  cost improvement, a risk to watch, or a concrete next action.
- `Label these as your own inference, not the source.` They are extrapolations —
  keep them clearly separate from the faithful summary so they are never
  mistaken for claims the source made.
- Tailor to the user's context when known (role, company, goals); otherwise
  keep ideas broadly applicable.
- If the content has no plausible business angle, say so briefly instead of
  forcing ideas.

## Output format

Lead with the source so the summary is always traceable. Keep formatting light.
Match length to type: Articles short; Index pages as long as the item count requires.

```md
<title> or file: filename.pdf

source: [title](url) or pdf attachment

TL;DR: ...

<br><br>

## Key points

- ...

(Index only) By item

- Item title — 2-4 sentence summary

(if evaluative) Pros / Cons

-   - ...
-   - ...

<br><br>

(if technical / operational) Configuration & fine print

- setting / flag / option — what it does, default, caveat

<br><br>

## Business value (ideas / inference, not from the source)

1. ...
2. ...
   (1-3 items)

(if any) Dates & action items
```

## Step 6 — Calendar & Notes follow-ups

Offer these whenever they apply; also perform them on direct request. Always
embed the source link or file reference in whatever gets saved.

`Calendar` — if the content contains dated events (webinar, deadline, meetup,
launch, conference session):

- Use your available calendar capability, preferring `Apple Calendar` when an
  Apple calendar integration is present. If you have no calendar capability,
  say so.
- Always confirm date, time, and timezone before creating. Put the source URL
  in the event notes.

`Notes` — to save the summary itself:

- Use your available notes capability, preferring `Apple Notes` when an Apple
  Notes integration is present.
- Title the note clearly (source title + date) and place the source link / file
  name at the top.
- If you have no notes capability, say so and fall back to: returning the
  summary as a clearly formatted, copy-ready block, or saving it as a file.
  Never silently drop the request.
