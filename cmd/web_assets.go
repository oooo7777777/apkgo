package cmd

const webAppsHTML = `<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>apkgo Apps</title>
  <style>
    :root {
      --bg: #090b10;
      --panel: rgba(14,18,25,0.88);
      --line: rgba(130,164,255,0.14);
      --line-strong: rgba(130,164,255,0.28);
      --text: #edf3ff;
      --muted: #8d9ab3;
      --accent: #4de2c5;
      --accent-2: #6aa9ff;
      --ok: #46d39a;
      --bad: #ff6b7a;
      --shadow: 0 28px 80px rgba(0,0,0,0.45);
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      min-height: 100vh;
      color: var(--text);
      font-family: "SF Mono", "JetBrains Mono", "IBM Plex Sans", "PingFang SC", "Microsoft YaHei", monospace, sans-serif;
      background:
        radial-gradient(circle at 12% 10%, rgba(77,226,197,0.12), transparent 22%),
        radial-gradient(circle at 88% 14%, rgba(106,169,255,0.12), transparent 20%),
        linear-gradient(180deg, #07090d 0%, #090b10 100%);
    }
    .shell {
      max-width: 1120px;
      margin: 0 auto;
      padding: 28px 20px 52px;
    }
    .hero, .panel {
      border-radius: 28px;
      background: var(--panel);
      border: 1px solid var(--line-strong);
      box-shadow: var(--shadow);
    }
    .hero {
      padding: 24px;
      margin-bottom: 20px;
    }
    .hero-top {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 16px;
    }
    h1 {
      margin: 0 0 10px;
      font-size: clamp(34px, 6vw, 58px);
      line-height: .94;
      letter-spacing: -0.06em;
    }
    .hero-desc {
      max-width: 720px;
      color: var(--muted);
      font-size: 14px;
      line-height: 1.8;
    }
    .page-actions {
      display: flex;
      gap: 12px;
      flex-wrap: wrap;
    }
    .hero-link, button, .primary-link {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      min-height: 44px;
      padding: 0 18px;
      border-radius: 999px;
      border: 1px solid rgba(255,255,255,0.08);
      background: rgba(255,255,255,0.05);
      color: var(--text);
      text-decoration: none;
      font: inherit;
      font-size: 12px;
      font-weight: 800;
      letter-spacing: 0.08em;
      text-transform: uppercase;
      cursor: pointer;
    }
    .primary-link, button.primary {
      background: linear-gradient(135deg, rgba(77,226,197,0.18), rgba(106,169,255,0.2));
      border-color: rgba(77,226,197,0.28);
    }
    .panel-head {
      padding: 18px 22px 14px;
      border-bottom: 1px solid rgba(255,255,255,0.05);
    }
    .panel-head-row {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 14px;
      margin-bottom: 8px;
    }
    .panel-kicker {
      color: var(--accent-2);
      font-size: 11px;
      font-weight: 800;
      letter-spacing: 0.08em;
      text-transform: uppercase;
      margin-bottom: 10px;
    }
    .panel h2 {
      margin: 0 0 8px;
      font-size: 24px;
      letter-spacing: -0.04em;
    }
    .note {
      margin: 0;
      color: var(--muted);
      font-size: 13px;
      line-height: 1.7;
    }
    .panel-body {
      padding: 22px;
    }
    .app-list {
      display: grid;
      gap: 14px;
    }
    .app-card {
      display: grid;
      grid-template-columns: minmax(0, 1fr) auto;
      gap: 12px;
      align-items: center;
      padding: 16px 18px;
      border: 1px solid rgba(130,164,255,0.14);
      border-radius: 20px;
      background: rgba(8,12,18,0.9);
    }
    .app-card.selected {
      border-color: rgba(77,226,197,0.4);
      box-shadow: 0 0 0 3px rgba(77,226,197,0.08);
    }
    .app-meta {
      display: grid;
      gap: 6px;
    }
    .app-name {
      font-size: 18px;
      font-weight: 800;
    }
    .app-subtitle {
      color: var(--muted);
      font-size: 13px;
      line-height: 1.6;
    }
    .badge {
      display: inline-flex;
      width: fit-content;
      padding: 4px 10px;
      border-radius: 999px;
      font-size: 11px;
      font-weight: 800;
      letter-spacing: 0.06em;
      text-transform: uppercase;
      background: rgba(255,255,255,0.06);
      color: var(--muted);
    }
    .badge.ok {
      background: rgba(70,211,154,0.16);
      color: #88efc3;
    }
    .app-actions {
      display: flex;
      gap: 10px;
      flex-wrap: wrap;
      justify-content: flex-end;
    }
    .ghost {
      background: rgba(255,255,255,0.04);
    }
    .danger {
      border-color: rgba(255,107,122,0.24);
      color: #ffc2c8;
      background: rgba(255,107,122,0.08);
    }
    .status {
      margin-top: 14px;
      padding: 12px 14px;
      border-radius: 16px;
      border: 1px solid transparent;
      display: none;
      font-size: 13px;
      line-height: 1.6;
    }
    .status.show { display: block; }
    .status.ok {
      background: rgba(70,211,154,0.12);
      border-color: rgba(70,211,154,0.24);
      color: #88efc3;
    }
    .status.bad {
      background: rgba(255,107,122,0.12);
      border-color: rgba(255,107,122,0.24);
      color: #ffc2c8;
    }
    .empty {
      padding: 18px;
      border-radius: 18px;
      background: rgba(255,255,255,0.03);
      color: var(--muted);
      text-align: center;
    }
    @media (max-width: 760px) {
      .hero-top, .app-card {
        grid-template-columns: 1fr;
        display: grid;
      }
      .app-actions {
        justify-content: flex-start;
      }
    }
  </style>
</head>
<body>
  <div class="shell">
    <section class="hero">
      <div class="hero-top">
        <div>
          <h1>APKGO</h1>
          <div class="hero-desc">一个面向多应用分发场景的 APK 发布工具。你可以先选择当前要操作的 App，系统会将对应配置切换为主 config，再继续完成发布、审核查询和历史记录管理。</div>
        </div>
      </div>
    </section>

    <section class="panel">
      <div class="panel-head">
        <div class="panel-head-row">
          <div class="panel-kicker">app center</div>
          <a class="primary-link" href="/config?mode=create-app">新增 App</a>
        </div>
        <h2>App 列表</h2>
      </div>
      <div class="panel-body">
        <div id="apps-status" class="status"></div>
        <div id="apps-list" class="app-list"></div>
      </div>
    </section>
  </div>

  <script>
    function escapeHtml(value) {
      return String(value || '')
        .replaceAll('&', '&amp;')
        .replaceAll('<', '&lt;')
        .replaceAll('>', '&gt;')
        .replaceAll('"', '&quot;')
        .replaceAll("'", '&#39;');
    }

    function setStatus(message, kind) {
      const el = document.getElementById('apps-status');
      if (!message) {
        el.textContent = '';
        el.className = 'status';
        return;
      }
      el.textContent = message;
      el.className = 'status show ' + kind;
    }

    async function loadApps() {
      const resp = await fetch('/api/apps');
      const data = await resp.json();
      if (!resp.ok) throw new Error(data.error || '读取 App 失败');
      renderApps(data.apps || []);
    }

    function refreshAppsOnReturn() {
      loadApps().catch((err) => setStatus(String(err), 'bad'));
    }

    window.addEventListener('pageshow', refreshAppsOnReturn);
    window.addEventListener('focus', refreshAppsOnReturn);
    document.addEventListener('visibilitychange', () => {
      if (!document.hidden) refreshAppsOnReturn();
    });

    function renderApps(apps) {
      const root = document.getElementById('apps-list');
      if (!apps.length) {
        root.innerHTML = '<div class="empty">还没有可用的 App，请先新增一个。</div>';
        return;
      }
      root.innerHTML = apps.map((app) => (
        '<section class="app-card' + (app.selected ? ' selected' : '') + '">' +
          '<div class="app-meta">' +
            '<div class="app-name">' + escapeHtml(app.name) + '</div>' +
            '<div class="app-subtitle">包名：' + escapeHtml(app.package_name || '未配置') + '</div>' +
            (app.selected ? '<span class="badge ok">当前主配置</span>' : '<span class="badge">可切换</span>') +
          '</div>' +
          '<div class="app-actions">' +
            (app.selected ? '<a class="ghost primary-link" href="/upload">进入发布页</a>' : '<button type="button" class="primary" data-select="' + escapeHtml(app.id) + '">选择并进入发布页</button>') +
            '<a class="ghost primary-link" href="/config?app=' + encodeURIComponent(app.id) + '">编辑</a>' +
            '<button type="button" class="danger" data-delete="' + escapeHtml(app.id) + '">删除</button>' +
          '</div>' +
        '</section>'
      )).join('');
    }

    async function selectApp(id) {
      setStatus('', '');
      const resp = await fetch('/api/apps/select', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id }),
      });
      const data = await resp.json();
      if (!resp.ok) throw new Error(data.error || '切换 App 失败');
      window.location.href = '/upload';
    }

    async function deleteApp(id) {
      if (!window.confirm('确认删除这个 App 配置吗？')) return;
      setStatus('', '');
      const resp = await fetch('/api/apps/delete', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id }),
      });
      const data = await resp.json();
      if (!resp.ok) throw new Error(data.error || '删除 App 失败');
      setStatus('App 已删除。', 'ok');
      renderApps(data.apps || []);
    }

    document.addEventListener('click', (event) => {
      const selectBtn = event.target.closest('[data-select]');
      if (selectBtn) {
        selectApp(selectBtn.getAttribute('data-select')).catch((err) => setStatus(String(err), 'bad'));
        return;
      }
      const deleteBtn = event.target.closest('[data-delete]');
      if (deleteBtn) {
        deleteApp(deleteBtn.getAttribute('data-delete')).catch((err) => setStatus(String(err), 'bad'));
      }
    });

    loadApps().catch((err) => setStatus(String(err), 'bad'));
  </script>
</body>
</html>
`

