<template>
    <svg :width="size" :height="size * 0.8" :viewBox="`0 0 30 20`" xmlns="http://www.w3.org/2000/svg"
        :aria-label="`${resolvedExt} file icon`" role="img">
        <!-- <rect x="0" y="0" :width="size" :height="size" rx="4" :fill="config.bg" /> -->
        <!-- <polygon points="28,0 40,12 40,48 0,48 0,0" :fill="config.bg" /> -->
        <polyline points="28,0 28,12 40,12" fill="none" :stroke="config.fold" stroke-width="1" />
        <text x="14" y="14" text-anchor="middle" :font-size="labelFontSize" font-weight="700"
            font-family="'JetBrains Mono', 'Fira Code', 'Courier New', monospace" :fill="config.color">
            {{ config.label  }}
        </text>
    </svg>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface IconConfig {
    bg: string
    fold: string
    color: string
    label: string
}

interface Props {
    /**
     * File extension (e.g. 'ts', 'js', 'py') or full filename (e.g. 'Dockerfile', '.env')
     */
    lang?: string;
    /**
     * Category hint: 'language' | 'web' | 'data' | 'doc' | 'image' | 'media' | 'build'
     * Optional — auto-resolved from lang if omitted.
     */
    type?: 'file' | 'directory';
    /**
     * Icon size in pixels (width). Height is automatically 1.2× width. Default: 40
     */
    size?: number;
}

const props = withDefaults(defineProps<Props>(), {
    lang: 'txt',
    size: 24,
})

