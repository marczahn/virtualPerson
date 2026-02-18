'use strict';

// ── Bio variable metadata ────────────────────────────────────────────────────

const BIO_META = {
  body_temp:        { label: 'Body Temp',        min: 34,  max: 43,  unit: '°C',     color: '#ff7b72' },
  heart_rate:       { label: 'Heart Rate',        min: 40,  max: 200, unit: 'bpm',    color: '#ffa657' },
  blood_pressure:   { label: 'Blood Pressure',    min: 80,  max: 200, unit: 'mmHg',   color: '#f0883e' },
  respiratory_rate: { label: 'Resp. Rate',        min: 8,   max: 40,  unit: 'br/min', color: '#e3b341' },
  spo2:             { label: 'SpO2',              min: 70,  max: 100, unit: '%',      color: '#3fb950' },
  blood_sugar:      { label: 'Blood Sugar',       min: 50,  max: 200, unit: 'mg/dL',  color: '#58a6ff' },
  hunger:           { label: 'Hunger',            min: 0,   max: 1,   unit: '',       color: '#79c0ff' },
  thirst:           { label: 'Thirst',            min: 0,   max: 1,   unit: '',       color: '#a5d6ff' },
  hydration:        { label: 'Hydration',         min: 0,   max: 1,   unit: '',       color: '#56d364' },
  glycogen:         { label: 'Glycogen',          min: 0,   max: 1,   unit: '',       color: '#3fb950' },
  cortisol:         { label: 'Cortisol',          min: 0,   max: 1,   unit: '',       color: '#f85149' },
  adrenaline:       { label: 'Adrenaline',        min: 0,   max: 1,   unit: '',       color: '#ffa657' },
  serotonin:        { label: 'Serotonin',         min: 0,   max: 1,   unit: '',       color: '#bc8cff' },
  dopamine:         { label: 'Dopamine',          min: 0,   max: 1,   unit: '',       color: '#d2a8ff' },
  endorphins:       { label: 'Endorphins',        min: 0,   max: 1,   unit: '',       color: '#ff7b72' },
  fatigue:          { label: 'Fatigue',           min: 0,   max: 1,   unit: '',       color: '#e3b341' },
  pain:             { label: 'Pain',              min: 0,   max: 1,   unit: '',       color: '#f85149' },
  muscle_tension:   { label: 'Muscle Tension',    min: 0,   max: 1,   unit: '',       color: '#ffa657' },
  immune_response:  { label: 'Immune Response',   min: 0,   max: 1,   unit: '',       color: '#56d364' },
  circadian_phase:  { label: 'Circadian Phase',   min: 0,   max: 24,  unit: 'hr',    color: '#58a6ff' },
};

const BIO_GROUPS = [
  { id: 'vital',    label: 'Vital Signs',             vars: ['body_temp', 'heart_rate', 'blood_pressure', 'respiratory_rate', 'spo2'] },
  { id: 'metabolic',label: 'Metabolic',               vars: ['blood_sugar', 'hunger', 'thirst', 'hydration', 'glycogen'] },
  { id: 'hormonal', label: 'Hormonal / Neurochemical', vars: ['cortisol', 'adrenaline', 'serotonin', 'dopamine', 'endorphins'] },
  { id: 'physical', label: 'Physical State',           vars: ['fatigue', 'pain', 'muscle_tension', 'immune_response'] },
  { id: 'circadian',label: 'Circadian',               vars: ['circadian_phase'] },
];

const HISTORY_SIZE = 60;

// Normalize a raw value to [0, 1] within its natural range.
function normalize(varKey, value) {
  const m = BIO_META[varKey];
  return (value - m.min) / (m.max - m.min);
}

// ── State ────────────────────────────────────────────────────────────────────

const bioHistory   = {};   // { varKey: number[] } normalized
const bioRawHistory= {};   // { varKey: number[] } raw values for tooltips
const bioLabels    = [];   // shared timestamp labels (HH:MM:SS)

