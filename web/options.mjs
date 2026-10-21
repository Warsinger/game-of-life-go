export function validateOptions(values, kind) {
  if (!['life', 'rover'].includes(kind)) throw new Error('Unknown application');
  const integer = (name, min, max) => {
    const raw = String(values[name] ?? '').trim();
    const n = Number(raw);
    if (!raw || !Number.isInteger(n) || n < min || n > max) throw new Error(name + ' must be an integer between ' + min + ' and ' + max + '.');
    return n;
  };
  const width = integer('width', 1, kind === 'life' ? 4096 : 8192);
  const height = integer('height', 1, kind === 'life' ? 4096 : 8192);
  const speed = integer('speed', 0, 60);
  const result = new URLSearchParams({ width, height, speed, debug: String(values.debug === true) });
  if (kind === 'life') {
    const cell = integer('cell', 1, 8192);
    const raw = String(values.variance ?? '').trim();
    const variance = Number(raw);
    if (!raw || !Number.isFinite(variance) || variance < 0 || variance > 1) throw new Error('Live-cell probability must be between 0 and 1.');
    if (width * height > 4_000_000) throw new Error('Use at most 4,000,000 cells.');
    if (width * cell > 8192 || height * cell > 8192) throw new Error('The board must fit within 8192 pixels in each direction.');
    result.set('cell', cell); result.set('variance', variance); result.set('lines', String(values.lines === true));
  }
  return result;
}
