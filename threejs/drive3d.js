// Simple browser port of the drive3D demo using Three.js
// This is a simplified recreation of the Go game's logic and visuals.

const CHUNK_SIZE = 50;
const VISIBLE_RADIUS = 2; // generate 5x5

let scene, camera, renderer;
let car, carState;
let chunks = new Map();
let collisionBoxes = [];
let lastTime = performance.now();
let fps = 0;
let fpsCounter = { last: performance.now(), frames: 0, value: 0 };
// Per-session randomness so the map differs on each page load
const SESSION_SEED = Math.floor(Math.random() * 0xFFFFFFFF);

init();

function init() {
  // Scene and renderer
  scene = new THREE.Scene();
  renderer = new THREE.WebGLRenderer({ antialias: true });
  renderer.setSize(window.innerWidth, window.innerHeight);
  document.body.appendChild(renderer.domElement);

  // Camera
  camera = new THREE.PerspectiveCamera(60, window.innerWidth / window.innerHeight, 0.1, 1000);
  camera.position.set(0, 5, -10);

  // Light
  const hemi = new THREE.HemisphereLight(0xffffff, 0x444444, 1.0);
  scene.add(hemi);
  const dir = new THREE.DirectionalLight(0xffffff, 0.8);
  dir.position.set(-3, 10, -10);
  scene.add(dir);

  // Car
  const carGeom = new THREE.BoxGeometry(1, 0.5, 2);
  const carMat = new THREE.MeshStandardMaterial({ color: 0xff3333 });
  car = new THREE.Mesh(carGeom, carMat);
  scene.add(car);
  carState = { position: new THREE.Vector3(CHUNK_SIZE / 2, 0.25, CHUNK_SIZE / 2), yaw: 0, speed: 0, steering: 0 };

  // Controls
  setupInput();

  // UI
  setupUI();

  // Generate initial chunks
  for (let i = -VISIBLE_RADIUS; i <= VISIBLE_RADIUS; i++) {
    for (let j = -VISIBLE_RADIUS; j <= VISIBLE_RADIUS; j++) {
      generateChunk(i, j);
    }
  }

  window.addEventListener('resize', onResize);
  requestAnimationFrame(loop);
}

function onResize() {
  camera.aspect = window.innerWidth / window.innerHeight;
  camera.updateProjectionMatrix();
  renderer.setSize(window.innerWidth, window.innerHeight);
}

function key(k) { return keys[k]; }
const keys = {};
function setupInput() {
  window.addEventListener('keydown', (e) => { keys[e.code] = true; });
  window.addEventListener('keyup', (e) => { keys[e.code] = false; });
}

function setupUI() {
  const menu = document.getElementById('menu');
  const menuPlay = document.getElementById('menu-play');
  const btnPlay = document.getElementById('btn-play');
  const btnSettings = document.getElementById('btn-settings');
  const hud = document.getElementById('hud');
  const settings = document.getElementById('settings');
  const btnReturn = document.getElementById('btn-return');
  const toggleFPS = document.getElementById('toggle-fps');
  const toggleSpeed = document.getElementById('toggle-speed');

  menuPlay.onclick = () => { menu.style.display = 'none'; hud.style.display = 'block'; };
  btnPlay.onclick = () => { menu.style.display = 'none'; hud.style.display = 'block'; };
  btnSettings.onclick = () => { settings.style.display = settings.style.display === 'none' ? 'block' : 'none'; };
  btnReturn.onclick = () => { settings.style.display = 'none'; menu.style.display = 'flex'; hud.style.display = 'none'; }

  // Link toggles to HUD
  const fpsElem = document.getElementById('fps');
  const speedElem = document.getElementById('speed');
  toggleFPS.onchange = () => { fpsElem.style.display = toggleFPS.checked ? 'block' : 'none'; };
  toggleSpeed.onchange = () => { speedElem.style.display = toggleSpeed.checked ? 'block' : 'none'; };
}

function loop(now) {
  const dt = Math.min((now - lastTime) / 1000, 0.05);
  lastTime = now;

  update(dt);
  render();

  // fps calc
  fpsCounter.frames++;
  if (now - fpsCounter.last >= 500) {
    fpsCounter.value = Math.round((fpsCounter.frames * 1000) / (now - fpsCounter.last));
    fpsCounter.last = now;
    fpsCounter.frames = 0;
    const fpsElem = document.getElementById('fps');
    if (fpsElem) fpsElem.textContent = `FPS: ${fpsCounter.value}`;
  }

  // update speed HUD and debug counts
  const speedElem = document.getElementById('speed');
  const toggleSpeed = document.getElementById('toggle-speed');
  if (speedElem) {
    if (toggleSpeed && toggleSpeed.checked) {
      // convert internal speed to a readable km/h-like value (tunable)
      // user requested speed divided by 10 -> previous factor 50 becomes 5
      const kmh = Math.round(Math.abs(carState.speed) * 5);
      speedElem.textContent = `Speed: ${kmh} km/h`;
    }
  }
  const chunksElem = document.getElementById('dbg-chunks');
  const collElem = document.getElementById('dbg-coll');
  if (chunksElem) chunksElem.textContent = `Chunks: ${chunks.size}`;
  if (collElem) collElem.textContent = `Collisions: ${collisionBoxes.length}`;

  requestAnimationFrame(loop);
}

