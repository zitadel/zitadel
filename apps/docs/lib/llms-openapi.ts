import path from 'path';
import { createOpenAPI } from 'fumadocs-openapi/server';

// Renders the <APIPage /> tags of generated API reference pages as plain
// Markdown, so that the LLM text output contains the endpoint details that the
// website renders from the OpenAPI spec (method, path, parameters, fields).

type Schema = Record<string, any>;
type ProcessedDocument = Awaited<ReturnType<ReturnType<typeof createOpenAPI>['getSchema']>>;

// Nested objects deeper than this are summarised instead of expanded, which
// keeps llms-full.txt at a reasonable size.
const MAX_DEPTH = 4;

const API_PAGE_TAG = /<APIPage\b([\s\S]*?)\/>/g;
const TAG_ATTRIBUTE = /(\w+)=(?:"([^"]*)"|\{([\s\S]*?)\})(?=\s+\w+=|\s*$)/g;

const documents = new Map<string, Promise<ProcessedDocument>>();

function loadDocument(document: string): Promise<ProcessedDocument> {
  const absolute = path.isAbsolute(document)
    ? document
    : path.join(process.cwd(), 'openapi', document.slice('openapi/'.length));

  let loaded = documents.get(absolute);
  if (!loaded) {
    loaded = createOpenAPI({ input: [absolute] }).getSchema(absolute);
    documents.set(absolute, loaded);
  }
  return loaded;
}

