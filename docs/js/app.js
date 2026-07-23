/**
 * Project Aether — Redesign v2
 * Cosmic theme · Animated orbs · Scroll progress · Stats
 */

document.addEventListener('DOMContentLoaded', () => {
  initOrbCanvas();
  initScrollProgress();
  renderContent();
  setupScrollAnimations();
});

// ─── SCROLL PROGRESS BAR ───
function initScrollProgress() {
  const bar = document.getElementById('scroll-progress');
  window.addEventListener('scroll', () => {
    const h = document.documentElement.scrollHeight - window.innerHeight;
    bar.style.width = h > 0 ? (window.scrollY / h * 100) + '%' : '0%';
  });
}

// ─── ORB BACKGROUND ───
function initOrbCanvas() {
  const canvas = document.getElementById('orb-canvas');
  if (!canvas) return;
  const ctx = canvas.getContext('2d');
  let w, h;

  function resize() {
    w = canvas.width = window.innerWidth;
    h = canvas.height = window.innerHeight;
  }
  resize();
  window.addEventListener('resize', resize);

  // Floating orbs
  const orbs = [
    { x: 0.2, y: 0.3, r: 220, dx: 0.0003, dy: 0.0002, c: '99,102,241', a: 0.06 },
    { x: 0.8, y: 0.6, r: 300, dx: -0.0002, dy: 0.0004, c: '6,182,212', a: 0.05 },
    { x: 0.5, y: 0.8, r: 180, dx: 0.0004, dy: -0.0003, c: '16,185,129', a: 0.04 },
    { x: 0.7, y: 0.2, r: 160, dx: -0.00035, dy: -0.0002, c: '244,114,182', a: 0.035 },
  ];

  function draw() {
    ctx.clearRect(0, 0, w, h);
    for (const o of orbs) {
      o.x += o.dx * w;
      o.y += o.dy * h;
      if (o.x < -0.2) o.x = 1.2; if (o.x > 1.2) o.x = -0.2;
      if (o.y < -0.2) o.y = 1.2; if (o.y > 1.2) o.y = -0.2;

      const cx = o.x * w, cy = o.y * h;
      const grad = ctx.createRadialGradient(cx, cy, 0, cx, cy, o.r);
      grad.addColorStop(0, `rgba(${o.c},${o.a})`);
      grad.addColorStop(0.5, `rgba(${o.c},${o.a * 0.4})`);
      grad.addColorStop(1, `rgba(${o.c},0)`);
      ctx.fillStyle = grad;
      ctx.fillRect(cx - o.r, cy - o.r, o.r * 2, o.r * 2);
    }
    requestAnimationFrame(draw);
  }
  draw();
}

