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
- external URLs go to "## External Sources" as standard links, never wikilinks
- noise exclusion: EXCLUDE test files (*_test.go, *.spec.ts), mock/stub files, local script utilities, standard library wrappers, and generic helper/utility packages (e.g. logger, json, config-parser, db-driver). Do not model them as entities or edges unless explicitly directed.
- topology accuracy: ONLY link for substantial functional/architectural interactions (RPC, API call, DB write, message pub/sub). DO NOT link for trivial helper imports or infrastructure dependencies.
- alias merging: Verify evidence context; never merge different logical entities in separate zones just because they share a name (e.g. "auth" middleware vs "auth" service) — name them uniquely with a zone prefix (e.g. "payments-handler" and "auth-handler") if necessary.
- edge grounding: Each edge must have direct source evidence from the initiator (e.g., code call, import, SQL query, config read, or message content) noted in the relation description (e.g. "file:line", API, or table name). Do not link based on naming similarities, conventions, or hypothesis.
- mentions limit: Keep mentions relation edges to a maximum of 1~2 per dimension; if more are needed, define a more specific relation or remove them.`



phase('Discover')
const found = await parallel(sources.map((s) => () =>
    agent(
        `Scan this source for knowledge-graph entity candidates. Be highly selective to avoid noise: EXCLUDE helpers, utility scripts, tests (*_test.go, *.spec.*), mocks/fakes, and standard/third-party package wrappers. Focus only on substantial, independent units like services, databases, or main modules.\n` +
        `Source type: ${s.type}\nTarget: ${s.target}\n${RULES}\n` +
        `Return candidates with evidence (one line: where you saw it). ` +
        `Do NOT write any files.`,
        { label: `discover:${s.type}`, phase: 'Discover', schema: CANDIDATES_SCHEMA },
    )))

const candidates = found.filter(Boolean).flatMap((r) => r.candidates)
log(`${candidates.length} entity candidates from ${sources.length} sources`)

phase('Identify')
const identified = await agent(
    `Merge these entity candidates into canonical entities. Prevent false-positive merges: only merge if name/alias overlap indicates the EXACT SAME logical component in the same business context. If two candidates have the same name but operate in different contexts (e.g. "auth" helper vs "auth" service), keep them separate with distinct names/zones (e.g., payments-handler vs auth-handler). Verify candidate evidence before merging.\n` +
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
        `sentence description and an empty "References:" list. Exclude trivial internal helper functions or boilerplate getter/setters. Focus on key business concepts, API interfaces, main lifecycles, and core architectural components. No edges yet.`,
        { label: `extract:${e.name}`, phase: 'Extract' },
    )))


phase('Connect')
const registry = entities.map((e) => `${e.zone}/${e.name}`).join(', ')
await parallel(entities.map((e) => () =>
    agent(
        `Add directed edges for entity ${e.name} (${root}/${e.zone}/${e.name}.md). Be extremely accurate: ONLY link when there is a significant, real-world connection (such as RPC invocation, DB read/write, or message publishing/subscribing) between entities. DO NOT build links for logger calls, config loading, or helper utility imports. For each link, provide a descriptive relation and a brief note explaining the business context.\n` +
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
    `are exempt), (6) edge grounding audit — randomly audit N forward edges to ensure they contain direct source evidence (file:line, table name, API), remove or downgrade those without evidence and list them in the report, (7) hub-noise detection — identify entities with exceptionally high inbound+outbound edges primarily using "mentions" relations, list them as hub-noise candidates for manual review/dismantling. ` +
    `Fix what is mechanical (heading typos, inverted directions), then rebuild every entity's "## Backlinks" section ` +
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
