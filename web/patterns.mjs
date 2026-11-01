export const MAX_SIDE = 64;
export const MAX_BYTES = 8192;
export const STORE_KEY = 'golgo.starting-patterns.v1';

export function parseGrid(value) {
  const raw = String(value ?? '');
  if (new TextEncoder().encode(raw).length > MAX_BYTES) throw new Error('Starting grid must be at most 8192 bytes.');
  const rows = raw.replaceAll('\r\n', '\n').trim().split('\n');
  const width = rows[0].length;
  if (!width || width > MAX_SIDE || rows.length > MAX_SIDE) throw new Error('Starting grid must be 1–64 columns and rows.');
  for (const row of rows) {
    if (row.length !== width) throw new Error('All starting-grid rows must have the same width.');
    if (!/^[.O]+$/.test(row)) throw new Error('Use O for live cells and . for dead cells.');
  }
  return {width, height: rows.length, text: rows.join('\n')};
}

export function startingParams(params, values, library) {
  const id = String(values.pattern || 'random');
  params.set('pattern', id);
  if (id === 'random') return params;
  const pattern = id === 'custom' ? values.seed : library.find(p => p.id === id)?.grid;
  if (pattern === undefined) throw new Error('Choose an available starting pattern.');
  const grid = parseGrid(pattern);
  if (grid.width > Number(params.get('width')) || grid.height > Number(params.get('height'))) {
    throw new Error(`Starting pattern is ${grid.width}×${grid.height}; increase the board dimensions.`);
  }
  if (id === 'custom') params.set('seed', grid.text);
  return params;
}

export class PatternStore {
  constructor(getStorage = () => globalThis.localStorage) {
    this.getStorage = getStorage;
    this.entries = [];
    this.persistent = true;
    try {
      const raw = getStorage().getItem(STORE_KEY);
      if (raw && raw.length > 200000) throw new Error('Oversized library');
      const entries = JSON.parse(raw || '[]');
      if (!Array.isArray(entries) || entries.length > 20) throw new Error('Invalid library');
      const seen = new Set();
      for (const entry of entries) {
        if (!entry || typeof entry.name !== 'string' || !entry.name.trim() || entry.name.length > 40 || seen.has(entry.name)) throw new Error('Invalid saved pattern');
        seen.add(entry.name);
        this.entries.push({name:entry.name, grid:parseGrid(entry.grid).text});
      }
    } catch {
      this.entries = [];
      this.persistent = false;
    }
  }
  list() { return this.entries.map(entry => ({...entry})); }
  persist() {
    try { this.getStorage().setItem(STORE_KEY, JSON.stringify(this.entries)); this.persistent = true; }
    catch { this.persistent = false; }
    return this.persistent;
  }
  save(name, grid) {
    name = String(name).trim();
    if (!name || name.length > 40) throw new Error('Name the pattern using 1–40 characters.');
    grid = parseGrid(grid).text;
    if (this.entries.length >= 20 && !this.entries.some(e => e.name === name)) throw new Error('Your library holds 20 patterns; delete one before adding another.');
    this.entries = [{name,grid}, ...this.entries.filter(e => e.name !== name)];
    return this.persist();
  }
  remove(name) { this.entries = this.entries.filter(entry => entry.name !== name); return this.persist(); }
}
