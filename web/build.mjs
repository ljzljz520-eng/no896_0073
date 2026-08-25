import { mkdir, writeFile } from 'node:fs/promises';
await mkdir('dist', {recursive:true});
await writeFile('dist/index.html', '<!doctype html><html><body><main><h1>Online Education Game Rank</h1><p>Browse learning games by subject, age, and device.</p></main></body></html>');
