/**
 * Project Aether — GitHub Pages App
 * Single-page summary — Bahasa Indonesia
 */

document.addEventListener('DOMContentLoaded', () => {
  initBgCanvas();
  renderContent();
  setupScrollAnimations();
});

function renderContent() {
  document.getElementById('content').innerHTML = `
    <div class="content-inner">

      <!-- ═══ HERO ═══ -->
      <section class="hero">
        <div class="hero-badge">v0.1.0 · GraphOS</div>
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
          Runtime yang memahami perangkat lunak sebagai pengetahuan terstruktur — bukan sekadar kumpulan file.
          Dibangun untuk era baru kolaborasi manusia dan AI dalam rekayasa perangkat lunak.
        </p>
        <div class="hero-actions">
          <a class="btn-primary" href="#visi">Jelajahi Visi ↓</a>
        </div>
      </section>

      <!-- ═══ AKRONIM ═══ -->
      <section id="akronim" class="section">
        <h2 class="section-title">Makna Nama</h2>
        <div class="acronym-grid">
          <div class="acronym-card"><span class="acr-letter">A</span><span class="acr-word">Graph-Native</span><span class="acr-desc">Dibangun di atas knowledge graph — relasi adalah yang utama, file hanyalah representasi.</span></div>
          <div class="acronym-card"><span class="acr-letter">U</span><span class="acr-word">niversal</span><span class="acr-desc">Runtime yang agnostik terhadap bahasa pemrograman dan dapat memahami basis kode apa pun.</span></div>
          <div class="acronym-card"><span class="acr-letter">T</span><span class="acr-word">hinking</span><span class="acr-desc">Lapisan penalaran berbantuan AI di atas inti sistem yang deterministik.</span></div>
          <div class="acronym-card"><span class="acr-letter">O</span><span class="acr-word">rchestrated</span><span class="acr-desc">Koordinasi berbasis event di seluruh subsistem secara real-time.</span></div>
          <div class="acronym-card"><span class="acr-letter">M</span><span class="acr-word">odular</span><span class="acr-desc">Arsitektur microkernel — perluas kemampuan tanpa menyentuh inti sistem.</span></div>
          <div class="acronym-card"><span class="acr-letter">A</span><span class="acr-word">utonomous</span><span class="acr-desc">Runtime yang sadar diri dan berevolusi bersama proyek.</span></div>
          <div class="acronym-card"><span class="acr-letter">T</span><span class="acr-word">raceable</span><span class="acr-desc">Setiap keputusan tercatat, setiap perubahan dapat dijelaskan.</span></div>
          <div class="acronym-card"><span class="acr-letter">I</span><span class="acr-word">ntegrated</span><span class="acr-desc">Satu runtime untuk analisis, perencanaan, eksekusi, dan peninjauan.</span></div>
          <div class="acronym-card"><span class="acr-letter">C</span><span class="acr-word">onsistent</span><span class="acr-desc">Knowledge Model selalu sinkron dengan source code.</span></div>
          <div class="acronym-card"><span class="acr-letter">E</span><span class="acr-word">xtensible</span><span class="acr-desc">Sistem plugin untuk alat, penyedia layanan, dan bahasa.</span></div>
          <div class="acronym-card"><span class="acr-letter">R</span><span class="acr-word">untime-First</span><span class="acr-desc">Desktop, CLI, API — semua antarmuka menuju runtime yang sama.</span></div>
        </div>
      </section>

      <!-- ═══ FILOSOFI ═══ -->
      <section id="filosofi" class="section">
        <h2 class="section-title">Filosofi Inti</h2>
        <div class="phil-grid">
          <div class="phil-card">
            <div class="phil-icon">🧠</div>
            <h3>Pengetahuan Sebelum Kode</h3>
            <p>Source code adalah <em>representasi</em> dari pengetahuan, bukan kebenaran mutlak. Runtime bekerja berdasarkan Knowledge Model, bukan file mentah.</p>
          </div>
          <div class="phil-card">
            <div class="phil-icon">⚙️</div>
            <h3>Runtime Sebelum Antarmuka</h3>
            <p>Runtime adalah inti sistem. Aplikasi desktop, CLI, API, dan ekstensi editor hanyalah antarmuka terhadap runtime tersebut.</p>
          </div>
          <div class="phil-card">
            <div class="phil-icon">🔬</div>
            <h3>Kecerdasan Melalui Struktur</h3>
            <p>AI memahami proyek melalui pengetahuan terstruktur yang dibangun oleh runtime — bukan melalui potongan file mentah.</p>
          </div>
          <div class="phil-card">
            <div class="phil-icon">🎯</div>
            <h3>Inti yang Deterministik</h3>
            <p>Segala sesuatu yang dapat dihitung secara algoritmik tetap berada di runtime. AI hanya menangani penalaran dan kreativitas.</p>
          </div>
        </div>
      </section>

      <!-- ═══ VISI ═══ -->
      <section id="visi" class="section">
        <h2 class="section-title">Visi &amp; Misi</h2>
        <div class="vision-block">
          <div class="vision-card primary">
            <h3>Visi</h3>
            <p>Membangun <em>runtime</em> rekayasa perangkat lunak yang memahami proyek sebagai <strong>model pengetahuan hidup</strong>, sehingga manusia dan AI dapat bekerja bersama secara konsisten, dalam skala besar, dengan pertanggungjawaban penuh.</p>
          </div>
          <div class="vision-card secondary">
            <h3>Misi Utama</h3>
            <ul>
              <li>Menggeser paradigma dari <strong>berpusat pada file</strong> menjadi <strong>berpusat pada pengetahuan</strong></li>
              <li>Mengurangi ketergantungan pada <em>context window</em> LLM</li>
              <li>Menyediakan fondasi untuk <strong>rekayasa perangkat lunak otonom</strong></li>
              <li>Membangun arsitektur modular yang independen terhadap penyedia AI</li>
            </ul>
          </div>
        </div>
      </section>

      <!-- ═══ PRINSIP ═══ -->
      <section id="prinsip" class="section">
        <h2 class="section-title">Prinsip Arsitektur</h2>
        <div class="principles-scroll">
          <div class="principle-item"><span class="principle-tag">AP-001</span> Arsitektur Microkernel</div>
          <div class="principle-item"><span class="principle-tag">AP-002</span> Isolasi Layanan</div>
          <div class="principle-item"><span class="principle-tag">AP-003</span> Event Sebagai Prioritas</div>
          <div class="principle-item"><span class="principle-tag">AP-004</span> Berbasis API</div>
          <div class="principle-item"><span class="principle-tag">AP-005</span> Layanan Tanpa Status</div>
          <div class="principle-item"><span class="principle-tag">AP-006</span> Pembalikan Ketergantungan</div>
          <div class="principle-item"><span class="principle-tag">AP-007</span> Ketergantungan Eksplisit</div>
          <div class="principle-item"><span class="principle-tag">AP-008</span> State Transaksional</div>
          <div class="principle-item"><span class="principle-tag">AP-009</span> Versi pada Segala Hal</div>
          <div class="principle-item"><span class="principle-tag">AP-010</span> Terobservasi Secara Bawaan</div>
        </div>
      </section>

      <!-- ═══ STACK ═══ -->
      <section id="teknologi" class="section">
        <h2 class="section-title">Tumpukan Teknologi</h2>
        <div class="stack-grid">
          <div class="stack-item"><span class="stack-label">Bahasa</span><span class="stack-value">Go</span></div>
          <div class="stack-item"><span class="stack-label">Desktop</span><span class="stack-value">Wails</span></div>
          <div class="stack-item"><span class="stack-label">Penyimpanan</span><span class="stack-value">SQLite + Graph Storage</span></div>
          <div class="stack-item"><span class="stack-label">Parser</span><span class="stack-value">Tree-sitter</span></div>
          <div class="stack-item"><span class="stack-label">AI Provider</span><span class="stack-value">Cloud / Lokal / Enterprise</span></div>
          <div class="stack-item"><span class="stack-label">Arsitektur</span><span class="stack-value">Graph-Native Microkernel</span></div>
        </div>
      </section>

      <!-- ═══ ROADMAP ═══ -->
      <section id="roadmap" class="section">
        <h2 class="section-title">Peta Pengembangan</h2>
        <div class="roadmap">
          <div class="roadmap-phase">
            <div class="phase-dot"></div>
            <div class="phase-content">
              <h3>Fase 1 — Fondasi Runtime</h3>
              <p>Runtime Coordinator · Sistem Event · Lapisan Penyimpanan</p>
            </div>
          </div>
          <div class="roadmap-line"></div>
          <div class="roadmap-phase">
            <div class="phase-dot"></div>
            <div class="phase-content">
              <h3>Fase 2 — Mesin Pemahaman</h3>
              <p>Pemindai Workspace · Pipeline Parser · Knowledge Model · Knowledge Graph</p>
            </div>
          </div>
          <div class="roadmap-line"></div>
          <div class="roadmap-phase">
            <div class="phase-dot"></div>
            <div class="phase-content">
              <h3>Fase 3 — Lapisan Kecerdasan</h3>
              <p>Context Engine · AI Provider · Sistem Perencanaan</p>
            </div>
          </div>
          <div class="roadmap-line"></div>
          <div class="roadmap-phase">
            <div class="phase-dot"></div>
            <div class="phase-content">
              <h3>Fase 4 — Otomatisasi Terkendali</h3>
              <p>Action Processor · Validasi · Integrasi Git</p>
            </div>
          </div>
          <div class="roadmap-line"></div>
          <div class="roadmap-phase">
            <div class="phase-dot"></div>
            <div class="phase-content">
              <h3>Fase 5 — Ekosistem Ekstensi</h3>
              <p>Sistem Plugin · Integrasi Eksternal · Komunitas</p>
            </div>
          </div>
        </div>
      </section>

      <!-- ═══ FOOTER ═══ -->
      <footer class="footer">
        <p>Project Aether — <strong>A</strong> Graph-Native <strong>A</strong>utonomous <strong>S</strong>oftware <strong>E</strong>ngineering <strong>R</strong>untime</p>
        <p class="footer-meta">GraphOS · Tahap Draft · Dibangun dengan pengetahuan terstruktur</p>
      </footer>

    </div>
  `;
}

