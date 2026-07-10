<template>
    <div class="h-full flex w-full flex-col bg-slate-900/40 backdrop-blur-xl"
        :class="{'justify-center':newSession}"
    >
        <template v-if="!newSession">
            <!-- Header -->
            <div class="flex h-9 items-center space-x-2 border-b border-slate-800/60 px-4 bg-slate-900/60">
                <div class="flex h-4 w-4 items-center justify-center text-indigo-400">
                    <i class="ye ye-user-tie"></i>
                </div>
                <span class="text-xs font-bold tracking-wider text-slate-200 uppercase">Agent Assistant</span>
            </div>

            <!-- Chat history -->
            <div class="flex-1 overflow-y-auto p-3 space-y-4" ref="chatHistoryContainer">
                <div v-for="(msg, idx) in messages" :key="idx" :class="[
                    'flex flex-col max-w-[90%] rounded-xl p-3 text-lg leading-relaxed border',
                    msg.role === 'user'
                        ? 'ml-auto bg-slate-800 border-slate-700 text-slate-100'
                        : msg.role === 'tool'
                            ? 'mr-auto bg-slate-950/60 border-slate-800/40 text-slate-500 font-mono'
                            : 'mr-auto bg-slate-900/80 border-slate-800/60 text-slate-300',
                ]">
                    <span class="font-bold text-[9px] uppercase tracking-wider mb-1 block"
                        :class="msg.role === 'user' ? 'text-indigo-300' : msg.role === 'tool' ? 'text-amber-500/70' : 'text-emerald-400'">
                        {{ msg.role === 'user' ? 'You' : msg.role === 'tool' ? 'Tool' : 'Agent' }}
                    </span>
                    <div class="whitespace-pre-wrap text-sm" v-html="msg.content"></div>
                </div>

                <!-- Scaffold progress card -->
                <div v-if="scaffoldProgress.active"
                    class="mr-auto w-full max-w-[90%] rounded-xl border border-slate-800/60 bg-slate-900/80 p-3 space-y-2">
                    <div class="flex items-center justify-between">
                        <span class="text-[9px] font-bold uppercase tracking-wider text-emerald-400">Scaffolding</span>
                        <span class="text-[9px] text-slate-500">
                            {{ scaffoldProgress.index }} / {{ scaffoldProgress.total }}
                        </span>
                    </div>
                    <div class="h-1 w-full rounded-full bg-slate-800 overflow-hidden">
                        <div class="h-full rounded-full bg-indigo-500 transition-all duration-300"
                            :style="{ width: progressPercent + '%' }" />
                    </div>
                    <p class="truncate text-[10px] text-slate-400 font-mono">
                        {{ scaffoldProgress.file || '…' }}
                    </p>
                    <div v-if="scaffoldProgress.files.length"
                        class="max-h-28 overflow-y-auto space-y-0.5 pt-1 border-t border-slate-800/60">
                        <div v-for="(f, i) in scaffoldProgress.files" :key="i"
                            class="flex items-center gap-1.5 text-[10px] text-slate-500 font-mono">
                            <svg class="h-2.5 w-2.5 shrink-0 text-emerald-500" fill="none" viewBox="0 0 24 24"
                                stroke-width="3" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" d="m4.5 12.75 6 6 9-13.5" />
                            </svg>
                            <span class="truncate">{{ f }}</span>
                        </div>
                    </div>
                </div>

                <!-- Thinking indicator -->
                <div v-if="isAiThinking"
                    class="flex items-center space-x-1.5 mr-auto bg-slate-900/80 border border-slate-800/60 rounded-xl p-3 text-xs text-slate-500">
                    <span class="h-1 w-1 animate-bounce rounded-full bg-slate-400" style="animation-delay: 0ms" />
                    <span class="h-1 w-1 animate-bounce rounded-full bg-slate-400" style="animation-delay: 150ms" />
                    <span class="h-1 w-1 animate-bounce rounded-full bg-slate-400" style="animation-delay: 300ms" />
                </div>
            </div>
        </template>

        <!-- Input -->
        <div class="p-3 w-full">
            <form @submit.prevent="sendMessage" class="relative flex items-center">
                <div class="w-full flex flex-col gap-3">
                    <template v-if="newSession">
                        <div class="flex items-center text-lg text-gray-300 gap-2">
                            <span>New session in</span>
                            <button type="button" class="flex items-center gap-1 hover:text-white">
                                <span><i class="ye ye-folder-open text-blue-500 ye-sm px-2"></i> {{ store.active?.name }}</span>
                                <span class="text-xs">⌄</span>
                            </button>
                        </div>
                    </template>

                    <div class="border border-gray-400/20 rounded-xl bg-slate-900 shadow-lg overflow-visible flex flex-col">
                        <div class="p-3">
                            <textarea v-model="userInput" placeholder="What's your next milestone?"
                                class="w-full bg-transparent leading-none text-gray-200 placeholder-gray-600 resize-none outline-none min-h-[60px]"
                                @keydown.enter.exact.prevent="sendMessage"
                            ></textarea>
                        </div>

                        <div class="px-2 py-2 flex items-center justify-between border-t border-slate-800/60 gap-2">
                            <!-- Left: provider + model selectors -->
                            <div class="flex items-center gap-2 min-w-0">
                                <!-- Provider pill selector -->
                                
                                <!-- Model dropdown -->
                                <div class="relative min-w-0" ref="modelDropdownRef">
                                    <button
                                        type="button"
                                        class="flex items-center gap-1.5 px-2 py-1 rounded-md border border-slate-700/60 bg-slate-800/60 text-2 text-slate-300 hover:text-slate-100 hover:border-slate-600 transition-all duration-150 max-w-35"
                                        @click="providerModelDropdownOpen = !providerModelDropdownOpen"
                                    >
                                        <span class="truncate font-mono block leading-none">{{ selectedProvider }}</span>
                                        <svg class="shrink-0 text-slate-500 transition-transform duration-150"
                                            :class="{ 'rotate-180': providerModelDropdownOpen }"
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
                                        <div v-if="providerModelDropdownOpen"
                                            class="absolute bottom-full left-0 mb-1.5 z-50 min-w-45 rounded-lg border border-slate-700/80 bg-slate-900 shadow-xl shadow-black/40 overflow-hidden"
                                        >
                                            <div class="px-2 pt-2 pb-1 border-b border-slate-800/60">
                                                <p class=" font-semibold uppercase tracking-wider text-slate-500">
                                                    Select Provider
                                                </p>
                                            </div>
                                            <div class="py-1 max-h-48 overflow-y-auto">
                                                <template v-for="group in providerOptions" :key="group.value">
                                                    
                                                    <button
                                                        type="button"
                                                        :class="[
                                                            'w-full flex items-center justify-between gap-2 px-3 py-1.5 text-left transition-colors duration-100',
                                                            selectedProvider === group.value
                                                                ? 'bg-indigo-500/20 text-indigo-300'
                                                                : 'text-slate-300 hover:bg-slate-800 hover:text-slate-100'
                                                        ]"
                                                        @click="onProviderChange(group.value)"
                                                    >
                                                        <span class="flex flex-col min-w-0">
                                                            <span class=" truncate leading-none">{{ group.label }}</span>
                                                        </span>
                                                        <svg v-if="selectedProvider === group.value"
                                                            class="shrink-0 text-indigo-400" width="10" height="10" viewBox="0 0 10 10" fill="none">
                                                            <path d="M1.5 5L4 7.5L8.5 2.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                                                        </svg>
                                                    </button>
                                                </template>
                                            </div>
                                        </div>
                                    </Transition>
                                </div>

                                <!-- Model dropdown -->
                                <div class="relative min-w-0" ref="modelDropdownRef">
                                    <button
                                        type="button"
                                        class="flex items-center gap-1.5 px-2 py-1 rounded-md border border-slate-700/60 bg-slate-800/60 text-slate-300 hover:text-slate-100 hover:border-slate-600 transition-all duration-150 max-w-[140px]"
                                        @click="modelDropdownOpen = !modelDropdownOpen"
                                    >
                                        <span class="truncate font-mono leading-none">{{ selectedModel }}</span>
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
                                            <div class="px-2 pt-2 pb-2 border-b border-slate-800/60">
                                                <p class=" font-semibold uppercase leading-none tracking-wider text-slate-500">
                                                    {{ currentProvider?.label }} Models
                                                </p>
                                            </div>
                                            <div class="py-2 max-h-48 overflow-y-auto">
                                                <template v-for="group in groupedModels" :key="group.group">
                                                    <div v-if="group.group" class="px-2 pt-2 pb-0.5 leading-none  font-semibold uppercase tracking-wider text-slate-600">
                                                        {{ group.group }}
                                                    </div>
                                                    <button
                                                        v-for="m in group.models"
                                                        :key="m.value"
                                                        type="button"
                                                        :disabled="m.disabled"
                                                        :class="[
                                                            'w-full flex items-center justify-between gap-2 px-3 py-2 text-left transition-colors duration-100',
                                                            selectedModel === m.value
                                                                ? 'bg-indigo-500/20 text-indigo-300'
                                                                : m.disabled
                                                                    ? 'text-slate-600 cursor-not-allowed'
                                                                    : 'text-slate-300 hover:bg-slate-800 hover:text-slate-100'
                                                        ]"
                                                        @click="onModelChange(m.value)"
                                                    >
                                                        <span class="flex flex-col min-w-0 ">
                                                            <span class=" font-mono truncate leading-none mb-1">{{ m.label }}</span>
                                                            <span v-if="m.description" class="text-sm text-slate-500 truncate leading-none">{{ m.description }}</span>
                                                        </span>
                                                        <svg v-if="selectedModel === m.value"
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

                            <!-- Right: send -->
                            <button
                                type="submit"
                                :disabled="isAiThinking || !userInput.trim()"
                                class="shrink-0 flex items-center justify-center w-6 h-6 rounded-md transition-all duration-150"
                                :class="isAiThinking || !userInput.trim()
                                    ? 'text-slate-700 cursor-not-allowed'
                                    : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800'"
                            >
                                <i class="ye ye-paper-plane text-xs"></i>
                            </button>
                        </div>
                    </div>

                    <div class="flex items-center justify-between text-xs text-gray-500 px-1">
                        <div class="flex items-center gap-1">
                            <i class="ye ye-key ye-sm"></i> Default Approvals
                        </div>
                        <div class="flex items-center gap-4">
                            <button type="button" class="flex items-center gap-1 hover:text-gray-300"><i class="ye ye-folder-open ye-sm"></i> Folder</button>
                            <button type="button" class="flex items-center gap-1 hover:text-gray-300"><i class="ye ye-leaf ye-sm"></i> main</button>
                        </div>
                    </div>
                </div>
            </form>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { useWorkspace } from '@/store/workspace'
