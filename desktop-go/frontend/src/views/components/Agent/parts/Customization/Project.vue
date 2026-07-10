<template>
    <div class="h-full flex w-full flex-col bg-slate-900/40 backdrop-blur-xl overflow-y-auto">
        <!-- Header -->
        <div class="flex h-9 items-center space-x-2 border-b border-slate-800/60 px-4 bg-slate-900/60 shrink-0">
            <i class="ye ye-folder-open text-indigo-400 text-sm"></i>
            <span class="text-sm font-bold tracking-wider text-slate-200 uppercase">Generate Project</span>
        </div>

        <!-- Form -->
        <form @submit.prevent="onSubmit" class="flex-1 overflow-y-auto p-5 space-y-5">
            <!-- Project Name -->
            <div>
                <label class="block text-sm font-medium text-slate-400 mb-1.5">
                    Project Name <span class="text-red-400">*</span>
                </label>
                <input v-model="form.project_name" type="text" placeholder="e.g. my-awesome-app"
                    class="w-full h-10 bg-slate-900 border border-slate-700/60 rounded-lg px-3 py-2 text-sm text-slate-200 placeholder-slate-600 outline-none focus:border-indigo-500/60 transition-colors"
                    required />
            </div>

            <!-- Description -->
            <div>
                <label class="block text-sm font-medium text-slate-400 mb-1.5">
                    Description <span class="text-red-400">*</span>
                </label>
                <textarea v-model="form.description" rows="3" placeholder="Detailed description of what the project does"
                    class="w-full h-10 bg-slate-900 border border-slate-700/60 rounded-lg px-3 py-2 text-sm text-slate-200 placeholder-slate-600 outline-none focus:border-indigo-500/60 transition-colors resize-none"
                    required></textarea>
            </div>

            <!-- Size & Complexity row -->
            <div class="grid grid-cols-2 gap-4">
                <div>
                    <label class="block text-sm font-medium text-slate-400 mb-1.5">
                        Size <span class="text-red-400">*</span>
                    </label>
                    <select v-model="form.size"
                        class="w-full h-10 bg-slate-900 border border-slate-700/60 rounded-lg px-3 py-2 text-sm text-slate-200 outline-none focus:border-indigo-500/60 transition-colors appearance-none cursor-pointer"
                        required>
                        <option value="" disabled>Select size</option>
                        <option value="small">Small</option>
                        <option value="middle">Middle</option>
                        <option value="large">Large</option>
                    </select>
                </div>
                <div>
                    <label class="block text-sm font-medium text-slate-400 mb-1.5">
                        Complexity <span class="text-red-400">*</span>
                    </label>
                    <select v-model="form.complexity"
                        class="w-full h-10 bg-slate-900 border border-slate-700/60 rounded-lg px-3 py-2 text-sm text-slate-200 outline-none focus:border-indigo-500/60 transition-colors appearance-none cursor-pointer"
                        required>
                        <option value="" disabled>Select complexity</option>
                        <option value="low">Low</option>
                        <option value="medium">Medium</option>
                        <option value="high">High</option>
                    </select>
                </div>
            </div>

            <!-- Framework -->
            <div>
                <label class="block text-sm font-medium text-slate-400 mb-1.5">
                    Framework / Language <span class="text-slate-600">(optional)</span>
                </label>
                <input v-model="form.framework" type="text" placeholder="e.g. Go, Vue.js, Python Flask (auto-recommended if omitted)"
                    class="w-full h-10 bg-slate-900 border border-slate-700/60 rounded-lg px-3 py-2 text-sm text-slate-200 placeholder-slate-600 outline-none focus:border-indigo-500/60 transition-colors" />
            </div>

            <!-- Modules -->
            <div>
                <label class="block text-sm font-medium text-slate-400 mb-1.5">
                    Modules / Features <span class="text-red-400">*</span>
                </label>
                <input v-model="form.modules" type="text" placeholder="Comma-separated: auth, payments, dashboard, api"
                    class="w-full h-10 bg-slate-900 border border-slate-700/60 rounded-lg px-3 py-2 text-sm text-slate-200 placeholder-slate-600 outline-none focus:border-indigo-500/60 transition-colors"
                    required />
                <p class="mt-1 text-sm text-slate-600">Comma-separated list of modules/features to include</p>
            </div>

            <!-- Features (detailed) -->
            <div>
                <label class="block text-sm font-medium text-slate-400 mb-1.5">
                    Detailed Feature Requirements <span class="text-red-400">*</span>
                </label>
                <textarea v-model="form.features" rows="4" placeholder="Describe each module in detail: what it should do, key functionality, user interactions..."
                    class="w-full h-10 bg-slate-900 border border-slate-700/60 rounded-lg px-3 py-2 text-sm text-slate-200 placeholder-slate-600 outline-none focus:border-indigo-500/60 transition-colors resize-none"
                    required></textarea>
            </div>

            <!-- DB Tables row -->
            <div class="grid grid-cols-2 gap-4">
                <div>
                    <label class="block text-sm font-medium text-slate-400 mb-1.5">
                        Min DB Tables <span class="text-red-400">*</span>
                    </label>
                    <input v-model.number="form.db_tables_min" type="number" min="0" placeholder="1"
                        class="w-full h-10 bg-slate-900 border border-slate-700/60 rounded-lg px-3 py-2 text-sm text-slate-200 placeholder-slate-600 outline-none focus:border-indigo-500/60 transition-colors"
                        required />
                </div>
                <div>
                    <label class="block text-sm font-medium text-slate-400 mb-1.5">
                        Max DB Tables <span class="text-red-400">*</span>
                    </label>
                    <input v-model.number="form.db_tables_max" type="number" min="0" placeholder="10"
                        class="w-full h-10 bg-slate-900 border border-slate-700/60 rounded-lg px-3 py-2 text-sm text-slate-200 placeholder-slate-600 outline-none focus:border-indigo-500/60 transition-colors"
                        required />
                </div>
            </div>

            <!-- Root Path -->
            <div>
                <label class="block text-sm font-medium text-slate-400 mb-1.5">
                    Root Path <span class="text-red-400">*</span>
                </label>
                <div class="flex gap-2">
                    <input v-model="form.root_path" type="text" :placeholder="store.active?.path || '/path/to/project'"
                        class="flex-1 bg-slate-900 border border-slate-700/60 rounded-lg px-3 py-2 text-sm text-slate-200 placeholder-slate-600 outline-none focus:border-indigo-500/60 transition-colors font-mono"
                        required />
                    <button type="button" @click="pickFolder"
                        class="shrink-0 px-3 py-2 rounded-lg border border-slate-700/60 bg-slate-800 text-sm text-slate-400 hover:text-slate-200 hover:border-slate-600 transition-colors"
                        title="Browse…">
                        <i class="ye ye-folder-open"></i>
                    </button>
                </div>
                <p v-if="store.active?.path" class="mt-1 text-sm text-slate-600">
                    Current workspace: {{ store.active.path }}
                </p>
            </div>

            <!-- Submit -->
            <div class="flex items-center justify-between pt-2 border-t border-slate-800/60">
                <div class="flex items-center gap-2">
                    <!-- Provider + Model pill selectors (reusing ChatPanel types) -->
                    <div class="relative" ref="providerDropdownRef">
                        <button type="button"
                            class="flex items-center gap-1.5 px-2.5 py-1.5 leading-none rounded-md border border-slate-700/60 bg-slate-800/60 text-sm text-slate-300 hover:text-slate-100 hover:border-slate-600 transition-all"
                            @click="providerOpen = !providerOpen"
                        >
                            <span>{{ selectedProviderLabel }}</span>
                            <svg class="shrink-0 text-slate-500 transition-transform" :class="{ 'rotate-180': providerOpen }" width="7" height="7" viewBox="0 0 10 10" fill="none">
                                <path d="M2 3.5L5 6.5L8 3.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                            </svg>
                        </button>
                        <!-- Provider dropdown -->
                        <Transition name="fade-slide">
                            <div v-if="providerOpen"
                                class="absolute bottom-full left-0 mb-1.5 z-50 min-w-40 rounded-lg border border-slate-700/80 bg-slate-900 shadow-xl shadow-black/40 overflow-hidden">
                                <div class="py-1 max-h-48 overflow-y-auto">
                                    <button v-for="p in providerOptions" :key="p.value" type="button"
                                        :class="[
                                            'w-full flex items-center gap-2 px-3 py-1.5 text-left transition-colors text-sm',
                                            selectedProvider === p.value
                                                ? 'bg-indigo-500/20 text-indigo-300'
                                                : 'text-slate-300 hover:bg-slate-800'
                                        ]"
                                        @click="selectProvider(p.value)">
                                        <span>{{ p.icon }}</span>
                                        <span>{{ p.label }}</span>
                                    </button>
                                </div>
                            </div>
                        </Transition>
                    </div>

                    <div class="relative" ref="modelDropdownRef">
                        <button type="button"
                            class="flex items-center gap-1.5 px-2.5 py-1.5 rounded-md border border-slate-700/60 bg-slate-800/60 text-sm text-slate-300 hover:text-slate-100 hover:border-slate-600 transition-all"
                            @click="modelOpen = !modelOpen"
                        >
                            <span class="max-w-[100px] truncate font-mono leading-none ">{{ selectedModel }}</span>
                            <svg class="shrink-0 text-slate-500 transition-transform" :class="{ 'rotate-180': modelOpen }" width="7" height="7" viewBox="0 0 10 10" fill="none">
                                <path d="M2 3.5L5 6.5L8 3.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                            </svg>
                        </button>
                        <!-- Model dropdown -->
                        <Transition name="fade-slide">
                            <div v-if="modelOpen"
                                class="absolute bottom-full left-0 mb-1.5 z-50 min-w-[160px] rounded-lg border border-slate-700/80 bg-slate-900 shadow-xl shadow-black/40 overflow-hidden">
                                <div class="px-2 pt-2 pb-1 border-b border-slate-800/60">
                                    <p class="text-[9px] font-semibold uppercase tracking-wider text-slate-500">{{ currentProvider?.label }} Models</p>
                                </div>
                                <div class="py-1 max-h-48 overflow-y-auto">
                                    <template v-for="group in groupedModels" :key="group.group">
                                        <div v-if="group.group" class="px-2 pt-2 pb-0.5 text-[9px] font-semibold uppercase tracking-wider text-slate-600">
                                            {{ group.group }}
                                        </div>
                                        <button v-for="m in group.models" :key="m.value" type="button"
                                            :disabled="m.disabled"
                                            :class="[
                                                'w-full flex items-center justify-between gap-2 px-3 py-1.5 text-left transition-colors text-sm',
                                                selectedModel === m.value
                                                    ? 'bg-indigo-500/20 text-indigo-300'
                                                    : m.disabled ? 'text-slate-600 cursor-not-allowed' : 'text-slate-300 hover:bg-slate-800'
                                            ]"
                                            @click="selectModel(m.value)">
                                            <span class="flex flex-col">
                                                <span class="font-mono truncate">{{ m.label }}</span>
                                                <span v-if="m.description" class="text-[9px] text-slate-500">{{ m.description }}</span>
                                            </span>
                                        </button>
                                    </template>
                                </div>
                            </div>
                        </Transition>
                    </div>
                </div>

                <button type="submit" :disabled="submitting"
                    class="flex items-center gap-1.5 px-4 py-2 rounded-lg bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed text-sm font-medium text-white transition-colors">
                    <template v-if="submitting">
                        <svg class="animate-spin h-3.5 w-3.5" viewBox="0 0 24 24" fill="none">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                        </svg>
                        Generating…
                    </template>
                    <template v-else>
                        <i class="ye ye-paper-plane text-sm"></i>
                        Generate Project
                    </template>
                </button>
            </div>

            <!-- Feedback message -->
            <p v-if="feedback" :class="[
                'text-sm transition-opacity',
                feedback.type === 'success' ? 'text-emerald-400' : 'text-red-400'
            ]">{{ feedback.message }}</p>
        </form>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useWorkspace } from '@/store/workspace'