for (const g of BIO_GROUPS) {
  for (const v of g.vars) {
    bioHistory[v]    = [];
    bioRawHistory[v] = [];
  }
}

// ── Chart instances (one per group) ──────────────────────────────────────────

const charts = {};           // { groupId: Chart }
const hiddenVars = new Set(); // currently hidden variable keys

Chart.defaults.color = '#8b949e';
Chart.defaults.borderColor = '#30363d';
Chart.defaults.font.family = 'ui-monospace, SFMono-Regular, Consolas, monospace';
Chart.defaults.font.size = 11;

function makeChart(canvasId, group) {
  const ctx = document.getElementById(canvasId).getContext('2d');
  const datasets = group.vars.map(v => ({
    label:            BIO_META[v].label,
    data:             [],
    borderColor:      BIO_META[v].color,
    backgroundColor:  BIO_META[v].color + '22',
    borderWidth:      1.5,
    pointRadius:      0,
    tension:          0.3,
    fill:             false,
    hidden:           false,
  }));

  return new Chart(ctx, {
    type: 'line',
    data: { labels: [], datasets },
    options: {
      animation:   false,
      responsive:  true,
      maintainAspectRatio: false,
      interaction: { mode: 'index', intersect: false },
      plugins: {
        legend: { display: false },
        tooltip: {
          callbacks: {
            label: (ctx) => {
              const varKey = group.vars[ctx.datasetIndex];
              const raw    = bioRawHistory[varKey][ctx.dataIndex];
              const meta   = BIO_META[varKey];
              if (raw === undefined) return '';
              const val = meta.unit ? `${raw.toFixed(2)} ${meta.unit}` : raw.toFixed(3);
              return ` ${meta.label}: ${val}`;
            },
          },
        },
      },
      scales: {
        x: {
          ticks: { maxTicksLimit: 6, maxRotation: 0 },
          grid:  { color: '#21262d' },
        },
        y: {
          min:   0,
          max:   1,
          ticks: { maxTicksLimit: 5, callback: v => v.toFixed(1) },
          grid:  { color: '#21262d' },
        },
      },
    },
  });
}

// ── DOM construction ─────────────────────────────────────────────────────────

function buildBioPanel() {
  const container = document.getElementById('bio-groups');

  for (const group of BIO_GROUPS) {
    const details = document.createElement('details');
    details.className = 'bio-group';
    details.open = true;

    const summary = document.createElement('summary');
    summary.textContent = group.label;
    details.appendChild(summary);

    const body = document.createElement('div');
    body.className = 'bio-group-body';

    // Checkbox legend
    const legend = document.createElement('div');
    legend.className = 'var-legend';

    group.vars.forEach((v, i) => {
      const label = document.createElement('label');
      label.className = 'var-toggle';
      label.id = `toggle-${v}`;

      const cb = document.createElement('input');
      cb.type = 'checkbox';
      cb.checked = true;
      cb.addEventListener('change', () => toggleVar(group.id, v, i, cb.checked, label));

      const swatch = document.createElement('span');
      swatch.className = 'swatch';
      swatch.style.background = BIO_META[v].color;

      label.appendChild(cb);
      label.appendChild(swatch);
      label.appendChild(document.createTextNode(BIO_META[v].label));
      legend.appendChild(label);
    });

    body.appendChild(legend);

    // Chart
    const wrap = document.createElement('div');
    wrap.className = 'chart-wrap';
    const canvas = document.createElement('canvas');
    canvas.id = `chart-${group.id}`;
    wrap.appendChild(canvas);
    body.appendChild(wrap);

    details.appendChild(body);
    container.appendChild(details);

    charts[group.id] = makeChart(`chart-${group.id}`, group);
  }
}

function toggleVar(groupId, varKey, datasetIndex, visible, labelEl) {
  const chart = charts[groupId];
  if (!chart) return;
  chart.data.datasets[datasetIndex].hidden = !visible;
  chart.update('none');
  if (visible) {
    hiddenVars.delete(varKey);
    labelEl.classList.remove('hidden-var');
  } else {
    hiddenVars.add(varKey);
    labelEl.classList.add('hidden-var');
  }
}

