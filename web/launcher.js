import {startLauncher} from './shared/launcher.mjs';
import {validateOptions} from './options.mjs';
import {startingParams} from './patterns.mjs';
import {setupPatterns} from './pattern-ui.mjs';
let library = [];
// Start the random game immediately; loading the optional library must not block it.
setupPatterns(document).then(patterns => { library = patterns; });
startLauncher({
  validateOptions: values => startingParams(validateOptions(values, 'life'), values, library),
  readOptions: form => ({...Object.fromEntries(new FormData(form)), debug: form.elements.debug.checked, lines: form.elements.lines.checked}),
  readyMessage: 'Ready. R restarts the selected board. + / − adjust speed.'
});
