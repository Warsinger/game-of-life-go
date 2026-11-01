import {test} from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import {setupPatterns} from '../web/pattern-ui.mjs';
import {PatternStore} from '../web/patterns.mjs';
const library=JSON.parse(fs.readFileSync(new URL('../internal/patterns/library.json',import.meta.url)));
class Element {
 constructor(value=''){this.value=value;this.children=[];this.listeners={};this.textContent='';}
 addEventListener(name,fn){this.listeners[name]=fn;}
 append(node){this.children.push(node);if(!this.value)this.value=node.value;}
 replaceChildren(){this.children=[];this.value='';}
 emit(name){return this.listeners[name]?.();}
}
function fixture(){
 const nodes=new Map(['pattern','pattern-editor','seed','pattern-description','pattern-status','saved-pattern','pattern-name','load-pattern','delete-pattern','save-pattern','load-paused'].map(id=>[id,new Element()]));
 nodes.get('pattern').value='random';nodes.get('seed').value='.O.\n..O\nOOO';
 const form={elements:{speed:{value:'15'}},submits:0,requestSubmit(){this.submits++;}};
 const document={querySelector:()=>form,getElementById:id=>nodes.get(id),createElement:()=>new Element()};
 const data=new Map();const store=new PatternStore(()=>({getItem:k=>data.get(k)||null,setItem:(k,v)=>data.set(k,v)}));
 return {nodes,form,document,store};
}

test('preset editing, named save/load/delete and Load paused are connected to the real form',async()=>{
 const {nodes:n,form,document,store}=fixture();
 const loaded=await setupPatterns(document,async()=>({ok:true,json:async()=>library}),store);
 assert.equal(loaded.length,6);assert.equal(n.get('pattern-editor').hidden,true);assert.equal(n.get('load-pattern').disabled,true);
 n.get('pattern').value='glider';n.get('pattern').emit('change');assert.equal(n.get('seed').value,'.O.\n..O\nOOO');assert.equal(n.get('pattern-editor').hidden,false);
 n.get('seed').value='OO\nOO';n.get('seed').emit('input');assert.equal(n.get('pattern').value,'custom');
 n.get('pattern-name').value='<My block>';n.get('save-pattern').emit('click');assert.equal(store.list()[0].name,'<My block>');assert.equal(n.get('saved-pattern').children[0].textContent,'<My block>');
 n.get('seed').value='O';n.get('load-pattern').emit('click');assert.equal(n.get('seed').value,'OO\nOO');assert.equal(n.get('pattern').value,'custom');
 n.get('load-paused').emit('click');assert.equal(form.elements.speed.value,'0');assert.equal(form.submits,1);
 n.get('delete-pattern').emit('click');assert.deepEqual(store.list(),[]);assert.equal(n.get('load-pattern').disabled,true);
});

test('unavailable library and blocked storage leave Random/custom usable and report session saves',async()=>{
 const {nodes:n,document}=fixture();const store=new PatternStore(()=>{throw Error('blocked')});
 assert.deepEqual(await setupPatterns(document,async()=>{throw Error('offline')},store),[]);
 assert.match(n.get('pattern-status').textContent,/could not load/);
 n.get('pattern').value='custom';n.get('pattern').emit('change');n.get('pattern-name').value='local';n.get('save-pattern').emit('click');
 assert.match(n.get('pattern-status').textContent,/session only/);assert.equal(store.list().length,1);
});