const webUploadHTML = `<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>apkgo Web</title>
  <style>
    :root {
      --bg: #090b10;
      --bg-soft: #0e1219;
      --panel: rgba(14,18,25,0.88);
      --panel-2: rgba(18,24,34,0.96);
      --line: rgba(130,164,255,0.14);
      --line-strong: rgba(130,164,255,0.28);
      --text: #edf3ff;
      --muted: #8d9ab3;
      --accent: #4de2c5;
      --accent-2: #6aa9ff;
      --accent-3: #8d6bff;
      --ok: #46d39a;
      --warn: #ffb454;
      --bad: #ff6b7a;
      --shadow: 0 28px 80px rgba(0,0,0,0.45);
      --radius: 24px;
    }
    * { box-sizing: border-box; }
    html, body {
      min-height: 100%;
    }
    body {
      margin: 0;
      color: var(--text);
      font-family: "SF Mono", "JetBrains Mono", "IBM Plex Sans", "PingFang SC", "Microsoft YaHei", monospace, sans-serif;
      background:
        radial-gradient(circle at 12% 10%, rgba(77,226,197,0.12), transparent 22%),
        radial-gradient(circle at 88% 14%, rgba(106,169,255,0.12), transparent 20%),
        radial-gradient(circle at 50% 120%, rgba(141,107,255,0.1), transparent 30%),
        linear-gradient(180deg, #07090d 0%, var(--bg) 100%);
      overflow-x: hidden;
    }
    body::before {
      content: "";
      position: fixed;
      inset: 0;
      pointer-events: none;
      background-image:
        linear-gradient(rgba(255,255,255,0.02) 1px, transparent 1px),
        linear-gradient(90deg, rgba(255,255,255,0.02) 1px, transparent 1px);
      background-size: 32px 32px;
      mask-image: radial-gradient(circle at center, black 45%, transparent 92%);
    }
    .shell {
      max-width: 1360px;
      margin: 0 auto;
      padding: 28px 20px 52px;
    }
    .hero {
      position: relative;
      overflow: hidden;
      padding: 22px 24px;
      border-radius: 32px;
      background:
        linear-gradient(135deg, rgba(18,24,34,0.96), rgba(10,14,20,0.94)),
        linear-gradient(180deg, rgba(77,226,197,0.05), rgba(106,169,255,0.04));
      border: 1px solid var(--line-strong);
      box-shadow: var(--shadow);
      margin-bottom: 20px;
    }
    .hero::after {
      content: "";
      position: absolute;
      inset: -1px;
      border-radius: inherit;
      padding: 1px;
      background: linear-gradient(135deg, rgba(77,226,197,0.3), rgba(106,169,255,0.18), rgba(141,107,255,0.2));
      -webkit-mask:
        linear-gradient(#000 0 0) content-box,
        linear-gradient(#000 0 0);
      -webkit-mask-composite: xor;
      mask-composite: exclude;
      pointer-events: none;
    }
    .hero-top {
      display: grid;
      gap: 16px;
      position: relative;
      z-index: 1;
    }
    .hero-bar {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 16px;
    }
    .hero-actions {
      display: flex;
      align-items: center;
      justify-content: flex-end;
      flex: 0 0 auto;
      gap: 12px;
    }
    .hero-link {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      min-height: 44px;
      padding: 0 18px;
      border-radius: 999px;
      border: 1px solid rgba(255,255,255,0.08);
      background: rgba(255,255,255,0.05);
      color: var(--text);
      text-decoration: none;
      font-size: 12px;
      font-weight: 800;
      letter-spacing: 0.08em;
      text-transform: uppercase;
      transition: transform .16s ease, border-color .16s ease, background .16s ease;
    }
    .hero-link:hover {
      transform: translateY(-1px);
      border-color: rgba(77,226,197,0.28);
      background: rgba(77,226,197,0.08);
    }
    h1 {
      margin: 0;
      font-size: clamp(36px, 6vw, 64px);
      line-height: 0.94;
      letter-spacing: -0.06em;
    }
    .hero-side {
      display: grid;
      gap: 12px;
      grid-template-columns: repeat(3, minmax(0, 1fr));
    }
    .hero-stat {
      padding: 12px 14px;
      border-radius: 18px;
      background: rgba(255,255,255,0.03);
      border: 1px solid rgba(255,255,255,0.06);
      backdrop-filter: blur(12px);
    }
    .hero-stat strong {
      display: block;
      color: var(--muted);
      font-size: 11px;
      font-weight: 800;
      letter-spacing: 0.08em;
      text-transform: uppercase;
      margin-bottom: 8px;
    }
    .hero-stat span {
      font-size: 13px;
      line-height: 1.5;
      color: var(--text);
    }
    .workspace {
      display: grid;
      grid-template-columns: minmax(0, 1.05fr) minmax(360px, 0.95fr);
      gap: 20px;
      align-items: stretch;
    }
    .panel {
      background: var(--panel);
      border: 1px solid var(--line);
      border-radius: 28px;
      box-shadow: var(--shadow);
      overflow: hidden;
      height: 100%;
    }
    .panel-head {
      padding: 18px 22px 14px;
      border-bottom: 1px solid rgba(255,255,255,0.05);
      background: linear-gradient(180deg, rgba(255,255,255,0.02), rgba(255,255,255,0));
    }
    .panel-kicker {
      color: var(--accent-2);
      font-size: 11px;
      font-weight: 800;
      letter-spacing: 0.08em;
      text-transform: uppercase;
      margin-bottom: 10px;
    }
    .panel h2 {
      margin: 0 0 8px;
      font-size: 24px;
      letter-spacing: -0.04em;
    }
    .panel p.note {
      margin: 0;
      color: var(--muted);
      font-size: 13px;
      line-height: 1.7;
    }
    .panel-body {
      padding: 22px;
    }
    .field {
      display: grid;
      gap: 8px;
    }
    .field-card {
      position: relative;
      display: grid;
      gap: 8px;
      padding: 14px 16px;
      border: 1px solid rgba(130,164,255,0.14);
      border-radius: 18px;
      background: rgba(8,12,18,0.9);
      cursor: pointer;
      transition: border-color .18s ease, box-shadow .18s ease, transform .18s ease, background .18s ease;
    }
    .field-card:hover {
      transform: translateY(-1px);
      border-color: rgba(77,226,197,0.28);
      background: rgba(10,16,24,0.94);
    }
    .field-card:focus-within {
      border-color: rgba(77,226,197,0.45);
      box-shadow: 0 0 0 4px rgba(77,226,197,0.1);
    }
    .field-card .small {
      margin: 0;
    }
    .field-card-meta {
      color: var(--muted);
      font-size: 12px;
      line-height: 1.6;
    }
    .hidden {
      display: none;
    }
    label.small {
      font-size: 12px;
      color: var(--muted);
      font-weight: 700;
      letter-spacing: 0.04em;
      text-transform: uppercase;
    }
    input, textarea, select {
      width: 100%;
      border: 1px solid rgba(130,164,255,0.14);
      border-radius: 16px;
      background: rgba(8,12,18,0.9);
      padding: 14px 15px;
      font: inherit;
      color: var(--text);
      transition: border-color .18s ease, box-shadow .18s ease, transform .18s ease;
    }
    select {
      appearance: none;
      -webkit-appearance: none;
      padding-right: 64px;
      background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='14' height='9' viewBox='0 0 14 9' fill='none'%3E%3Cpath d='M1 1.5L7 7.5L13 1.5' stroke='%23edf3ff' stroke-width='1.8' stroke-linecap='round' stroke-linejoin='round'/%3E%3C/svg%3E");
      background-repeat: no-repeat;
      background-position: right 20px center;
      background-size: 14px 9px;
    }
    textarea {
      min-height: 120px;
      resize: vertical;
    }
    input:focus, textarea:focus, select:focus {
      outline: none;
      border-color: rgba(77,226,197,0.45);
      box-shadow: 0 0 0 4px rgba(77,226,197,0.1);
    }
    input[type="datetime-local"] {
      border: 1px solid rgba(255,255,255,0.06);
      background:
        linear-gradient(180deg, rgba(255,255,255,0.03), rgba(255,255,255,0.01)),
        rgba(6,10,16,0.96);
      padding: 16px 18px;
      border-radius: 16px;
      color: var(--text);
      font-size: 15px;
      font-weight: 700;
      letter-spacing: 0.02em;
      cursor: pointer;
      color-scheme: dark;
    }
    input[type="datetime-local"]::-webkit-calendar-picker-indicator {
      cursor: pointer;
      filter: invert(88%) sepia(12%) saturate(834%) hue-rotate(107deg) brightness(101%) contrast(93%);
      opacity: 0.9;
    }
    .upload-grid {
      display: grid;
      gap: 16px;
    }
    .market-box {
      border: 1px solid rgba(255,255,255,0.06);
      border-radius: 20px;
      background: rgba(255,255,255,0.03);
      padding: 14px;
    }
    .market-box.loading {
      color: var(--muted);
    }
    .market-list {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(190px, 1fr));
      gap: 12px;
    }
    .market-card {
      position: relative;
      border: 1px solid rgba(130,164,255,0.12);
      background: rgba(8,12,18,0.84);
      border-radius: 18px;
      padding: 14px;
      min-height: 98px;
      transition: border-color .18s ease, box-shadow .18s ease, transform .18s ease;
    }
    .market-card:hover {
      transform: translateY(-1px);
      border-color: rgba(77,226,197,0.26);
    }
    .market-card.disabled {
      opacity: .5;
      background: rgba(255,255,255,0.02);
    }
    .market-check {
      position: absolute;
      inset: 0;
      opacity: 0;
      cursor: pointer;
    }
    .market-card.selected {
      border-color: rgba(77,226,197,0.42);
      box-shadow: 0 0 0 4px rgba(77,226,197,0.08);
    }
    .market-title {
      font-size: 16px;
      font-weight: 800;
      margin-bottom: 6px;
    }
    .market-meta {
      font-size: 12px;
      color: var(--muted);
      line-height: 1.6;
    }
    .badge {
      position: absolute;
      top: 12px;
      right: 12px;
      border-radius: 999px;
      padding: 5px 9px;
      font-size: 11px;
      font-weight: 800;
      line-height: 1;
      letter-spacing: 0.02em;
    }
    .badge.ok {
      background: rgba(70,211,154,0.14);
      color: var(--ok);
    }
    .badge.warn {
      background: rgba(255,180,84,0.14);
      color: var(--warn);
    }
    .spinner {
      width: 18px;
      height: 18px;
      border: 2px solid rgba(77,226,197,0.16);
      border-top-color: var(--accent);
      border-radius: 50%;
      animation: spin .8s linear infinite;
      display: inline-block;
      vertical-align: middle;
      margin-right: 8px;
    }
    @keyframes spin {
      to { transform: rotate(360deg); }
    }
    .actions {
      display: flex;
      justify-content: center;
      gap: 12px;
      flex-wrap: wrap;
      margin-top: 10px;
    }
    button {
      border: 0;
      border-radius: 999px;
      padding: 16px 36px;
      min-width: 280px;
      font: inherit;
      font-weight: 800;
      cursor: pointer;
      transition: transform .16s ease, opacity .16s ease, box-shadow .16s ease, background .16s ease;
    }
    button:hover { transform: translateY(-1px); }
    button:disabled { opacity: .45; cursor: wait; transform: none; }
    .primary {
      color: #051014;
      background: linear-gradient(135deg, var(--accent), #8bf6df);
      box-shadow: 0 16px 34px rgba(77,226,197,0.18);
    }
    .secondary {
      min-width: 140px;
      color: var(--text);
      background: rgba(255,255,255,0.06);
      border: 1px solid rgba(255,255,255,0.08);
      box-shadow: none;
    }
    .status {
      margin-top: 16px;
      padding: 12px 14px;
      border-radius: 16px;
      display: none;
      font-size: 14px;
      line-height: 1.6;
    }
    .status.show { display: block; }
    .status.ok {
      background: rgba(70,211,154,0.12);
      color: var(--ok);
      border: 1px solid rgba(70,211,154,0.14);
    }
    .status.bad {
      background: rgba(255,107,122,0.1);
      color: var(--bad);
      border: 1px solid rgba(255,107,122,0.14);
    }
    .console {
      background: var(--panel-2);
      border: 1px solid var(--line);
      border-radius: 28px;
      overflow: hidden;
      box-shadow: var(--shadow);
      height: 100%;
      display: flex;
      flex-direction: column;
    }
    .console-head {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 16px 18px;
      background:
        linear-gradient(180deg, rgba(255,255,255,0.03), rgba(255,255,255,0)),
        linear-gradient(90deg, rgba(77,226,197,0.08), rgba(106,169,255,0.04));
      border-bottom: 1px solid rgba(255,255,255,0.05);
    }
    .console-title {
      color: var(--text);
      font-size: 13px;
      font-weight: 800;
      letter-spacing: 0.08em;
      text-transform: uppercase;
    }
    .console-meta {
      color: var(--muted);
      font-size: 12px;
    }
    .console-toolbar {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 12px;
      padding: 12px 16px;
      border-bottom: 1px solid rgba(255,255,255,0.05);
      background: rgba(255,255,255,0.02);
    }
    .console-ledger {
      display: flex;
      align-items: center;
      gap: 8px;
      flex-wrap: wrap;
    }
    .terminal-dot {
      width: 10px;
      height: 10px;
      border-radius: 50%;
      box-shadow: 0 0 12px currentColor;
    }
    .terminal-dot.red { color: #ff6b7a; background: #ff6b7a; }
    .terminal-dot.yellow { color: #ffb454; background: #ffb454; }
    .terminal-dot.green { color: #46d39a; background: #46d39a; }
    .console-hint {
      color: var(--muted);
      font-size: 11px;
      letter-spacing: 0.04em;
      text-transform: uppercase;
    }
    .console-body {
      padding: 14px 16px 0;
      max-height: 720px;
      min-height: 0;
      flex: 1 1 auto;
      overflow: auto;
    }
    .market-console {
      display: grid;
      gap: 12px;
      padding-bottom: 16px;
    }
    .market-console.hidden {
      display: none;
    }
    .market-stream-card {
      border: 1px solid rgba(130,164,255,0.12);
      border-radius: 20px;
      background: rgba(255,255,255,0.03);
      overflow: hidden;
    }
    .market-stream-card.collapsed .market-log-list {
      display: none;
    }
    .market-stream-card.ok {
      border-color: rgba(70,211,154,0.24);
    }
    .market-stream-card.bad {
      border-color: rgba(255,107,122,0.24);
    }
    .market-stream-card.active {
      border-color: rgba(77,226,197,0.28);
      box-shadow: 0 0 0 1px rgba(77,226,197,0.08);
    }
    .market-stream-head {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 16px;
      padding: 12px 14px 10px;
      border-bottom: 1px solid rgba(255,255,255,0.05);
      background: rgba(255,255,255,0.03);
    }
    .market-stream-title {
      display: grid;
      gap: 6px;
      flex: 1 1 auto;
    }
    .market-stream-name {
      font-size: 15px;
      font-weight: 800;
      letter-spacing: -0.02em;
    }
    .market-stream-actions {
      display: flex;
      align-items: center;
      gap: 8px;
      flex-wrap: wrap;
      justify-content: flex-end;
    }
    .market-stream-btn {
      min-width: 0;
      padding: 7px 10px;
      border-radius: 999px;
      border: 1px solid rgba(255,255,255,0.08);
      background: rgba(255,255,255,0.04);
      color: var(--text);
      font-size: 11px;
      font-weight: 800;
      letter-spacing: 0.04em;
      text-transform: uppercase;
      cursor: pointer;
      box-shadow: none;
    }
    .market-stream-btn:hover {
      border-color: rgba(77,226,197,0.28);
      background: rgba(77,226,197,0.08);
    }
    .market-stream-meta {
      font-size: 12px;
      color: var(--muted);
      line-height: 1.7;
    }
    .market-stream-body {
      padding: 12px 14px 14px;
      display: grid;
      gap: 10px;
    }
    .market-stream-card.collapsed .market-stream-body {
      gap: 10px;
    }
    .market-phase {
      font-size: 11px;
      color: var(--muted);
      line-height: 1.6;
    }
    .market-progress {
      display: grid;
      gap: 8px;
    }
    .market-progress-track {
      position: relative;
      height: 10px;
      border-radius: 999px;
      background: rgba(255,255,255,0.06);
      overflow: hidden;
    }
    .market-progress-fill {
      height: 100%;
      width: 0%;
      border-radius: inherit;
      background: linear-gradient(90deg, var(--accent), #8bf6df);
      transition: width .24s ease;
    }
    .market-progress-meta {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 12px;
      font-size: 10px;
      color: var(--muted);
      line-height: 1.6;
    }
    .market-log-list {
      display: grid;
      gap: 8px;
      max-height: 220px;
      overflow: auto;
      padding-right: 4px;
    }
    .market-log-empty {
      font-size: 12px;
      color: var(--muted);
      line-height: 1.8;
    }
    .market-log-line {
      display: grid;
      grid-template-columns: 54px 62px 1fr;
      gap: 10px;
      font-size: 12px;
      line-height: 1.7;
      color: #dbe6ff;
      padding: 8px 0;
      border-bottom: 1px solid rgba(255,255,255,0.04);
    }
    .market-log-line:last-child {
      border-bottom: 0;
    }
    .market-log-time {
      color: rgba(237,243,255,0.36);
      font-size: 11px;
      letter-spacing: 0.03em;
      padding-top: 2px;
    }
    .market-log-kind {
      font-size: 11px;
      font-weight: 800;
      letter-spacing: 0.08em;
      text-transform: uppercase;
      opacity: 0.94;
    }
    .market-log-text {
      white-space: pre-wrap;
      word-break: break-word;
    }
    .market-log-line.info .market-log-kind { color: #91b8ff; }
    .market-log-line.phase .market-log-kind { color: #6ee7ff; }
    .market-log-line.progress .market-log-kind { color: var(--accent); }
    .market-log-line.success .market-log-kind { color: var(--ok); }
    .market-log-line.error .market-log-kind { color: var(--bad); }
    .market-log-line.json .market-log-kind { color: #c5a5ff; }
    .task-empty {
      display: grid;
      gap: 10px;
      padding: 18px 2px 22px;
      color: var(--muted);
      border-bottom: 1px solid rgba(255,255,255,0.05);
    }
    .task-empty strong {
      color: var(--text);
      font-size: 14px;
      letter-spacing: 0.02em;
    }
    .task-empty span {
      font-size: 12px;
      line-height: 1.8;
    }
    .console-summary {
      margin: 0 16px 14px;
      padding: 12px 14px;
      border-radius: 18px;
      border: 1px solid rgba(130,164,255,0.16);
      background: rgba(255,255,255,0.03);
      display: none;
      gap: 8px;
    }
    .console-summary.show {
      display: grid;
    }
    .console-summary.ok {
      border-color: rgba(70,211,154,0.2);
      background: rgba(70,211,154,0.08);
    }
    .console-summary.bad {
      border-color: rgba(255,107,122,0.2);
      background: rgba(255,107,122,0.08);
    }
    .summary-title {
      font-size: 12px;
      font-weight: 800;
      letter-spacing: 0.08em;
      text-transform: uppercase;
    }
    .summary-text {
      font-size: 13px;
      line-height: 1.7;
      color: var(--text);
    }
    .result-panel {
      margin: 0 16px 16px;
      border: 1px solid rgba(255,255,255,0.08);
      border-radius: 16px;
      background: rgba(255,255,255,0.03);
      overflow: hidden;
    }
    .result-panel summary {
      cursor: pointer;
      list-style: none;
      padding: 14px 16px;
      color: var(--text);
      font-size: 13px;
      font-weight: 800;
      background: rgba(255,255,255,0.03);
    }
    .result-panel summary::-webkit-details-marker {
      display: none;
    }
    .result-json {
      margin: 0;
      padding: 0 16px 16px;
      color: #dbe6ff;
      font-size: 12px;
      line-height: 1.7;
      white-space: pre-wrap;
      word-break: break-word;
    }
    .modal {
      position: fixed;
      inset: 0;
      display: none;
      align-items: center;
      justify-content: center;
      padding: 20px;
      background: rgba(4,6,10,0.72);
      backdrop-filter: blur(12px);
      z-index: 1000;
    }
    .modal.show {
      display: flex;
    }
    .modal-card {
      width: min(760px, 100%);
      background: linear-gradient(180deg, rgba(14,18,25,0.98), rgba(11,15,21,0.96));
      border: 1px solid rgba(130,164,255,0.16);
      border-radius: 30px;
      box-shadow: 0 30px 90px rgba(0,0,0,0.55);
      padding: 26px;
    }
    .modal-head {
      display: flex;
      justify-content: space-between;
      gap: 16px;
      align-items: flex-start;
      margin-bottom: 18px;
    }
    .modal-title {
      margin: 0 0 6px;
      font-size: 30px;
      line-height: 1.03;
      letter-spacing: -0.05em;
    }
    .modal-sub {
      margin: 0;
      color: var(--muted);
      font-size: 14px;
      line-height: 1.7;
    }
    .confirm-grid {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 12px;
    }
    .confirm-item {
      border: 1px solid rgba(255,255,255,0.06);
      background: rgba(255,255,255,0.03);
      border-radius: 18px;
      padding: 14px;
    }
    .confirm-item.full {
      grid-column: 1 / -1;
    }
    .confirm-label {
      font-size: 12px;
      color: var(--muted);
      font-weight: 700;
      letter-spacing: 0.04em;
      text-transform: uppercase;
      margin-bottom: 6px;
    }
    .confirm-value {
      font-size: 15px;
      line-height: 1.7;
      word-break: break-word;
    }
    .confirm-markets {
      display: flex;
      flex-wrap: wrap;
      gap: 10px;
    }
    .confirm-market {
      border-radius: 999px;
      padding: 8px 12px;
      background: rgba(77,226,197,0.1);
      color: var(--accent);
      font-size: 13px;
      font-weight: 800;
    }
    .modal-actions {
      display: flex;
      justify-content: flex-end;
      gap: 12px;
      margin-top: 22px;
    }
    @media (max-width: 900px) {
      .workspace {
        grid-template-columns: 1fr;
      }
      .hero-side {
        grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
      }
      .console-body {
        min-height: 360px;
      }
    }
    @media (max-width: 720px) {
      .shell {
        padding: 20px 14px 32px;
      }
      .hero {
        padding: 18px;
      }
      .hero-bar {
        flex-direction: column;
        align-items: stretch;
      }
      .hero-actions {
        justify-content: flex-start;
      }
      .hero-side {
        grid-template-columns: 1fr;
      }
      .market-log-line {
        grid-template-columns: 1fr;
        gap: 4px;
      }
      .market-stream-head {
        flex-direction: column;
      }
      .market-stream-actions {
        justify-content: flex-start;
      }
      .feishu-meta-grid {
        grid-template-columns: 1fr;
      }
      .feishu-store-row {
        flex-direction: column;
      }
      .confirm-grid {
        grid-template-columns: 1fr;
      }
      .modal-card {
        padding: 20px;
      }
    }
  </style>
</head>
<body>
  <div class="shell">
    <section class="hero">
      <div class="hero-top">
        <div class="hero-bar">
          <h1 id="page-title">Upload</h1>
          <div class="hero-actions">
            <a class="hero-link" href="/apps">App 管理</a>
            <a class="hero-link" href="/config">配置中心</a>
            <a class="hero-link" href="/history">发布记录</a>
            <a class="hero-link" href="/audit">审核查询</a>
          </div>
        </div>
        <div class="hero-side">
          <div class="hero-stat">
            <strong>Artifact Input</strong>
            <span>zip 或 APK / 渠道包自动识别</span>
          </div>
          <div class="hero-stat">
            <strong>Dispatch Mode</strong>
            <span>审核后自动发布 / 定时发布</span>
          </div>
          <div class="hero-stat">
            <strong>Notify Chain</strong>
            <span>发布完成后自动走全局 Hook 与飞书通知</span>
          </div>
        </div>
      </div>
    </section>

    <section class="workspace">
      <section class="panel">
        <div class="panel-head">
          <div class="panel-kicker">deploy flow</div>
          <h2>发布执行</h2>
          <p class="note">支持上传 1 个 zip，或者直接上传 1 个或多个 APK。系统会自动遍历 zip 内所有层级的文件夹并识别渠道包，但不能在同一次里混传 zip 和 APK。</p>
        </div>
        <div class="panel-body">
          <form id="upload-form" class="upload-grid">
            <div class="field">
              <label class="small" for="archive">上传文件</label>
              <input id="archive" name="archive" type="file" accept=".zip,.apk" multiple required>
            </div>

            <div class="field hidden" id="market-field">
              <label class="small">识别到的市场</label>
              <div id="market-box" class="market-box">
                <div id="market-list" class="market-list"></div>
              </div>
            </div>

            <div class="field">
              <label class="small" for="publish_mode">发布方式</label>
              <select id="publish_mode" name="publish_mode">
                <option value="auto">审核后自动发布</option>
                <option value="scheduled">定时发布</option>
              </select>
            </div>

            <label class="field hidden field-card" id="publish-time-field" for="publish_time">
              <span class="small">发布时间</span>
              <span class="field-card-meta">点击整块区域选择发布时间，仅在定时发布时需要设置。</span>
              <input id="publish_time" name="publish_time" type="datetime-local">
            </label>

            <div class="field">
              <label class="small" for="notes">更新说明</label>
              <textarea id="notes" name="notes" placeholder="修复登录问题、提升稳定性"></textarea>
            </div>

            <div class="actions">
              <button type="submit" class="primary" id="upload-btn" disabled>开始发布</button>
            </div>
            <div id="upload-status" class="status"></div>
          </form>
        </div>
      </section>

      <section class="console">
        <div class="console-head">
          <div class="console-title">任务</div>
          <div class="console-meta hidden" id="console-meta"></div>
        </div>
        <div class="console-toolbar">
          <div class="console-ledger">
            <span class="terminal-dot red"></span>
            <span class="terminal-dot yellow"></span>
            <span class="terminal-dot green"></span>
          </div>
          <div class="console-hint">selected store tasks / progress / logs</div>
        </div>
        <div class="console-body">
          <div id="market-console" class="market-console hidden"></div>
          <div id="task-empty" class="task-empty">
            <strong>任务面板已就绪</strong>
            <span>上传文件并识别出渠道包后，只要你勾选市场，右侧就会立刻出现对应的任务卡片。后续发布进度、状态和日志都显示在各自卡片里。</span>
          </div>
        </div>
        <div id="console-summary" class="console-summary hidden">
          <div id="summary-title" class="summary-title"></div>
          <div id="summary-text" class="summary-text"></div>
        </div>
        <details id="result-panel" class="result-panel hidden">
          <summary>最终 JSON 结果</summary>
          <pre id="result-json" class="result-json"></pre>
        </details>
      </section>
    </section>
  </div>

  <div id="confirm-modal" class="modal" aria-hidden="true">
    <div class="modal-card">
      <div class="modal-head">
        <div>
          <h3 class="modal-title">确认本次发布</h3>
          <p class="modal-sub">请再确认一次这次要提交的文件、发布方式、更新说明和目标市场。确认后会立即开始执行发布。</p>
        </div>
      </div>

      <div class="confirm-grid">
        <div class="confirm-item">
          <div class="confirm-label">上传文件</div>
          <div id="confirm-archive" class="confirm-value"></div>
        </div>
        <div class="confirm-item">
          <div class="confirm-label">发布方式</div>
          <div id="confirm-mode" class="confirm-value"></div>
        </div>
        <div class="confirm-item">
          <div class="confirm-label">发布时间</div>
          <div id="confirm-time" class="confirm-value"></div>
        </div>
        <div class="confirm-item">
          <div class="confirm-label">目标市场数</div>
          <div id="confirm-count" class="confirm-value"></div>
        </div>
        <div class="confirm-item full">
          <div class="confirm-label">更新说明</div>
          <div id="confirm-notes" class="confirm-value"></div>
        </div>
        <div class="confirm-item full">
          <div class="confirm-label">目标市场</div>
          <div id="confirm-markets" class="confirm-markets"></div>
        </div>
      </div>

      <div class="modal-actions">
        <button type="button" id="confirm-cancel" class="secondary">返回修改</button>
        <button type="button" id="confirm-submit" class="primary">确认发布</button>
      </div>
    </div>
  </div>

	  <script>
	    let detectedArtifacts = [];
	    let inspectedUploadKey = null;
	    let selectedStores = new Set();
	    let pendingSubmit = null;
      let marketStreams = new Map();
      let currentSelectionMode = 'auto';

    function setStatus(id, message, kind) {
      const el = document.getElementById(id);
      if (!message) {
        el.textContent = '';
        el.className = 'status';
        return;
      }
      el.textContent = message;
      el.className = 'status show ' + kind;
    }

    function syncPublishModeUI() {
      const mode = document.getElementById('publish_mode').value;
      const wrap = document.getElementById('publish-time-field');
      const input = document.getElementById('publish_time');
      const scheduled = mode === 'scheduled';
      wrap.classList.toggle('hidden', !scheduled);
      input.required = scheduled;
      if (!scheduled) input.value = '';
    }

    function formatPublishTime(value) {
      if (!value) return '';
      return value.replace('T', ' ') + ':00';
    }

    function syncPublishTimeMin() {
      const input = document.getElementById('publish_time');
      const now = new Date();
      now.setSeconds(0, 0);
      const pad = (n) => String(n).padStart(2, '0');
      const local =
        now.getFullYear() + '-' +
        pad(now.getMonth() + 1) + '-' +
        pad(now.getDate()) + 'T' +
        pad(now.getHours()) + ':' +
        pad(now.getMinutes());
      input.min = local;
      if (input.value && input.value < local) input.value = local;
    }

    function escapeHtml(value) {
      return String(value || '')
        .replaceAll('&', '&amp;')
        .replaceAll('<', '&lt;')
        .replaceAll('>', '&gt;')
        .replaceAll('"', '&quot;')
        .replaceAll("'", '&#39;');
    }

    async function applyCurrentAppTitle() {
      try {
        const resp = await fetch('/api/app/current');
        const data = await resp.json();
        if (!resp.ok) return;
        const appName = data?.name || 'APKGO';
        document.getElementById('page-title').textContent = appName + ' Upload';
        document.title = appName + ' Upload';
      } catch (_) {}
    }

    function setMarketLoading() {
      const field = document.getElementById('market-field');
      const box = document.getElementById('market-box');
      const list = document.getElementById('market-list');
      field.classList.remove('hidden');
      box.classList.add('loading');
      list.innerHTML = '<span><span class="spinner"></span>正在识别上传文件中的渠道包...</span>';
    }

    function updateSubmitState() {
      const hasSelectable = detectedArtifacts.some(item => item.configured);
      const hasSelectedConfigured = detectedArtifacts.some(item => item.configured && selectedStores.has(item.store));
      document.getElementById('upload-btn').disabled = !(hasSelectable && hasSelectedConfigured);
    }

    function toggleStore(store, checked) {
      if (checked) selectedStores.add(store);
      else selectedStores.delete(store);
      renderMarkets(detectedArtifacts);
      syncTaskCardsFromSelection();
      updateSubmitState();
    }

    function renderMarkets(items) {
      const field = document.getElementById('market-field');
      const box = document.getElementById('market-box');
      const list = document.getElementById('market-list');
      field.classList.remove('hidden');
      box.classList.remove('loading');
      list.innerHTML = '';
      for (const item of items) {
        const checked = selectedStores.has(item.store);
        const card = document.createElement('label');
        card.className = 'market-card' + (checked ? ' selected' : '') + (item.configured ? '' : ' disabled');
        card.innerHTML =
          '<span class="badge ' + (item.configured ? 'ok' : 'warn') + '">' + (item.configured ? '已配置 key' : '未配置 key') + '</span>' +
          '<input class="market-check" type="checkbox" ' + (item.configured ? '' : 'disabled') + (checked ? ' checked' : '') + '>' +
          '<div class="market-title">' + escapeHtml(item.display_name || item.store) + '</div>' +
          '<div class="market-meta">市场：' + escapeHtml(item.store) + '</div>' +
          '<div class="market-meta">渠道：' + escapeHtml(item.channel || '手动选择') + '</div>';
        const input = card.querySelector('input');
        input.addEventListener('change', () => toggleStore(item.store, input.checked));
        list.appendChild(card);
      }
    }

    function clearMarkets() {
      detectedArtifacts = [];
      inspectedUploadKey = null;
      selectedStores = new Set();
      currentSelectionMode = 'auto';
      document.getElementById('market-field').classList.add('hidden');
      document.getElementById('market-list').innerHTML = '';
      document.getElementById('upload-btn').disabled = true;
      resetLogs();
    }

    function getSelectedFiles() {
      return Array.from(document.getElementById('archive').files || []);
    }

    function buildUploadKey(files) {
      return files
        .map(file => file.name + ':' + file.size + ':' + file.lastModified)
        .sort()
        .join('|');
    }

    function describeUploadFiles(files) {
      if (!files.length) return '';
      if (files.length === 1) return files[0].name;
      return files.length + ' 个文件';
    }

    function validateSelectedFiles(files) {
      if (!files.length) return '请先选择 zip 或 APK 文件。';
      const zipFiles = files.filter(file => file.name.toLowerCase().endsWith('.zip'));
      const apkFiles = files.filter(file => file.name.toLowerCase().endsWith('.apk'));
      if (zipFiles.length && apkFiles.length) return '不能同时上传 zip 和 APK，请二选一。';
      if (zipFiles.length > 1) return '暂时只支持上传 1 个 zip，或上传多个 APK。';
      if (!zipFiles.length && apkFiles.length !== files.length) return '当前只支持 zip 或 APK 文件。';
      return '';
    }

	    function resetLogs() {
        marketStreams = new Map();
        document.getElementById('market-console').innerHTML = '';
        document.getElementById('market-console').classList.add('hidden');
        document.getElementById('task-empty').style.display = '';
	      const summary = document.getElementById('console-summary');
	      summary.className = 'console-summary hidden';
	      document.getElementById('summary-title').textContent = '';
	      document.getElementById('summary-text').textContent = '';
	      document.getElementById('result-panel').classList.add('hidden');
	      document.getElementById('result-json').textContent = '';
	    }

	    function nowTimeLabel() {
	      const now = new Date();
	      const pad = (n) => String(n).padStart(2, '0');
	      return pad(now.getHours()) + ':' + pad(now.getMinutes()) + ':' + pad(now.getSeconds());
	    }

	    function showSummary(kind, title, text) {
	      const summary = document.getElementById('console-summary');
	      summary.className = 'console-summary show ' + kind;
	      document.getElementById('summary-title').textContent = title;
	      document.getElementById('summary-text').textContent = text;
	    }

    function storeName(store) {
      const named = {
        huawei: '华为',
        xiaomi: '小米',
        oppo: 'OPPO',
        vivo: 'vivo',
        honor: '荣耀',
        tencent: '应用宝',
        samsung: 'Samsung',
        googleplay: 'Google Play',
        pgyer: '蒲公英',
        fir: 'fir.im',
        script: 'Script',
      };
      return named[store] || store || '-';
    }

    function ensureMarketCard(store) {
      if (!store) return null;
      if (marketStreams.has(store)) return marketStreams.get(store);

      const root = document.getElementById('market-console');
      root.classList.remove('hidden');
      document.getElementById('task-empty').style.display = 'none';
      const card = document.createElement('section');
      card.className = 'market-stream-card active collapsed';
      card.innerHTML =
        '<div class="market-stream-head">' +
          '<div class="market-stream-title">' +
            '<div class="market-stream-name">' + escapeHtml(storeName(store)) + '</div>' +
            '<div class="market-stream-meta">市场标识：' + escapeHtml(store) + '</div>' +
          '</div>' +
          '<div class="market-stream-actions">' +
            '<button type="button" class="market-stream-btn" data-role="copy">复制日志</button>' +
            '<button type="button" class="market-stream-btn" data-role="toggle">展开日志</button>' +
            '<span class="badge info" data-role="badge">等待中</span>' +
          '</div>' +
        '</div>' +
        '<div class="market-stream-body">' +
          '<div class="market-phase" data-role="phase">等待任务开始...</div>' +
          '<div class="market-progress">' +
            '<div class="market-progress-track"><div class="market-progress-fill" data-role="fill"></div></div>' +
            '<div class="market-progress-meta"><span data-role="percent">0%</span><span data-role="progress-text">尚未开始上传</span></div>' +
          '</div>' +
          '<div class="market-log-list" data-role="logs">' +
            '<div class="market-log-empty">当前市场的阶段日志、进度日志和接口返回会显示在这里。</div>' +
          '</div>' +
        '</div>';
      root.appendChild(card);
      const entry = {
        card,
        badge: card.querySelector('[data-role="badge"]'),
        phase: card.querySelector('[data-role="phase"]'),
        fill: card.querySelector('[data-role="fill"]'),
        percent: card.querySelector('[data-role="percent"]'),
        progressText: card.querySelector('[data-role="progress-text"]'),
        logs: card.querySelector('[data-role="logs"]'),
        toggle: card.querySelector('[data-role="toggle"]'),
        copy: card.querySelector('[data-role="copy"]'),
        lines: [],
        collapsed: true,
      };
      entry.toggle.addEventListener('click', () => toggleMarketCard(store));
      entry.copy.addEventListener('click', async () => {
        const payload = entry.lines.join('\n').trim();
        if (!payload) {
          setStatus('upload-status', storeName(store) + ' 当前没有可复制的日志。', 'bad');
          return;
        }
        try {
          await navigator.clipboard.writeText(payload);
          setStatus('upload-status', storeName(store) + ' 日志已复制到剪贴板。', 'ok');
        } catch (err) {
          setStatus('upload-status', '复制失败：' + String(err), 'bad');
        }
      });
      marketStreams.set(store, entry);
      return entry;
    }

    function removeMarketCard(store) {
      const entry = marketStreams.get(store);
      if (!entry) return;
      entry.card.remove();
      marketStreams.delete(store);
      if (!marketStreams.size) {
        document.getElementById('market-console').classList.add('hidden');
        document.getElementById('task-empty').style.display = '';
      }
    }

    function syncTaskCardsFromSelection() {
      const selected = new Set(selectedStores);
      for (const store of selected) {
        ensureMarketCard(store);
      }
      for (const store of Array.from(marketStreams.keys())) {
        if (!selected.has(store)) {
          removeMarketCard(store);
        }
      }
    }

    function setMarketCollapsed(entry, collapsed) {
      if (!entry) return;
      entry.collapsed = collapsed;
      entry.card.classList.toggle('collapsed', collapsed);
      entry.toggle.textContent = collapsed ? '展开日志' : '收起日志';
    }

    function toggleMarketCard(store) {
      const entry = ensureMarketCard(store);
      if (!entry) return;
      setMarketCollapsed(entry, !entry.collapsed);
    }

    function setMarketBadge(entry, label, cls) {
      if (!entry) return;
      entry.badge.className = 'badge ' + cls;
      entry.badge.textContent = label;
    }

    function updateMarketProgress(store, sent, total) {
      const entry = ensureMarketCard(store);
      if (!entry) return;
      const pct = total > 0 ? Math.max(0, Math.min(100, Math.round((sent / total) * 100))) : 0;
      entry.fill.style.width = pct + '%';
      entry.percent.textContent = pct + '%';
      entry.progressText.textContent = total > 0
        ? (sent + ' / ' + total + ' bytes')
        : String(sent || 0);
    }

    function appendMarketLog(store, kind, line) {
      const entry = ensureMarketCard(store);
      if (!entry) return;
      const empty = entry.logs.querySelector('.market-log-empty');
      if (empty) empty.remove();
      const time = nowTimeLabel();
      entry.lines.push('[' + time + '] [' + kind.toUpperCase() + '] ' + line);
      const row = document.createElement('div');
      row.className = 'market-log-line ' + kind;
      row.innerHTML =
        '<div class="market-log-time">' + time + '</div>' +
        '<div class="market-log-kind">' + escapeHtml(kind) + '</div>' +
        '<div class="market-log-text">' + escapeHtml(line) + '</div>';
      entry.logs.appendChild(row);
      entry.logs.scrollTop = entry.logs.scrollHeight;
    }

    function formatBytes(sent, total) {
      if (!total) return String(sent || 0);
      const pct = total > 0 ? Math.min(100, Math.round((sent / total) * 100)) : 0;
      return pct + '% (' + sent + '/' + total + ' bytes)';
    }

    function renderStreamEvent(event) {
      switch (event.type) {
        case 'start':
          const startStores = event.stores || event.data?.stores || [];
          for (const store of startStores) {
            const entry = ensureMarketCard(store);
            setMarketBadge(entry, '排队中', 'info');
          }
          break;
        case 'phase':
          if (event.store) {
            const entry = ensureMarketCard(event.store);
            entry.phase.textContent = '当前阶段：' + (event.phase || event.data?.phase || '处理中');
            setMarketBadge(entry, '进行中', 'warn');
            entry.card.classList.add('active');
            appendMarketLog(event.store, 'phase', '阶段：' + (event.phase || event.data?.phase || ''));
          }
          break;
        case 'total':
          if (event.store) {
            const totalBytes = event.total_bytes || event.data?.total_bytes || 0;
            updateMarketProgress(event.store, 0, totalBytes);
            appendMarketLog(event.store, 'info', '文件大小：' + totalBytes + ' bytes');
          }
          break;
        case 'bytes':
          if (event.store) {
            const sent = event.sent ?? event.data?.sent ?? 0;
            const total = event.total ?? event.data?.total ?? 0;
            updateMarketProgress(event.store, sent, total);
            appendMarketLog(event.store, 'progress', '上传进度：' + formatBytes(sent, total));
          }
          break;
        case 'result':
          if (event.store) {
            const success = event.success ?? event.data?.success ?? false;
            const err = event.error || event.data?.error || '';
            const durationMs = event.duration_ms || event.data?.duration_ms || 0;
            const entry = ensureMarketCard(event.store);
            entry.phase.textContent = success ? '当前阶段：已完成' : '当前阶段：执行失败';
            setMarketBadge(entry, success ? '成功' : '失败', success ? 'ok' : 'bad');
            entry.card.classList.remove('active');
            entry.card.classList.add(success ? 'ok' : 'bad');
            updateMarketProgress(event.store, success ? 1 : 0, 1);
            appendMarketLog(event.store, success ? 'success' : 'error', '结果：' + (success ? '成功' : '失败') + (err ? '，' + err : '') + (durationMs ? '，耗时 ' + durationMs + ' ms' : ''));
            break;
          }
          document.getElementById('result-panel').classList.remove('hidden');
          document.getElementById('result-json').textContent = JSON.stringify(event.data, null, 2);
          break;
        case 'done':
            const doneResults = event.results || event.data?.results;
	          if (doneResults) {
	            const successCount = doneResults.filter(item => item.success).length;
	            const failCount = doneResults.length - successCount;
	            showSummary(
	              failCount > 0 ? 'bad' : 'ok',
	              failCount > 0 ? '发布完成，但有失败市场' : '全部市场发布完成',
	              '成功 ' + successCount + ' 个，失败 ' + failCount + ' 个。'
	            );
	          }
	          if (doneResults) {
	            document.getElementById('result-panel').classList.remove('hidden');
	            document.getElementById('result-json').textContent = JSON.stringify({ apk: event.apk || event.data?.apk, results: doneResults }, null, 2);
	          }
	          break;
        case 'store.start':
          if (event.store) {
            const entry = ensureMarketCard(event.store);
            setMarketBadge(entry, '启动中', 'warn');
            entry.phase.textContent = '当前阶段：准备开始';
            appendMarketLog(event.store, 'info', event.message || JSON.stringify(event));
          }
          break;
        case 'store.done':
          if (event.store) {
            const entry = ensureMarketCard(event.store);
            setMarketBadge(entry, '已完成', 'ok');
            entry.phase.textContent = '当前阶段：发布完成';
            entry.card.classList.remove('active');
            entry.card.classList.add('ok');
            updateMarketProgress(event.store, 1, 1);
            appendMarketLog(event.store, 'success', event.message || JSON.stringify(event));
          }
          break;
	        case 'store.error':
            if (event.store) {
              const entry = ensureMarketCard(event.store);
              setMarketBadge(entry, '失败', 'bad');
              entry.phase.textContent = '当前阶段：执行失败';
              entry.card.classList.remove('active');
              entry.card.classList.add('bad');
              appendMarketLog(event.store, 'error', event.message || JSON.stringify(event));
            }
	          showSummary('bad', '市场发布失败', event.message || '请查看错误日志');
	          break;
        case 'log':
          if (event.store) {
            appendMarketLog(event.store, 'info', event.message || JSON.stringify(event));
          }
          break;
        case 'feishu.sent':
          if (event.data?.ok) {
            setStatus('upload-status', event.message || '飞书通知已发送。', 'ok');
          } else {
            setStatus('upload-status', event.message || '飞书通知发送失败。', 'bad');
          }
          break;
	        default:
	          if (event.type === 'error') {
	            showSummary('bad', '发布执行失败', event.message || 'unknown error');
          }
        }
    }

    async function readStreamAsEvents(resp) {
      const reader = resp.body.getReader();
      const decoder = new TextDecoder();
      let buffer = '';
      while (true) {
        const { value, done } = await reader.read();
        if (done) break;
        buffer += decoder.decode(value, { stream: true });
        const parts = buffer.split('\n');
        buffer = parts.pop() || '';
        for (const line of parts) {
          const trimmed = line.trim();
          if (!trimmed) continue;
          try {
            renderStreamEvent(JSON.parse(trimmed));
          } catch (err) {
            setStatus('upload-status', '收到无法解析的流式事件，请检查接口输出。', 'bad');
          }
        }
      }
      if (buffer.trim()) {
        try {
          renderStreamEvent(JSON.parse(buffer.trim()));
        } catch (err) {
          setStatus('upload-status', '收到无法解析的流式尾数据，请检查接口输出。', 'bad');
        }
      }
    }

    function openConfirmModal(payload) {
      pendingSubmit = payload;
      document.getElementById('confirm-archive').textContent = payload.archiveName;
      document.getElementById('confirm-mode').textContent = payload.publishModeText;
      document.getElementById('confirm-time').textContent = payload.publishTimeText;
      document.getElementById('confirm-count').textContent = String(payload.selectedItems.length) + ' 个市场';
      document.getElementById('confirm-notes').textContent = payload.notesText;

      const marketRoot = document.getElementById('confirm-markets');
      marketRoot.innerHTML = '';
      for (const item of payload.selectedItems) {
        const chip = document.createElement('div');
        chip.className = 'confirm-market';
        chip.textContent = item.display_name || item.store;
        marketRoot.appendChild(chip);
      }

      const modal = document.getElementById('confirm-modal');
      modal.classList.add('show');
      modal.setAttribute('aria-hidden', 'false');
    }

    function closeConfirmModal() {
      pendingSubmit = null;
      const modal = document.getElementById('confirm-modal');
      modal.classList.remove('show');
      modal.setAttribute('aria-hidden', 'true');
    }

    async function doSubmit(payload) {
      const btn = document.getElementById('upload-btn');
      btn.disabled = true;
      setStatus('upload-status', '正在发布，请稍候...', 'ok');
      resetLogs();
      for (const store of payload.stores) {
        const entry = ensureMarketCard(store);
        setMarketBadge(entry, '排队中', 'info');
      }

      const fd = new FormData();
      for (const file of payload.files) {
        fd.append('archive', file);
      }
      fd.append('notes', payload.notes);
      fd.append('publish_mode', payload.publishMode);
      fd.append('publish_time', payload.publishTime);
      fd.append('stores', payload.stores.join(','));

      const resp = await fetch('/api/upload', { method: 'POST', body: fd });
      await readStreamAsEvents(resp);
	      if (!resp.ok) {
	        setStatus('upload-status', '发布失败', 'bad');
	        showSummary('bad', '发布执行失败', '接口已返回失败，请查看上方日志和最终 JSON 结果。');
	        btn.disabled = false;
	        return;
	      }
      setStatus('upload-status', '发布完成。', 'ok');
      btn.disabled = false;
    }

    async function inspectArchive(files) {
      const fileError = validateSelectedFiles(files);
      if (fileError) {
        clearMarkets();
        setStatus('upload-status', fileError, 'bad');
        return;
      }
      setMarketLoading();
      document.getElementById('upload-btn').disabled = true;
      const fd = new FormData();
      for (const file of files) {
        fd.append('archive', file);
      }
      const resp = await fetch('/api/inspect', { method: 'POST', body: fd });
      const data = await resp.json();
      if (!resp.ok) {
        clearMarkets();
        document.getElementById('result-panel').classList.remove('hidden');
        document.getElementById('result-json').textContent = JSON.stringify(data, null, 2);
        setStatus('upload-status', data.error || '识别失败，请检查配置和上传文件。', 'bad');
        return;
      }
      detectedArtifacts = data.upload?.artifacts || [];
      currentSelectionMode = data.upload?.selection_mode || 'auto';
      inspectedUploadKey = buildUploadKey(files);
      if (currentSelectionMode === 'manual') {
        selectedStores = new Set();
      } else {
        selectedStores = new Set(detectedArtifacts.filter(item => item.configured).map(item => item.store));
      }
      renderMarkets(detectedArtifacts);
      syncTaskCardsFromSelection();
      updateSubmitState();
      if (currentSelectionMode === 'manual') {
        setStatus('upload-status', '未匹配到市场别名，已展示所有已配置市场，请手动选择。', 'bad');
      } else {
        setStatus('upload-status', '', 'ok');
      }
    }

	    document.getElementById('archive').addEventListener('change', async (event) => {
	      const files = Array.from(event.target.files || []);
	      closeConfirmModal();
	      clearMarkets();
	      if (!files.length) {
	        return;
      }
      try {
        await inspectArchive(files);
      } catch (err) {
        clearMarkets();
        setStatus('upload-status', String(err), 'bad');
	      }
	    });

	    document.getElementById('publish_mode').addEventListener('change', syncPublishModeUI);
	    document.getElementById('publish-time-field').addEventListener('click', (event) => {
	      const input = document.getElementById('publish_time');
	      if (event.target !== input) {
	        input.focus();
	        if (typeof input.showPicker === 'function') input.showPicker();
	      }
	    });
	    document.getElementById('confirm-cancel').addEventListener('click', closeConfirmModal);
    document.getElementById('confirm-submit').addEventListener('click', async () => {
      if (!pendingSubmit) return;
      const payload = pendingSubmit;
      closeConfirmModal();
      await doSubmit(payload);
    });
    document.getElementById('confirm-modal').addEventListener('click', (event) => {
      if (event.target.id === 'confirm-modal') closeConfirmModal();
    });

    document.getElementById('upload-form').addEventListener('submit', async (event) => {
      event.preventDefault();
      const files = getSelectedFiles();
      const fileError = validateSelectedFiles(files);
      if (fileError) {
        setStatus('upload-status', fileError, 'bad');
        return;
      }
      const currentKey = buildUploadKey(files);
      if (!detectedArtifacts.length || inspectedUploadKey !== currentKey) {
        setStatus('upload-status', '请先等待文件识别完成。', 'bad');
        return;
      }
      if (!selectedStores.size) {
        setStatus('upload-status', '请至少选择一个已配置 key 的市场。', 'bad');
        return;
      }

      const selectedItems = detectedArtifacts.filter(item => selectedStores.has(item.store));
      const publishModeText = document.getElementById('publish_mode').value === 'scheduled'
        ? '定时发布'
        : '审核后自动发布';
      const publishTimeText = document.getElementById('publish_mode').value === 'scheduled'
        ? (formatPublishTime(document.getElementById('publish_time').value) || '未填写')
        : '无需设置';
      const notesText = (document.getElementById('notes').value || '').trim() || '未填写';

      openConfirmModal({
        files,
        archiveName: describeUploadFiles(files),
        publishMode: document.getElementById('publish_mode').value,
        publishModeText,
        publishTime: formatPublishTime(document.getElementById('publish_time').value),
        publishTimeText,
        notes: document.getElementById('notes').value,
        notesText,
        stores: Array.from(selectedStores),
        selectedItems,
      });
    });

    syncPublishTimeMin();
    applyCurrentAppTitle();
    syncPublishModeUI();
    resetLogs();
  </script>
</body>
</html>
`