// ── Bio state update ──────────────────────────────────────────────────────────

function updateBioPanel(payload) {
  const bio        = payload.bio_state;
  const thresholds = bio.thresholds || [];
  const ts         = new Date(payload.timestamp).toLocaleTimeString('en-GB');

  // Roll labels
  if (bioLabels.length >= HISTORY_SIZE) bioLabels.shift();
  bioLabels.push(ts);

  // Identify all variables in threshold breach
  const inThreshold = new Set();
  for (const th of thresholds) {
    // Mark entire groups by system name prefix match — best-effort
    inThreshold.add(th.system);
  }

  // Update per-group charts
  for (const group of BIO_GROUPS) {
    const chart = charts[group.id];
    if (!chart) continue;

    let groupHasThreshold = false;

    group.vars.forEach((v, i) => {
      const raw  = bio[v] ?? 0;
      const norm = normalize(v, raw);

      if (bioHistory[v].length >= HISTORY_SIZE)    bioHistory[v].shift();
      if (bioRawHistory[v].length >= HISTORY_SIZE)  bioRawHistory[v].shift();
      bioHistory[v].push(norm);
      bioRawHistory[v].push(raw);

      chart.data.datasets[i].data = bioHistory[v];
    });

    chart.data.labels = bioLabels;

    // Highlight chart border if any threshold is active for this group
    for (const th of thresholds) {
      if (thresholdMatchesGroup(th.system, group.id)) {
        groupHasThreshold = true;
        break;
      }
    }
    const canvasEl = document.getElementById(`chart-${group.id}`);
    if (canvasEl) {
      canvasEl.style.outline = groupHasThreshold ? '1px solid #f85149' : 'none';
    }

    chart.update('none');
  }

  // Update threshold alert banner
  const alertEl = document.getElementById('threshold-alerts');
  const listEl  = document.getElementById('threshold-list');
  listEl.innerHTML = '';
  if (thresholds.length > 0) {
    alertEl.classList.add('active');
    for (const th of thresholds) {
      const li = document.createElement('li');
      li.textContent = `[${th.condition.toUpperCase()}] ${th.system}: ${th.description}`;
      listEl.appendChild(li);
    }
  } else {
    alertEl.classList.remove('active');
  }
}

function thresholdMatchesGroup(system, groupId) {
  const map = {
    vital:    ['thermoregulation', 'cardiovascular', 'respiratory', 'glycemic'],
    metabolic:['glycemic', 'dehydration', 'starvation'],
    hormonal: ['stress', 'cortisol'],
    physical: ['pain', 'fatigue', 'immune'],
    circadian:['circadian'],
  };
  const keywords = map[groupId] || [];
  return keywords.some(k => system.toLowerCase().includes(k));
}

// ── Psych state update ────────────────────────────────────────────────────────

let personalityChart = null;

function updatePsychPanel(payload) {
  const ps = payload.psych_state;

  setBar('bar-arousal',     ps.arousal);
  setBar('bar-energy',      ps.energy);
  setBar('bar-cogload',     ps.cognitive_load);
  setBar('bar-regulation',  ps.regulation_capacity);
  setValenceBar(ps.valence);

  setVal('val-arousal',    ps.arousal.toFixed(2));
  setVal('val-energy',     ps.energy.toFixed(2));
  setVal('val-cogload',    ps.cognitive_load.toFixed(2));
  setVal('val-regulation', ps.regulation_capacity.toFixed(2));
  setVal('val-valence',    ps.valence.toFixed(2));

  renderTags('coping-list',     ps.active_coping      || [], 'coping');
  renderTags('distortion-list', ps.active_distortions || [], 'distortion');
  renderIsolation(ps.isolation_phase, ps.loneliness_level);

  if (ps.personality) {
    updateRadar(ps.personality);
  }
}

