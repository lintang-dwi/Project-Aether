/**
 * Project AETHER — Futuristic Dark Theme
 * Graph particle bg · Node acronym visualizer · Counter animations
 */

document.addEventListener('DOMContentLoaded', () => {
  initParticleNetwork();
  initScrollProgress();
  renderContent();
  initCounters();
  initAcronymGraph();
  setupScrollAnimations();
});

// ─── SCROLL PROGRESS ───
function initScrollProgress() {
  const bar = document.getElementById('scroll-progress');
  if (!bar) return;
  window.addEventListener('scroll', () => {
    const h = document.documentElement.scrollHeight - window.innerHeight;
    bar.style.width = h > 0 ? (window.scrollY / h * 100) + '%' : '0%';
  }, { passive: true });
}

// ─── PARTICLE NETWORK CANVAS ───
function initParticleNetwork() {
  const c = document.getElementById('particle-canvas');
  if (!c) return;
  const ctx = c.getContext('2d');
  let w, h, particles = [];

  const res = () => { w = c.width = window.innerWidth; h = c.height = window.innerHeight; };
  res(); window.addEventListener('resize', res);

  const P = 60;
  for (let i = 0; i < P; i++) {
    particles.push({
      x: Math.random() * w, y: Math.random() * h,
      vx: (Math.random() - 0.5) * 0.2, vy: (Math.random() - 0.5) * 0.2,
      r: Math.random() * 1.5 + 0.5,
    });
  }

  let frame = 0;

  function draw() {
    ctx.clearRect(0, 0, w, h);
    frame++;

    // Draw connections
    for (let i = 0; i < particles.length; i++) {
      for (let j = i + 1; j < particles.length; j++) {
        const dx = particles[i].x - particles[j].x;
        const dy = particles[i].y - particles[j].y;
        const dist = Math.sqrt(dx * dx + dy * dy);
        if (dist < 160) {
          const alpha = 0.06 * (1 - dist / 160);
          ctx.beginPath();
          ctx.moveTo(particles[i].x, particles[i].y);
          ctx.lineTo(particles[j].x, particles[j].y);
          ctx.strokeStyle = `rgba(99, 102, 241, ${alpha})`;
          ctx.lineWidth = 0.5;
          ctx.stroke();
        }
      }
    }

    // Draw nodes
    for (const p of particles) {
      p.x += p.vx; p.y += p.vy;
      if (p.x < 0) p.x = w; if (p.x > w) p.x = 0;
      if (p.y < 0) p.y = h; if (p.y > h) p.y = 0;

      ctx.beginPath();
      ctx.arc(p.x, p.y, p.r, 0, Math.PI * 2);
      ctx.fillStyle = 'rgba(99, 102, 241, 0.15)';
      ctx.fill();
    }

    requestAnimationFrame(draw);
  }
  draw();
}

// ─── COUNTER ANIMATION ───
function initCounters() {
  const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (!entry.isIntersecting) return;
      const el = entry.target;
      const target = parseInt(el.dataset.count);
      if (!target || el.dataset.done) return;
      el.dataset.done = '1';
      animateCounter(el, target);
    });
  }, { threshold: 0.3 });

  document.querySelectorAll('[data-count]').forEach(el => observer.observe(el));
}

function animateCounter(el, target) {
  const dur = 1200;
  const start = performance.now();

  function tick(now) {
    const t = Math.min((now - start) / dur, 1);
    const eased = 1 - Math.pow(1 - t, 3); // ease-out cubic
    el.textContent = Math.floor(eased * target);
    if (t < 1) requestAnimationFrame(tick);
    else el.textContent = target;
  }
  requestAnimationFrame(tick);
}