import { EventsOn, EventsOff } from '@wails/runtime/runtime'
import { AgentChat, LoadConfigFromFile } from '@wails/go/main/App'
import { providerOptions, type ProviderConfig, type ModelGroup } from '@/views/components/Chat/ChatPanel'

// ── Store ────────────────────────────────────────────────────────────────────

const store = useWorkspace()

// ── Form state ───────────────────────────────────────────────────────────────

const form = reactive({
    project_name: '',
    description: '',
    size: '' as 'small' | 'middle' | 'large' | '',
    complexity: '' as 'low' | 'medium' | 'high' | '',
    modules: '',
    features: '',
    db_tables_min: 1,
    db_tables_max: 10,
    framework: '',
    root_path: store.active?.path || '',
})

const submitting = ref(false)
const feedback = ref<{ type: 'success' | 'error'; message: string } | null>(null)

// ── Provider / Model selectors ───────────────────────────────────────────────

const selectedProvider = ref<string>('ollama')
const selectedModel    = ref<string>('qwen3:8b')
const providerOpen = ref(false)
const modelOpen    = ref(false)
const providerDropdownRef = ref<HTMLElement | null>(null)
const modelDropdownRef    = ref<HTMLElement | null>(null)

const selectedProviderLabel = computed(() =>
    providerOptions.find(p => p.value === selectedProvider.value)?.label ?? selectedProvider.value
)

