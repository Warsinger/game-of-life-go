"use strict";
const session = new URLSearchParams(location.search).get('session');
const input = window.appInput = {left:false, right:false, jump:false};
const report = (type, message) => parent.postMessage({type, session, message}, location.origin);
window.appReady = () => report('app-ready');
window.addEventListener('error', event => report('app-error', String(event.error?.message || event.message)));
window.addEventListener('unhandledrejection', event => report('app-error', String(event.reason?.message || event.reason)));
window.addEventListener('message', event => {
  if (event.origin !== location.origin || event.source !== parent || event.data?.session !== session || event.data?.type !== 'app-input') return;
  const {key, down} = event.data;
  if (!['left','right','jump'].includes(key) || typeof down !== 'boolean') return;
  if (key !== 'jump' || down) input[key] = down;
});
function release() {input.left = input.right = input.jump = false;}
window.addEventListener('blur', release);
document.addEventListener('visibilitychange', () => { if (document.hidden) release(); });
window.addEventListener('keydown', event => { if (['Space','ArrowUp','ArrowDown','ArrowLeft','ArrowRight'].includes(event.code)) event.preventDefault(); });
(async () => {
  try {
    const go = new Go();
    const response = await fetch("game.wasm");
    if (!response.ok) throw new Error('Game download failed: HTTP ' + response.status);
    let bytes = await response.arrayBuffer();
    const magic = new Uint8Array(bytes, 0, Math.min(2, bytes.byteLength));
    if (magic[0] === 0x1f && magic[1] === 0x8b) {
      if (typeof DecompressionStream === 'undefined') throw new Error('This browser cannot decompress the game. Try a current browser.');
      bytes = await new Response(new Blob([bytes]).stream().pipeThrough(new DecompressionStream('gzip'))).arrayBuffer();
    }
    const result = await WebAssembly.instantiate(bytes, go.importObject);
    await go.run(result.instance);
    throw new Error('Game stopped. Use Apply / restart to run again.');
  } catch (error) { report('app-error', String(error.message || error)); }
})();
