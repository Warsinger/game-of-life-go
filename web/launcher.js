import { validateOptions } from './options.mjs';
const form = document.querySelector('form');
const frame = document.getElementById('game');
const status = document.getElementById('status');
const kind = document.body.dataset.app;
let session = '';
let timeout;
let ready = false;
function focusGame() {
  frame.focus({preventScroll: true});
  const canvas = frame.contentDocument?.querySelector('canvas');
  if (canvas) { canvas.tabIndex = 0; canvas.focus({preventScroll: true}); }
}
function sendInput(key, down) {
  if (ready) frame.contentWindow.postMessage({type:'app-input', session, key, down}, location.origin);
}
function releaseInput() { for (const key of ['left','right']) sendInput(key, false); }
function launch(event) {
  event?.preventDefault();
  const values = Object.fromEntries(new FormData(form));
  values.debug = form.elements.debug.checked;
  if (kind === 'life') values.lines = form.elements.lines.checked;
  let query;
  try { query = validateOptions(values, kind); }
  catch (error) { status.textContent = error.message; status.dataset.state = 'error'; return; }
  releaseInput(); ready = false;
  session = crypto.randomUUID(); query.set('session', session);
  frame.style.aspectRatio = query.get('width') + '/' + query.get('height');
  const url = new URL(frame.dataset.src, location.href); url.search = query.toString();
  clearTimeout(timeout);
  status.textContent = 'Loading… The first visit downloads the game.'; status.dataset.state = 'loading';
  frame.src = url.href;
  timeout = setTimeout(() => { if (!ready) { status.textContent = 'Still loading. If the game does not appear, use Apply / restart to retry.'; } }, 60000);
}
window.addEventListener('message', event => {
  if (event.origin !== location.origin || event.source !== frame.contentWindow || event.data?.session !== session) return;
  if (event.data.type === 'app-ready' && status.dataset.state !== 'runtime-error') {
    clearTimeout(timeout); ready = true; status.dataset.state = 'ready';
    status.textContent = kind === 'life' ? 'Ready. Settings apply when you restart the board.' : 'Ready. Use arrow keys and Space, or the on-screen controls.';
    focusGame();
  } else if (event.data.type === 'app-error') {
    clearTimeout(timeout); ready = false; status.dataset.state = 'runtime-error';
    status.textContent = 'Could not run. Apply / restart to retry. ' + event.data.message;
  }
});
form.addEventListener('submit', launch);
document.getElementById('focus').addEventListener('click', focusGame);
document.getElementById('fullscreen').addEventListener('click', async () => {
  try { await frame.requestFullscreen(); focusGame(); }
  catch { status.textContent = 'Fullscreen is unavailable in this browser. You can still play in the page.'; }
});
for (const button of document.querySelectorAll('[data-control]')) {
  button.addEventListener('pointerdown', event => {
    if (event.button !== 0) return;
    event.preventDefault(); button.setPointerCapture(event.pointerId); sendInput(button.dataset.control, true);
  });
  for (const name of ['pointerup','pointercancel','lostpointercapture']) button.addEventListener(name, () => sendInput(button.dataset.control, false));
  button.addEventListener('keydown', event => { if (['Enter',' '].includes(event.key)) {event.preventDefault(); if (!event.repeat) sendInput(button.dataset.control, true);} });
  button.addEventListener('keyup', event => { if (['Enter',' '].includes(event.key)) {event.preventDefault(); sendInput(button.dataset.control, false);} });
}
window.addEventListener('blur', releaseInput);
document.addEventListener('visibilitychange', () => { if (document.hidden) releaseInput(); });
launch();
