import {parseGrid, PatternStore} from './patterns.mjs';

// App-specific starting-grid UI; the shared launcher still owns the game lifecycle.
export async function setupPatterns(document, fetchLibrary = () => fetch(new URL('./patterns.json', import.meta.url), {signal: AbortSignal.timeout(5000)}), store = new PatternStore()) {
  const form = document.querySelector('form');
  const source = document.getElementById('pattern');
  const editor = document.getElementById('pattern-editor');
  const seed = document.getElementById('seed');
  const description = document.getElementById('pattern-description');
  const status = document.getElementById('pattern-status');
  const saved = document.getElementById('saved-pattern');
  const name = document.getElementById('pattern-name');
  let library = [];
  try {
    const response = await fetchLibrary();
    if (!response.ok) throw new Error('Library unavailable');
    const data = await response.json();
    if (!Array.isArray(data) || data.length > 50) throw new Error('Invalid library');
    const seen = new Set(['random','custom']);
    for (const pattern of data) {
      if (typeof pattern.id !== 'string' || !/^[a-z0-9-]+$/.test(pattern.id) || seen.has(pattern.id) || typeof pattern.name !== 'string' || typeof pattern.description !== 'string') throw new Error('Invalid pattern');
      parseGrid(pattern.grid); seen.add(pattern.id);
    }
    library = data;
    for (const pattern of library) {
      const option = document.createElement('option'); option.value = pattern.id; option.textContent = pattern.name;
      source.append(option);
    }
  } catch { status.textContent = 'Built-in library could not load. Random and custom grids are still available.'; }
  if (!store.persistent) status.textContent += ' Saved patterns are unavailable; new saves may last only this session.';
  function refreshSaved() {
    saved.replaceChildren();
    for (const entry of store.list()) { const option = document.createElement('option'); option.value = entry.name; option.textContent = entry.name; saved.append(option); }
    document.getElementById('load-pattern').disabled = !store.list().length;
    document.getElementById('delete-pattern').disabled = !store.list().length;
  }
  function updateSelection() {
    const pattern = library.find(p => p.id === source.value);
    editor.hidden = source.value === 'random';
    description.textContent = pattern?.description || (source.value === 'random' ? 'A fresh random board on each restart.' : 'Define a rectangular starting grid. It is centered on an otherwise empty board.');
    if (pattern) seed.value = pattern.grid;
  }
  function report(persistent, action) { status.textContent = persistent ? `${action} in this browser.` : `${action} for this session only; browser storage is unavailable.`; }
  source.addEventListener('change',updateSelection);
  seed.addEventListener('input',()=>{source.value='custom';updateSelection();});
  document.getElementById('load-paused').addEventListener('click',()=>{form.elements.speed.value='0';form.requestSubmit();});
  document.getElementById('save-pattern').addEventListener('click',()=>{
    try {
      if (source.value === 'random') throw new Error('Choose or enter a pattern to save; the running random board is not captured.');
      report(store.save(name.value,seed.value),'Saved'); refreshSaved();
    } catch(error) {status.textContent=error.message;}
  });
  document.getElementById('load-pattern').addEventListener('click',()=>{
    const entry=store.list().find(p=>p.name===saved.value);if(!entry)return;
    source.value='custom';seed.value=entry.grid;name.value=entry.name;updateSelection();
    status.textContent='Pattern loaded into the editor. Apply / restart or Load paused to use it.';
  });
  document.getElementById('delete-pattern').addEventListener('click',()=>{report(store.remove(saved.value),'Deleted');refreshSaved();});
  refreshSaved(); updateSelection();
  return library;
}
