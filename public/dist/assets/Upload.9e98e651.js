import{b5 as G,_ as J,a0 as K,a as V,b as N,k as P,c5 as Q,e as o,S as T,ad as _,B as M,E as X,I as O,a3 as v,ak as R,n as Y,bM as j,a6 as Z,b8 as E,aj as H,c3 as ee,b1 as U,c6 as te,av as ae,c7 as oe,c8 as ne,c9 as re,p as se}from"./index.2d320a5e.js";import{u as le,c as ce,d as ie}from"./index.a9c7f046.js";async function*de(a,l,p){const n=new Set;async function i(){const[m,c]=await Promise.race(n);return n.delete(m),c}for(const m of l){const c=(async()=>await p(m,l))().then(r=>[c,r]);n.add(c),n.size>=a&&(yield await i())}for(;n.size;)yield await i()}const ue={pending:"neutral",uploading:"info",backending:"info",success:"success",error:"danger"},pe=async a=>{let l=[];const p=async(n,i)=>{await new Promise((c,r)=>{const d=s=>{console.error(s),r(s)};if(n.isFile)n.file(s=>{const u=new File([s],i+s.name,{type:s.type});l.push(u),console.log(u),c({})},d);else if(n.isDirectory){const s=n.createReader(),u=()=>{s.readEntries(async w=>{for(let g=0;g<w.length;g++)await p(w[g],i+n.name+"/");c({}),w.length>0&&u()},d)};u()}})};return await p(a,""),l},me=a=>({name:a.name,path:a.webkitRelativePath?a.webkitRelativePath:a.name,size:a.size,progress:0,speed:0,status:"pending"}),ge=async(a,l,p,n=!1)=>{let i=new Date().valueOf(),m=0;const c=new FormData;c.append("file",l);const r=await G.put("/fs/form",c,{headers:{"File-Path":encodeURIComponent(a),"As-Task":n,"Content-Type":"multipart/form-data","Last-Modified":l.lastModified,Password:J()},onUploadProgress:d=>{if(d.total){const s=d.loaded/d.total*100|0;p("progress",s);const u=new Date().valueOf(),w=(u-i)/1e3;if(w>1){const y=(d.loaded-m)/w,h=(d.total-d.loaded)/y;p("speed",y),console.log(h),i=u,m=d.loaded}s===100&&p("status","backending")}}});if(r.code!==200)return new Error(r.message)},fe=async(a,l,p,n=!1)=>{let i=new Date().valueOf(),m=0;const c=await G.put("/fs/put",l,{headers:{"File-Path":encodeURIComponent(a),"As-Task":n,"Content-Type":l.type||"application/octet-stream","Last-Modified":l.lastModified,Password:J()},onUploadProgress:r=>{if(r.total){const d=r.loaded/r.total*100|0;p("progress",d);const s=new Date().valueOf(),u=(s-i)/1e3;if(u>1){const g=(r.loaded-m)/u,$=(r.total-r.loaded)/g;p("speed",g),console.log($),i=s,m=r.loaded}d===100&&p("status","backending")}}});if(c.code!==200)return new Error(c.message)},he=[{name:"Stream",upload:fe,provider:/.*/},{name:"Form",upload:ge,provider:/.*/}],we=()=>he.filter(a=>a.provider.test(K.provider)),ke=a=>{const l=V();return o(v,{w:"$full",spacing:"$1",rounded:"$lg",border:"1px solid $neutral7",alignItems:"start",p:"$2",get _hover(){return{border:`1px solid ${R()}`}},get children(){return[o(U,{css:{wordBreak:"break-all"},get children(){return a.path}}),o(_,{spacing:"$2",get children(){return[o(te,{get colorScheme(){return ue[a.status]},get children(){return l(`home.upload.${a.status}`)}}),o(U,{get children(){return[ae(()=>oe(a.speed)),"/s"]}})]}}),o(ne,{w:"$full",trackColor:"$info3",rounded:"$full",get value(){return a.progress},size:"sm",get children(){return o(re,{get color(){return R()},rounded:"$md"})}}),o(U,{color:"$danger10",get children(){return a.msg}})]}})},$e=()=>{const a=V(),{pathname:l}=N(),{refresh:p}=le(),[n,i]=P(!1),[m,c]=P(!1),[r,d]=P(!1),[s,u]=Q({uploads:[]}),w=()=>s.uploads.every(({status:e})=>["success","error"].includes(e));let g,y;const $=async e=>{if(e.length!==0){c(!0);for(const t of e){const k=me(t);u("uploads",f=>[...f,k])}for await(const t of de(3,e,q))console.log(t);p(void 0,!0)}},h=(e,t,k)=>{u("uploads",f=>f.path===e,t,k)},D=we(),[x,W]=P(D[0]),q=async e=>{const t=e.webkitRelativePath?e.webkitRelativePath:e.name;h(t,"status","uploading");const k=se(l(),t);try{const f=await x().upload(k,e,(S,C)=>{h(t,S,C)},r());f?(h(t,"status","error"),h(t,"msg",f.message)):(h(t,"status","success"),h(t,"progress",100))}catch(f){console.error(f),h(t,"status","error"),h(t,"msg",f.message)}};return o(v,{w:"$full",pb:"$2",spacing:"$2",get children(){return o(T,{get when(){return!m()},get fallback(){return[o(_,{spacing:"$2",get children(){return[o(M,{colorScheme:"accent",onClick:()=>{u("uploads",e=>e.filter(({status:t})=>!["success","error"].includes(t))),console.log(s.uploads)},get children(){return a("home.upload.clear_done")}}),o(T,{get when(){return w()},get children(){return o(M,{onClick:()=>{c(!1)},get children(){return a("home.upload.back")}})}})]}}),o(X,{get each(){return s.uploads},children:e=>o(ke,e)})]},get children(){return[o(O,{type:"file",multiple:!0,ref(e){const t=g;typeof t=="function"?t(e):g=e},display:"none",onChange:e=>{var t;$(Array.from((t=e.target.files)!=null?t:[]))}}),o(O,{type:"file",multiple:!0,webkitdirectory:!0,ref(e){const t=y;typeof t=="function"?t(e):y=e},display:"none",onChange:e=>{var t;$(Array.from((t=e.target.files)!=null?t:[]))}}),o(v,{w:"$full",justifyContent:"center",get border(){return`2px dashed ${n()?R():"$neutral8"}`},rounded:"$lg",onDragOver:e=>{e.preventDefault(),i(!0)},onDragLeave:()=>{i(!1)},onDrop:async e=>{var A,I,z,B;e.preventDefault(),e.stopPropagation(),i(!1);const t=[],k=Array.from((I=(A=e.dataTransfer)==null?void 0:A.items)!=null?I:[]),f=Array.from((B=(z=e.dataTransfer)==null?void 0:z.files)!=null?B:[]);let S=k.length;const C=[];for(let F=0;F<S;F++){const b=k[F].webkitGetAsEntry();b!=null&&b.isFile?t.push(f[F]):b!=null&&b.isDirectory&&C.push(b)}for(const F of C){const L=await pe(F);t.push(...L)}t.length===0&&Y.warning(a("home.upload.no_files_drag")),$(t)},spacing:"$4",h:"$56",get children(){return o(T,{get when(){return!n()},get fallback(){return o(j,{get children(){return a("home.upload.release")}})},get children(){return[o(j,{get children(){return a("home.upload.upload-tips")}}),o(Z,{w:"30%",get children(){return o(E,{get value(){return x().name},onChange:e=>{W(D.find(t=>t.name===e))},get options(){return D.map(e=>({label:e.name,value:e.name}))}})}}),o(_,{spacing:"$4",get children(){return[o(H,{compact:!0,size:"xl",get["aria-label"](){return a("home.upload.upload_folder")},colorScheme:"accent",get icon(){return o(ce,{})},onClick:()=>{y.click()}}),o(H,{compact:!0,size:"xl",get["aria-label"](){return a("home.upload.upload_files")},get icon(){return o(ie,{})},onClick:()=>{g.click()}})]}}),o(ee,{get checked(){return r()},onChange:()=>{d(!r())},get children(){return a("home.upload.add_as_task")}})]}})}})]}})}})};export{$e as default};