const ICON_MAP: Record<string, IconConfig> = {
    // ── Languages ────────────────────────────────────────────
    ts: { bg: '#1A1A2E', fold: '#3A3A5C', color: '#3178C6', label: 'TS' },
    tsx: { bg: '#1A1F35', fold: '#2A3050', color: '#3178C6', label: 'TSX' },
    js: { bg: '#1A1A2E', fold: '#3A3A5C', color: '#F7DF1E', label: 'JS' },
    jsx: { bg: '#1A1F35', fold: '#2A3050', color: '#61DAFB', label: 'JSX' },
    mjs: { bg: '#1A1A2E', fold: '#3A3A5C', color: '#F7DF1E', label: 'MJS' },
    cjs: { bg: '#1A1A2E', fold: '#3A3A5C', color: '#F7DF1E', label: 'CJS' },
    py: { bg: '#1A2E1A', fold: '#2A4A2A', color: '#4CAF50', label: 'PY' },
    pyi: { bg: '#1A2E1A', fold: '#2A4A2A', color: '#4CAF50', label: 'PYI' },
    rb: { bg: '#2E1A1A', fold: '#4A2A2A', color: '#CC342D', label: 'RB' },
    go: { bg: '#1A2A35', fold: '#2A3A48', color: '#00ADD8', label: 'GO' },
    rs: { bg: '#2E1A1A', fold: '#4A2A2A', color: '#DEA584', label: 'RS' },
    c: { bg: '#1A1A2A', fold: '#2A2A3C', color: '#9B59B6', label: 'C' },
    h: { bg: '#1A1A2A', fold: '#2A2A3C', color: '#9B59B6', label: 'H' },
    cpp: { bg: '#1A1A2A', fold: '#2A2A3C', color: '#9B59B6', label: 'C++' },
    cc: { bg: '#1A1A2A', fold: '#2A2A3C', color: '#9B59B6', label: 'C++' },
    cxx: { bg: '#1A1A2A', fold: '#2A2A3C', color: '#9B59B6', label: 'C++' },
    hpp: { bg: '#1A1A2A', fold: '#2A2A3C', color: '#9B59B6', label: 'H++' },
    java: { bg: '#1A2E1A', fold: '#2A4A2A', color: '#B07219', label: 'JAV' },
    class: { bg: '#1A2E1A', fold: '#2A4A2A', color: '#B07219', label: 'CLS' },
    kt: { bg: '#1A2A35', fold: '#2A3A48', color: '#7F52FF', label: 'KT' },
    kts: { bg: '#1A2A35', fold: '#2A3A48', color: '#7F52FF', label: 'KTS' },
    swift: { bg: '#1A2E2A', fold: '#2A4A3A', color: '#5AC8FA', label: 'SW' },
    dart: { bg: '#1A2A35', fold: '#2A3A48', color: '#00B4AB', label: 'DRT' },
    lua: { bg: '#2E2A1A', fold: '#4A3A2A', color: '#E67E22', label: 'LUA' },
    cs: { bg: '#2A1A2E', fold: '#3C2A4A', color: '#A97BFF', label: 'CS' },
    fs: { bg: '#1A2A35', fold: '#2A3A48', color: '#378ADD', label: 'FS' },
    zig: { bg: '#2E2A1A', fold: '#4A3A2A', color: '#FFD700', label: 'ZIG' },
    ex: { bg: '#2A1A2E', fold: '#3C2A4A', color: '#9B30FF', label: 'EX' },
    exs: { bg: '#2A1A2E', fold: '#3C2A4A', color: '#9B30FF', label: 'EXS' },
    erl: { bg: '#2E1A1A', fold: '#4A2A2A', color: '#A90533', label: 'ERL' },
    hs: { bg: '#2A1A2E', fold: '#3C2A4A', color: '#5D4F85', label: 'HS' },
    clj: { bg: '#1A2E1A', fold: '#2A4A2A', color: '#5881D8', label: 'CLJ' },
    scala: { bg: '#2E1A1A', fold: '#4A2A2A', color: '#DC322F', label: 'SC' },
    r: { bg: '#1A2A35', fold: '#2A3A48', color: '#276DC3', label: 'R' },
    php: { bg: '#2A1A2E', fold: '#3C2A4A', color: '#8993BE', label: 'PHP' },
    m: { bg: '#1A2A35', fold: '#2A3A48', color: '#438EFF', label: 'M' },
    mm: { bg: '#1A2A35', fold: '#2A3A48', color: '#438EFF', label: 'MM' },

    // ── Web & Styles ─────────────────────────────────────────
    html: { bg: '#2E1A1A', fold: '#4A2A2A', color: '#E34C26', label: 'HTM' },
    htm: { bg: '#2E1A1A', fold: '#4A2A2A', color: '#E34C26', label: 'HTM' },
    css: { bg: '#1A1A2E', fold: '#2A2A4A', color: '#264DE4', label: 'CSS' },
    scss: { bg: '#2A1A2E', fold: '#3C2A4A', color: '#CD6799', label: 'SCSS' },
    sass: { bg: '#2A1A2E', fold: '#3C2A4A', color: '#CD6799', label: 'SASS' },
    less: { bg: '#1A2A35', fold: '#2A3A48', color: '#1D365D', label: 'LESS' },
    styl: { bg: '#1A2E1A', fold: '#2A4A2A', color: '#4CAF50', label: 'STL' },
    vue: { bg: '#2A1A1A', fold: '#3C2A2A', color: '#42B883', label: 'VUE' },
    svelte: { bg: '#1A2E2A', fold: '#2A3C30', color: '#FF3E00', label: 'SVE' },
    svg: { bg: '#1A2E1A', fold: '#2A4A2A', color: '#41B883', label: 'SVG' },
    astro: { bg: '#1A1A2E', fold: '#2A2A3C', color: '#FF5D01', label: 'AST' },

    // ── Data & Config ─────────────────────────────────────────
    json: { bg: '#2A2A1A', fold: '#3C3C2A', color: '#CBCB41', label: 'JSON' },
    jsonc: { bg: '#2A2A1A', fold: '#3C3C2A', color: '#CBCB41', label: 'JSNC' },
    yaml: { bg: '#2A2A1A', fold: '#3C3C2A', color: '#CC3E44', label: 'YAML' },
    yml: { bg: '#2A2A1A', fold: '#3C3C2A', color: '#CC3E44', label: 'YML' },
    xml: { bg: '#1A2E1A', fold: '#2A4A2A', color: '#4CAF50', label: 'XML' },
    csv: { bg: '#2A2A1A', fold: '#3C3C2A', color: '#E5C07B', label: 'CSV' },
    toml: { bg: '#2A2E1A', fold: '#3A4A2A', color: '#6A9955', label: 'TOML' },
    ini: { bg: '#2A2A1A', fold: '#3C3C2A', color: '#888', label: 'INI' },
    env: { bg: '#1A1A1A', fold: '#2A2A2A', color: '#888', label: 'ENV' },
    prisma: { bg: '#1A2A1A', fold: '#2A3A2A', color: '#4EC9B0', label: 'PRM' },
    sql: { bg: '#1A2A2A', fold: '#2A3A3A', color: '#3B82F6', label: 'SQL' },
    graphql: { bg: '#2A1A2E', fold: '#3C2A4A', color: '#E535AB', label: 'GQL' },
    gql: { bg: '#2A1A2E', fold: '#3C2A4A', color: '#E535AB', label: 'GQL' },
    proto: { bg: '#1A2A35', fold: '#2A3A48', color: '#3178C6', label: 'PRT' },

    // ── Docs & Scripts ────────────────────────────────────────
    md: { bg: '#1A2A35', fold: '#2A3A48', color: '#519ABA', label: 'MD' },
    mdx: { bg: '#1A2A35', fold: '#2A3A48', color: '#61DAFB', label: 'MDX' },
    rst: { bg: '#1A2A35', fold: '#2A3A48', color: '#519ABA', label: 'RST' },
    txt: { bg: '#1A2A2A', fold: '#2A3A3A', color: '#4EC9B0', label: 'TXT' },
    pdf: { bg: '#2E1A1A', fold: '#4A2A2A', color: '#F44747', label: 'PDF' },
    docx: { bg: '#1A2A2E', fold: '#2A3A44', color: '#40A9FF', label: 'DOC' },
    doc: { bg: '#1A2A2E', fold: '#2A3A44', color: '#40A9FF', label: 'DOC' },
    xlsx: { bg: '#1A2A1A', fold: '#2A4A2A', color: '#1D6F42', label: 'XLS' },
    xls: { bg: '#1A2A1A', fold: '#2A4A2A', color: '#1D6F42', label: 'XLS' },
    pptx: { bg: '#2E1A1A', fold: '#4A2A2A', color: '#D04423', label: 'PPT' },
    sh: { bg: '#1A1A1A', fold: '#2A2A2A', color: '#89D185', label: 'SH' },
    bash: { bg: '#1A1A1A', fold: '#2A2A2A', color: '#89D185', label: 'BSH' },
    zsh: { bg: '#1A1A1A', fold: '#2A2A2A', color: '#89D185', label: 'ZSH' },
    fish: { bg: '#1A1A1A', fold: '#2A2A2A', color: '#89D185', label: 'FSH' },
    ps1: { bg: '#1A1A2A', fold: '#2A2A3C', color: '#C586C0', label: 'PS1' },
    bat: { bg: '#1A1A1A', fold: '#2A2A2A', color: '#C1C1C1', label: 'BAT' },

    // ── Images & Media ────────────────────────────────────────
    png: { bg: '#2E2A1A', fold: '#4A3A2A', color: '#E9A145', label: 'PNG' },
    jpg: { bg: '#2E2A1A', fold: '#4A3A2A', color: '#E9A145', label: 'JPG' },
    jpeg: { bg: '#2E2A1A', fold: '#4A3A2A', color: '#E9A145', label: 'JPG' },
    gif: { bg: '#1A1A2A', fold: '#2A2A3C', color: '#7F77DD', label: 'GIF' },
    webp: { bg: '#1A2E2A', fold: '#2A4A3A', color: '#5DCAA5', label: 'WBP' },
    ico: { bg: '#2E2A1A', fold: '#4A3A2A', color: '#E9A145', label: 'ICO' },
    bmp: { bg: '#2E2A1A', fold: '#4A3A2A', color: '#E9A145', label: 'BMP' },
    tiff: { bg: '#2E2A1A', fold: '#4A3A2A', color: '#E9A145', label: 'TIF' },
    mp4: { bg: '#1A2A2E', fold: '#2A3A44', color: '#56B6C2', label: 'MP4' },
    mov: { bg: '#1A2A2E', fold: '#2A3A44', color: '#56B6C2', label: 'MOV' },
    webm: { bg: '#1A2A2E', fold: '#2A3A44', color: '#56B6C2', label: 'WBM' },
    mp3: { bg: '#2E1A2A', fold: '#4A2A3A', color: '#C678DD', label: 'MP3' },
    wav: { bg: '#2E1A2A', fold: '#4A2A3A', color: '#C678DD', label: 'WAV' },
    ogg: { bg: '#2E1A2A', fold: '#4A2A3A', color: '#C678DD', label: 'OGG' },

    // ── Build & Package ───────────────────────────────────────
    'package.json': { bg: '#2A2A1A', fold: '#3C3C2A', color: '#CBCB41', label: 'PKG' },
    'package-lock.json': { bg: '#1A1A1A', fold: '#2A2A2A', color: '#4EC9B0', label: 'LCK' },
    'yarn.lock': { bg: '#1A1A1A', fold: '#2A2A2A', color: '#2C8EBB', label: 'LCK' },
    'pnpm-lock.yaml': { bg: '#1A1A1A', fold: '#2A2A2A', color: '#F69220', label: 'LCK' },
    'tsconfig.json': { bg: '#1A2A2E', fold: '#2A3A44', color: '#3178C6', label: 'TSC' },
    'vite.config': { bg: '#1A1A2E', fold: '#2A2A3C', color: '#646CFF', label: 'VIT' },
    'webpack.config': { bg: '#1A2A35', fold: '#2A3A48', color: '#8DD6F9', label: 'WPK' },
    'rollup.config': { bg: '#2E1A1A', fold: '#4A2A2A', color: '#FF3333', label: 'RUP' },
    makefile: { bg: '#1A1A1A', fold: '#2A2A2A', color: '#D4D4D4', label: 'MKF' },
    dockerfile: { bg: '#1A2A35', fold: '#2A3A48', color: '#2496ED', label: 'DCK' },
    'docker-compose': { bg: '#1A2A35', fold: '#2A3A48', color: '#2496ED', label: 'DCO' },
    '.gitignore': { bg: '#2E1A1A', fold: '#4A2A2A', color: '#F05032', label: 'GIT' },
    '.eslintrc': { bg: '#2A1A2E', fold: '#3C2A4A', color: '#4B32C3', label: 'ESL' },
    '.prettierrc': { bg: '#1A2E2A', fold: '#2A4A3A', color: '#F7B93E', label: 'PRT' },
    '.babelrc': { bg: '#2A2A1A', fold: '#3C3C2A', color: '#F5DA55', label: 'BBL' },
    '.editorconfig': { bg: '#1A1A1A', fold: '#2A2A2A', color: '#FEFEFE', label: 'EDC' },
    'cargo.toml': { bg: '#2E1A1A', fold: '#4A2A2A', color: '#DEA584', label: 'CGO' },
    'go.mod': { bg: '#1A2A35', fold: '#2A3A48', color: '#00ADD8', label: 'MOD' },
    'requirements.txt': { bg: '#1A2E1A', fold: '#2A4A2A', color: '#4CAF50', label: 'REQ' },
    'pyproject.toml': { bg: '#1A2E1A', fold: '#2A4A2A', color: '#4CAF50', label: 'PPY' },
    'gemfile': { bg: '#2E1A1A', fold: '#4A2A2A', color: '#CC342D', label: 'GEM' },
    lock: { bg: '#1A1A1A', fold: '#2A2A2A', color: '#4EC9B0', label: 'LCK' },
    zip: { bg: '#2A2A1A', fold: '#3C3C2A', color: '#E5C07B', label: 'ZIP' },
    tar: { bg: '#2A2A1A', fold: '#3C3C2A', color: '#E5C07B', label: 'TAR' },
    gz: { bg: '#2A2A1A', fold: '#3C3C2A', color: '#E5C07B', label: 'GZ' },
    wasm: { bg: '#2A1A2E', fold: '#3C2A4A', color: '#654FF0', label: 'WSM' },
    cert: { bg: '#1A2A2A', fold: '#2A3A3A', color: '#4EC9B0', label: 'CRT' },
    pem: { bg: '#1A2A2A', fold: '#2A3A3A', color: '#4EC9B0', label: 'PEM' },
}

