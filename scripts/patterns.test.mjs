import {test} from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import {parseGrid, startingParams, PatternStore, STORE_KEY} from '../web/patterns.mjs';
const library=JSON.parse(fs.readFileSync(new URL('../internal/patterns/library.json',import.meta.url)));
const params=()=>new URLSearchParams({width:'20',height:'20'});

test('portable library, custom grid limits and canonical query parameters',()=>{
 for(const p of library){parseGrid(p.grid);assert.equal(startingParams(params(),{pattern:p.id},library).get('pattern'),p.id);}
 assert.equal(startingParams(params(),{},library).get('pattern'),'random');
 const query=startingParams(params(),{pattern:'custom',seed:'\r\n.O.\r\n..O\r\nOOO\r\n'},library);
 assert.equal(query.get('seed'),'.O.\n..O\nOOO');assert.equal(new URLSearchParams(query.toString()).get('seed'),query.get('seed'));
 for(const text of ['', 'O.\nO', 'O \nOO','O\n\nO','o','😀','O'.repeat(65),'O\n'.repeat(65),'O'.repeat(8193)]) assert.throws(()=>parseGrid(text));
 assert.throws(()=>startingParams(new URLSearchParams({width:2,height:2}),{pattern:'glider'},library),/increase/);
 assert.throws(()=>startingParams(params(),{pattern:'missing'},library),/available/);
});

test('named library persists exact starting grids, overwrites names and deletes safely',()=>{
 const data=new Map();const storage={getItem:k=>data.get(k)||null,setItem:(k,v)=>data.set(k,v)};
 const store=new PatternStore(()=>storage);
 assert.equal(store.save('<my glider>',library.find(p=>p.id==='glider').grid),true);
 assert.deepEqual(new PatternStore(()=>storage).list(),store.list());
 store.save('<my glider>','OO\nOO');assert.equal(store.list().length,1);assert.equal(store.list()[0].grid,'OO\nOO');
 const copy=store.list();copy[0].grid='bad';assert.notEqual(store.list()[0].grid,'bad');
 assert.throws(()=>store.save('', 'O'));assert.throws(()=>store.save('x'.repeat(41),'O'));
 for(let i=0;i<19;i++)store.save('pattern '+i,'O');assert.throws(()=>store.save('too many','O'),/20/);
 store.remove('<my glider>');assert.equal(new PatternStore(()=>storage).list().length,19);
 assert.ok(data.has(STORE_KEY));
});

test('storage faults/corruption retain a usable bounded session library',()=>{
 for(const getter of [()=>{throw Error('blocked')},()=>({getItem(){throw Error('get blocked')},setItem(){throw Error('blocked')}}),()=>({getItem(){return null},setItem(){throw Error('quota')}})]){
  const store=new PatternStore(getter);assert.equal(store.save('local','O'),false);assert.deepEqual(store.list(),[{name:'local',grid:'O'}]);assert.equal(store.remove('local'),false);assert.deepEqual(store.list(),[]);
 }
 for(const raw of ['bad','{}','null',JSON.stringify([{name:'x',grid:'bad'}]),' '.repeat(200001)]){
  const store=new PatternStore(()=>({getItem:()=>raw,setItem(){throw Error('blocked')}}));assert.equal(store.persistent,false);assert.deepEqual(store.list(),[]);
 }
});