import { EventsOn, EventsOff } from '@wails/runtime/runtime'
import { AgentChat, LoadConfigFromFile, SaveConfigToFile } from '@wails/go/main/App'
import { providerOptions, type ModelGroup, type ProviderConfig } from './ChatPanel'


// ── Types ────────────────────────────────────────────────────────────────────

interface ScaffoldProgress {
    file: string
    index: number
    total: number
    done: boolean
    error?: string
}

const props = defineProps<{
    sessionId?: string | null
    small: boolean
}>()

const store = useWorkspace()

const selectedProvider = ref<string>('ollama')
const selectedModel    = ref<string>('qwen3:1.7b')
const providerModelDropdownOpen = ref(false)
const modelDropdownOpen = ref(false)
const modelDropdownRef  = ref<HTMLElement | null>(null)

const chatHistoryContainer = ref<HTMLElement | null>(null)
const userInput = ref('')
const isAiThinking = ref(false)
const currentSessionID = ref<string>(props.sessionId || 'session-' + Math.random().toString(36).substring(2, 9))

// Use the store's session messages when a session is selected
const messages = ref<{ role: string; content: string }[]>([])

const newSession = computed<boolean>(() => (props.small? false : messages.value.length === 0))

// ── Scaffold progress ────────────────────────────────────────────────────────