const webConfigHTML = `<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>apkgo Config</title>
  <style>
    :root {
      --bg: #090b10;
      --panel: rgba(14,18,25,0.88);
      --line: rgba(130,164,255,0.14);
      --line-strong: rgba(130,164,255,0.28);
      --text: #edf3ff;
      --muted: #8d9ab3;
      --accent: #4de2c5;
      --ok: #46d39a;
      --warn: #ffb454;
      --bad: #ff6b7a;
      --info: #91b8ff;
      --shadow: 0 28px 80px rgba(0,0,0,0.45);
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      min-height: 100vh;
      color: var(--text);
      font-family: "SF Mono", "JetBrains Mono", "IBM Plex Sans", "PingFang SC", "Microsoft YaHei", monospace, sans-serif;
      background:
        radial-gradient(circle at 12% 10%, rgba(77,226,197,0.12), transparent 22%),
        radial-gradient(circle at 88% 14%, rgba(106,169,255,0.12), transparent 20%),
        linear-gradient(180deg, #07090d 0%, var(--bg) 100%);
    }
    .shell {
      max-width: 1120px;
      margin: 0 auto;
      padding: 28px 20px 52px;
    }
    .hero, .panel, .modal-card {
      background: var(--panel);
      border: 1px solid var(--line);
      box-shadow: var(--shadow);
    }
    .hero {
      position: relative;
      overflow: hidden;
      padding: 22px 24px;
      border-radius: 32px;
      margin-bottom: 20px;
      background:
        linear-gradient(135deg, rgba(18,24,34,0.96), rgba(10,14,20,0.94)),
        linear-gradient(180deg, rgba(77,226,197,0.05), rgba(106,169,255,0.04));
      border-color: var(--line-strong);
    }
    .hero::after {
      content: "";
      position: absolute;
      inset: -1px;
      border-radius: inherit;
      padding: 1px;
      background: linear-gradient(135deg, rgba(77,226,197,0.3), rgba(106,169,255,0.18), rgba(141,107,255,0.2));
      -webkit-mask: linear-gradient(#000 0 0) content-box, linear-gradient(#000 0 0);
      -webkit-mask-composite: xor;
      mask-composite: exclude;
      pointer-events: none;
    }
    .hero-top {
      display: grid;
      gap: 16px;
      position: relative;
      z-index: 1;
    }
    .hero-bar {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 16px;
    }
    .hero-link {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      min-height: 44px;
      padding: 0 18px;
      border-radius: 999px;
      border: 1px solid rgba(255,255,255,0.08);
      background: rgba(255,255,255,0.05);
      color: var(--text);
      text-decoration: none;
      font-size: 12px;
      font-weight: 800;
      letter-spacing: 0.08em;
      text-transform: uppercase;
    }
    h1 {
      margin: 0;
      font-size: clamp(34px, 6vw, 58px);
      line-height: 0.94;
      letter-spacing: -0.06em;
    }
    .hero-desc {
      max-width: 760px;
      color: var(--muted);
      font-size: 14px;
      line-height: 1.8;
    }
    .panel {
      border-radius: 28px;
      overflow: hidden;
      margin-bottom: 20px;
    }
    .panel-head {
      padding: 18px 22px 14px;
      border-bottom: 1px solid rgba(255,255,255,0.05);
      background: linear-gradient(180deg, rgba(255,255,255,0.02), rgba(255,255,255,0));
    }
    .panel-kicker {
      color: var(--accent);
      font-size: 11px;
      font-weight: 800;
      letter-spacing: 0.08em;
      text-transform: uppercase;
      margin-bottom: 10px;
    }
    .panel h2 {
      margin: 0;
      font-size: 24px;
      letter-spacing: -0.04em;
    }
    .panel-body {
      padding: 10px 22px 22px;
    }
    .config-group {
      margin-top: 20px;
      padding: 14px 14px 12px;
      border-radius: 24px;
      background: rgba(255,255,255,0.02);
      border: 1px solid rgba(130,164,255,0.12);
    }
    .group-title {
      margin: 0 0 10px;
      padding-bottom: 10px;
      border-bottom: 1px solid rgba(255,255,255,0.06);
      font-size: 14px;
      font-weight: 800;
      letter-spacing: 0.02em;
      color: var(--accent);
      text-transform: uppercase;
    }
    .config-list {
      display: grid;
      gap: 8px;
    }
    .config-row {
      display: grid;
      grid-template-columns: minmax(180px, 220px) minmax(0, 1fr) auto;
      gap: 12px;
      align-items: center;
      padding: 10px 12px;
      border-radius: 14px;
      background: rgba(255,255,255,0.025);
      border: 1px solid rgba(130,164,255,0.1);
    }
    .config-row.is-configured {
      background: rgba(70,211,154,0.05);
      border-color: rgba(70,211,154,0.18);
    }
    .config-row.is-unconfigured {
      background: rgba(255,255,255,0.02);
      border-color: rgba(255,255,255,0.07);
    }
    .config-meta {
      display: grid;
      gap: 2px;
      align-items: center;
    }
    .config-name {
      font-size: 14px;
      font-weight: 800;
    }
    .config-subtitle {
      color: var(--muted);
      font-size: 11px;
      line-height: 1.4;
      letter-spacing: 0.03em;
      text-transform: uppercase;
    }
    .config-summary {
      color: var(--text);
      font-size: 13px;
      line-height: 1.5;
      font-weight: 500;
      opacity: 0.92;
      word-break: break-word;
    }
    .config-summary.ok {
      color: var(--ok);
    }
    .config-summary.pending {
      color: var(--muted);
    }
    .config-action {
      border: 1px solid rgba(77,226,197,0.18);
      border-radius: 999px;
      padding: 8px 14px;
      min-width: 72px;
      font: inherit;
      font-size: 12px;
      font-weight: 800;
      cursor: pointer;
      color: var(--accent);
      background: rgba(77,226,197,0.08);
      box-shadow: none;
    }
    .config-action[disabled],
    .ghost-btn[disabled] {
      opacity: 0.62;
      cursor: wait;
    }
    .status {
      padding: 12px 14px;
      border-radius: 16px;
      display: none;
      font-size: 14px;
      line-height: 1.6;
      margin-bottom: 16px;
    }
    .status.show { display: block; }
    .status.ok {
      background: rgba(70,211,154,0.12);
      color: var(--ok);
      border: 1px solid rgba(70,211,154,0.14);
    }
    .status.bad {
      background: rgba(255,107,122,0.1);
      color: var(--bad);
      border: 1px solid rgba(255,107,122,0.14);
    }
    .status.info {
      background: rgba(145,184,255,0.1);
      color: var(--info);
      border: 1px solid rgba(145,184,255,0.14);
    }
    .modal {
      position: fixed;
      inset: 0;
      display: none;
      align-items: center;
      justify-content: center;
      padding: 24px;
      background: rgba(3,6,10,0.72);
      backdrop-filter: blur(10px);
      z-index: 40;
    }
    .modal.show {
      display: flex;
    }
    .modal-card {
      width: min(760px, 100%);
      border-radius: 28px;
      overflow: hidden;
    }
    .modal-head {
      padding: 20px 22px 14px;
      border-bottom: 1px solid rgba(255,255,255,0.05);
    }
    .modal-body {
      padding: 20px 22px 22px;
      display: grid;
      gap: 14px;
    }
    .modal-status {
      display: none;
      padding: 12px 14px;
      border-radius: 16px;
      border: 1px solid rgba(255,255,255,0.08);
      font-size: 13px;
      line-height: 1.6;
    }
    .modal-status.show {
      display: block;
    }
    .modal-status.bad {
      color: #ffd5da;
      background: rgba(255,107,122,0.12);
      border-color: rgba(255,107,122,0.28);
    }
    .modal-status.ok {
      color: #d6ffef;
      background: rgba(70,211,154,0.12);
      border-color: rgba(70,211,154,0.28);
    }
    .modal-stack {
      display: grid;
      gap: 14px;
    }
    .doc-card {
      display: grid;
      grid-template-columns: 40px minmax(0, 1fr) 18px;
      gap: 10px;
      align-items: center;
      padding: 12px 14px;
      border-radius: 18px;
      background: rgba(21,94,108,0.34);
      border: 1px solid #35d0ff;
      color: var(--text);
      text-decoration: none;
    }
    .doc-icon {
      width: 40px;
      height: 40px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #54c9f0;
      color: #08212a;
      font-size: 16px;
      font-weight: 800;
    }
    .doc-title {
      font-size: 13px;
      font-weight: 800;
      margin-bottom: 2px;
    }
    .doc-link {
      color: rgba(237,243,255,0.58);
      font-size: 11px;
      line-height: 1.45;
      word-break: break-all;
    }
    .doc-arrow {
      font-size: 16px;
      color: var(--text);
      text-align: right;
    }
    .modal-grid {
      display: grid;
      gap: 14px;
      grid-template-columns: 1fr;
    }
    .field {
      display: grid;
      gap: 8px;
    }
    label.small {
      font-size: 12px;
      color: var(--muted);
      font-weight: 700;
      letter-spacing: 0.04em;
      text-transform: uppercase;
    }
    input, textarea {
      width: 100%;
      border: 1px solid rgba(130,164,255,0.14);
      border-radius: 16px;
      background: rgba(8,12,18,0.9);
      padding: 14px 15px;
      font: inherit;
      color: var(--text);
    }
    textarea {
      min-height: 120px;
      resize: vertical;
    }
    input[type="file"] {
      padding: 12px 14px;
      line-height: 1.4;
    }
    .field-value {
      color: var(--muted);
      font-size: 12px;
      line-height: 1.6;
      word-break: break-all;
    }
    .hint {
      color: var(--muted);
      font-size: 12px;
      line-height: 1.6;
    }
    .modal-actions {
      display: flex;
      justify-content: flex-end;
      gap: 12px;
      flex-wrap: wrap;
      margin-top: 6px;
    }
    .ghost-btn {
      color: var(--text);
      background: rgba(255,255,255,0.05);
      border: 1px solid rgba(255,255,255,0.08);
      border-radius: 999px;
      padding: 12px 18px;
      min-width: 96px;
      font: inherit;
      font-weight: 800;
      cursor: pointer;
    }
    @media (max-width: 860px) {
      .config-row {
        grid-template-columns: 1fr;
      }
      .modal-grid {
        grid-template-columns: 1fr;
      }
    }
  </style>
</head>
<body>
  <div class="shell">
    <section class="hero">
      <div class="hero-top">
        <div class="hero-bar">
          <h1 id="page-title">Config</h1>
          <a class="hero-link" href="/upload">返回发布页</a>
        </div>
        <div class="hero-desc">集中维护发布流程需要的核心配置。每项配置按条目展示，可单独编辑与保存。</div>
      </div>
    </section>

    <section class="panel">
      <div class="panel-head">
        <div class="panel-kicker">config center</div>
        <h2 id="config-title">配置列表</h2>
      </div>
      <div class="panel-body">
        <div id="config-groups"></div>
      </div>
    </section>
  </div>

  <div id="config-modal" class="modal" aria-hidden="true">
    <div class="modal-card">
      <div class="modal-head">
        <div class="panel-kicker">edit item</div>
        <h2 id="modal-title">编辑配置</h2>
      </div>
      <div class="modal-body">
        <div id="modal-status" class="modal-status"></div>
        <div id="modal-desc" class="hint"></div>
        <div id="modal-fields" class="modal-grid"></div>
        <div class="modal-actions">
          <button type="button" id="modal-cancel" class="ghost-btn">取消</button>
          <button type="button" id="modal-save" class="config-action">保存</button>
        </div>
      </div>
    </div>
  </div>

  <script>
    let configPayload = null;
    let modalSection = null;
    let modalSaving = false;
    const createAppMode = new URLSearchParams(window.location.search).get('mode') === 'create-app';

    function escapeHtml(value) {
      return String(value || '')
        .replaceAll('&', '&amp;')
        .replaceAll('<', '&lt;')
        .replaceAll('>', '&gt;')
        .replaceAll('"', '&quot;')
        .replaceAll("'", '&#39;');
    }

    async function applyCurrentAppTitle() {
      try {
        const resp = await fetch('/api/config' + window.location.search);
        const data = await resp.json();
        if (!resp.ok) return;
        const appName = data?.ui?.app_name || data?.current_app || 'APKGO';
        document.getElementById('page-title').textContent = appName + ' Config';
        document.title = appName + ' Config';
      } catch (_) {}
    }

    function setModalStatus(message, kind) {
      const el = document.getElementById('modal-status');
      if (!message) {
        el.textContent = '';
        el.className = 'modal-status';
        return;
      }
      el.textContent = message;
      el.className = 'modal-status show ' + kind;
    }

    function groupTitle(key) {
      const map = {
        app: 'App 名称',
        stores: '商店凭证',
        ui: '包名配置',
        hooks: '飞书机器人配置',
        aliases: '市场渠道名称配置',
      };
      return map[key] || key;
    }

    function getSection(sectionKey) {
      return (configPayload?.sections || []).find((item) => item.key === sectionKey) || null;
    }

    function renderGroups() {
      const groups = document.getElementById('config-groups');
      const allItems = configPayload?.items || [];
      const ordered = ['app', 'ui', 'hooks', 'stores', 'aliases'];
      groups.innerHTML = ordered.map((groupKey) => {
        const items = allItems
          .filter((item) => item.group_key === groupKey)
          .sort((a, b) => {
            const orderDiff = Number(a.order || 0) - Number(b.order || 0);
            if (orderDiff !== 0) return orderDiff;
            if (Boolean(a.configured) !== Boolean(b.configured)) {
              return a.configured ? -1 : 1;
            }
            return String(a.display_name || '').localeCompare(String(b.display_name || ''), 'zh-CN');
          });
        if (!items.length) return '';
        const list = items.map((item) => {
          const rowClass = item.configured ? 'config-row is-configured' : 'config-row is-unconfigured';
          const summaryClass = item.configured ? 'config-summary ok' : 'config-summary pending';
          return (
            '<div class="' + rowClass + '">' +
              '<div class="config-meta">' +
                '<div class="config-name">' + escapeHtml(item.display_name) + '</div>' +
                (item.subtitle ? '<div class="config-subtitle">' + escapeHtml(item.subtitle) + '</div>' : '') +
              '</div>' +
              '<div class="' + summaryClass + '">' + escapeHtml(item.summary || '未配置') + '</div>' +
              '<button type="button" class="config-action" data-group="' + escapeHtml(item.group_key) + '" data-section="' + escapeHtml(item.section_key) + '" data-key="' + escapeHtml(item.key) + '">' + escapeHtml(item.edit_label || '编辑') + '</button>' +
            '</div>'
          );
        }).join('');
        return '<section class="config-group"><div class="group-title">' + escapeHtml(groupTitle(groupKey)) + '</div><div class="config-list">' + list + '</div></section>';
      }).join('');
    }

    function sectionValue(sectionKey, fieldKey) {
      if (sectionKey === 'app') return configPayload?.ui?.app_name || '';
      if (sectionKey === 'hooks') return configPayload?.hooks?.feishu_webhook || '';
      if (sectionKey === 'ui') return configPayload?.ui?.default_audit_package || '';
      if (sectionKey === 'aliases') return '';
      return configPayload?.stores_config?.[sectionKey]?.[fieldKey] || '';
    }

    function openModal(groupKey, sectionKey, itemKey) {
      modalSection = { groupKey, sectionKey, itemKey };
      const title = document.getElementById('modal-title');
      const desc = document.getElementById('modal-desc');
      const fields = document.getElementById('modal-fields');
      setModalStatus('', '');
      title.textContent = '编辑 ' + itemKey;
      desc.textContent = '';

      if (groupKey === 'aliases' && configPayload?.market_aliases?.[itemKey]) {
        title.textContent = '编辑 ' + itemKey + ' 渠道名称';
        desc.textContent = '使用换行分隔多个别名，例如 tencent 和 qq。';
        fields.innerHTML =
          '<div class="field full">' +
            '<label class="small">别名列表</label>' +
            '<textarea id="alias-editor">' + escapeHtml((configPayload.market_aliases[itemKey] || []).join('\n')) + '</textarea>' +
          '</div>';
      } else {
        const section = getSection(sectionKey);
        if (!section) return;
        title.textContent = '编辑 ' + (section.display_name || itemKey);
        desc.textContent = section.description || '';
        const docCard = section.doc_url
          ? '<a class="doc-card" href="' + escapeHtml(section.doc_url) + '" target="_blank" rel="noreferrer">' +
              '<div class="doc-icon">◫</div>' +
              '<div>' +
                '<div class="doc-title">查看获取凭证文档</div>' +
                '<div class="doc-link">' + escapeHtml(section.doc_url) + '</div>' +
              '</div>' +
              '<div class="doc-arrow">↗</div>' +
            '</a>'
          : '';
        fields.innerHTML = '<div class="modal-stack">' + docCard + '<div class="modal-grid">' + (section.fields || []).map((field) => {
          const value = sectionValue(sectionKey, field.key);
          const control = field.file
            ? '<input type="file" data-field="' + escapeHtml(field.key) + '" ' + (field.accept ? 'accept="' + escapeHtml(field.accept) + '"' : '') + '>'
            : (field.multiline
              ? '<textarea data-field="' + escapeHtml(field.key) + '" placeholder="' + escapeHtml(field.placeholder || '') + '">' + escapeHtml(value) + '</textarea>'
              : '<input ' + (field.secret ? 'type="password"' : 'type="text"') + ' data-field="' + escapeHtml(field.key) + '" value="' + escapeHtml(value) + '" placeholder="' + escapeHtml(field.placeholder || '') + '">');
          return (
            '<div class="field">' +
              '<label class="small">' + escapeHtml(field.label) + '</label>' +
              control +
              (field.file ? '<div class="field-value">当前文件：' + escapeHtml(value || '未上传') + '</div>' : '') +
              (field.advanced ? '<div class="hint">高级字段</div>' : '') +
            '</div>'
          );
        }).join('') + '</div></div>';
      }

      document.getElementById('config-modal').classList.add('show');
      document.getElementById('config-modal').setAttribute('aria-hidden', 'false');
    }

    function closeModal() {
      if (modalSaving) return;
      if (createAppMode) {
        window.location.href = '/apps';
        return;
      }
      setModalStatus('', '');
      modalSection = null;
      document.getElementById('config-modal').classList.remove('show');
      document.getElementById('config-modal').setAttribute('aria-hidden', 'true');
    }

    function buildSavePayload() {
      return {
        ui: configPayload.ui || {},
        hooks: configPayload.hooks || {},
        stores: configPayload.stores_config || {},
        market_aliases: configPayload.market_aliases || {},
        target_group: modalSection?.groupKey || '',
        target_section: modalSection?.sectionKey || '',
        app_mode: createAppMode ? 'create' : '',
        app_id: configPayload.app_id || '',
      };
    }

    function setModalSaving(saving) {
      modalSaving = saving;
      const saveButton = document.getElementById('modal-save');
      const cancelButton = document.getElementById('modal-cancel');
      if (saveButton) {
        saveButton.disabled = saving;
        if (createAppMode) {
          saveButton.textContent = saving ? '确定中...' : '确定';
        } else {
          saveButton.textContent = saving ? '保存中...' : '保存';
        }
      }
      if (cancelButton) {
        cancelButton.disabled = saving;
      }
    }

    async function saveModal() {
      if (!modalSection || modalSaving) return;
      const payload = buildSavePayload();
      const formData = new FormData();
      let hasFile = false;
      if (modalSection.groupKey === 'aliases' && payload.market_aliases[modalSection.itemKey] != null) {
        payload.market_aliases[modalSection.itemKey] = String(document.getElementById('alias-editor').value || '')
          .split('\n')
          .map((item) => item.trim())
          .filter(Boolean);
      } else if (modalSection.sectionKey === 'app') {
        payload.ui.app_name = String(document.querySelector('[data-field="app_name"]').value || '').trim();
      } else if (modalSection.sectionKey === 'ui') {
        payload.ui.default_audit_package = String(document.querySelector('[data-field="default_audit_package"]').value || '').trim();
      } else if (modalSection.sectionKey === 'hooks') {
        payload.hooks.feishu_webhook = String(document.querySelector('[data-field="feishu_webhook"]').value || '').trim();
      } else {
        if (!payload.stores[modalSection.sectionKey]) payload.stores[modalSection.sectionKey] = {};
        for (const field of document.querySelectorAll('[data-field]')) {
          const fieldKey = field.getAttribute('data-field');
          if (field.type === 'file') {
            const file = field.files && field.files[0];
            if (file) {
              hasFile = true;
              formData.append('store_file_' + modalSection.sectionKey + '__' + fieldKey, file);
            }
            continue;
          }
          payload.stores[modalSection.sectionKey][fieldKey] = String(field.value || '').trim();
        }
      }

      setModalStatus('', '');
      setModalSaving(true);
      try {
        let resp;
        if (hasFile) {
          formData.append('payload', JSON.stringify(payload));
          resp = await fetch('/api/config/save', {
            method: 'POST',
            body: formData,
          });
        } else {
          resp = await fetch('/api/config/save', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload),
          });
        }
        const data = await resp.json();
        if (!resp.ok) {
          setModalStatus(data.error || '保存失败', 'bad');
          return;
        }
        if (createAppMode) {
          window.location.href = '/config?app=' + encodeURIComponent(data.app_id || '');
          return;
        }
        await loadConfig();
        setModalStatus('配置保存成功。', 'ok');
        setTimeout(() => {
          if (!modalSaving) closeModal();
        }, 500);
      } finally {
        setModalSaving(false);
      }
    }

    async function loadConfig() {
      const resp = await fetch('/api/config' + window.location.search);
      const data = await resp.json();
      if (!resp.ok) throw new Error(data.error || '读取配置失败');
      configPayload = data;
      if (!configPayload.market_aliases || !Object.keys(configPayload.market_aliases).length) {
        configPayload.market_aliases = {};
      }
      if (createAppMode) {
        document.getElementById('config-title').textContent = '新增 App';
        const desc = document.querySelector('.hero-desc');
        if (desc) desc.textContent = '这里会创建一个空白 App。请先填写 App 名称，再按顺序补充包名、飞书配置、市场凭证和渠道配置。保存后会自动成为当前主配置。';
        const saveButton = document.getElementById('modal-save');
        if (saveButton) saveButton.textContent = '确定';
      }
      renderGroups();
      if (createAppMode) {
        openModal('app', 'app', 'app_name');
      }
    }

    document.addEventListener('click', (event) => {
      const trigger = event.target.closest('[data-group][data-section][data-key]');
      if (trigger) {
        openModal(trigger.getAttribute('data-group'), trigger.getAttribute('data-section'), trigger.getAttribute('data-key'));
      }
    });
    document.getElementById('modal-cancel').addEventListener('click', closeModal);
    document.getElementById('modal-save').addEventListener('click', saveModal);
    document.getElementById('config-modal').addEventListener('click', (event) => {
      if (event.target.id === 'config-modal') closeModal();
    });

    applyCurrentAppTitle();
    loadConfig().catch((err) => {
      document.getElementById('config-groups').innerHTML =
        '<div class="status show bad">读取配置失败：' + escapeHtml(String(err)) + '</div>';
    });
  </script>
</body>
</html>
`

