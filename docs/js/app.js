/**
 * Project Aether — GitHub Pages App
 * Loads & renders all existing markdown files dynamically
 */

// ─── File Index (mirrors existing structure, no modifications) ───
const FILES = {
  "Project Aether": {
    type: "section",
    items: [
      { name: "Overview", path: "../README.md" },
      { name: "Foundation", path: "../FOUNDATION.md" },
      { name: "Glossary", path: "../Project%20Aether%20Glossary.md" },
    ]
  },
  "Product Requirements Document": {
    type: "folder",
    base: "../Product%20Requirements%20Document/",
    items: [
      "1. Executive Summary.md",
      "2. Product Tenets.md",
      "3. Product Overview.md",
      "4. Problem Statement.md",
      "5. Goals.md",
      "6. Success Metrics.md",
      "7. Target Users.md",
      "8. Product Scope.md",
      "9. Product Architecture Overview.md",
      "10. Core Capabilities.md",
      "11. Functional Requirements.md",
      "12. Non-Functional Requirements (NFR).md",
      "13. System Constraints & Assumptions.md",
      "14. Development Roadmap.md",
      "15. Success Metrics & Acceptance Criteria.md",
    ]
  },
  "Architecture Decision Record": {
    type: "folder",
    base: "../Architecture%20Decision%20Record/",
    items: [
      "ADR-001 — Runtime Coordinator as Core Orchestrator.md",
      "ADR-002 — Knowledge Graph as Source of Understanding.md",
      "ADR-003 — Workspace → Knowledge → Task → Context → AI Dependency Model.md",
      "ADR-004 — Event-Driven Runtime Architecture.md",
      "ADR-005 — Deterministic Core vs AI Responsibility Boundary.md",
      "ADR-006 — Local First Desktop Runtime.md",
      "ADR-007 — Extension atau Plugin Architecture.md",
      "ADR-008 — Security Boundary dan Permission Model.md",
      "ADR-010 — AI Context Retrieval Strategy.md",
      "ADR-011 — Go as Runtime Implementation Language.md",
      "ADR-012 — Desktop Framework Selection (Wails).md",
      "ADR-013 — Storage Strategy (SQLite + Graph Storage).md",
      "ADR-014 — Event Bus Implementation.md",
      "ADR-015 — Parser Architecture (Tree-sitter atau AST Pipeline).md",
      "ADR-016 — Git Integration Strategy.md",
      "ADR-017 — AI Provider Abstraction Layer.md",
      "ADR-018 — Local Model Support Strategy.md",
      "ADR-019 — Configuration Management.md",
      "ADR-020 — Build & Release Architecture.md",
    ]
  },
  "Software Architecture Document": {
    type: "folder",
    base: "../Software%20Architecture%20Document/",
    items: [
      "1. Architecture Overview.md",
      "2. Architectural Principles.md",
      "3. System Context Diagram.md",
      "4. Container Architecture.md",
      "5. Runtime Architecture Detail.md",
      "6. Graph Engine Architecture.md",
      "7. Knowledge Engine Architecture.md",
      "8. Task Engine Architecture.md",
      "9. Context Engine Architecture.md",
      "10. Agent Runtime Architecture.md",
      "11. Tool Runtime Architecture.md",
      "12. Validation Engine Architecture.md",
      "13. Memory Engine Architecture.md",
      "14. Event Architecture.md",
      "15. Storage Architecture.md",
      "16. Workspace Architecture.md",
      "17. Git Engine Architecture.md",
      "18. Parser & Knowledge Extraction Architecture.md",
      "19. Knowledge Graph Architecture.md",
      "20. Runtime Coordinator Architecture.md",
      "21. UI Engine Architecture.md",
      "22. Deployment Architecture.md",
      "23. Security Architecture.md",
      "24. Observability Architecture.md",
      "25. Extension & Plugin Architecture.md",
      "26. API & External Interface Architecture.md",
      "27. Data Flow Architecture.md",
      "28. Performance & Scalability Architecture.md",
      "29. Testing & Quality Architecture.md",
      "30. Final System Architecture Summary.md",
    ]
  }
};

// ─── State ───
let allDocIndex = [];

// ─── DOM refs ───
const sidebar = document.getElementById('sidebar');
const content = document.getElementById('content');
const searchInput = document.getElementById('search-input');
const searchResults = document.getElementById('search-results');

// ─── Init ───
document.addEventListener('DOMContentLoaded', () => {
  renderSidebar();
  setupSearch();
  initBgCanvas();
  loadWelcome();
});

// ─── Render Sidebar Navigation ───
function renderSidebar() {
  sidebar.innerHTML = '';
  allDocIndex = [];

  for (const [sectionName, section] of Object.entries(FILES)) {
    const sectionEl = document.createElement('div');
    sectionEl.className = 'nav-section';

    const title = document.createElement('div');
    title.className = 'nav-section-title';
    title.innerHTML = `<span class="arrow">▼</span> ${sectionName} <span class="count">${section.items.length}</span>`;
    title.addEventListener('click', () => {
      title.classList.toggle('collapsed');
    });
    sectionEl.appendChild(title);

    const items = document.createElement('div');
    items.className = 'nav-items';

    section.items.forEach((item, idx) => {
      let displayName, fullPath;
      if (section.type === 'section') {
        displayName = item.name;
        fullPath = item.path;
      } else {
        displayName = item.replace(/\.md$/, '');
        fullPath = section.base + encodeURIComponent(item);
      }

      allDocIndex.push({ name: displayName, path: fullPath, section: sectionName });

      const a = document.createElement('a');
      a.className = 'nav-item';
      a.textContent = displayName;
      a.dataset.path = fullPath;
      a.addEventListener('click', (e) => {
        e.preventDefault();
        document.querySelectorAll('.nav-item').forEach(n => n.classList.remove('active'));
        a.classList.add('active');
        loadDoc(fullPath);
      });
      items.appendChild(a);
    });

    sectionEl.appendChild(items);
    sidebar.appendChild(sectionEl);

    // Collapse sections with > 10 items initially
    if (section.items.length > 10) {
      title.classList.add('collapsed');
    }
  }
}