// ─── ACRONYM GRAPH ───
function initAcronymGraph() {
  // Big word reveal
  const bwObserver = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (!entry.isIntersecting) return;
      const letters = entry.target.querySelectorAll('.bw-letter');
      letters.forEach((l, i) => {
        setTimeout(() => l.classList.add('lit'), i * 120);
      });
      bwObserver.unobserve(entry.target);
    });
  }, { threshold: 0.5 });

  const bw = document.querySelector('.big-word');
  if (bw) bwObserver.observe(bw);

  // Node click/hover
  let activeNode = null;
  const nodes = document.querySelectorAll('.acr-node');
  nodes.forEach(node => {
    // Hover -> highlight lines
    node.addEventListener('mouseenter', () => {
      const idx = parseInt(node.dataset.idx);
      document.querySelectorAll('.acr-line').forEach(line => {
        if (parseInt(line.dataset.from) === idx || parseInt(line.dataset.to) === idx) {
          line.classList.add('highlight');
        }
      });
    });
    node.addEventListener('mouseleave', () => {
      document.querySelectorAll('.acr-line').forEach(line => line.classList.remove('highlight'));
    });

    // Click -> expand card
    node.addEventListener('click', () => {
      const expandId = node.dataset.expand;
      const expanded = document.getElementById(expandId);
      if (!expanded) return;

      if (activeNode && activeNode !== node) {
        activeNode.classList.remove('active');
        const prevExpand = document.getElementById(activeNode.dataset.expand);
        if (prevExpand) prevExpand.classList.remove('show');
      }

      if (node.classList.contains('active')) {
        node.classList.remove('active');
        expanded.classList.remove('show');
        activeNode = null;
      } else {
        node.classList.add('active');
        expanded.classList.add('show');
        activeNode = node;
      }
    });
  });
}