const webAuditHTML = `<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>apkgo Audit</title>
  <style>
    :root {
      --bg: #090b10;
      --panel: rgba(14,18,25,0.88);
      --line: rgba(130,164,255,0.14);
      --line-strong: rgba(130,164,255,0.28);
      --text: #edf3ff;
      --muted: #8d9ab3;
      --accent: #4de2c5;
      --ok: #46d39a;
      --warn: #ffb454;
      --bad: #ff6b7a;
      --shadow: 0 28px 80px rgba(0,0,0,0.45);
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      min-height: 100vh;
      color: var(--text);
      font-family: "SF Mono", "JetBrains Mono", "IBM Plex Sans", "PingFang SC", "Microsoft YaHei", monospace, sans-serif;
      background:
        radial-gradient(circle at 12% 10%, rgba(77,226,197,0.12), transparent 22%),
        radial-gradient(circle at 88% 14%, rgba(106,169,255,0.12), transparent 20%),
        linear-gradient(180deg, #07090d 0%, var(--bg) 100%);
    }
    .shell {
      max-width: 1120px;
      margin: 0 auto;
      padding: 28px 20px 52px;
    }
    .hero, .panel {
      background: var(--panel);
      border: 1px solid var(--line);
      box-shadow: var(--shadow);
    }
    .hero {
      position: relative;
      overflow: hidden;
      padding: 22px 24px;
      border-radius: 32px;
      margin-bottom: 20px;
      background:
        linear-gradient(135deg, rgba(18,24,34,0.96), rgba(10,14,20,0.94)),
        linear-gradient(180deg, rgba(77,226,197,0.05), rgba(106,169,255,0.04));
      border-color: var(--line-strong);
    }
    .hero::after {
      content: "";
      position: absolute;
      inset: -1px;
      border-radius: inherit;
      padding: 1px;
      background: linear-gradient(135deg, rgba(77,226,197,0.3), rgba(106,169,255,0.18), rgba(141,107,255,0.2));
      -webkit-mask: linear-gradient(#000 0 0) content-box, linear-gradient(#000 0 0);
      -webkit-mask-composite: xor;
      mask-composite: exclude;
      pointer-events: none;
    }
    .hero-top {
      position: relative;
      z-index: 1;
      display: grid;
      gap: 14px;
    }
    .hero-bar {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 16px;
    }
    h1 {
      margin: 0;
      font-size: clamp(34px, 6vw, 58px);
      line-height: 0.94;
      letter-spacing: -0.06em;
    }
    .hero-desc {
      max-width: 760px;
      color: var(--muted);
      font-size: 14px;
      line-height: 1.8;
    }
    .hero-link {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      min-height: 44px;
      padding: 0 18px;
      border-radius: 999px;
      border: 1px solid rgba(255,255,255,0.08);
      background: rgba(255,255,255,0.05);
      color: var(--text);
      text-decoration: none;
      font-size: 12px;
      font-weight: 800;
      letter-spacing: 0.08em;
      text-transform: uppercase;
    }
    .hero-side {
      display: grid;
      gap: 12px;
      grid-template-columns: repeat(3, minmax(0, 1fr));
    }
    .hero-stat {
      padding: 12px 14px;
      border-radius: 18px;
      background: rgba(255,255,255,0.03);
      border: 1px solid rgba(255,255,255,0.06);
    }
    .hero-stat strong {
      display: block;
      color: var(--muted);
      font-size: 11px;
      font-weight: 800;
      letter-spacing: 0.08em;
      text-transform: uppercase;
      margin-bottom: 8px;
    }
    .hero-stat span {
      font-size: 13px;
      line-height: 1.6;
    }
    .panel {
      border-radius: 28px;
      overflow: hidden;
    }
    .panel-head {
      padding: 18px 22px 14px;
      border-bottom: 1px solid rgba(255,255,255,0.05);
      background: linear-gradient(180deg, rgba(255,255,255,0.02), rgba(255,255,255,0));
    }
    .panel-kicker {
      color: #6aa9ff;
      font-size: 11px;
      font-weight: 800;
      letter-spacing: 0.08em;
      text-transform: uppercase;
      margin-bottom: 10px;
    }
    .panel h2 {
      margin: 0 0 8px;
      font-size: 24px;
      letter-spacing: -0.04em;
    }
    .note {
      margin: 0;
      color: var(--muted);
      font-size: 13px;
      line-height: 1.8;
    }
    .panel-body {
      padding: 22px;
    }
    .head-row {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 16px;
      margin-bottom: 2px;
    }
    .head-actions {
      display: flex;
      align-items: center;
      gap: 10px;
      flex-wrap: wrap;
      justify-content: flex-end;
    }
    button {
      border: 0;
      border-radius: 999px;
      padding: 14px 24px;
      min-width: 220px;
      font: inherit;
      font-weight: 800;
      cursor: pointer;
      color: #051014;
      background: linear-gradient(135deg, var(--accent), #8bf6df);
      box-shadow: 0 16px 34px rgba(77,226,197,0.18);
    }
    button:disabled {
      opacity: .45;
      cursor: wait;
    }
    .secondary-btn {
      color: var(--text);
      background: rgba(255,255,255,0.06);
      border: 1px solid rgba(255,255,255,0.08);
      box-shadow: none;
    }
    .secondary-btn.hidden {
      display: none;
    }
    .status {
      margin-top: 16px;
      padding: 12px 14px;
      border-radius: 16px;
      display: none;
      font-size: 14px;
      line-height: 1.6;
    }
    .status.show { display: block; }
    .status.ok {
      background: rgba(70,211,154,0.12);
      color: var(--ok);
      border: 1px solid rgba(70,211,154,0.14);
    }
    .status.bad {
      background: rgba(255,107,122,0.1);
      color: var(--bad);
      border: 1px solid rgba(255,107,122,0.14);
    }
    .status.info {
      background: rgba(106,169,255,0.1);
      color: #91b8ff;
      border: 1px solid rgba(106,169,255,0.14);
    }
    .audit-results {
      display: grid;
      gap: 12px;
      margin-top: 16px;
    }
    .audit-card {
      border: 1px solid rgba(130,164,255,0.12);
      border-radius: 18px;
      background: rgba(8,12,18,0.84);
      padding: 16px;
      display: grid;
      gap: 10px;
    }
    .audit-card.reviewing {
      border-color: rgba(255,180,84,0.24);
      box-shadow: inset 0 0 0 1px rgba(255,180,84,0.08);
    }
    .audit-card.approved {
      border-color: rgba(70,211,154,0.24);
      box-shadow: inset 0 0 0 1px rgba(70,211,154,0.08);
    }
    .audit-card.rejected {
      border-color: rgba(255,107,122,0.24);
      box-shadow: inset 0 0 0 1px rgba(255,107,122,0.08);
    }
    .audit-card.withdrawn,
    .audit-card.unknown {
      border-color: rgba(106,169,255,0.18);
    }
    .audit-card.unsupported {
      opacity: 0.72;
      background: rgba(255,255,255,0.02);
    }
    .audit-head {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 12px;
    }
    .audit-store {
      font-size: 16px;
      font-weight: 800;
      letter-spacing: -0.02em;
    }
    .audit-detail {
      font-size: 12px;
      line-height: 1.7;
      color: var(--muted);
      white-space: pre-wrap;
      word-break: break-word;
    }
    .audit-actions {
      display: flex;
      justify-content: flex-end;
      margin-top: 2px;
    }
    .audit-link {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      min-height: 34px;
      padding: 0 14px;
      border-radius: 999px;
      border: 1px solid rgba(255,255,255,0.08);
      background: rgba(255,255,255,0.05);
      color: var(--text);
      text-decoration: none;
      font-size: 11px;
      font-weight: 800;
      letter-spacing: 0.04em;
      transition: transform .16s ease, border-color .16s ease, background .16s ease;
    }
    .audit-link:hover {
      transform: translateY(-1px);
      border-color: rgba(77,226,197,0.28);
      background: rgba(77,226,197,0.08);
    }
    .badge {
      border-radius: 999px;
      padding: 5px 9px;
      font-size: 11px;
      font-weight: 800;
      line-height: 1;
      letter-spacing: 0.02em;
    }
    .badge.ok {
      background: rgba(70,211,154,0.14);
      color: var(--ok);
    }
    .badge.warn {
      background: rgba(255,180,84,0.14);
      color: var(--warn);
    }
    .badge.bad {
      background: rgba(255,107,122,0.14);
      color: var(--bad);
    }
    .badge.info {
      background: rgba(106,169,255,0.14);
      color: #91b8ff;
    }
    @media (max-width: 900px) {
      .hero-side {
        grid-template-columns: 1fr;
      }
      .head-row {
        flex-direction: column;
        align-items: stretch;
      }
      button {
        min-width: 0;
        width: 100%;
      }
    }
    @media (max-width: 720px) {
      .shell {
        padding: 20px 14px 32px;
      }
      .hero {
        padding: 18px;
      }
      .hero-bar {
        flex-direction: column;
      }
    }
  </style>
</head>
<body>
  <div class="shell">
    <section class="hero">
      <div class="hero-top">
        <div class="hero-bar">
          <h1 id="page-title">Audit</h1>
          <a class="hero-link" href="/upload">返回发布页</a>
        </div>
        <div class="hero-desc">这里单独查询所有已配置市场的审核状态。默认直接使用配置中的包名，不需要额外输入。</div>
        <div class="hero-side">
          <div class="hero-stat">
            <strong>Default Package</strong>
            <span id="default-package-label">读取中...</span>
          </div>
          <div class="hero-stat">
            <strong>Scope</strong>
            <span>所有已配置市场</span>
          </div>
          <div class="hero-stat">
            <strong>Result</strong>
            <span>统一映射为审核中 / 通过 / 驳回 / 撤回 / 未知</span>
          </div>
        </div>
      </div>
    </section>

    <section class="panel">
      <div class="panel-head">
        <div class="panel-kicker">audit flow</div>
        <div class="head-row">
          <h2>审核状态查询</h2>
          <div class="head-actions">
            <button type="button" id="sync-feishu-btn" class="secondary-btn hidden">同步至飞书</button>
            <button type="button" id="audit-btn">查询所有市场审核状态</button>
          </div>
        </div>
      </div>
      <div class="panel-body">
        <div id="audit-status" class="status"></div>
        <div id="audit-results" class="audit-results"></div>
      </div>
    </section>
  </div>

  <script>
    let defaultPackage = '';
    let latestAuditData = null;

    function escapeHtml(value) {
      return String(value || '')
        .replaceAll('&', '&amp;')
        .replaceAll('<', '&lt;')
        .replaceAll('>', '&gt;')
        .replaceAll('"', '&quot;')
        .replaceAll("'", '&#39;');
    }

    async function applyCurrentAppTitle() {
      try {
        const resp = await fetch('/api/app/current');
        const data = await resp.json();
        if (!resp.ok) return;
        const appName = data?.name || 'APKGO';
        document.getElementById('page-title').textContent = appName + ' Audit';
        document.title = appName + ' Audit';
      } catch (_) {}
    }

    function setStatus(message, kind) {
      const el = document.getElementById('audit-status');
      if (!message) {
        el.textContent = '';
        el.className = 'status';
        return;
      }
      el.textContent = message;
      el.className = 'status show ' + kind;
    }

    function stateMeta(item) {
      if (!item.supported) return { label: '暂不支持', badge: 'info', card: 'unsupported' };
      if (item.error) return { label: '查询失败', badge: 'bad', card: 'rejected' };
      switch (item.state) {
        case 'approved': return { label: '审核通过', badge: 'ok', card: 'approved' };
        case 'rejected': return { label: '审核驳回', badge: 'bad', card: 'rejected' };
        case 'withdrawn': return { label: '已撤回', badge: 'info', card: 'withdrawn' };
        case 'reviewing': return { label: '审核中', badge: 'warn', card: 'reviewing' };
        default: return { label: '状态未知', badge: 'info', card: 'unknown' };
      }
    }

    function storeName(store) {
      const named = {
        huawei: '华为',
        xiaomi: '小米',
        oppo: 'OPPO',
        vivo: 'vivo',
        honor: '荣耀',
        tencent: '应用宝',
        samsung: 'Samsung',
        googleplay: 'Google Play',
        pgyer: '蒲公英',
        fir: 'fir.im',
        script: 'Script',
      };
      return named[store] || store;
    }

    function auditSortRank(item) {
      if (!item.supported) return 90;
      if (item.error) return 80;
      switch (item.state) {
        case 'approved': return 10;
        case 'reviewing': return 20;
        case 'rejected': return 30;
        case 'withdrawn': return 40;
        case 'unknown': return 50;
        default: return 60;
      }
    }

    function manualViewURL(store) {
      const urls = {
        huawei: 'https://developer.huawei.com/consumer/cn/',
        tencent: 'https://open.tencent.com/',
        oppo: 'https://open.oppomobile.com/',
        honor: 'https://developer.honor.com/cn/',
        vivo: 'https://developer.vivo.com.cn/',
        xiaomi: 'https://dev.mi.com/xiaomihyperos',
        pgyer: 'https://www.pgyer.com/',
      };
      return urls[store] || '';
    }

    function renderResults(data) {
      latestAuditData = data;
      const root = document.getElementById('audit-results');
      document.getElementById('sync-feishu-btn').classList.remove('hidden');
      root.innerHTML = '';
      const stores = [...(data.stores || [])].sort((a, b) => {
        const rankDiff = auditSortRank(a) - auditSortRank(b);
        if (rankDiff !== 0) return rankDiff;
        return storeName(a.store).localeCompare(storeName(b.store), 'zh-CN');
      });
      for (const item of stores) {
        const meta = stateMeta(item);
        let versionValue = '';
        if (item.versionCode != null) {
          versionValue = String(item.versionCode);
        }
        if (!versionValue && item.version_code != null) {
          versionValue = String(item.version_code);
        }
        const detail = !item.supported
          ? (item.store === 'xiaomi'
            ? '官方未提供审核状态查询 API，请前往小米后台查看。'
            : '该市场当前还没有接入审核状态查询接口。')
          : item.error
            ? item.error
            : (item.detail || '');
        const manualURL = manualViewURL(item.store);
        const action = manualURL
          ? '<div class="audit-actions"><a class="audit-link" href="' + escapeHtml(manualURL) + '" target="_blank" rel="noreferrer">手动查看</a></div>'
          : '';
        const card = document.createElement('div');
        card.className = 'audit-card ' + meta.card;
        card.innerHTML =
          '<div class="audit-head">' +
            '<div class="audit-store">' + escapeHtml(storeName(item.store)) + '</div>' +
            '<span class="badge ' + meta.badge + '">' + escapeHtml(meta.label) + '</span>' +
          '</div>' +
          '<div class="audit-detail">版本号：' + escapeHtml(versionValue) + '</div>' +
          (detail ? '<div class="audit-detail">' + escapeHtml(detail) + '</div>' : '') +
          action;
        root.appendChild(card);
      }
    }

    async function loadConfig() {
      const resp = await fetch('/api/config');
      const data = await resp.json();
      defaultPackage = data?.ui?.default_audit_package || '';
      const value = defaultPackage || '未配置';
      document.getElementById('default-package-label').textContent = value;
    }

    async function queryAudit() {
      const btn = document.getElementById('audit-btn');
      btn.disabled = true;
      setStatus('正在查询所有已配置市场的审核状态...', 'info');
      try {
        const fd = new FormData();
        const resp = await fetch('/api/audit', { method: 'POST', body: fd });
        const data = await resp.json();
        if (!resp.ok) {
          latestAuditData = null;
          document.getElementById('sync-feishu-btn').classList.add('hidden');
          setStatus(data.error || '查询失败', 'bad');
          return;
        }
        renderResults(data);
        const total = (data.stores || []).length;
        const resolved = (data.stores || []).filter(item => item.supported && !item.error && ['approved', 'rejected', 'withdrawn'].includes(item.state)).length;
        setStatus('查询完成，共返回 ' + total + ' 个市场，其中已到终态 ' + resolved + ' 个。', 'ok');
      } catch (err) {
        setStatus('查询失败：' + String(err), 'bad');
      } finally {
        btn.disabled = false;
      }
    }

    async function syncAuditToFeishu() {
      const btn = document.getElementById('sync-feishu-btn');
      if (!latestAuditData || !latestAuditData.stores || !latestAuditData.stores.length) {
        setStatus('请先查询到审核结果，再同步到飞书。', 'bad');
        return;
      }
      btn.disabled = true;
      setStatus('正在同步审核状态到飞书...', 'info');
      try {
        const resp = await fetch('/api/audit/sync-feishu', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(latestAuditData),
        });
        const data = await resp.json();
        if (!resp.ok) {
          setStatus(data.error || '同步失败', 'bad');
          return;
        }
        setStatus('已同步到飞书机器人。', 'ok');
      } catch (err) {
        setStatus('同步失败：' + String(err), 'bad');
      } finally {
        btn.disabled = false;
      }
    }

    document.getElementById('audit-btn').addEventListener('click', queryAudit);
    document.getElementById('sync-feishu-btn').addEventListener('click', syncAuditToFeishu);

    applyCurrentAppTitle();
    loadConfig()
      .then(() => queryAudit())
      .catch((err) => {
        const msg = '读取配置失败：' + String(err);
        document.getElementById('default-package-label').textContent = msg;
        setStatus(msg, 'bad');
      });
  </script>
</body>
</html>
`

