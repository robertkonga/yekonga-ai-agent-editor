<template>
    <div class="h-full bg-slate-900 text-gray-200 p-8 overflow-y-auto">
        <div class="max-w-2xl mx-auto">
            <h1 class="text-2xl font-bold mb-8 flex items-center gap-3">
                <i class="ye ye-gear"></i> Settings
            </h1>

            <!-- ── API Keys ─────────────────────────────────────────────── -->
            <section class="mb-6 bg-slate-800/50 p-6 rounded-xl border border-slate-700/50">
                <h2 class="text-lg font-semibold mb-5 flex items-center gap-2">
                    <i class="ye ye-key text-yellow-500"></i> API Keys
                </h2>
                <div class="space-y-5">

                    <!-- Anthropic -->
                    <div>
                        <div class="flex items-center gap-2 mb-1">
                            <span class="text-sm font-medium text-gray-300">◆ Anthropic</span>
                            <span class="text-[10px] px-1.5 py-0.5 rounded bg-indigo-500/15 text-indigo-400 border border-indigo-500/20">claude-opus · sonnet · haiku</span>
                        </div>
                        <div class="flex gap-2">
                            <div class="relative flex-grow">
                                <input
                                    v-model="keys.anthropic"
                                    :type="showKey.anthropic ? 'text' : 'password'"
                                    placeholder="sk-ant-…"
                                    class="w-full bg-slate-900 border border-slate-700 rounded px-3 py-2 pr-9 text-sm focus:outline-none focus:border-blue-500 transition-colors font-mono"
                                />
                                <button type="button" tabindex="-1"
                                    class="absolute right-2.5 top-1/2 -translate-y-1/2 text-slate-500 hover:text-slate-300"
                                    @click="showKey.anthropic = !showKey.anthropic">
                                    <i :class="showKey.anthropic ? 'ye ye-eye-slash' : 'ye ye-eye'" class="ye-sm"></i>
                                </button>
                            </div>
                            <SaveButton :saving="savingKey === 'anthropic'" @click="saveKey('anthropic')" />
                        </div>
                        <KeyStatus :status="keyStatus.anthropic" />
                    </div>

                    <div class="border-t border-slate-700/40"></div>

                    <!-- Gemini -->
                    <div>
                        <div class="flex items-center gap-2 mb-1">
                            <span class="text-sm font-medium text-gray-300">✦ Google Gemini</span>
                            <span class="text-[10px] px-1.5 py-0.5 rounded bg-blue-500/15 text-blue-400 border border-blue-500/20">gemini-2.5-flash · pro</span>
                        </div>
                        <div class="flex gap-2">
                            <div class="relative flex-grow">
                                <input
                                    v-model="keys.gemini"
                                    :type="showKey.gemini ? 'text' : 'password'"
                                    placeholder="AIza…"
                                    class="w-full bg-slate-900 border border-slate-700 rounded px-3 py-2 pr-9 text-sm focus:outline-none focus:border-blue-500 transition-colors font-mono"
                                />
                                <button type="button" tabindex="-1"
                                    class="absolute right-2.5 top-1/2 -translate-y-1/2 text-slate-500 hover:text-slate-300"
                                    @click="showKey.gemini = !showKey.gemini">
                                    <i :class="showKey.gemini ? 'ye ye-eye-slash' : 'ye ye-eye'" class="ye-sm"></i>
                                </button>
                            </div>
                            <SaveButton :saving="savingKey === 'gemini'" @click="saveKey('gemini')" />
                        </div>
                        <KeyStatus :status="keyStatus.gemini" />
                    </div>

                    <div class="border-t border-slate-700/40"></div>

                    <!-- DeepSeek -->
                    <div>
                        <div class="flex items-center gap-2 mb-1">
                            <span class="text-sm font-medium text-gray-300">🐋 DeepSeek</span>
                            <span class="text-[10px] px-1.5 py-0.5 rounded bg-cyan-500/15 text-cyan-400 border border-cyan-500/20">deepseek-chat · reasoner</span>
                        </div>
                        <div class="flex gap-2">
                            <div class="relative flex-grow">
                                <input
                                    v-model="keys.deepseek"
                                    :type="showKey.deepseek ? 'text' : 'password'"
                                    placeholder="sk-…"
                                    class="w-full bg-slate-900 border border-slate-700 rounded px-3 py-2 pr-9 text-sm focus:outline-none focus:border-blue-500 transition-colors font-mono"
                                />
                                <button type="button" tabindex="-1"
                                    class="absolute right-2.5 top-1/2 -translate-y-1/2 text-slate-500 hover:text-slate-300"
                                    @click="showKey.deepseek = !showKey.deepseek">
                                    <i :class="showKey.deepseek ? 'ye ye-eye-slash' : 'ye ye-eye'" class="ye-sm"></i>
                                </button>
                            </div>
                            <SaveButton :saving="savingKey === 'deepseek'" @click="saveKey('deepseek')" />
                        </div>
                        <KeyStatus :status="keyStatus.deepseek" />
                    </div>
                </div>
            </section>

            <!-- ── Ollama ───────────────────────────────────────────────── -->
            <section class="mb-6 bg-slate-800/50 p-6 rounded-xl border border-slate-700/50">
                <h2 class="text-lg font-semibold mb-1 flex items-center gap-2">
                    <span>🦙</span> Ollama <span class="text-sm font-normal text-slate-500">· local</span>
                </h2>
                <p class="text-[11px] text-slate-500 mb-4">No API key required. Ollama must be running locally.</p>

                <div class="space-y-4">
                    <div>
                        <label class="block text-xs font-medium text-gray-400 mb-1">Host URL</label>
                        <div class="flex gap-2">
                            <input
                                v-model="ollamaHost"
                                type="text"
                                placeholder="http://localhost:11434"
                                class="flex-grow bg-slate-900 border border-slate-700 rounded px-3 py-2 text-sm focus:outline-none focus:border-blue-500 transition-colors font-mono"
                            />
                            <SaveButton :saving="savingKey === 'ollamaHost'" @click="saveOllamaHost" />
                        </div>
                    </div>

                    <!-- Ollama status ping -->
                    <div class="flex items-center gap-2">
                        <button type="button"
                            class="text-[11px] px-3 py-1 rounded border border-slate-700 text-slate-400 hover:text-slate-200 hover:border-slate-500 transition-colors"
                            :disabled="pingingOllama"
                            @click="pingOllama">
                            {{ pingingOllama ? 'Checking…' : 'Test connection' }}
                        </button>
                        <span v-if="ollamaStatus" class="text-[11px]"
                            :class="ollamaStatus === 'ok' ? 'text-emerald-400' : 'text-red-400'">
                            {{ ollamaStatus === 'ok' ? '● Running' : '● Not reachable' }}
                        </span>
                    </div>
                </div>
            </section>

            <!-- ── Agent defaults ──────────────────────────────────────── -->
            <section class="mb-6 bg-slate-800/50 p-6 rounded-xl border border-slate-700/50">
                <h2 class="text-lg font-semibold mb-4 flex items-center gap-2">
                    <i class="ye ye-robot text-blue-500"></i> Agent Defaults
                </h2>
                <div class="flex space-y-4">
                    <div class="flex-1">
                        <label class="block text-xs font-medium text-gray-400 mb-1">Default Provider</label>
                        <div class="flex flex-wrap gap-2">
                            <button
                                v-for="p in providerOptions"
                                :key="p.value"
                                type="button"
                                :class="[
                                    'flex items-center gap-1.5 px-3 py-1.5 rounded-lg border text-sm transition-all duration-150',
                                    defaultProvider === p.value
                                        ? 'border-indigo-500 bg-indigo-500/15 text-indigo-300'
                                        : 'border-slate-700 text-slate-400 hover:border-slate-500 hover:text-slate-200'
                                ]"
                                @click="onProviderChange(p.value)"

                            >
                                <span>{{ p.icon }}</span>
                                <span>{{ p.label }}</span>
                            </button>
                        </div>
                    </div>
                    <div class="w-auto">
                        <label class="block text-xs font-medium text-gray-400 mb-1">Default Model</label>
                        <!-- Model dropdown -->
                        <div class="relative min-w-0 m-full" ref="modelDropdownRef">
                            <button
                                type="button"
                                class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-slate-700/60 bg-slate-800/60 text-sm text-slate-300 hover:text-slate-100 hover:border-slate-600 transition-all duration-150 max-w-[140px]"
                                @click="modelDropdownOpen = !modelDropdownOpen"
                            >
                                <span class="truncate">{{ defaultModel }}</span>
                                <svg class="shrink-0 text-slate-500 transition-transform duration-150"
                                    :class="{ 'rotate-180': modelDropdownOpen }"
                                    width="8" height="8" viewBox="0 0 10 10" fill="none">
                                    <path d="M2 3.5L5 6.5L8 3.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                                </svg>
                            </button>

                            <!-- Model menu -->
                            <Transition
                                enter-active-class="transition duration-100 ease-out"
                                enter-from-class="opacity-0 -translate-y-1"
                                leave-active-class="transition duration-75 ease-in"
                                leave-to-class="opacity-0 -translate-y-1"
                            >
                                <div v-if="modelDropdownOpen"
                                    class="absolute bottom-full left-0 mb-1.5 z-50 min-w-[180px] rounded-lg border border-slate-700/80 bg-slate-900 shadow-xl shadow-black/40 overflow-hidden"
                                >
                                    <div class="px-2 pt-2 pb-1 border-b border-slate-800/60">
                                        <p class="text-[9px] font-semibold uppercase tracking-wider text-slate-500">
                                            {{ defaultModel }} Models
                                        </p>
                                    </div>
                                    <div class="py-1 max-h-48 overflow-y-auto">
                                        <template v-for="group in currentModels" :key="group.group">
                                            <div v-if="group.group" class="px-2 pt-2 pb-0.5 text-[9px] font-semibold uppercase tracking-wider text-slate-600">
                                                {{ group.group }}
                                            </div>
                                            <button
                                                v-for="m in group.models"
                                                :key="m.value"
                                                type="button"
                                                :disabled="m.disabled"
                                                :class="[
                                                    'w-full flex items-center justify-between gap-2 px-3 py-1.5 text-left transition-colors duration-100',
                                                    defaultModel === m.value
                                                        ? 'bg-indigo-500/20 text-indigo-300'
                                                        : m.disabled
                                                            ? 'text-slate-600 cursor-not-allowed'
                                                            : 'text-slate-300 hover:bg-slate-800 hover:text-slate-100'
                                                ]"
                                                @click="onModelChange(m.value)"
                                            >
                                                <span class="flex flex-col min-w-0">
                                                    <span class="text-[11px] font-mono truncate">{{ m.label }}</span>
                                                    <span v-if="m.description" class="text-[9px] text-slate-500 truncate">{{ m.description }}</span>
                                                </span>
                                                <svg v-if="defaultModel === m.value"
                                                    class="shrink-0 text-indigo-400" width="10" height="10" viewBox="0 0 10 10" fill="none">
                                                    <path d="M1.5 5L4 7.5L8.5 2.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                                                </svg>
                                            </button>
                                        </template>
                                    </div>
                                </div>
                            </Transition>
                        </div>
                    </div>
                </div>
            </section>

            <!-- ── Footer ──────────────────────────────────────────────── -->
            <div class="flex items-center justify-between mt-8">
                <p v-if="globalStatus" class="text-xs transition-all"
                    :class="globalStatus.ok ? 'text-emerald-400' : 'text-red-400'">
                    {{ globalStatus.msg }}
                </p>
                <div v-else></div>
                <button @click="resetSettings"
                    class="text-sm text-gray-500 hover:text-red-400 transition-colors">
                    Reset all defaults
                </button>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, defineComponent, h, watch } from 'vue'