function update(dt) {
  // Controls: ArrowUp: accelerate, ArrowDown: brake, ArrowLeft/Right: steer
  const accelBase = 5.0;
  if (keys['ArrowUp']) {
    carState.speed += accelBase * dt;
  } else if (keys['ArrowDown']) {
    carState.speed -= accelBase * dt;
  } else {
    carState.speed *= 0.995;
    if (Math.abs(carState.speed) < 0.01) carState.speed = 0;
  }

  if (keys['ArrowLeft']) { carState.steering = Math.max(carState.steering - 2 * dt, -1); }
  else if (keys['ArrowRight']) { carState.steering = Math.min(carState.steering + 2 * dt, 1); }
  else { carState.steering *= 0.9; }

  carState.yaw += carState.steering * dt;

  const forward = new THREE.Vector3(Math.cos(carState.yaw), 0, Math.sin(carState.yaw));
  carState.position.addScaledVector(forward, carState.speed * dt);

  // Collision: simple bounding-sphere against collisionBoxes
  for (const b of collisionBoxes) {
    const dx = carState.position.x - b.center.x;
    const dz = carState.position.z - b.center.z;
    const dist2 = dx * dx + dz * dz;
    const minDist = 1 + b.radius;
    if (dist2 < minDist * minDist) {
      // push back to avoid overlap (add a negative scaled forward vector)
      carState.position.addScaledVector(forward, -0.25);
      carState.speed = 0;
      break;
    }
  }

  // Update car mesh
  car.position.copy(carState.position);
  car.rotation.y = -carState.yaw;

  // Update camera
  const camTarget = carState.position.clone();
  const camPos = carState.position.clone().add(new THREE.Vector3(-Math.cos(carState.yaw) * 5, 3, -Math.sin(carState.yaw) * 5));
  camera.position.lerp(camPos, 0.1);
  camera.lookAt(camTarget);

  // Update visible chunks around player
  const playerChunkX = Math.floor(carState.position.x / CHUNK_SIZE);
  const playerChunkZ = Math.floor(carState.position.z / CHUNK_SIZE);
  for (let i = playerChunkX - VISIBLE_RADIUS; i <= playerChunkX + VISIBLE_RADIUS; i++) {
    for (let j = playerChunkZ - VISIBLE_RADIUS; j <= playerChunkZ + VISIBLE_RADIUS; j++) {
      const key = `${i},${j}`;
      if (!chunks.has(key)) {
        try { generateChunk(i, j); }
        catch (e) { console.error('generateChunk failed', e); }
      }
    }
  }

  // prune chunks outside of visible radius to avoid memory/logic blowups
  pruneChunks(playerChunkX, playerChunkZ);
}

function render() { renderer.render(scene, camera); }