const currentProvider = computed<ProviderConfig | null>(() =>
    providerOptions.find(p => p.value === selectedProvider.value) ?? null
)

const groupedModels = computed<ModelGroup[]>(() =>
    currentProvider.value?.modelGroups ?? []
)

function selectProvider(value: string) {
    selectedProvider.value = value
    const p = providerOptions.find(p => p.value === value)
    if (p) selectedModel.value = p.defaultModel
    providerOpen.value = false
}

function selectModel(value: string) {
    selectedModel.value = value
    modelOpen.value = false
}

// ── Close dropdowns on outside click ─────────────────────────────────────────

function onOutsideClick(e: MouseEvent) {
    if (providerOpen.value && !providerDropdownRef.value?.contains(e.target as Node)) {
        providerOpen.value = false
    }
    if (modelOpen.value && !modelDropdownRef.value?.contains(e.target as Node)) {
        modelOpen.value = false
    }
}

// ── Folder picker ────────────────────────────────────────────────────────────

async function pickFolder() {
    try {
        const { OpenWorkspaceDialog } = await import('@wails/go/main/App')
        const path = await OpenWorkspaceDialog()
        if (path) form.root_path = path
    } catch {
        // Fallback: use current workspace path
        if (store.active?.path) form.root_path = store.active.path
    }
}