import { SaveConfigToFile, LoadConfigFromFile, SaveAPIKey, SaveOllamaHost } from '@wails/go/main/App'
import { providerOptions, type ModelGroup } from './components/Agent/parts/Customization/ChatPanel'
import { computed } from 'vue'

// ── Inline sub-components ────────────────────────────────────────────────────

const SaveButton = defineComponent({
    props: { saving: Boolean },
    emits: ['click'],
    setup(props, { emit }) {
        return () => h('button', {
            type: 'button',
            disabled: props.saving,
            onClick: () => emit('click'),
            class: 'shrink-0 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white px-4 py-2 rounded text-sm font-medium transition-colors min-w-[64px]',
        }, props.saving ? 'Saving…' : 'Save')
    },
})

const KeyStatus = defineComponent({
    props: { status: String as () => 'saved' | 'error' | null },
    setup(props) {
        return () => props.status
            ? h('p', {
                class: `mt-1 text-[10px] ${props.status === 'saved' ? 'text-emerald-400' : 'text-red-400'}`,
              }, props.status === 'saved' ? '✓ Saved' : '✕ Failed to save')
            : null
    },
})

// ── State ────────────────────────────────────────────────────────────────────

const keys = ref({ anthropic: '', gemini: '', deepseek: '' })
const showKey = ref({ anthropic: false, gemini: false, deepseek: false })
const keyStatus = ref<Record<string, 'saved' | 'error' | null>>({
    anthropic: null, gemini: null, deepseek: null,
})
const modelDropdownOpen = ref<boolean>(false);
const defaultProvider = ref('ollama')
const defaultModel = ref('qwen3:8b')
const ollamaHost      = ref('http://localhost:11434')
const savingKey       = ref<string | null>(null)
const pingingOllama   = ref(false)
const ollamaStatus    = ref<'ok' | 'error' | null>(null)
const globalStatus    = ref<{ ok: boolean; msg: string } | null>(null)