const scaffoldProgress = ref({
    active: false,
    file: '',
    index: 0,
    total: 0,
    files: [] as string[],
})

const progressPercent = computed(() => {
    if (!scaffoldProgress.value.total) return 0
    return Math.round((scaffoldProgress.value.index / scaffoldProgress.value.total) * 100)
})


const currentProvider = computed(() =>
    providerOptions.find(p => p.value === selectedProvider.value) ?? null
)

const groupedModels = computed<ModelGroup[]>(() =>
    currentProvider.value?.modelGroups ?? []
)

function onProviderChange(value: string) {
    selectedProvider.value = value;

    const p = providerOptions.find(p => p.value === value)
    if (p) selectedModel.value = p.defaultModel
    providerModelDropdownOpen.value = false

    try {
        SaveConfigToFile('DefaultProvider', value)
    } catch {}
}

function onModelChange(value: string) {
    selectedModel.value = value
    modelDropdownOpen.value = false

    try {
        SaveConfigToFile('DefaultModel', value)
    } catch {}
}

// Close model dropdown on outside click
function onOutsideClick(e: MouseEvent) {
    if (modelDropdownOpen.value && !modelDropdownRef.value?.contains(e.target as Node)) {
        modelDropdownOpen.value = false
    }
}


// ── Store / refs ─────────────────────────────────────────────────────────────

