export const meta = {
    name: 'topology-builder',
    description: 'Build a wikilink knowledge topology from multiple sources',
    phases: [
        { title: 'Discover', detail: 'one agent per source, collect entity candidates' },
        { title: 'Identify', detail: 'merge candidates into canonical entities' },
        { title: 'Extract', detail: 'one agent per entity, write dimensions' },
        { title: 'Connect', detail: 'one agent per entity, write wikilink edges' },
        { title: 'Verify', detail: 'link integrity, backlinks, _index.md' },
    ],
}

// args: { root: string, sources: [{ type, target }], seedZones?: string[] }
const root = (args && args.root) || '~/projects/product/topologies'
const sources = (args && args.sources) || []
if (!sources.length) throw new Error('args.sources is required: [{ type, target }]')

const SKILL = 'plugins/general/skills/topology-builder/SKILL.md (read the nearest copy; ' +
    'if unreadable, follow the rules embedded in this prompt)'

const CANDIDATES_SCHEMA = {
    type: 'object',
    required: ['candidates'],
    properties: {
        candidates: {
            type: 'array',
            items: {
                type: 'object',
                required: ['name', 'type', 'evidence'],
                properties: {
                    name: { type: 'string' },
                    type: { type: 'string' },
                    zone: { type: 'string' },
                    aliases: { type: 'array', items: { type: 'string' } },
                    evidence: { type: 'string' },
                },
            },
        },
    },
}

const ENTITIES_SCHEMA = {
    type: 'object',
    required: ['entities'],
    properties: {
        entities: {
            type: 'array',
            items: {
                type: 'object',
                required: ['name', 'type', 'zone'],
                properties: {
                    name: { type: 'string' },
                    type: { type: 'string' },
                    zone: { type: 'string' },
                    aliases: { type: 'array', items: { type: 'string' } },
                    sources: { type: 'array', items: { type: 'string' } },
                },
            },
        },
    },
}

const RULES = `Rules (from topology-builder skill):
- entity = deployable/operable/readable unit (service, system, datastore, corpus, channel, team); never a single code file or function
- file: <zone>/<name>.md under ${root}; kebab-case; filename globally unique across zones
- YAML frontmatter: name, type, zone, tags, aliases, sources
- 2~12 dimension sections per entity; first line under each "## " heading is a kind annotation "kind: concept|method|state|interface", then description, then a "References:" list
- kind is judged by the section's subject matter, not its implementation form: a rule/policy is "concept" even when implemented as a function; an action is "method" even when it acts on state data ("state" is reserved for states and lifecycles themselves)
- "## External Sources" and "## Backlinks" are fixed sections, not dimensions: no kind line, no References edges, excluded from the 3~12 cap
- edge: "- <relation> [[entity#Section]] — note"; relations: calls, uses, reads-from, writes-to, publishes-to, subscribes-to, depends-on, mentions, owned-by
- edge direction = initiator -> receiver; NEVER encode "invoked by X" as a forward edge — reverse relations belong to the auto-generated Backlinks section
- Section must match the target file's "## " heading verbatim (Read the file and copy, never spell from memory); same-file edges use [[#Section]]
- external URLs go to "## External Sources" as standard links, never wikilinks`

phase('Discover')
const found = await parallel(sources.map((s) => () =>
    agent(
        `Scan this source for knowledge-graph entity candidates.\n` +
        `Source type: ${s.type}\nTarget: ${s.target}\n${RULES}\n` +
        `Return candidates with evidence (one line: where you saw it). ` +
        `Do NOT write any files.`,
        { label: `discover:${s.type}`, phase: 'Discover', schema: CANDIDATES_SCHEMA },
    )))
const candidates = found.filter(Boolean).flatMap((r) => r.candidates)
log(`${candidates.length} entity candidates from ${sources.length} sources`)

phase('Identify')
const identified = await agent(
    `Merge these entity candidates into canonical entities.\n` +
    `Candidates: ${JSON.stringify(candidates)}\n${RULES}\n` +
    `Merge duplicates by name/alias overlap (accumulate aliases and sources), ` +
    `assign each entity a zone folder name (kebab-case), and for each entity ` +
    `write a skeleton file ${root}/<zone>/<name>.md containing only the YAML ` +
    `frontmatter and the "# Title" line. Create directories as needed. ` +
    `Return the final entity list.`,
    { label: 'identify:merge', phase: 'Identify', schema: ENTITIES_SCHEMA },
)
const entities = identified.entities
log(`${entities.length} canonical entities`)

phase('Extract')
// Barrier between Extract and Connect: edges must point at headings that
// already exist on disk, so every dimension section is written first.
await parallel(entities.map((e) => () =>
    agent(
        `Fill in the entity file ${root}/${e.zone}/${e.name}.md (it already has ` +
        `frontmatter). Sources to consult: ${JSON.stringify(e.sources || [])}.\n${RULES}\n` +
        `Add a one-line positioning sentence under the title, then 2~12 "## " ` +
        `dimension sections (concepts/methods/states/interfaces) each starting ` +
        `with its "kind: <concept|method|state|interface>" line, then a 1-3 ` +
        `sentence description and an empty "References:" list. No edges yet.`,
        { label: `extract:${e.name}`, phase: 'Extract' },
    )))

phase('Connect')
const registry = entities.map((e) => `${e.zone}/${e.name}`).join(', ')
await parallel(entities.map((e) => () =>
    agent(
        `Add directed edges for entity ${e.name} (${root}/${e.zone}/${e.name}.md).\n` +
        `Entity registry: ${registry}\n${RULES}\n` +
        `Read sibling entity files to copy their exact "## " headings into ` +
        `wikilinks. Fill each dimension's "References:" list and add ` +
        `"## External Sources" when applicable. Exploration depth: max 2 hops ` +
        `from this entity; entities beyond that go into a frontier note in your ` +
        `final message, NOT into new files. Return the frontier list (may be empty).`,
        { label: `connect:${e.name}`, phase: 'Connect' },
    )))

phase('Verify')
const report = await agent(
    `Run the Verification step of the topology-builder skill on ${root}: ` +
    `(1) duplicate filenames across zones, (2) broken wikilinks — target file ` +
    `missing or "## Section" heading not matching verbatim, (3) unlinked ` +
    `entities — counting only cross-entity forward edges in "References:" ` +
    `lists (never Backlinks sections or same-file [[#Section]] links): those ` +
    `with zero inbound edges and those with zero outbound edges, (4) inverted ` +
    `edges — forward edges whose note says ` +
    `"invoked/called/used by" the target, (5) dimension sections whose first ` +
    `line is not a valid "kind:" annotation (External Sources and Backlinks ` +
    `are exempt). Fix what is mechanical (heading typos, ` +
    `inverted directions), then rebuild every entity's "## Backlinks" section ` +
    `strictly by recomputing from forward edges across the whole graph (keep the ` +
    `"<!-- auto-generated: do not hand-edit -->" marker) and write ` +
    `${root}/_index.md with the entity registry table, a Mermaid overview ` +
    `(zones as subgraphs, edges aggregated to entity level), the Frontier list, ` +
    `and — as the LAST section of the file — "## Unlinked" with two lists of ` +
    `[[entity]] wikilinks: "no inbound" (no entity links to it) and ` +
    `"no outbound" (it links to no entity); mark entities on both lists ` +
    `"(orphan)", write "None" for empty lists. Return a short integrity report.`,
    { label: 'verify:integrity', phase: 'Verify' },
)

return { entities: entities.length, report }