function initBgCanvas() {
  const canvas = document.getElementById('bg-canvas');
  if (!canvas) return;
  const ctx = canvas.getContext('2d');
  let w, h, particles = [];

  function resize() {
    w = canvas.width = window.innerWidth;
    h = canvas.height = window.innerHeight;
  }
  resize();
  window.addEventListener('resize', resize);

  const PCOUNT = 90;
  for (let i = 0; i < PCOUNT; i++) {
    particles.push({
      x: Math.random() * w, y: Math.random() * h,
      vx: (Math.random() - 0.5) * 0.25, vy: (Math.random() - 0.5) * 0.25,
      r: Math.random() * 2 + 0.5,
    });
  }

  function draw() {
    ctx.clearRect(0, 0, w, h);
    for (const p of particles) {
      p.x += p.vx; p.y += p.vy;
      if (p.x < 0) p.x = w; if (p.x > w) p.x = 0;
      if (p.y < 0) p.y = h; if (p.y > h) p.y = 0;
      ctx.beginPath();
      ctx.arc(p.x, p.y, p.r, 0, Math.PI * 2);
      ctx.fillStyle = 'rgba(56, 189, 248, 0.3)';
      ctx.fill();
    }
    for (let i = 0; i < particles.length; i++) {
      for (let j = i + 1; j < particles.length; j++) {
        const dx = particles[i].x - particles[j].x;
        const dy = particles[i].y - particles[j].y;
        const dist = Math.sqrt(dx * dx + dy * dy);
        if (dist < 180) {
          ctx.beginPath();
          ctx.moveTo(particles[i].x, particles[i].y);
          ctx.lineTo(particles[j].x, particles[j].y);
          ctx.strokeStyle = `rgba(56, 189, 248, ${0.08 * (1 - dist / 180)})`;
          ctx.lineWidth = 0.5;
          ctx.stroke();
        }
      }
    }
    requestAnimationFrame(draw);
  }
  draw();
}

function setupScrollAnimations() {
  const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) entry.target.classList.add('visible');
    });
  }, { threshold: 0.1 });
  document.querySelectorAll('.section, .hero, .footer').forEach(el => {
    el.classList.add('fade-section');
    observer.observe(el);
  });
}