// ─── RENDER CONTENT ───
function renderContent() {
  const acronyms = [
    { letter: 'A', word: 'Graph-Native', desc: 'Dibangun di atas knowledge graph — relasi adalah yang utama, file hanyalah representasi turunan.' },
    { letter: 'U', word: 'niversal', desc: 'Agnostik terhadap bahasa pemrograman. Runtime yang mampu memahami basis kode apa pun.' },
    { letter: 'T', word: 'hinking', desc: 'Lapisan penalaran berbantuan AI di atas inti sistem yang deterministik dan dapat diprediksi.' },
    { letter: 'O', word: 'rchestrated', desc: 'Koordinasi berbasis event secara real-time di seluruh subsistem runtime.' },
    { letter: 'M', word: 'odular', desc: 'Arsitektur microkernel — perkuat sistem tanpa menyentuh inti.' },
    { letter: 'A', word: 'utonomous', desc: 'Runtime yang berevolusi bersama proyek, belajar dari setiap perubahan.' },
    { letter: 'T', word: 'raceable', desc: 'Setiap keputusan terekam. Setiap perubahan memiliki jejak yang dapat ditelusuri.' },
    { letter: 'I', word: 'ntegrated', desc: 'Satu runtime untuk analisis, perencanaan, eksekusi, dan peninjauan kode.' },
    { letter: 'C', word: 'onsistent', desc: 'Knowledge Model selalu sinkron dengan source code, tidak pernah tertinggal.' },
    { letter: 'E', word: 'xtensible', desc: 'Sistem plugin untuk tools, provider AI, dan bahasa pemrograman baru.' },
    { letter: 'R', word: 'untime-First', desc: 'Desktop, CLI, API, ekstensi editor — semua berbasis runtime yang identik.' },
  ];

  // Build big word letters
  const bigWord = 'AUTOMATIC'.split('').map(l => `<span class="bw-letter">${l}</span>`).join('');

  // Build nodes + SVG lines
  let nodesHtml = '';
  let linesHtml = '';
  const gap = 100; // percentage-based spacing

  acronyms.forEach((a, i) => {
    const x = 5 + (i / (acronyms.length - 1)) * 90;
    nodesHtml += `
      <div class="acr-node" data-idx="${i}" data-expand="exp-${i}">
        <div class="acr-node-circle">${a.letter}</div>
        <div class="acr-node-label">${a.word}</div>
      </div>`;
    if (i < acronyms.length - 1) {
      linesHtml += `<line class="acr-line" data-from="${i}" data-to="${i+1}" x1="${x}%" y1="26" x2="${5 + ((i+1) / (acronyms.length-1)) * 90}%" y2="26"/>`;
    }
  });

  // Expanded cards
  let expandedHtml = '';
  acronyms.forEach((a, i) => {
    expandedHtml += `
      <div id="exp-${i}" class="acr-expanded">
        <div class="exp-word">${a.word}</div>
        <div class="exp-desc">${a.desc}</div>
      </div>`;
  });

  const content = `<!---●≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡ HERO ≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡●--->
      <section class="hero">
        <div class="hero-badge"><span class="live-dot"></span> GRAPHOS · V0.1.0 · FOUNDATION STAGE</div>
        <h1 class="hero-title">
          <span class="hero-letter" style="--i:0">A</span>
          <span class="hero-letter" style="--i:1">E</span>
          <span class="hero-letter" style="--i:2">T</span>
          <span class="hero-letter" style="--i:3">H</span>
          <span class="hero-letter" style="--i:4">E</span>
          <span class="hero-letter" style="--i:5">R</span>
        </h1>
        <p class="hero-subtitle">
          <span class="hl">A</span> Graph-Native · <span class="hl2">A</span>utonomous ·
          <span class="hl3">S</span>oftware · <span class="hl">E</span>ngineering · <span class="hl2">R</span>untime
        </p>
        <p class="hero-desc">
          Runtime generasi baru yang memahami perangkat lunak<br>
          sebagai <strong style="color:var(--accent3)">pengetahuan terstruktur</strong> — bukan sekadar file.
        </p>
        <div class="hero-actions">
          <a class="btn-primary" href="#nama"><span>Jelajahi Visi ↓</span></a>
        </div>
      </section>

      <!---●≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡ STATS ≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡●--->
      <div class="stats-grid">
        <div class="stat-card"><div class="stat-icon">📄</div><div class="stat-num" data-count="65">0</div><div class="stat-label">Dokumen</div></div>
        <div class="stat-card"><div class="stat-icon">⚖️</div><div class="stat-num" data-count="55">0</div><div class="stat-label">Keputusan Arsitektur</div></div>
        <div class="stat-card"><div class="stat-icon">🧩</div><div class="stat-num" data-count="30">0</div><div class="stat-label">Komponen Sistem</div></div>
        <div class="stat-card"><div class="stat-icon">🚀</div><div class="stat-num" data-count="5">0</div><div class="stat-label">Fase Pengembangan</div></div>
      </div>

      <!---●≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡ MAKNA NAMA ≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡●--->
      <section id="nama" class="section">
        <h2 class="section-title">Makna Dibalik Nama</h2>
        <p class="acronym-intro">
          <strong>AETHER</strong> adalah akronim berlapis. Setiap huruf mewakili kata kunci arsitektur —
          dan bersama-sama mereka membentuk kata <em>AUTOMATIC</em>: sifat inti dari runtime ini.
          <br>Klik atau sentuh setiap <strong>node</strong> untuk melihat arti lengkapnya.
        </p>

        <div class="big-word">${bigWord}</div>

        <div class="acronym-scope">
          <svg class="acronym-lines-svg" viewBox="0 0 100 52" preserveAspectRatio="none">${linesHtml}</svg>
          <div class="acronym-nodes">${nodesHtml}</div>
        </div>

        <div style="display:flex; flex-wrap:wrap; justify-content:center; gap:8px; margin-top:16px;">
          ${expandedHtml}
        </div>
      </section>

      <!---●≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡ FILOSOFI ≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡●--->
      <section id="filosofi" class="section">
        <h2 class="section-title">Filosofi Inti</h2>
        <div class="phil-grid">
          <div class="phil-card">
            <div class="phil-icon-wrap"><svg viewBox="0 0 24 24"><path d="M12 2l2 7h7l-5.5 4 2 7L12 16l-5.5 4 2-7L3 9h7z"/></svg></div>
            <h3>Pengetahuan &gt; Kode</h3>
            <p>Source code adalah <em>representasi</em> dari pengetahuan, bukan kebenaran mutlak. Runtime bekerja dari Knowledge Model, bukan file mentah.</p>
          </div>
          <div class="phil-card">
            <div class="phil-icon-wrap"><svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg></div>
            <h3>Runtime adalah Sistem</h3>
            <p>Desktop, CLI, API, ekstensi — semuanya hanyalah antarmuka menuju runtime yang sama. Tidak ada logika bisnis di luar runtime.</p>
          </div>
          <div class="phil-card">
            <div class="phil-icon-wrap"><svg viewBox="0 0 24 24"><path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"/><path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"/></svg></div>
            <h3>Kecerdasan dari Struktur</h3>
            <p>AI tidak membaca file mentah. AI menerima konteks terstruktur yang sudah disiapkan oleh runtime — lebih akurat, lebih efisien.</p>
          </div>
          <div class="phil-card">
            <div class="phil-icon-wrap"><svg viewBox="0 0 24 24"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg></div>
            <h3>Deterministik adalah Dasar</h3>
            <p>Semua yang bisa dihitung secara algoritmik dilakukan runtime. AI hanya untuk penalaran dan kreativitas — bukan keputusan kritis.</p>
          </div>
        </div>
      </section>

      <!---●≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡ VISI ≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡●--->
      <section id="visi" class="section">
        <h2 class="section-title">Visi &amp; Misi</h2>
        <div class="vision-block">
          <div class="vision-card primary">
            <h3>Visi</h3>
            <p>Membangun <em>runtime</em> yang memahami proyek sebagai <strong style="color:var(--accent3)">model pengetahuan hidup</strong> — manusia dan AI bekerja bersama secara konsisten, dalam skala besar, dengan pertanggungjawaban penuh.</p>
          </div>
          <div class="vision-card">
            <h3>Misi Utama</h3>
            <ul>
              <li><svg viewBox="0 0 24 24"><polyline points="20 6 9 17 4 12"/></svg> Menggeser paradigma dari <strong style="color:var(--text)">file-centric</strong> ke <strong style="color:var(--text)">knowledge-centric</strong></li>
              <li><svg viewBox="0 0 24 24"><polyline points="20 6 9 17 4 12"/></svg> Memutus ketergantungan pada <em>context window</em> LLM yang terbatas</li>
              <li><svg viewBox="0 0 24 24"><polyline points="20 6 9 17 4 12"/></svg> Fondasi untuk <strong style="color:var(--text)">Autonomous Software Engineering</strong> yang nyata</li>
              <li><svg viewBox="0 0 24 24"><polyline points="20 6 9 17 4 12"/></svg> Arsitektur modular, <strong style="color:var(--text)">provider-independent</strong>, dan extensible</li>
            </ul>
          </div>
        </div>
      </section>

      <!---●≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡ PRINSIP ≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡●--->
      <section id="prinsip" class="section">
        <h2 class="section-title">10 Prinsip Arsitektur</h2>
        <div class="principles-grid">
          <div class="principle-card"><div class="pcode">AP-01</div><div class="pname">Arsitektur Microkernel</div></div>
          <div class="principle-card"><div class="pcode">AP-02</div><div class="pname">Isolasi Layanan</div></div>
          <div class="principle-card"><div class="pcode">AP-03</div><div class="pname">Event First</div></div>
          <div class="principle-card"><div class="pcode">AP-04</div><div class="pname">API Driven</div></div>
          <div class="principle-card"><div class="pcode">AP-05</div><div class="pname">Stateless Services</div></div>
          <div class="principle-card"><div class="pcode">AP-06</div><div class="pname">Dependency Inversion</div></div>
          <div class="principle-card"><div class="pcode">AP-07</div><div class="pname">Eksplisit</div></div>
          <div class="principle-card"><div class="pcode">AP-08</div><div class="pname">Transaksional</div></div>
          <div class="principle-card"><div class="pcode">AP-09</div><div class="pname">Version Everything</div></div>
          <div class="principle-card"><div class="pcode">AP-10</div><div class="pname">Observable by Default</div></div>
        </div>
      </section>

      <!---●≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡ STACK ≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡●--->
      <section id="stack" class="section">
        <h2 class="section-title">Tumpukan Teknologi</h2>
        <div class="stack-grid">
          <div class="stack-card"><div class="stack-icon">🐹</div><div class="stack-label">Bahasa Inti</div><div class="stack-value">Go</div></div>
          <div class="stack-card"><div class="stack-icon">🖥️</div><div class="stack-label">Desktop</div><div class="stack-value">Wails</div></div>
          <div class="stack-card"><div class="stack-icon">🗄️</div><div class="stack-label">Penyimpanan</div><div class="stack-value">SQLite + Graph</div></div>
          <div class="stack-card"><div class="stack-icon">🌳</div><div class="stack-label">Parser</div><div class="stack-value">Tree-sitter</div></div>
          <div class="stack-card"><div class="stack-icon">🤖</div><div class="stack-label">AI Provider</div><div class="stack-value">Cloud · Lokal · Enterprise</div></div>
          <div class="stack-card"><div class="stack-icon">🔌</div><div class="stack-label">Event Bus</div><div class="stack-value">Internal Pub/Sub</div></div>
          <div class="stack-card"><div class="stack-icon">🔷</div><div class="stack-label">Arsitektur</div><div class="stack-value">Graph-Native</div></div>
          <div class="stack-card"><div class="stack-icon">📜</div><div class="stack-label">Lisensi</div><div class="stack-value">TBD</div></div>
        </div>
      </section>

      <!---●≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡ ROADMAP ≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡●--->
      <section id="roadmap" class="section">
        <h2 class="section-title">Peta Pengembangan</h2>
        <div class="roadmap">
          <div class="roadmap-phase"><div class="phase-dot"></div><div class="phase-content"><h3>Fase 1 · Fondasi Runtime</h3><p>Runtime Coordinator · Event System · Storage Layer · Core Lifecycle</p></div></div>
          <div class="roadmap-line"></div>
          <div class="roadmap-phase"><div class="phase-dot"></div><div class="phase-content"><h3>Fase 2 · Mesin Pemahaman</h3><p>Workspace Scanner · Parser Pipeline · Knowledge Model · Knowledge Graph</p></div></div>
          <div class="roadmap-line"></div>
          <div class="roadmap-phase"><div class="phase-dot"></div><div class="phase-content"><h3>Fase 3 · Lapisan Kecerdasan</h3><p>Context Engine · AI Provider Abstraction · Planning System · Task Engine</p></div></div>
          <div class="roadmap-line"></div>
          <div class="roadmap-phase"><div class="phase-dot"></div><div class="phase-content"><h3>Fase 4 · Otomatisasi Terkendali</h3><p>Action Processor · Validation Engine · Git Engine · Memory Engine</p></div></div>
          <div class="roadmap-line"></div>
          <div class="roadmap-phase"><div class="phase-dot"></div><div class="phase-content"><h3>Fase 5 · Ekosistem &amp; Skalabilitas</h3><p>Plugin System · Extension API · Performance Optimization · Collaborative Runtime</p></div></div>
        </div>
      </section>

      <!---●≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡ FOOTER ≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡●--->
      <footer class="footer">
        <p><strong>Project Aether</strong> — <strong>A</strong> Graph-Native <strong>A</strong>utonomous <strong>S</strong>oftware <strong>E</strong>ngineering <strong>R</strong>untime</p>
        <p class="footer-meta">GraphOS · Tahap Draft · Dibangun dengan pengetahuan terstruktur</p>
      </footer>`;

  document.getElementById('content').innerHTML = content;
}

// ─── SCROLL REVEAL ───
function setupScrollAnimations() {
  const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) entry.target.classList.add('visible');
    });
  }, { threshold: 0.1, rootMargin: '0px 0px -40px 0px' });

  document.querySelectorAll('.section, .hero, .footer, .stats-grid').forEach(el => {
    el.classList.add('fade-section');
    observer.observe(el);
  });
}