const FALLBACK: IconConfig = {
    bg: '#1E1E1E',
    fold: '#2A2A2A',
    color: '#858585',
    label: '?',
}

const resolvedExt = computed(() => {
    const raw = (props.lang ?? '').trim().toLowerCase()

    const fullNameMatches = [
        'package.json', 'package-lock.json', 'yarn.lock', 'pnpm-lock.yaml',
        'tsconfig.json', 'vite.config', 'webpack.config', 'rollup.config',
        'makefile', 'dockerfile', 'docker-compose', '.gitignore', '.eslintrc',
        '.prettierrc', '.babelrc', '.editorconfig', 'cargo.toml', 'go.mod',
        'requirements.txt', 'pyproject.toml', 'gemfile',
    ]

    for (const name of fullNameMatches) {
        if (raw === name || raw.startsWith(name + '.') || raw.endsWith('/' + name)) {
            return name
        }
    }

    const ext = raw.startsWith('.') ? raw.slice(1) : raw.split('.').pop() ?? raw
    
    return ext
})

const config = computed<IconConfig>(() => {
    return ICON_MAP[resolvedExt.value] ?? FALLBACK
})

const labelFontSize = computed(() => {
    const len = config.value.label.length;

    if (len <= 2) return 14;
    if (len <= 3) return 13.5;

    return 13.4
})
</script>