function resetScaffold() {
    scaffoldProgress.value = { active: false, file: '', index: 0, total: 0, files: [] }
}

// ── Message helpers ──────────────────────────────────────────────────────────

const appendOrUpdateAssistant = (text: string) => {
    const last = messages.value[messages.value.length - 1]
    if (last && last.role === 'assistant') {
        last.content += text
    } else {
        messages.value.push({ role: 'assistant', content: text })
    }
}

// ── Watch for session selection changes ──────────────────────────────────────

watch(() => store.active?.sessionId, (newId) => {
    if (newId) {
        currentSessionID.value = newId;
        // Load the session's messages from the store (populated by selectSession)
        messages.value = [...store.activeSessionMessages];
    }
});

// When the store loads new messages (from selectSession), update local messages
watch(() => store.activeSessionMessages, (newMessages) => {
    if (store.active?.sessionId && newMessages.length > 0) {
        messages.value = [...newMessages];
    }
}, { deep: true });

// ── Wails events ─────────────────────────────────────────────────────────────
 
onMounted(async () => {
    document.addEventListener('mousedown', onOutsideClick, true)

    try {
        selectedProvider.value = await LoadConfigFromFile("DefaultProvider")
    } catch (error) {
        console.warn(error);
    }
    try {
        selectedModel.value = await LoadConfigFromFile("DefaultModel")
    } catch (error) {
        console.warn(error);        
    }

    store.selectSession(currentSessionID.value)
    messages.value = [...store.activeSessionMessages];

    EventsOn('scaffold:progress', (p: ScaffoldProgress) => {
        if (p.error) {
            isAiThinking.value = false
            messages.value.push({ role: 'assistant', content: `❌ Scaffold failed:\n${p.error}` })
            resetScaffold()
            scrollToBottom()
            return
        }
        if (p.done) {
            isAiThinking.value = false
            scaffoldProgress.value.active = false
            messages.value.push({
                role: 'assistant',
                content: `✅ Done — ${p.total} file${p.total !== 1 ? 's' : ''} written successfully.`,
            })
            scrollToBottom() 
            return
        }
        scaffoldProgress.value.active = true
        scaffoldProgress.value.file   = p.file
        scaffoldProgress.value.index  = p.index
        scaffoldProgress.value.total  = p.total
        if (p.file && p.file !== 'Contacting LLM…' && p.file !== 'Parsing plan…') {
            scaffoldProgress.value.files.push(p.file)
        }
        scrollToBottom() 
    })

    EventsOn('agent:message', (text: string) => {
        appendOrUpdateAssistant(text)
        scrollToBottom()
    })

    EventsOn('agent:tool', (call: { name: string; input: string }) => {
        // console.log('Tool call:', call)
        messages.value.push({ role: 'tool', content: `🔧 ${call.name} => ${call.input}` })
        scrollToBottom()
    })

    EventsOn('agent:done', () => {
        isAiThinking.value = false
        store.fetchSessions()
    })

    EventsOn('agent:error', (msg: string) => {
        isAiThinking.value = false
        messages.value.push({ role: 'assistant', content: `❌ ${msg}` })
    })
})

onUnmounted(() => {
    document.removeEventListener('mousedown', onOutsideClick, true)
    EventsOff('scaffold:progress')
    EventsOff('agent:message')
    EventsOff('agent:tool')
    EventsOff('agent:done')
    EventsOff('agent:error')
})

// ── Scroll ───────────────────────────────────────────────────────────────────

const scrollToBottom = async () => {
    await nextTick()
    if (chatHistoryContainer.value) {
        chatHistoryContainer.value.scrollTop = chatHistoryContainer.value.scrollHeight
    }
}

// ── Send ─────────────────────────────────────────────────────────────────────

const sendMessage = async () => {
    if (!userInput.value.trim() || isAiThinking.value) return

    const text = userInput.value.trim()
    userInput.value = ''
    messages.value.push({ role: 'user', content: text })
    scrollToBottom()

    isAiThinking.value = true;
    resetScaffold()

    try {
        await AgentChat(currentSessionID.value, text, selectedProvider.value, selectedModel.value)
    } catch (err: any) {
        isAiThinking.value = false
        messages.value.push({
            role: 'assistant',
            content: `❌ ${err?.message ?? String(err)}`,
        })
        scrollToBottom()
    }
}
</script>

<style lang="scss">
think {
    font-style: italic;
    opacity: 0.5;
}
</style>