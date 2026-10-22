import {test} from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import {gunzipSync} from 'node:zlib';
import {packagePages} from './package-pages.mjs';
function fixture(t) {
 const root=fs.mkdtempSync(path.join(os.tmpdir(),'game-pages-test-'));
 t.after(()=>fs.rmSync(root,{recursive:true,force:true}));
 const source=path.join(root,'source');fs.mkdirSync(source);
 for(const file of ['index.html','game.html','launcher.js','options.mjs','game.js','styles.css']) fs.copyFileSync(new URL('../web/'+file,import.meta.url),path.join(source,file));
 fs.cpSync(new URL('../node_modules/@warsinger/web-shared/browser',import.meta.url),path.join(source,'shared'),{recursive:true});
 for(const name of ['favicon.ico','icon-180.png','icon-192.png','icon-512.png']) fs.copyFileSync(new URL('../assets/icons/'+name,import.meta.url),path.join(source,name));
 fs.writeFileSync(path.join(source,'wasm_exec.js'),'// matching Go runtime');
 fs.writeFileSync(path.join(source,'game.wasm'),Buffer.from([0,97,115,109,1,0,0,0]));
 return {root,source};
}
test('packages actual HTML/module references, gzip bytes, revision and cache policy',t=>{
 const {root,source}=fixture(t);const output=path.join(root,'pages');
 const hash=packagePages(source,output,'abcdef');const assets=path.join(output,'assets',hash);
 assert.deepEqual(gunzipSync(fs.readFileSync(path.join(assets,'game.wasm.gz'))),fs.readFileSync(path.join(source,'game.wasm')));
 assert.match(fs.readFileSync(path.join(assets,'game.js'),'utf8'),/new URL\('game.wasm.gz', import.meta.url\)/);
 const index=fs.readFileSync(path.join(output,'index.html'),'utf8');
 for(const file of ['styles.css','launcher.js','game.html']) assert.ok(index.includes('assets/'+hash+'/'+file));
 for(const match of fs.readFileSync(path.join(assets,'game.html'),'utf8').matchAll(/src="([^"]+)"/g)) assert.ok(fs.existsSync(path.join(assets,match[1])));
 assert.ok(fs.readFileSync(path.join(assets,'launcher.js'),'utf8').includes("'./options.mjs'"));
 for (const file of ['launcher.mjs','frame.mjs','status.mjs','go-runtime.mjs','shell.css','header.css','footer.css']) assert.ok(fs.existsSync(path.join(assets,'shared',file)));
 assert.deepEqual(JSON.parse(fs.readFileSync(path.join(output,'version.json'))),{revision:'abcdef',assets:'assets/'+hash});
 const headers=fs.readFileSync(path.join(output,'_headers'),'utf8');assert.match(headers,/no-cache/);assert.match(headers,/immutable/);assert.match(headers,/application\/gzip/);
 assert.ok(fs.existsSync(path.join(output,'404.html')));
 assert.throws(()=>packagePages(source,output),/EEXIST/);
});
test('asset paths change when implementation changes and remain deterministic',t=>{
 const {root,source}=fixture(t);
 const a=packagePages(source,path.join(root,'a'),'one');const b=packagePages(source,path.join(root,'b'),'two');assert.equal(a,b);
 fs.appendFileSync(path.join(source,'options.mjs'),'\n// changed');
 const c=packagePages(source,path.join(root,'c'));assert.notEqual(c,a);
 fs.appendFileSync(path.join(source,'shared/shell.css'),'\n/* shared change */');
 assert.notEqual(packagePages(source,path.join(root,'d')),c);
});
test('refuses missing references and overlarge assets before creating output',t=>{
 const {root,source}=fixture(t);fs.writeFileSync(path.join(source,'game.js'),'bad loader');
 assert.throws(()=>packagePages(source,path.join(root,'bad')),/fetch reference/);assert.equal(fs.existsSync(path.join(root,'bad')),false);
 fs.copyFileSync(new URL('../web/game.js',import.meta.url),path.join(source,'game.js'));
 fs.writeFileSync(path.join(source,'wasm_exec.js'),Buffer.alloc(25*1024*1024+1));
 assert.throws(()=>packagePages(source,path.join(root,'large')),/25 MiB/);
});

test('icons resolve from both entry points and change the immutable asset hash',t=>{
 const {root,source}=fixture(t);
 const a=packagePages(source,path.join(root,'before'));
 const assets=path.join(root,'before','assets',a);
 for(const [file,base] of [['index.html',path.join(root,'before')],['game.html',assets]]) {
  for(const match of fs.readFileSync(path.join(base,file),'utf8').matchAll(/<link[^>]+href="([^"]+)"/g)) {
   assert.ok(fs.existsSync(path.join(base,match[1])),file+' missing '+match[1]);
  }
 }
 assert.deepEqual(fs.readFileSync(path.join(root,'before','favicon.ico')),fs.readFileSync(path.join(source,'favicon.ico')));
 fs.appendFileSync(path.join(source,'icon-192.png'),Buffer.from([1]));
 assert.notEqual(packagePages(source,path.join(root,'after')),a);
});