// ─── Load & Render Markdown ───
async function loadDoc(path) {
  content.innerHTML = `<div class="loading"><div class="spinner"></div>Loading...</div>`;
  try {
    const resp = await fetch(path);
    if (!resp.ok) throw new Error(`HTTP ${resp.status}`);
    let md = await resp.text();

    // Strip outermost code fences if present (some files are wrapped)
    md = md.replace(/^```(?:markdown\s+)?(?:id="[^"]*")?\s*\n?/i, '');
    md = md.replace(/\n```\s*$/, '');

    const html = marked.parse(md, { breaks: true, gfm: true });
    content.innerHTML = `<div class="content-inner markdown-body">${html}</div>`;

    // Highlight code blocks
    document.querySelectorAll('pre code').forEach(block => {
      hljs.highlightElement(block);
    });

    // Update URL hash
    const docName = path.split('/').pop().replace(/\.md$/, '');
    history.replaceState(null, '', '#' + encodeURIComponent(path));
  } catch (err) {
    content.innerHTML = `<div class="error">Failed to load document: ${err.message}</div>`;
  }
}

// ─── Welcome Page ───
function loadWelcome() {
  const sections = Object.entries(FILES);
  const total = allDocIndex.length;
  content.innerHTML = `
    <div class="content-inner welcome">
      <h1>Project Aether</h1>
      <p class="subtitle">A Graph-Native Autonomous Software Engineering Runtime — explore the full documentation below.</p>
      <div class="cards">
        ${sections.map(([name, s]) => `
          <div class="card" onclick="document.querySelector('.nav-section-title').click()">
            <div class="icon">${getSectionIcon(name)}</div>
            <div class="name">${name}</div>
            <div class="desc">${s.items.length} documents</div>
          </div>
        `).join('')}
      </div>
      <div style="margin-top: 28px; font-size: 0.85rem; color: var(--text-muted);">
        ${total} documents · Foundation → PRD → ADR → SAD
      </div>
    </div>
  `;

  // Check hash for direct doc link
  const hash = decodeURIComponent(location.hash.slice(1));
  if (hash) {
    const item = document.querySelector(`.nav-item[data-path="${hash}"]`);
    if (item) item.click();
  }
}

function getSectionIcon(name) {
  if (name.includes('Foundation') || name === 'Project Aether') return '🚀';
  if (name.includes('Product')) return '📋';
  if (name.includes('Architecture Decision')) return '⚖️';
  if (name.includes('Software')) return '🏗️';
  return '📄';
}

// ─── Search ───
function setupSearch() {
  searchInput.addEventListener('focus', () => {
    if (searchInput.value.trim()) showSearch();
  });
  searchInput.addEventListener('input', () => {
    if (searchInput.value.trim()) showSearch();
    else searchResults.classList.remove('show');
  });
  document.addEventListener('click', (e) => {
    if (!e.target.closest('.header-search')) searchResults.classList.remove('show');
  });
}

function showSearch() {
  const q = searchInput.value.toLowerCase().trim();
  if (!q) return;
  const results = allDocIndex.filter(d => d.name.toLowerCase().includes(q) || d.section.toLowerCase().includes(q));
  searchResults.innerHTML = results.slice(0, 20).map(r => `
    <div class="search-result-item" data-path="${r.path}">
      <div class="path">${r.section}</div>
      <div class="match">${highlight(r.name, q)}</div>
    </div>
  `).join('');
  searchResults.classList.add('show');
  searchResults.querySelectorAll('.search-result-item').forEach(el => {
    el.addEventListener('click', () => {
      const path = el.dataset.path;
      document.querySelectorAll('.nav-item').forEach(n => n.classList.remove('active'));
      const nav = document.querySelector(`.nav-item[data-path="${path}"]`);
      if (nav) { nav.classList.add('active'); nav.scrollIntoView({ block: 'nearest' }); }
      loadDoc(path);
      searchResults.classList.remove('show');
      searchInput.value = '';
    });
  });
}

function highlight(text, q) {
  const idx = text.toLowerCase().indexOf(q);
  if (idx === -1) return text;
  return text.slice(0, idx) + '<mark>' + text.slice(idx, idx + q.length) + '</mark>' + text.slice(idx + q.length);
}

// ─── Animated Background (Particle Network) ───
function initBgCanvas() {
  const canvas = document.getElementById('bg-canvas');
  const ctx = canvas.getContext('2d');
  let w, h, particles = [];

  function resize() {
    w = canvas.width = window.innerWidth;
    h = canvas.height = window.innerHeight;
  }
  resize();
  window.addEventListener('resize', resize);

  const PCOUNT = 80;
  for (let i = 0; i < PCOUNT; i++) {
    particles.push({
      x: Math.random() * w,
      y: Math.random() * h,
      vx: (Math.random() - 0.5) * 0.3,
      vy: (Math.random() - 0.5) * 0.3,
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