const modelDropdownRef  = ref<HTMLElement | null>(null)

const currentProvider = computed(() =>
    providerOptions.find(p => p.value === defaultProvider.value) ?? null
)

const currentModels = computed<ModelGroup[]>(() =>
    currentProvider.value?.modelGroups ?? []
)

function onProviderChange(value: string) {
    defaultProvider.value = value;

    const p = providerOptions.find(p => p.value === value)
    if (p) defaultModel.value = p.defaultModel

    modelDropdownOpen.value = false
}

function onModelChange(value: string) {
    defaultModel.value = value;
    modelDropdownOpen.value = false
}

// Close model dropdown on outside click
function onOutsideClick(e: MouseEvent) {
    if (modelDropdownOpen.value && !modelDropdownRef.value?.contains(e.target as Node)) {
        modelDropdownOpen.value = false
    }
}

// ── Helpers ───────────────────────────────────────────────────────────────────

function flashGlobal(ok: boolean, msg: string) {
    globalStatus.value = { ok, msg }
    setTimeout(() => { globalStatus.value = null }, 3000)
}

// ── Save / load ───────────────────────────────────────────────────────────────

async function saveKey(type: 'anthropic' | 'gemini' | 'deepseek') {
    savingKey.value = type
    keyStatus.value[type] = null
    try {
        await SaveAPIKey(type, keys.value[type])
        keyStatus.value[type] = 'saved'
        flashGlobal(true, 'Key saved.')
    } catch {
        keyStatus.value[type] = 'error'
        flashGlobal(false, 'Failed to save key.')
    } finally {
        savingKey.value = null
    } 
}