// ── Scaffold progress listener ───────────────────────────────────────────────

interface ScaffoldProgress {
    file: string
    index: number
    total: number
    done: boolean
    error?: string
}

onMounted(async () => {
    document.addEventListener('mousedown', onOutsideClick, true)

    // Load saved defaults
    try {
        selectedProvider.value = await LoadConfigFromFile("DefaultProvider") || 'ollama'
        selectedModel.value = await LoadConfigFromFile("DefaultModel") || 'qwen3:8b'
    } catch { /* use defaults */ }

    // Listen for scaffold progress
    EventsOn('scaffold:progress', (p: ScaffoldProgress) => {
        if (p.error) {
            submitting.value = false
            feedback.value = { type: 'error', message: `Scaffold failed: ${p.error}` }
            return
        }
        if (p.done) {
            submitting.value = false
            feedback.value = {
                type: 'success',
                message: `✅ Project created — ${p.total} file${p.total !== 1 ? 's' : ''} written successfully.`,
            }
            return
        }
        // Still in progress – keep submitting true
    })

    EventsOn('agent:error', (msg: string) => {
        submitting.value = false
        feedback.value = { type: 'error', message: `Agent error: ${msg}` }
    })

    EventsOn('agent:done', () => {
        submitting.value = false
    })
})

