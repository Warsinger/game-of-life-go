import {test} from 'node:test';
import assert from 'node:assert/strict';
import {validateOptions} from '../web/options.mjs';
const life = {width:80,height:80,cell:10,variance:0.5,speed:15,debug:false,lines:false};
const rover = {width:800,height:600,speed:30,debug:false};
test('all startup options survive validation, including booleans and zero',()=>{
  const values=validateOptions({...life,speed:0,debug:true,lines:true,variance:0},'life');
  assert.deepEqual(Object.fromEntries(values),{width:'80',height:'80',speed:'0',debug:'true',cell:'10',variance:'0',lines:'true'});
  assert.equal(validateOptions(rover,'rover').toString(),'width=800&height=600&speed=30&debug=false');
});
test('reject invalid numbers, limits and allocation/viewport overflow',()=>{
  for(const patch of [{width:0},{height:-1},{width:4097},{width:1.5},{speed:61},{speed:-1},{cell:0},{cell:8193},{variance:NaN},{variance:Infinity},{variance:-0.1},{variance:1.1},{variance:''},{width:''},{width:4096,height:4096,cell:1},{width:1000,cell:10}]) {
    assert.throws(()=>validateOptions({...life,...patch},'life'),undefined,JSON.stringify(patch));
  }
  for(const patch of [{width:8193},{height:0},{speed:Infinity}]) assert.throws(()=>validateOptions({...rover,...patch},'rover'));
  assert.throws(()=>validateOptions(life,'unknown'));
});
test('boundary settings and canonical values are accepted',()=>{
  assert.equal(validateOptions({...life,width:4096,height:1,cell:2,variance:1,speed:60},'life').get('width'),'4096');
  assert.equal(validateOptions({...rover,width:8192,height:8192},'rover').get('width'),'8192');
  assert.equal(validateOptions({...life,width:'080'},'life').get('width'),'80');
});