async function saveOllamaHost() {
    savingKey.value = 'ollamaHost'
    try {
        await SaveOllamaHost(ollamaHost.value)
        flashGlobal(true, 'Ollama host saved.')
    } catch {
        flashGlobal(false, 'Failed to save Ollama host.')
    } finally {
        savingKey.value = null
    }
}

async function saveDefaultProvider() {
    try {
        await SaveConfigToFile('DefaultProvider', defaultProvider.value)
    } catch {}
}

async function saveDefaultModel() {
    try {
        await SaveConfigToFile('DefaultModel', defaultModel.value)
    } catch {}
}

async function pingOllama() {
    pingingOllama.value = true
    ollamaStatus.value  = null
    try {
        const res = await fetch(`${ollamaHost.value}/api/tags`, { signal: AbortSignal.timeout(4000) })
        ollamaStatus.value = res.ok ? 'ok' : 'error'
    } catch {
        ollamaStatus.value = 'error'
    } finally {
        pingingOllama.value = false
    }
}

async function loadSettings() {
    const load = async (k: string, fallback = '') => {
        try { return await LoadConfigFromFile(k) } catch(e) { console.log(e); return fallback }
    }

    keys.value.anthropic    = await load('APIKey_anthropic')
    keys.value.gemini       = await load('APIKey_gemini')
    keys.value.deepseek     = await load('APIKey_deepseek')
    ollamaHost.value        = await load('OllamaHost', 'http://localhost:11434')
    defaultProvider.value   = await load('DefaultProvider', 'ollama')
    defaultModel.value      = await load('DefaultModel', 'qwen3:8b')    
}

function resetSettings() {
    if (!confirm('Reset all defaults? Saved API keys will not be deleted.')) return
    defaultProvider.value = 'ollama'
    defaultModel.value = 'qwen3:8b'
    ollamaHost.value      = 'http://localhost:11434'
    saveDefaultProvider()
    saveOllamaHost()
    flashGlobal(true, 'Defaults reset.')
}

watch(defaultProvider, (v2,v1)=>{
    if(v2 != v1) {
        if(currentModels.value[0]) {
            defaultModel.value = currentModels.value[0].models[0].value;
        }
    }

    saveDefaultProvider();
})

watch(defaultModel, (v2,v1)=>{
    saveDefaultModel();
})

// ── Lifecycle ─────────────────────────────────────────────────────────────────
onMounted(loadSettings)
</script>