const webHistoryHTML = `<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>apkgo History</title>
  <style>
    :root {
      --bg: #090b10;
      --panel: rgba(14,18,25,0.88);
      --line: rgba(130,164,255,0.14);
      --line-strong: rgba(130,164,255,0.28);
      --text: #edf3ff;
      --muted: #8d9ab3;
      --accent: #4de2c5;
      --ok: #46d39a;
      --warn: #ffb454;
      --bad: #ff6b7a;
      --info: #91b8ff;
      --shadow: 0 28px 80px rgba(0,0,0,0.45);
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      min-height: 100vh;
      color: var(--text);
      font-family: "SF Mono", "JetBrains Mono", "IBM Plex Sans", "PingFang SC", "Microsoft YaHei", monospace, sans-serif;
      background:
        radial-gradient(circle at 12% 10%, rgba(77,226,197,0.12), transparent 22%),
        radial-gradient(circle at 88% 14%, rgba(106,169,255,0.12), transparent 20%),
        linear-gradient(180deg, #07090d 0%, var(--bg) 100%);
    }
    .shell {
      max-width: 1200px;
      margin: 0 auto;
      padding: 28px 20px 52px;
    }
    .hero, .panel {
      background: var(--panel);
      border: 1px solid var(--line);
      box-shadow: var(--shadow);
    }
    .hero {
      position: relative;
      overflow: hidden;
      padding: 22px 24px;
      border-radius: 32px;
      margin-bottom: 20px;
      background:
        linear-gradient(135deg, rgba(18,24,34,0.96), rgba(10,14,20,0.94)),
        linear-gradient(180deg, rgba(77,226,197,0.05), rgba(106,169,255,0.04));
      border-color: var(--line-strong);
    }
    .hero::after {
      content: "";
      position: absolute;
      inset: -1px;
      border-radius: inherit;
      padding: 1px;
      background: linear-gradient(135deg, rgba(77,226,197,0.3), rgba(106,169,255,0.18), rgba(141,107,255,0.2));
      -webkit-mask: linear-gradient(#000 0 0) content-box, linear-gradient(#000 0 0);
      -webkit-mask-composite: xor;
      mask-composite: exclude;
      pointer-events: none;
    }
    .hero-top {
      position: relative;
      z-index: 1;
      display: grid;
      gap: 14px;
    }
    .hero-bar {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 16px;
    }
    h1 {
      margin: 0;
      font-size: clamp(34px, 6vw, 58px);
      line-height: 0.94;
      letter-spacing: -0.06em;
    }
    .hero-desc {
      max-width: 780px;
      color: var(--muted);
      font-size: 14px;
      line-height: 1.8;
    }
    .hero-link {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      min-height: 44px;
      padding: 0 18px;
      border-radius: 999px;
      border: 1px solid rgba(255,255,255,0.08);
      background: rgba(255,255,255,0.05);
      color: var(--text);
      text-decoration: none;
      font-size: 12px;
      font-weight: 800;
      letter-spacing: 0.08em;
      text-transform: uppercase;
    }
    .hero-side {
      display: grid;
      gap: 12px;
      grid-template-columns: repeat(3, minmax(0, 1fr));
    }
    .hero-stat {
      padding: 12px 14px;
      border-radius: 18px;
      background: rgba(255,255,255,0.03);
      border: 1px solid rgba(255,255,255,0.06);
    }
    .hero-stat strong {
      display: block;
      color: var(--muted);
      font-size: 11px;
      font-weight: 800;
      letter-spacing: 0.08em;
      text-transform: uppercase;
      margin-bottom: 8px;
    }
    .hero-stat span {
      font-size: 13px;
      line-height: 1.6;
    }
    .panel {
      border-radius: 28px;
      overflow: hidden;
    }
    .panel-head {
      padding: 18px 22px 14px;
      border-bottom: 1px solid rgba(255,255,255,0.05);
      background: linear-gradient(180deg, rgba(255,255,255,0.02), rgba(255,255,255,0));
    }
    .panel-kicker {
      color: #6aa9ff;
      font-size: 11px;
      font-weight: 800;
      letter-spacing: 0.08em;
      text-transform: uppercase;
      margin-bottom: 10px;
    }
    .head-row {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 16px;
    }
    .head-side {
      color: var(--muted);
      font-size: 12px;
      line-height: 1.7;
      text-align: right;
      word-break: break-all;
    }
    .panel h2 {
      margin: 0;
      font-size: 24px;
      letter-spacing: -0.04em;
    }
    .panel-body {
      padding: 22px;
    }
    .status {
      margin-bottom: 16px;
      padding: 12px 14px;
      border-radius: 16px;
      display: none;
      font-size: 14px;
      line-height: 1.6;
    }
    .status.show { display: block; }
    .status.bad {
      background: rgba(255,107,122,0.1);
      color: var(--bad);
      border: 1px solid rgba(255,107,122,0.14);
    }
    .status.info {
      background: rgba(106,169,255,0.1);
      color: var(--info);
      border: 1px solid rgba(106,169,255,0.14);
    }
    .empty-state {
      border: 1px dashed rgba(130,164,255,0.16);
      border-radius: 22px;
      padding: 34px 20px;
      text-align: center;
      color: var(--muted);
      background: rgba(255,255,255,0.02);
      line-height: 1.8;
    }
    .history-list {
      display: grid;
      gap: 14px;
    }
    .history-card {
      border: 1px solid rgba(130,164,255,0.12);
      border-radius: 22px;
      background: rgba(8,12,18,0.84);
      padding: 18px;
      display: grid;
      gap: 8px;
      cursor: pointer;
      transition: transform .16s ease, border-color .16s ease, background .16s ease;
    }
    .history-card:hover {
      transform: translateY(-1px);
      border-color: rgba(77,226,197,0.26);
      background: rgba(10,16,24,0.9);
    }
    .history-head {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 16px;
    }
    .history-app {
      font-size: 18px;
      font-weight: 800;
      letter-spacing: -0.03em;
    }
    .badge {
      border-radius: 999px;
      padding: 6px 10px;
      font-size: 11px;
      font-weight: 800;
      line-height: 1;
      letter-spacing: 0.04em;
      white-space: nowrap;
    }
    .badge.ok {
      background: rgba(70,211,154,0.14);
      color: var(--ok);
    }
    .badge.warn {
      background: rgba(255,180,84,0.14);
      color: var(--warn);
    }
    .badge.bad {
      background: rgba(255,107,122,0.14);
      color: var(--bad);
    }
    .badge.info {
      background: rgba(106,169,255,0.14);
      color: var(--info);
    }
    .history-time {
      color: var(--muted);
      font-size: 14px;
      line-height: 1.7;
    }
    .history-hint {
      color: var(--muted);
      font-size: 12px;
      line-height: 1.6;
    }
    @media (max-width: 900px) {
      .hero-side {
        grid-template-columns: 1fr;
      }
      .head-row, .history-head {
        flex-direction: column;
        align-items: stretch;
      }
      .head-side {
        text-align: left;
      }
    }
    @media (max-width: 720px) {
      .shell {
        padding: 20px 14px 32px;
      }
      .hero {
        padding: 18px;
      }
      .hero-bar {
        flex-direction: column;
      }
    }
  </style>
</head>
<body>
  <div class="shell">
    <section class="hero">
      <div class="hero-top">
        <div class="hero-bar">
          <h1 id="page-title">History</h1>
          <a class="hero-link" href="/upload">返回发布页</a>
        </div>
        <div class="hero-side">
          <div class="hero-stat">
            <strong>Storage</strong>
            <span>本地 JSONL 持久化，轻量、可读、零额外依赖</span>
          </div>
          <div class="hero-stat">
            <strong>Contains</strong>
            <span>发布时间 / 更新文案 / 版本 / 渠道结果</span>
          </div>
          <div class="hero-stat">
            <strong>Scope</strong>
            <span>当前 app 独立持久化，切换 app 时自动同步到主 History</span>
          </div>
        </div>
      </div>
    </section>

    <section class="panel">
      <div class="panel-head">
        <div class="panel-kicker">release ledger</div>
        <div class="head-row">
          <h2>发布记录</h2>
        </div>
      </div>
      <div class="panel-body">
        <div id="history-status" class="status"></div>
        <div id="history-list" class="history-list"></div>
      </div>
    </section>
  </div>

  <script>
    function escapeHtml(value) {
      return String(value || '')
        .replaceAll('&', '&amp;')
        .replaceAll('<', '&lt;')
        .replaceAll('>', '&gt;')
        .replaceAll('"', '&quot;')
        .replaceAll("'", '&#39;');
    }

    async function applyCurrentAppTitle() {
      try {
        const resp = await fetch('/api/app/current');
        const data = await resp.json();
        if (!resp.ok) return;
        const appName = data?.name || 'APKGO';
        document.getElementById('page-title').textContent = appName + ' History';
        document.title = appName + ' History';
      } catch (_) {}
    }

    function setStatus(message, kind) {
      const el = document.getElementById('history-status');
      if (!message) {
        el.textContent = '';
        el.className = 'status';
        return;
      }
      el.textContent = message;
      el.className = 'status show ' + kind;
    }

    function statusMeta(status) {
      switch (status) {
        case 'success': return { label: '全部成功', cls: 'ok' };
        case 'failed': return { label: '全部失败', cls: 'bad' };
        case 'partial': return { label: '部分成功', cls: 'warn' };
        default: return { label: '状态未知', cls: 'info' };
      }
    }

    function renderEmpty() {
      document.getElementById('history-list').innerHTML =
        '<div class="empty-state">本地还没有发布记录。<br>你在发布页完成一次真实发布后，这里会自动出现记录。</div>';
    }

    function renderRecords(records) {
      const root = document.getElementById('history-list');
      if (!records || !records.length) {
        renderEmpty();
        return;
      }
      root.innerHTML = '';
      for (const record of records) {
        const meta = statusMeta(record.status);
        const card = document.createElement('article');
        card.className = 'history-card';

        card.innerHTML =
          '<div class="history-head">' +
            '<div class="history-app">' + escapeHtml(record.version_name || '-') + '</div>' +
            '<span class="badge ' + meta.cls + '">' + escapeHtml(meta.label) + '</span>' +
          '</div>' +
          '<div class="history-time">' + escapeHtml(record.published_at || record.timestamp || '-') + '</div>';
        card.addEventListener('click', () => {
          window.location.href = '/history/detail?ts=' + encodeURIComponent(record.timestamp || '');
        });
        root.appendChild(card);
      }
    }

    async function loadHistory() {
      setStatus('正在读取本地发布记录...', 'info');
      try {
        const resp = await fetch('/api/history');
        const data = await resp.json();
        if (!resp.ok) {
          setStatus(data.error || '读取失败', 'bad');
          renderEmpty();
          return;
        }
        renderRecords(data.records || []);
        setStatus('', 'info');
      } catch (err) {
        setStatus('读取失败：' + String(err), 'bad');
        renderEmpty();
      }
    }

    applyCurrentAppTitle();
    loadHistory();
  </script>
</body>
</html>
`

