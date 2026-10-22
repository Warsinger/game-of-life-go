import {startGoGame} from './shared/go-runtime.mjs';
startGoGame({wasmURL: new URL('game.wasm', import.meta.url), preventKeys: ['Space', 'ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight']});