// ─── CONTENT ───
function renderContent() {
  const totalDocs = 65; // ADR + PRD + SAD + Foundation + Glossary + README

  document.getElementById('content').innerHTML = `
    <div class="content-inner">

      <!-- HERO -->
      <section class="hero">
        <div class="hero-badge"><span class="live"></span> GraphOS · v0.1.0 · Foundation Stage</div>
        <h1 class="hero-title">
          <span class="hero-letter" style="--i:0">A</span>
          <span class="hero-letter" style="--i:1">E</span>
          <span class="hero-letter" style="--i:2">T</span>
          <span class="hero-letter" style="--i:3">H</span>
          <span class="hero-letter" style="--i:4">E</span>
          <span class="hero-letter" style="--i:5">R</span>
        </h1>
        <p class="hero-subtitle">
          <strong>A</strong> Graph-Native
          <strong>A</strong>utonomous
          <strong>S</strong>oftware
          <strong>E</strong>ngineering
          <strong>R</strong>untime
        </p>
        <p class="hero-desc">
          Runtime generasi baru yang memahami perangkat lunak<br class="hide-mobile">
          sebagai <strong>pengetahuan terstruktur</strong> — bukan sekadar file.
        </p>
        <div class="hero-actions">
          <a class="btn-primary" href="#visi"><span>Jelajahi Visi ↓</span></a>
        </div>
      </section>

      <!-- STATS -->
      <div class="stats-bar">
        <div class="stat-item"><div class="stat-num">${totalDocs}</div><div class="stat-label">Dokumen</div></div>
        <div class="stat-item"><div class="stat-num">55+</div><div class="stat-label">Keputusan Arsitektur</div></div>
        <div class="stat-item"><div class="stat-num">30</div><div class="stat-label">Komponen Sistem</div></div>
        <div class="stat-item"><div class="stat-num">5</div><div class="stat-label">Fase Pengembangan</div></div>
      </div>

      <!-- NAMA -->
      <section id="nama" class="section">
        <h2 class="section-title">Makna Dibalik Nama</h2>
        <div class="acronym-list">
          <div class="acronym-row"><span class="acr-badge">A</span><div class="acr-body"><span class="acr-word">Graph-Native</span><span class="acr-desc">Dibangun di atas knowledge graph — relasi adalah yang utama, file hanyalah representasi turunan.</span></div></div>
          <div class="acronym-row"><span class="acr-badge">U</span><div class="acr-body"><span class="acr-word">niversal</span><span class="acr-desc">Agnostik terhadap bahasa pemrograman. Runtime yang mampu memahami basis kode apa pun.</span></div></div>
          <div class="acronym-row"><span class="acr-badge">T</span><div class="acr-body"><span class="acr-word">hinking</span><span class="acr-desc">Lapisan penalaran berbantuan AI di atas inti sistem yang deterministik dan dapat diprediksi.</span></div></div>
          <div class="acronym-row"><span class="acr-badge">O</span><div class="acr-body"><span class="acr-word">rchestrated</span><span class="acr-desc">Koordinasi berbasis event secara real-time di seluruh subsistem runtime.</span></div></div>
          <div class="acronym-row"><span class="acr-badge">M</span><div class="acr-body"><span class="acr-word">odular</span><span class="acr-desc">Arsitektur microkernel — perkuat sistem tanpa menyentuh inti.</span></div></div>
          <div class="acronym-row"><span class="acr-badge">A</span><div class="acr-body"><span class="acr-word">utonomous</span><span class="acr-desc">Runtime yang berevolusi bersama proyek, belajar dari setiap perubahan.</span></div></div>
          <div class="acronym-row"><span class="acr-badge">T</span><div class="acr-body"><span class="acr-word">raceable</span><span class="acr-desc">Setiap keputusan terekam. Setiap perubahan memiliki jejak yang dapat ditelusuri.</span></div></div>
          <div class="acronym-row"><span class="acr-badge">I</span><div class="acr-body"><span class="acr-word">ntegrated</span><span class="acr-desc">Satu runtime untuk analisis, perencanaan, eksekusi, dan peninjauan kode.</span></div></div>
          <div class="acronym-row"><span class="acr-badge">C</span><div class="acr-body"><span class="acr-word">onsistent</span><span class="acr-desc">Knowledge Model selalu sinkron dengan source code, tidak pernah tertinggal.</span></div></div>
          <div class="acronym-row"><span class="acr-badge">E</span><div class="acr-body"><span class="acr-word">xtensible</span><span class="acr-desc">Sistem plugin untuk tools, provider AI, dan bahasa pemrograman baru.</span></div></div>
          <div class="acronym-row"><span class="acr-badge">R</span><div class="acr-body"><span class="acr-word">untime-First</span><span class="acr-desc">Desktop, CLI, API, ekstensi editor — semua berbasis runtime yang identik.</span></div></div>
        </div>
      </section>

      <!-- FILOSOFI -->
      <section id="filosofi" class="section">
        <h2 class="section-title">Filosofi Inti</h2>
        <div class="phil-grid">
          <div class="phil-card">
            <span class="phil-icon">🧠</span>
            <h3>Pengetahuan &gt; Kode</h3>
            <p>Source code adalah <em>representasi</em> dari pengetahuan, bukan kebenaran mutlak. Runtime bekerja dari Knowledge Model, bukan file mentah.</p>
          </div>
          <div class="phil-card">
            <span class="phil-icon">⚙️</span>
            <h3>Runtime adalah Sistem</h3>
            <p>Desktop, CLI, API, ekstensi — semuanya hanyalah antarmuka menuju runtime yang sama. Tidak ada logika bisnis di luar runtime.</p>
          </div>
          <div class="phil-card">
            <span class="phil-icon">🔬</span>
            <h3>Kecerdasan dari Struktur</h3>
            <p>AI tidak membaca file mentah. AI menerima konteks terstruktur yang sudah disiapkan oleh runtime — lebih akurat, lebih efisien.</p>
          </div>
          <div class="phil-card">
            <span class="phil-icon">🎯</span>
            <h3>Deterministik adalah Dasar</h3>
            <p>Semua yang bisa dihitung secara algoritmik dilakukan runtime. AI hanya untuk penalaran dan kreativitas — bukan keputusan kritis.</p>
          </div>
        </div>
      </section>

      <!-- VISI -->
      <section id="visi" class="section">
        <h2 class="section-title">Visi &amp; Misi</h2>
        <div class="vision-block">
          <div class="vision-card">
            <h3>Visi</h3>
            <p>Membangun <em>runtime</em> yang memahami proyek sebagai <strong>model pengetahuan hidup</strong> — manusia dan AI bekerja bersama secara konsisten, dalam skala besar, dengan pertanggungjawaban penuh.</p>
          </div>
          <div class="vision-card">
            <h3>Misi Utama</h3>
            <ul>
              <li>Menggeser paradigma dari <strong>file-centric</strong> ke <strong>knowledge-centric</strong></li>
              <li>Memutus ketergantungan pada <em>context window</em> LLM yang terbatas</li>
              <li>Fondasi untuk <strong>Autonomous Software Engineering</strong> yang nyata</li>
              <li>Arsitektur modular, <strong>provider-independent</strong>, dan extensible</li>
            </ul>
          </div>
        </div>
      </section>

      <!-- PRINSIP -->
      <section id="prinsip" class="section">
        <h2 class="section-title">10 Prinsip Arsitektur</h2>
        <div class="principles-scroll">
          <div class="principle-item"><span class="principle-tag">AP-01</span> Microkernel</div>
          <div class="principle-item"><span class="principle-tag">AP-02</span> Isolasi Layanan</div>
          <div class="principle-item"><span class="principle-tag">AP-03</span> Event First</div>
          <div class="principle-item"><span class="principle-tag">AP-04</span> API Driven</div>
          <div class="principle-item"><span class="principle-tag">AP-05</span> Stateless</div>
          <div class="principle-item"><span class="principle-tag">AP-06</span> Dependency Inversion</div>
          <div class="principle-item"><span class="principle-tag">AP-07</span> Eksplisit</div>
          <div class="principle-item"><span class="principle-tag">AP-08</span> Transaksional</div>
          <div class="principle-item"><span class="principle-tag">AP-09</span> Versioning</div>
          <div class="principle-item"><span class="principle-tag">AP-10</span> Observable</div>
        </div>
      </section>

      <!-- STACK -->
      <section id="stack" class="section">
        <h2 class="section-title">Tumpukan Teknologi</h2>
        <div class="stack-grid">
          <div class="stack-item"><span class="stack-label">Bahasa Inti</span><span class="stack-value">Go</span></div>
          <div class="stack-item"><span class="stack-label">Desktop UI</span><span class="stack-value">Wails</span></div>
          <div class="stack-item"><span class="stack-label">Penyimpanan</span><span class="stack-value">SQLite + Graph</span></div>
          <div class="stack-item"><span class="stack-label">Parser</span><span class="stack-value">Tree-sitter</span></div>
          <div class="stack-item"><span class="stack-label">AI Provider</span><span class="stack-value">Cloud · Lokal · Enterprise</span></div>
          <div class="stack-item"><span class="stack-label">Event Bus</span><span class="stack-value">Internal Pub/Sub</span></div>
          <div class="stack-item"><span class="stack-label">Arsitektur</span><span class="stack-value">Graph-Native</span></div>
          <div class="stack-item"><span class="stack-label">Lisensi</span><span class="stack-value">TBD</span></div>
        </div>
      </section>

      <!-- ROADMAP -->
      <section id="roadmap" class="section">
        <h2 class="section-title">Peta Pengembangan</h2>
        <div class="roadmap">
          <div class="roadmap-phase">
            <div class="phase-dot"></div>
            <div class="phase-content">
              <h3>Fase 1 · Fondasi Runtime</h3>
              <p>Runtime Coordinator · Event System · Storage Layer · Core Lifecycle</p>
            </div>
          </div>
          <div class="roadmap-line"></div>
          <div class="roadmap-phase">
            <div class="phase-dot"></div>
            <div class="phase-content">
              <h3>Fase 2 · Mesin Pemahaman</h3>
              <p>Workspace Scanner · Parser Pipeline · Knowledge Model · Knowledge Graph</p>
            </div>
          </div>
          <div class="roadmap-line"></div>
          <div class="roadmap-phase">
            <div class="phase-dot"></div>
            <div class="phase-content">
              <h3>Fase 3 · Lapisan Kecerdasan</h3>
              <p>Context Engine · AI Provider Abstraction · Planning System · Task Engine</p>
            </div>
          </div>
          <div class="roadmap-line"></div>
          <div class="roadmap-phase">
            <div class="phase-dot"></div>
            <div class="phase-content">
              <h3>Fase 4 · Otomatisasi Terkendali</h3>
              <p>Action Processor · Validation Engine · Git Engine · Memory Engine</p>
            </div>
          </div>
          <div class="roadmap-line"></div>
          <div class="roadmap-phase">
            <div class="phase-dot"></div>
            <div class="phase-content">
              <h3>Fase 5 · Ekosistem &amp; Skalabilitas</h3>
              <p>Plugin System · Extension API · Performance Optimization · Collaborative Runtime</p>
            </div>
          </div>
        </div>
      </section>

      <!-- FOOTER -->
      <footer class="footer">
        <p><strong>Project Aether</strong> — <strong>A</strong> Graph-Native <strong>A</strong>utonomous <strong>S</strong>oftware <strong>E</strong>ngineering <strong>R</strong>untime</p>
        <p class="footer-meta">GraphOS · Tahap Draft · Dibangun dengan pengetahuan terstruktur</p>
      </footer>

    </div>
  `;
}

// ─── SCROLL ANIMATIONS ───
function setupScrollAnimations() {
  const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) entry.target.classList.add('visible');
    });
  }, { threshold: 0.1, rootMargin: '0px 0px -40px 0px' });

  document.querySelectorAll('.section, .hero, .footer, .stats-bar').forEach(el => {
    el.classList.add('fade-section');
    observer.observe(el);
  });
}