// Simple chunk generator inspired by the Go code
function generateChunk(i, j) {
  const key = `${i},${j}`;
  // Simple seeded randomness using a hash
  // mix the session seed with chunk coordinates so different page loads have different worlds
  const seed = hash(`${SESSION_SEED}:${i},${j}`);
  const rnd = mulberry32(seed);

  // Determine chunk type
  const chunkType = Math.floor(rnd()*6);

  const posX = i * CHUNK_SIZE;
  const posZ = j * CHUNK_SIZE;

  // Ground plane
  const color = (function(){
    switch(chunkType){
      case 1: return 0x555555; // City
      case 2: return 0x888888; // Commercial
      case 3: return 0xE0C07A; // Desert
      case 4: return 0x347A1F; // Forest
      case 5: return 0xFFFFFF; // Snow
      default: return 0x999999; // Highway
    }
  })();
  const groundGeom = new THREE.PlaneGeometry(CHUNK_SIZE, CHUNK_SIZE);
  const groundMat = new THREE.MeshStandardMaterial({ color: color, side: THREE.DoubleSide });
  const ground = new THREE.Mesh(groundGeom, groundMat);
  ground.rotation.x = -Math.PI/2;
  ground.position.set(posX + CHUNK_SIZE/2, 0, posZ + CHUNK_SIZE/2);
  scene.add(ground);

  // Road (plus shape)
  const roadColor = (chunkType===4)?0x693f24: (chunkType===5?0x99ddff:0x333333);
  const roadW = 5;
  const roadMat = new THREE.MeshStandardMaterial({ color: roadColor });
  const roadH = new THREE.Mesh(new THREE.PlaneGeometry(CHUNK_SIZE, roadW), roadMat);
  roadH.rotation.x = -Math.PI/2;
  roadH.position.set(posX + CHUNK_SIZE/2, 0.01, posZ + CHUNK_SIZE/2);
  scene.add(roadH);
  const roadV = new THREE.Mesh(new THREE.PlaneGeometry(roadW, CHUNK_SIZE), roadMat);
  roadV.rotation.x = -Math.PI/2;
  roadV.position.set(posX + CHUNK_SIZE/2, 0.01, posZ + CHUNK_SIZE/2);
  scene.add(roadV);

  // Add objects based on type and track chunk-local collision boxes
  const objects = [];
  const chunkCollisionBoxes = [];
  const objCount = (chunkType===4?20:(chunkType===3?10:(chunkType===1?5:(chunkType===2?3:2))));
  let placed = 0;
  let attempts = 0;
  while (placed < objCount && attempts < objCount * 6) {
    attempts++;
    const bx = posX + rnd()*CHUNK_SIZE;
    const bz = posZ + rnd()*CHUNK_SIZE;
    // Skip road center (+ shape)
    const rx = Math.abs((bx - posX) - CHUNK_SIZE/2);
    const rz = Math.abs((bz - posZ) - CHUNK_SIZE/2);
    if (rx <= roadW || rz <= roadW) { continue; }
    const h = (chunkType===1)?(10 + rnd()*40):(chunkType===2?10: (chunkType===4?8+ rnd()*6: 3+ rnd()*5));
    const geom = new THREE.BoxGeometry( Math.max(1, h/3), h, Math.max(1, h/3));
    const mat = new THREE.MeshStandardMaterial({ color: (chunkType===3?0x8B4513:(chunkType===4?0x2E8B57:(chunkType===5?0xFFFFFF:0x8888ff))) });
    const m = new THREE.Mesh(geom, mat);
    m.position.set(bx, h/2, bz);
    scene.add(m);
    objects.push(m);
    // collision box (keep a reference so we can remove it when chunk is pruned)
    const box = { center: { x: bx, z: bz }, radius: Math.max(1, Math.max(h/3,h/3)) };
    collisionBoxes.push(box);
    chunkCollisionBoxes.push(box);
    placed++;
  }

  chunks.set(key, { ground, roadH, roadV, objects, collisionBoxes: chunkCollisionBoxes });
}

// Remove chunks outside the visible radius and dispose of their geometry/materials
function pruneChunks(centerX, centerZ){
  for (const [key, data] of chunks.entries()){
    const parts = key.split(',').map(Number);
    const i = parts[0], j = parts[1];
    if (Math.abs(i - centerX) > VISIBLE_RADIUS || Math.abs(j - centerZ) > VISIBLE_RADIUS){
      // remove meshes
      const removeMesh = (m) => {
        try {
          scene.remove(m);
          if (m.geometry) m.geometry.dispose();
          if (m.material) {
            if (Array.isArray(m.material)) m.material.forEach(mat=>mat.dispose());
            else m.material.dispose();
          }
        } catch(e){ /* defensive */ }
      };

      if (data.ground) removeMesh(data.ground);
      if (data.roadH) removeMesh(data.roadH);
      if (data.roadV) removeMesh(data.roadV);
      if (data.objects) data.objects.forEach(removeMesh);

      // remove collision boxes that were registered for this chunk
      if (data.collisionBoxes && data.collisionBoxes.length){
        for (const cb of data.collisionBoxes){
          const idx = collisionBoxes.indexOf(cb);
          if (idx !== -1) collisionBoxes.splice(idx, 1);
        }
      }

      chunks.delete(key);
    }
  }
}

// tiny seeded hash
function hash(s){
  let h = 2166136261 >>> 0;
  for (let i=0;i<s.length;i++) h = Math.imul(h ^ s.charCodeAt(i), 16777619);
  return h >>> 0;
}

function mulberry32(a) {
  return function() {
    var t = a += 0x6D2B79F5;
    t = Math.imul(t ^ t >>> 15, t | 1);
    t ^= t + Math.imul(t ^ t >>> 7, t | 61);
    return ((t ^ t >>> 14) >>> 0) / 4294967296;
  }
}
