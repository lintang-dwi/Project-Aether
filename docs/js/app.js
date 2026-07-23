/**
 * Project Aether — GitHub Pages App
 * Futuristic single-page summary. No file navigation.
 */

document.addEventListener('DOMContentLoaded', () => {
  initBgCanvas();
  renderContent();
  setupScrollAnimations();
});

// ─── Content ───
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
          A runtime that understands software as structured knowledge — not just files.
          Built for the next era of human-AI collaboration in engineering.
        </p>
        <div class="hero-actions">
          <a class="btn-primary" href="#vision">Explore the Vision ↓</a>
        </div>
      </section>

      <!-- ═══ ACRONYM EXPLODER ═══ -->
      <section id="acronym" class="section">
        <h2 class="section-title">The Name</h2>
        <div class="acronym-grid">
          <div class="acronym-card"><span class="acr-letter">A</span><span class="acr-word">Graph-Native</span><span class="acr-desc">Built on a knowledge graph — relationships first, files second.</span></div>
          <div class="acronym-card"><span class="acr-letter">U</span><span class="acr-word">niversal</span><span class="acr-desc">Language-agnostic runtime that understands any codebase.</span></div>
          <div class="acronym-card"><span class="acr-letter">T</span><span class="acr-word">hinking</span><span class="acr-desc">AI-assisted reasoning layer over a deterministic core.</span></div>
          <div class="acronym-card"><span class="acr-letter">O</span><span class="acr-word">rchestrated</span><span class="acr-desc">Event-driven coordination across all subsystems.</span></div>
          <div class="acronym-card"><span class="acr-letter">M</span><span class="acr-word">odular</span><span class="acr-desc">Microkernel architecture — extend without modifying core.</span></div>
          <div class="acronym-card"><span class="acr-letter">A</span><span class="acr-word">utonomous</span><span class="acr-desc">Self-aware runtime that evolves with the project.</span></div>
          <div class="acronym-card"><span class="acr-letter">T</span><span class="acr-word">raceable</span><span class="acr-desc">Every decision logged, every change explainable.</span></div>
          <div class="acronym-card"><span class="acr-letter">I</span><span class="acr-word">ntegrated</span><span class="acr-desc">One runtime for analysis, planning, execution, review.</span></div>
          <div class="acronym-card"><span class="acr-letter">C</span><span class="acr-word">onsistent</span><span class="acr-desc">Knowledge Model stays in sync with source code.</span></div>
          <div class="acronym-card"><span class="acr-letter">E</span><span class="acr-word">xtensible</span><span class="acr-desc">Plugin system for tools, providers, and languages.</span></div>
          <div class="acronym-card"><span class="acr-letter">R</span><span class="acr-word">untime-First</span><span class="acr-desc">Desktop, CLI, API — all interfaces to the same runtime.</span></div>
        </div>
      </section>

      <!-- ═══ PHILOSOPHY ═══ -->
      <section id="philosophy" class="section">
        <h2 class="section-title">Core Philosophy</h2>
        <div class="phil-grid">
          <div class="phil-card">
            <div class="phil-icon">🧠</div>
            <h3>Knowledge Before Code</h3>
            <p>Source code is a <em>representation</em> of knowledge, not the truth itself. The runtime works from the Knowledge Model, not raw files.</p>
          </div>
          <div class="phil-card">
            <div class="phil-icon">⚙️</div>
            <h3>Runtime Before Interface</h3>
            <p>The runtime is the system. Desktop apps, CLIs, APIs, and editor extensions are just interfaces to it.</p>
          </div>
          <div class="phil-card">
            <div class="phil-icon">🔬</div>
            <h3>Intelligence Through Structure</h3>
            <p>AI understands projects through structured knowledge built by the runtime — not through raw file dumps.</p>
          </div>
          <div class="phil-card">
            <div class="phil-icon">🎯</div>
            <h3>Deterministic Core</h3>
            <p>Everything that can be computed algorithmically stays in the runtime. AI handles reasoning and creativity only.</p>
          </div>
        </div>
      </section>

      <!-- ═══ VISION ═══ -->
      <section id="vision" class="section">
        <h2 class="section-title">Vision &amp; Mission</h2>
        <div class="vision-block">
          <div class="vision-card primary">
            <h3>Vision</h3>
            <p>Build a software engineering runtime that understands projects as <strong>living knowledge models</strong>, enabling humans and AI to work together consistently, at scale, with full accountability.</p>
          </div>
          <div class="vision-card secondary">
            <h3>Key Missions</h3>
            <ul>
              <li>Shift from <strong>file-centric</strong> to <strong>knowledge-centric</strong> development</li>
              <li>Reduce dependency on LLM context windows</li>
              <li>Provide a foundation for <strong>autonomous software engineering</strong></li>
              <li>Build a modular, provider-independent architecture</li>
            </ul>
          </div>
        </div>
      </section>

      <!-- ═══ PRINCIPLES ═══ -->
      <section id="principles" class="section">
        <h2 class="section-title">Architectural Principles</h2>
        <div class="principles-scroll">
          <div class="principle-item"><span class="principle-tag">AP-001</span> Microkernel Architecture</div>
          <div class="principle-item"><span class="principle-tag">AP-002</span> Service Isolation</div>
          <div class="principle-item"><span class="principle-tag">AP-003</span> Event First</div>
          <div class="principle-item"><span class="principle-tag">AP-004</span> API Driven</div>
          <div class="principle-item"><span class="principle-tag">AP-005</span> Stateless Services</div>
          <div class="principle-item"><span class="principle-tag">AP-006</span> Dependency Inversion</div>
          <div class="principle-item"><span class="principle-tag">AP-007</span> Explicit Dependencies</div>
          <div class="principle-item"><span class="principle-tag">AP-008</span> Transactional State</div>
          <div class="principle-item"><span class="principle-tag">AP-009</span> Version Everything</div>
          <div class="principle-item"><span class="principle-tag">AP-010</span> Observable by Default</div>
        </div>
      </section>

      <!-- ═══ STACK ═══ -->
      <section id="stack" class="section">
        <h2 class="section-title">Technology Stack</h2>
        <div class="stack-grid">
          <div class="stack-item"><span class="stack-label">Language</span><span class="stack-value">Go</span></div>
          <div class="stack-item"><span class="stack-label">Desktop</span><span class="stack-value">Wails</span></div>
          <div class="stack-item"><span class="stack-label">Storage</span><span class="stack-value">SQLite + Graph Storage</span></div>
          <div class="stack-item"><span class="stack-label">Parser</span><span class="stack-value">Tree-sitter</span></div>
          <div class="stack-item"><span class="stack-label">AI Providers</span><span class="stack-value">Cloud / Local / Enterprise</span></div>
          <div class="stack-item"><span class="stack-label">Architecture</span><span class="stack-value">Graph-Native Microkernel</span></div>
        </div>
      </section>

      <!-- ═══ ROADMAP ═══ -->
      <section id="roadmap" class="section">
        <h2 class="section-title">Development Roadmap</h2>
        <div class="roadmap">
          <div class="roadmap-phase">
            <div class="phase-dot"></div>
            <div class="phase-content">
              <h3>Phase 1 — Runtime Foundation</h3>
              <p>Runtime Coordinator · Event System · Storage Layer</p>
            </div>
          </div>
          <div class="roadmap-line"></div>
          <div class="roadmap-phase">
            <div class="phase-dot"></div>
            <div class="phase-content">
              <h3>Phase 2 — Understanding Engine</h3>
              <p>Workspace Scanner · Parser Pipeline · Knowledge Model · Knowledge Graph</p>
            </div>
          </div>
          <div class="roadmap-line"></div>
          <div class="roadmap-phase">
            <div class="phase-dot"></div>
            <div class="phase-content">
              <h3>Phase 3 — Intelligence Layer</h3>
              <p>Context Engine · AI Provider · Planning System</p>
            </div>
          </div>
          <div class="roadmap-line"></div>
          <div class="roadmap-phase">
            <div class="phase-dot"></div>
            <div class="phase-content">
              <h3>Phase 4 — Controlled Automation</h3>
              <p>Action Processor · Validation · Git Integration</p>
            </div>
          </div>
          <div class="roadmap-line"></div>
          <div class="roadmap-phase">
            <div class="phase-dot"></div>
            <div class="phase-content">
              <h3>Phase 5 — Extension Ecosystem</h3>
              <p>Plugin System · External Integration · Community</p>
            </div>
          </div>
        </div>
      </section>

      <!-- ═══ FOOTER ═══ -->
      <footer class="footer">
        <p>Project Aether — <strong>A</strong> Graph-Native <strong>A</strong>utonomous <strong>S</strong>oftware <strong>E</strong>ngineering <strong>R</strong>untime</p>
        <p class="footer-meta">GraphOS · Draft Stage · Built with structured knowledge</p>
      </footer>

    </div>
  `;
}

// ─── Animated Background (Particle Network) ───
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
      vx: (Math.random() - 0.5) * 0.25,
      vy: (Math.random() - 0.5) * 0.25,
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

// ─── Scroll Reveal ───
function setupScrollAnimations() {
  const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        entry.target.classList.add('visible');
      }
    });
  }, { threshold: 0.1 });

  document.querySelectorAll('.section, .hero, .footer').forEach(el => {
    el.classList.add('fade-section');
    observer.observe(el);
  });
}