onUnmounted(() => {
    document.removeEventListener('mousedown', onOutsideClick, true)
    EventsOff('scaffold:progress')
    EventsOff('agent:error')
    EventsOff('agent:done')
})

// ── Validation ───────────────────────────────────────────────────────────────

function validate(): boolean {
    if (!form.project_name.trim()) {
        feedback.value = { type: 'error', message: 'Project name is required.' }
        return false
    }
    if (!form.description.trim()) {
        feedback.value = { type: 'error', message: 'Description is required.' }
        return false
    }
    if (!form.size) {
        feedback.value = { type: 'error', message: 'Please select a project size.' }
        return false
    }
    if (!form.complexity) {
        feedback.value = { type: 'error', message: 'Please select a complexity level.' }
        return false
    }
    if (!form.modules.trim()) {
        feedback.value = { type: 'error', message: 'Please list at least one module or feature.' }
        return false
    }
    if (!form.features.trim()) {
        feedback.value = { type: 'error', message: 'Please describe the features in detail.' }
        return false
    }
    if (form.db_tables_min < 0) {
        feedback.value = { type: 'error', message: 'Minimum DB tables must be 0 or more.' }
        return false
    }
    if (form.db_tables_max < form.db_tables_min) {
        feedback.value = { type: 'error', message: 'Max DB tables must be >= min DB tables.' }
        return false
    }
    if (!form.root_path.trim()) {
        feedback.value = { type: 'error', message: 'Root path is required.' }
        return false
    }
    return true
}

// ── Submit ───────────────────────────────────────────────────────────────────

async function onSubmit() {
    feedback.value = null

    if (!validate()) return

    submitting.value = true

    // Build a structured prompt from the form data
    const prompt = `Create a project with the following specifications:

Project Name: ${form.project_name}
Description: ${form.description}
Size: ${form.size}
Complexity: ${form.complexity}
Modules: ${form.modules}
Detailed Features: ${form.features}
Min DB Tables: ${form.db_tables_min}
Max DB Tables: ${form.db_tables_max}
Framework: ${form.framework || '(auto-recommend)'}
Root Path: ${form.root_path}

Please scaffold the entire project with all necessary files.`

    // Generate a session ID for this project creation
    const sessionID = 'proj-' + Math.random().toString(36).substring(2, 9)

    try {
        await AgentChat(sessionID, prompt, selectedProvider.value, selectedModel.value)
    } catch (err: any) {
        submitting.value = false
        feedback.value = {
            type: 'error',
            message: `Failed to send: ${err?.message ?? String(err)}`,
        }
    }
}
</script>

<style scoped>
/* Custom select arrow */
select {
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='10' height='10' viewBox='0 0 10 10' fill='none'%3E%3Cpath d='M2 3.5L5 6.5L8 3.5' stroke='%236b7280' stroke-width='1.5' stroke-linecap='round' stroke-linejoin='round'/%3E%3C/svg%3E");
    background-repeat: no-repeat;
    background-position: right 10px center;
    padding-right: 28px;
}

/* Number input spinner hide */
input[type="number"]::-webkit-inner-spin-button,
input[type="number"]::-webkit-outer-spin-button {
    opacity: 0.5;
}

/* Transition for dropdowns */
.fade-slide-enter-active {
    transition: opacity 120ms, transform 120ms cubic-bezier(.2, 0, 0, 1);
}
.fade-slide-leave-active {
    transition: opacity 100ms, transform 100ms cubic-bezier(.4, 0, 1, 1);
}
.fade-slide-enter-from,
.fade-slide-leave-to {
    opacity: 0;
    transform: translateY(4px);
}
</style>
