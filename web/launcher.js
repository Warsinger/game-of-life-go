import {startLauncher} from './shared/launcher.mjs';
import {validateOptions} from './options.mjs';
startLauncher({
  validateOptions: values => validateOptions(values, 'life'),
  readOptions: form => ({...Object.fromEntries(new FormData(form)), debug: form.elements.debug.checked, lines: form.elements.lines.checked}),
  readyMessage: 'Ready. Settings apply when you restart the board.'
});