function setBar(id, value) {
  const el = document.getElementById(id);
  if (el) el.style.width = `${Math.max(0, Math.min(1, value)) * 100}%`;
}

function setValenceBar(valence) {
  const el = document.getElementById('bar-valence');
  if (!el) return;
  // valence is in [-1, 1]. Center = 50%. Positive → expands right. Negative → expands left.
  const pct    = (valence / 2) * 100; // -50% to +50%
  const isPos  = valence >= 0;
  const width  = Math.abs(pct);
  el.style.width      = `${width}%`;
  el.style.left       = isPos ? '50%' : `${50 + pct}%`;
  el.style.background = isPos ? '#3fb950' : '#f85149';
}

function setVal(id, text) {
  const el = document.getElementById(id);
  if (el) el.textContent = text;
}

function renderTags(containerId, items, cls) {
  const el = document.getElementById(containerId);
  if (!el) return;
  el.innerHTML = '';
  if (items.length === 0) {
    el.innerHTML = '<span class="empty-label">none</span>';
    return;
  }
  for (const item of items) {
    const tag = document.createElement('span');
    tag.className = `tag ${cls}`;
    tag.textContent = item.replace(/_/g, ' ');
    el.appendChild(tag);
  }
}

function renderIsolation(phase, loneliness) {
  const el = document.getElementById('isolation-tag');
  if (!el) return;
  el.className = `tag isolation-${phase}`;
  el.textContent = `${phase.replace(/_/g, ' ')} (${(loneliness * 100).toFixed(0)}%)`;
}

function updateRadar(personality) {
  const data = [
    personality.openness,
    personality.conscientiousness,
    personality.extraversion,
    personality.agreeableness,
    personality.neuroticism,
  ];

  if (!personalityChart) {
    const ctx = document.getElementById('personality-radar').getContext('2d');
    personalityChart = new Chart(ctx, {
      type: 'radar',
      data: {
        labels: ['Openness', 'Conscientiousness', 'Extraversion', 'Agreeableness', 'Neuroticism'],
        datasets: [{
          label: 'Personality',
          data,
          backgroundColor: 'rgba(88, 166, 255, 0.15)',
          borderColor:     '#58a6ff',
          pointBackgroundColor: '#58a6ff',
          pointRadius: 3,
          borderWidth: 1.5,
        }],
      },
      options: {
        animation: false,
        responsive: true,
        maintainAspectRatio: false,
        plugins: { legend: { display: false } },
        scales: {
          r: {
            min: 0, max: 1,
            ticks: { stepSize: 0.25, backdropColor: 'transparent', color: '#8b949e', font: { size: 10 } },
            grid:       { color: '#30363d' },
            angleLines:  { color: '#30363d' },
            pointLabels: { color: '#c9d1d9', font: { size: 11 } },
          },
        },
      },
    });
  } else {
    personalityChart.data.datasets[0].data = data;
    personalityChart.update('none');
  }
}

// ── WebSocket connection ──────────────────────────────────────────────────────

function setStatus(state) {
  const dot   = document.getElementById('status-dot');
  const label = document.getElementById('status-label');
  dot.className   = state;
  label.textContent = state === 'connected' ? 'Connected' : 'Disconnected';
}

function connect() {
  const ws = new WebSocket(`ws://${location.host}/ws`);

  ws.onopen  = () => setStatus('connected');
  ws.onclose = () => {
    setStatus('disconnected');
    setTimeout(connect, 3000);
  };
  ws.onerror = () => { /* onclose handles reconnect */ };

  ws.onmessage = (event) => {
    let msg;
    try { msg = JSON.parse(event.data); } catch { return; }

    if (msg.type === 'bio_state')   updateBioPanel(msg);
    if (msg.type === 'psych_state') updatePsychPanel(msg);
    // 'thought' messages are silently ignored in this read-only dashboard
  };
}

// ── Init ─────────────────────────────────────────────────────────────────────

document.addEventListener('DOMContentLoaded', () => {
  buildBioPanel();
  connect();
});