function decodeEntities(value: string): string {
  return value
    .replace(/&#x([0-9a-f]+);/gi, (_, hex) => String.fromCodePoint(parseInt(hex, 16)))
    .replace(/&#(\d+);/g, (_, dec) => String.fromCodePoint(parseInt(dec, 10)))
    .replace(/&quot;/g, '"')
    .replace(/&lt;/g, '<')
    .replace(/&gt;/g, '>')
    .replace(/&amp;/g, '&');
}

// Processed markdown serialises JSX attributes as `name="&#x22;...&#x22;"`,
// raw MDX as `name={...}`. Both hold JSON.
function parseAttributes(source: string): Record<string, unknown> {
  const attributes: Record<string, unknown> = {};
  for (const [, name, quoted, braced] of source.matchAll(TAG_ATTRIBUTE)) {
    const value = braced ?? decodeEntities(quoted);
    try {
      attributes[name] = JSON.parse(value);
    } catch {
      attributes[name] = value;
    }
  }
  return attributes;
}

function oneLine(text: string | undefined): string {
  return (text ?? '').replace(/\s+/g, ' ').trim();
}

function isNull(schema: Schema): boolean {
  return schema.type === 'null';
}

// `oneOf` is used for protobuf oneofs (several alternatives with properties),
// for optional fields (`[X, { type: 'null' }]`) and for plain unions such as
// `google.protobuf.Value`.
function isChoice(alternatives: Schema[] | undefined): alternatives is Schema[] {
  const present = alternatives?.filter((alternative) => !isNull(alternative)) ?? [];
  return present.length > 1 && present.some((alternative) => alternative.properties);
}

// Unwraps an optional field to its actual schema. The inner schema is returned
// as is, so it can still be matched by identity.
function unwrap(schema: Schema): Schema {
  const alternatives = schema.oneOf?.filter((alternative: Schema) => !isNull(alternative));
  return alternatives?.length === 1 ? alternatives[0] : schema;
}

// Merges `allOf` parts into a single object schema and collects `oneOf`
// groups (protobuf oneofs), which are rendered as "exactly one of".
function flatten(schema: Schema): { properties: Record<string, Schema>; required: Set<string>; oneOf: Schema[][] } {
  const properties: Record<string, Schema> = { ...(schema.properties ?? {}) };
  const required = new Set<string>(schema.required ?? []);
  const oneOf: Schema[][] = isChoice(schema.oneOf) ? [schema.oneOf] : [];

  for (const part of schema.allOf ?? []) {
    const nested = flatten(part);
    Object.assign(properties, nested.properties);
    nested.required.forEach((name) => required.add(name));
    oneOf.push(...nested.oneOf);
  }
  return { properties, required, oneOf };
}

function isObject(schema: Schema): boolean {
  return Boolean(schema.properties || schema.allOf || isChoice(schema.oneOf));
}

function typeLabel(input: Schema): string {
  const schema = unwrap(input);
  if (schema.oneOf && !isChoice(schema.oneOf)) {
    return schema.oneOf
      .filter((alternative: Schema) => !isNull(alternative))
      .map(typeLabel)
      .join(' | ');
  }
  if (schema.enum) return `enum: ${schema.enum.join(', ')}`;
  if (schema.const !== undefined) return `const: ${schema.const}`;

  const types = ([] as string[]).concat(schema.type ?? []).filter((type) => type !== 'null');
  if (types.includes('array')) return `array of ${schema.items ? typeLabel(schema.items) : 'any'}`;
  if (types.includes('object') && schema.additionalProperties && typeof schema.additionalProperties === 'object') {
    return `map of string to ${typeLabel(schema.additionalProperties)}`;
  }
  if (types.length === 0 && isObject(schema)) return 'object';

  const label = types.join(' | ') || 'any';
  return schema.format ? `${label} (${schema.format})` : label;
}

// The schema whose fields should be listed below a property, if any.
function expandable(input: Schema): Schema | undefined {
  const schema = unwrap(input);
  if (isObject(schema)) return schema;
  if (schema.items && isObject(schema.items)) return schema.items;
  const values = schema.additionalProperties;
  if (values && typeof values === 'object' && isObject(values)) return values;
  return undefined;
}

function renderFields(
  schema: Schema,
  doc: ProcessedDocument,
  depth: number,
  ancestors: Set<Schema>,
  lines: string[],
): void {
  const indent = '  '.repeat(depth);
  const { properties, required, oneOf } = flatten(schema);

  for (const [name, property] of Object.entries(properties)) {
    renderProperty(name, property, required.has(name), doc, depth, ancestors, lines);
  }

  for (const alternatives of oneOf) {
    lines.push(`${indent}- Exactly one of:`);
    for (const alternative of alternatives) {
      const nested = flatten(alternative);
      for (const [name, property] of Object.entries(nested.properties)) {
        renderProperty(name, property, false, doc, depth + 1, ancestors, lines);
      }
    }
  }
}

function renderProperty(
  name: string,
  property: Schema,
  required: boolean,
  doc: ProcessedDocument,
  depth: number,
  ancestors: Set<Schema>,
  lines: string[],
): void {
  const indent = '  '.repeat(depth);
  const description = oneLine(property.description ?? unwrap(property).description);
  const flags = [typeLabel(property), required ? 'required' : undefined, property.deprecated ? 'deprecated' : undefined]
    .filter(Boolean)
    .join(', ');
  lines.push(`${indent}- \`${name}\` (${flags})${description ? `: ${description}` : ''}`);

  const nested = expandable(property);
  if (!nested) return;

  if (ancestors.has(nested) || depth + 1 >= MAX_DEPTH) {
    const ref = doc.getRawRef(nested)?.split('/').pop();
    if (ref) lines.push(`${indent}  - Nested fields of \`${ref}\` not shown`);
    return;
  }

  ancestors.add(nested);
  renderFields(nested, doc, depth + 1, ancestors, lines);
  ancestors.delete(nested);
}

function renderBody(title: string, schema: Schema | undefined, doc: ProcessedDocument, lines: string[]): void {
  if (!schema) return;
  const fields: string[] = [];
  renderFields(schema, doc, 0, new Set([schema]), fields);

  lines.push('', title, '');
  lines.push(...(fields.length > 0 ? fields : ['Empty object (`{}`).']));
}

function renderOperation(doc: ProcessedDocument, pathName: string, method: string): string | undefined {
  const operation: Schema | undefined = (doc.dereferenced.paths as Schema | undefined)?.[pathName]?.[method];
  if (!operation) return undefined;

  const lines = [`## ${method.toUpperCase()} ${pathName}`];
  if (operation.operationId) lines.push('', `Operation ID: \`${operation.operationId}\``);

  const parameters: Schema[] = operation.parameters ?? [];
  for (const location of ['path', 'query', 'header']) {
    const matching = parameters.filter((parameter) => parameter.in === location);
    if (matching.length === 0) continue;

    lines.push('', `### ${location[0].toUpperCase()}${location.slice(1)} parameters`, '');
    for (const parameter of matching) {
      const flags = [typeLabel(parameter.schema ?? {}), parameter.required ? 'required' : undefined]
        .filter(Boolean)
        .join(', ');
      const description = oneLine(parameter.description ?? parameter.schema?.description);
      lines.push(`- \`${parameter.name}\` (${flags})${description ? `: ${description}` : ''}`);
    }
  }

  for (const [type, media] of Object.entries<Schema>(operation.requestBody?.content ?? {})) {
    renderBody(`### Request body (${type})`, media.schema, doc, lines);
  }

  for (const [status, response] of Object.entries<Schema>(operation.responses ?? {})) {
    const media = Object.entries<Schema>(response.content ?? {})[0];
    if (!status.startsWith('2')) {
      // The error shape is the same for every endpoint, so only its top-level fields are listed.
      const fields = Object.keys(flatten(media?.[1].schema ?? {}).properties);
      const summary = fields.length > 0 ? ` with fields ${fields.map((field) => `\`${field}\``).join(', ')}` : '';
      lines.push('', `### Response ${status}`, '', `${oneLine(response.description) || 'Error'}${summary}.`);
      continue;
    }
    renderBody(`### Response ${status}${media ? ` (${media[0]})` : ''}`, media?.[1].schema, doc, lines);
  }

  return lines.join('\n');
}

async function renderApiPage(source: string): Promise<string | undefined> {
  const { document, operations } = parseAttributes(source);
  if (typeof document !== 'string' || !Array.isArray(operations)) return undefined;

  const doc = await loadDocument(document);
  const rendered = operations
    .map(({ path: pathName, method }) => renderOperation(doc, pathName, method))
    .filter((section): section is string => Boolean(section));

  return rendered.length > 0 ? rendered.join('\n\n') : undefined;
}

export async function expandApiPages(text: string): Promise<string> {
  const tags = [...text.matchAll(API_PAGE_TAG)];
  if (tags.length === 0) return text;

  let result = text;
  for (const [tag, source] of tags) {
    try {
      const rendered = await renderApiPage(source);
      if (rendered) result = result.replace(tag, () => rendered);
    } catch (error) {
      console.warn(`[llms] Could not render ${tag}:`, error);
    }
  }
  return result;
}