const webHistoryDetailHTML = `<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>发布记录详情</title>
  <style>
    :root {
      --bg: #090b10;
      --panel: rgba(14,18,25,0.88);
      --line: rgba(130,164,255,0.14);
      --line-strong: rgba(130,164,255,0.28);
      --text: #edf3ff;
      --muted: #8d9ab3;
      --accent: #4de2c5;
      --ok: #46d39a;
      --warn: #ffb454;
      --bad: #ff6b7a;
      --info: #6aa9ff;
      --shadow: 0 28px 80px rgba(0,0,0,0.45);
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      min-height: 100vh;
      color: var(--text);
      font-family: "SF Mono", "JetBrains Mono", "IBM Plex Sans", "PingFang SC", "Microsoft YaHei", monospace, sans-serif;
      background:
        radial-gradient(circle at 12% 10%, rgba(77,226,197,0.12), transparent 22%),
        radial-gradient(circle at 88% 14%, rgba(106,169,255,0.12), transparent 20%),
        linear-gradient(180deg, #07090d 0%, var(--bg) 100%);
    }
    .shell {
      max-width: 1120px;
      margin: 0 auto;
      padding: 28px 20px 52px;
    }
    .hero, .panel {
      background: var(--panel);
      border: 1px solid var(--line);
      box-shadow: var(--shadow);
    }
    .hero {
      position: relative;
      overflow: hidden;
      padding: 22px 24px;
      border-radius: 32px;
      margin-bottom: 20px;
      background:
        linear-gradient(135deg, rgba(18,24,34,0.96), rgba(10,14,20,0.94)),
        linear-gradient(180deg, rgba(77,226,197,0.05), rgba(106,169,255,0.04));
      border-color: var(--line-strong);
    }
    .hero-bar, .head-row, .store-row {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 16px;
    }
    .hero-link {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      min-height: 44px;
      padding: 0 18px;
      border-radius: 999px;
      border: 1px solid rgba(255,255,255,0.08);
      background: rgba(255,255,255,0.05);
      color: var(--text);
      text-decoration: none;
      font-size: 12px;
      font-weight: 800;
      letter-spacing: 0.08em;
      text-transform: uppercase;
    }
    .hero-desc, .head-side, .notes-text, .store-detail {
      color: var(--muted);
      line-height: 1.7;
    }
    .hero-desc { margin-top: 12px; font-size: 14px; }
    .panel {
      border-radius: 28px;
      overflow: hidden;
    }
    .panel-head {
      padding: 18px 22px 14px;
      border-bottom: 1px solid rgba(255,255,255,0.05);
    }
    .panel-kicker, .section-label, .metric strong {
      color: var(--info);
      font-size: 11px;
      font-weight: 800;
      letter-spacing: 0.08em;
      text-transform: uppercase;
    }
    .panel-body {
      padding: 22px;
      display: grid;
      gap: 16px;
    }
    .status {
      padding: 12px 14px;
      border-radius: 16px;
      display: none;
      font-size: 14px;
      line-height: 1.6;
    }
    .status.show { display: block; }
    .status.bad {
      background: rgba(255,107,122,0.1);
      color: var(--bad);
      border: 1px solid rgba(255,107,122,0.14);
    }
    .status.info {
      background: rgba(106,169,255,0.1);
      color: var(--info);
      border: 1px solid rgba(106,169,255,0.14);
    }
    .detail-card, .notes-box, .store-table {
      border: 1px solid rgba(255,255,255,0.06);
      background: rgba(255,255,255,0.03);
      border-radius: 18px;
      padding: 16px;
    }
    .notes-box {
      margin-top: 18px;
      padding-top: 20px;
      padding-bottom: 20px;
    }
    .store-table {
      margin-top: 22px;
    }
    .detail-head {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 16px;
      margin-bottom: 14px;
    }
    .detail-title {
      font-size: 24px;
      font-weight: 800;
      letter-spacing: -0.04em;
    }
    .detail-sub {
      color: var(--muted);
      font-size: 13px;
      line-height: 1.7;
      margin-top: 6px;
    }
    .badge {
      border-radius: 999px;
      padding: 6px 10px;
      font-size: 11px;
      font-weight: 800;
      line-height: 1;
      letter-spacing: 0.04em;
      white-space: nowrap;
    }
    .badge.ok { background: rgba(70,211,154,0.14); color: var(--ok); }
    .badge.warn { background: rgba(255,180,84,0.14); color: var(--warn); }
    .badge.bad { background: rgba(255,107,122,0.14); color: var(--bad); }
    .badge.info { background: rgba(106,169,255,0.14); color: var(--info); }
    .detail-grid {
      display: grid;
      grid-template-columns: repeat(3, minmax(0, 1fr));
      gap: 12px;
    }
    .metric {
      border: 1px solid rgba(255,255,255,0.06);
      background: rgba(255,255,255,0.03);
      border-radius: 18px;
      padding: 12px 14px;
    }
    .metric strong {
      display: block;
      margin-bottom: 8px;
    }
    .metric span {
      font-size: 14px;
      line-height: 1.7;
      word-break: break-word;
    }
    .section-label {
      margin-bottom: 14px;
    }
    .notes-text {
      font-size: 13px;
      line-height: 1.9;
      white-space: pre-wrap;
      word-break: break-word;
    }
    .store-list {
      display: grid;
      gap: 10px;
    }
    .store-row {
      padding: 10px 12px;
      border-radius: 14px;
      background: rgba(255,255,255,0.03);
    }
    .store-name {
      font-size: 13px;
      font-weight: 700;
    }
    .store-detail {
      font-size: 12px;
      text-align: right;
      word-break: break-word;
    }
    .modal {
      position: fixed;
      inset: 0;
      background: rgba(2,6,12,0.72);
      backdrop-filter: blur(10px);
      display: none;
      align-items: center;
      justify-content: center;
      padding: 18px;
      z-index: 50;
    }
    .modal.show { display: flex; }
    .modal-card {
      width: min(520px, 100%);
      border-radius: 26px;
      border: 1px solid rgba(130,164,255,0.18);
      background:
        linear-gradient(180deg, rgba(18,24,34,0.98), rgba(10,14,20,0.98)),
        rgba(10,14,20,0.98);
      box-shadow: var(--shadow);
      padding: 22px;
      display: grid;
      gap: 18px;
    }
    .modal-head {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 16px;
    }
    .modal-title {
      margin: 0 0 8px;
      font-size: 22px;
      letter-spacing: -0.04em;
    }
    .modal-sub {
      margin: 0;
      color: var(--muted);
      font-size: 13px;
      line-height: 1.7;
    }
    .modal-actions {
      display: flex;
      justify-content: flex-end;
      gap: 12px;
    }
    .secondary, .danger {
      min-height: 44px;
      padding: 0 18px;
      border-radius: 999px;
      font: inherit;
      font-size: 13px;
      font-weight: 800;
      letter-spacing: 0.04em;
      cursor: pointer;
      border: 1px solid rgba(255,255,255,0.08);
      color: var(--text);
    }
    .secondary {
      background: rgba(255,255,255,0.05);
    }
    .danger {
      background: rgba(255,107,122,0.14);
      border-color: rgba(255,107,122,0.24);
      color: var(--bad);
    }
    @media (max-width: 900px) {
      .detail-grid {
        grid-template-columns: 1fr;
      }
      .hero-bar, .head-row, .detail-head, .store-row {
        flex-direction: column;
        align-items: stretch;
      }
      .head-side, .store-detail {
        text-align: left;
      }
    }
  </style>
</head>
<body>
  <div class="shell">
    <section class="hero">
      <div class="hero-bar">
        <h1 id="page-title">History Detail</h1>
        <a class="hero-link" href="/history">返回发布记录</a>
      </div>
    </section>

    <section class="panel">
      <div class="panel-head">
        <div class="panel-kicker">release detail</div>
        <div class="head-row">
          <h2></h2>
          <button type="button" id="delete-record-btn" class="danger">删除记录</button>
        </div>
      </div>
      <div class="panel-body">
        <div id="detail-status" class="status"></div>
        <div id="detail-root"></div>
      </div>
    </section>
  </div>

  <div id="delete-modal" class="modal" aria-hidden="true">
    <div class="modal-card">
      <div class="modal-head">
        <div>
          <h3 class="modal-title">确认删除这条发布记录</h3>
          <p class="modal-sub">删除后将从本地发布记录中移除，操作无法撤销。确认后会返回列表页面。</p>
        </div>
      </div>
      <div class="modal-actions">
        <button type="button" id="delete-cancel" class="secondary">取消</button>
        <button type="button" id="delete-confirm" class="danger">确认删除</button>
      </div>
    </div>
  </div>

  <script>
    let currentRecord = null;

    function escapeHtml(value) {
      return String(value || '')
        .replaceAll('&', '&amp;')
        .replaceAll('<', '&lt;')
        .replaceAll('>', '&gt;')
        .replaceAll('"', '&quot;')
        .replaceAll("'", '&#39;');
    }

    async function applyCurrentAppTitle() {
      try {
        const resp = await fetch('/api/app/current');
        const data = await resp.json();
        if (!resp.ok) return;
        const appName = data?.name || 'APKGO';
        document.getElementById('page-title').textContent = appName + ' History';
        document.title = appName + ' History Detail';
      } catch (_) {}
    }

    function setStatus(message, kind) {
      const el = document.getElementById('detail-status');
      if (!message) {
        el.textContent = '';
        el.className = 'status';
        return;
      }
      el.textContent = message;
      el.className = 'status show ' + kind;
    }

    function statusMeta(status) {
      switch (status) {
        case 'success': return { label: '全部成功', cls: 'ok' };
        case 'failed': return { label: '全部失败', cls: 'bad' };
        case 'partial': return { label: '部分成功', cls: 'warn' };
        default: return { label: '状态未知', cls: 'info' };
      }
    }

    function publishModeLabel(mode) {
      switch (mode) {
        case 'scheduled': return '定时发布';
        case 'auto': return '审核后自动发布';
        default: return mode || '未记录';
      }
    }

    function storeName(store) {
      const named = {
        huawei: '华为',
        xiaomi: '小米',
        oppo: 'OPPO',
        vivo: 'vivo',
        honor: '荣耀',
        tencent: '应用宝',
        samsung: 'Samsung',
        googleplay: 'Google Play',
        pgyer: '蒲公英',
        fir: 'fir.im',
        script: 'Script',
      };
      return named[store] || store || '-';
    }

    function renderNotFound() {
      currentRecord = null;
      document.getElementById('delete-record-btn').disabled = true;
      document.getElementById('detail-root').innerHTML =
        '<div class="detail-card">没有找到这条发布记录。</div>';
    }

    function renderDetail(record) {
      if (!record) {
        renderNotFound();
        return;
      }
      currentRecord = record;
      document.getElementById('delete-record-btn').disabled = false;
      const meta = statusMeta(record.status);
      const stores = (record.results || []).map((item) => {
        const badge = item.success
          ? '<span class="badge ok">成功</span>'
          : '<span class="badge bad">失败</span>';
        const detail = item.success
          ? '耗时 ' + String(item.duration_ms || 0) + ' ms'
          : escapeHtml(item.error || '未知错误');
        return '' +
          '<div class="store-row">' +
            '<div class="store-name">' + escapeHtml(storeName(item.store)) + '</div>' +
            '<div class="store-detail">' + badge + '<div>' + detail + '</div></div>' +
          '</div>';
      }).join('');

      document.getElementById('detail-root').innerHTML =
        '<div class="detail-card">' +
          '<div class="detail-head">' +
            '<div>' +
              '<div class="detail-title">版本 ' + escapeHtml(record.version_name || '-') + '</div>' +
              '<div class="detail-sub">' + escapeHtml(record.app_name || record.package_name || '未命名应用') + '</div>' +
            '</div>' +
            '<span class="badge ' + meta.cls + '">' + escapeHtml(meta.label) + '</span>' +
          '</div>' +
          '<div class="detail-grid">' +
            '<div class="metric"><strong>发布时间</strong><span>' + escapeHtml(record.published_at || record.timestamp || '-') + '</span></div>' +
            '<div class="metric"><strong>包名</strong><span>' + escapeHtml(record.package_name || '-') + '</span></div>' +
            '<div class="metric"><strong>版本号</strong><span>' + escapeHtml((record.version_name || '-') + ' (' + String(record.version_code || 0) + ')') + '</span></div>' +
            '<div class="metric"><strong>发布方式</strong><span>' + escapeHtml(publishModeLabel(record.publish_mode)) + '</span></div>' +
            '<div class="metric"><strong>定时发布时间</strong><span>' + escapeHtml(record.publish_time || '即时发布') + '</span></div>' +
            '<div class="metric"><strong>结果概览</strong><span>成功 ' + String(record.success_count || 0) + ' / 失败 ' + String(record.failure_count || 0) + '</span></div>' +
          '</div>' +
        '</div>' +
        '<div class="notes-box">' +
          '<div class="section-label">更新文案</div>' +
          '<div class="notes-text">' + escapeHtml(record.notes || '本次未填写更新文案') + '</div>' +
        '</div>' +
        '<div class="store-table">' +
          '<div class="section-label">市场结果</div>' +
          '<div class="store-list">' + stores + '</div>' +
        '</div>';
    }

    async function loadDetail() {
      const ts = new URLSearchParams(window.location.search).get('ts');
      if (!ts) {
        setStatus('缺少记录标识。', 'bad');
        renderNotFound();
        return;
      }
      setStatus('正在读取发布记录详情...', 'info');
      try {
        const resp = await fetch('/api/history');
        const data = await resp.json();
        if (!resp.ok) {
          setStatus(data.error || '读取失败', 'bad');
          renderNotFound();
          return;
        }
        const record = (data.records || []).find((item) => item.timestamp === ts);
        renderDetail(record);
        setStatus('', 'info');
      } catch (err) {
        setStatus('读取失败：' + String(err), 'bad');
        renderNotFound();
      }
    }

    function openDeleteModal() {
      const modal = document.getElementById('delete-modal');
      modal.classList.add('show');
      modal.setAttribute('aria-hidden', 'false');
    }

    function closeDeleteModal() {
      const modal = document.getElementById('delete-modal');
      modal.classList.remove('show');
      modal.setAttribute('aria-hidden', 'true');
    }

    async function deleteRecord() {
      if (!currentRecord?.timestamp) return;
      setStatus('正在删除发布记录...', 'info');
      const resp = await fetch('/api/history/delete', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ timestamp: currentRecord.timestamp }),
      });
      const data = await resp.json();
      if (!resp.ok) {
        setStatus(data.error || '删除失败', 'bad');
        return;
      }
      window.location.href = '/history?refresh=' + Date.now();
    }

    document.getElementById('delete-record-btn').addEventListener('click', () => {
      if (!currentRecord) return;
      openDeleteModal();
    });
    document.getElementById('delete-cancel').addEventListener('click', closeDeleteModal);
    document.getElementById('delete-confirm').addEventListener('click', async () => {
      closeDeleteModal();
      await deleteRecord();
    });
    document.getElementById('delete-modal').addEventListener('click', (event) => {
      if (event.target.id === 'delete-modal') closeDeleteModal();
    });

    applyCurrentAppTitle();
    loadDetail();
  </script>
</body>
</html>
`
