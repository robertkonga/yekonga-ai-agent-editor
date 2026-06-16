<template>
    <div class="h-full flex w-full flex-col bg-slate-900/40 backdrop-blur-xl"
        :class="{'justify-center':newSession}"
    >
        <template v-if="!newSession">
            <!-- Header -->
            <div class="flex h-9 items-center space-x-2 border-b border-slate-800/60 px-4 bg-slate-900/60">
                <div class="flex h-4 w-4 items-center justify-center  text-indigo-400 ">
                    <i class="ye ye-user-tie"></i>
                </div>
                <span class="text-xs font-bold tracking-wider text-slate-200 uppercase">Agent Assistant</span>
            </div>
    
            <!-- Chat history -->
            <div class="flex-1 overflow-y-auto p-3 space-y-4" ref="chatHistoryContainer">
                <div v-for="(msg, idx) in messages" :key="idx" :class="[
                    'flex flex-col max-w-[90%] rounded-xl p-3 text-xs leading-relaxed border',
                    msg.role === 'user'
                        ? 'ml-auto bg-slate-800 border-slate-700 text-slate-100'
                        : 'mr-auto bg-slate-900/80 border-slate-800/60 text-slate-300',
                ]">
                    <span class="font-bold text-[9px] uppercase tracking-wider mb-1 block"
                        :class="msg.role === 'user' ? 'text-indigo-300' : 'text-emerald-400'">
                        {{ msg.role === 'user' ? 'You' : 'Agent' }}
                    </span>
                    <p class="whitespace-pre-wrap text-[11px]">{{ msg.content }}</p>
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
    
                    <!-- Progress bar -->
                    <div class="h-1 w-full rounded-full bg-slate-800 overflow-hidden">
                        <div class="h-full rounded-full bg-indigo-500 transition-all duration-300"
                            :style="{ width: progressPercent + '%' }" />
                    </div>
    
                    <!-- Current file -->
                    <p class="truncate text-[10px] text-slate-400 font-mono">
                        {{ scaffoldProgress.file || '…' }}
                    </p>
    
                    <!-- Written files list -->
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
                            <button class="flex items-center gap-1 hover:text-white">
                                <span><i class="ye ye-folder-open text-blue-500 ye-sm px-2"></i> {{ store.active?.name }}</span>
                                <span class="text-xs">⌄</span>
                            </button>
                            <template v-if="false">
                                <span>with</span>
                                <button class="flex items-center gap-1 hover:text-white">
                                    <span><i class="ye ye-code ye-sm"></i> Copilot CLI</span>
                                    <span class="text-xs">⌄</span>
                                </button>
                            </template>
                        </div>
                    </template>
        
                    <div class="border border-gray-400/20 rounded-xl bg-slate-900 shadow-lg overflow-hidden flex flex-col">
                        <div v-if="false" class="bg-black/20 border-b border-slate-800/60 p-4 flex flex-col items-start justify-between">
                            <div class="flex w-full items-start justify-between gap-3">
                                <div class="flex items-center space-x-3">
                                    <span class="text-blue-400">ⓘ</span>
                                    <h3 class="text-sm font-semibold text-white">Credit Limit Reached</h3>
                                </div>
                                <button class="text-gray-500 hover:text-gray-300">✕</button>
                            </div>
                            <div class="flex  w-full items-start justify-between">
                                <p class="text-sm text-gray-400 mt-1">Upgrade to keep going.</p>
                                <!-- <button class="bg-blue-600 hover:bg-blue-500 text-white text-xs px-4 py-1 rounded-md">
                                    Upgrade
                                </button> -->
                            </div>
                        </div>
        
                        <div class="p-3">
                            <textarea v-model="userInput" placeholder="What's your next milestone?"
                                class="w-full bg-transparent text-gray-200 placeholder-gray-600 resize-none outline-none min-h-[60px]"></textarea>
                        </div>
        
                        <div class="px-4 py-2 flex items-center justify-between border-t border-slate-800/60">
                            <div class="flex items-center gap-4 text-xs text-gray-500">
                                <button class="hover:text-gray-300"><i class="ye ye-plus"></i></button>
                                <button class="hover:text-gray-300 flex items-center gap-1">
                                    <span>⧉</span> Agent
                                </button>
                                <!-- <span class="text-gray-600">|</span> -->
                                <!-- <button class="hover:text-gray-300">Claude Haiku 4.5</button> -->
                            </div>
                            <button class="text-gray-500 hover:text-gray-300"><i class="ye ye-paper-plane"></i></button>
                        </div>
                    </div>
        
                    <div class="flex items-center justify-between text-xs text-gray-500 px-1">
                        <div class="flex items-center gap-1">
                            <i class="ye ye-key ye-sm"></i> Default Approvals
                        </div>
                        <div class="flex items-center gap-4">
                            <button class="flex items-center gap-1 hover:text-gray-300"><i class="ye ye-folder-open ye-sm"></i> Folder</button>
                            <button class="flex items-center gap-1 hover:text-gray-300"><i class="ye ye-leaf ye-sm"></i> main</button>
                        </div>
                    </div>
                </div>

                <!-- <input v-model="userInput" type="text" :placeholder="`Ask about ${store.active?.activeFile?.name}...`"
                    class="w-full rounded-lg border border-slate-800 bg-slate-950 py-2 pl-3 pr-8 text-xs text-slate-100 placeholder-slate-600 outline-none transition focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500" />
                <button type="submit"
                    class="absolute right-1.5 rounded-md p-1 text-slate-500 transition hover:text-slate-200">
                    <svg class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round"
                            d="M6 12 3.269 3.125A59.769 59.769 0 0 1 21.485 12 59.768 59.768 0 0 1 3.27 20.875L5.999 12Zm0 0h7.5" />
                    </svg>
                </button> -->
            </form>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useWorkspace } from '@/store/workspace'
// import { GenerateProject } from '@wails/go/app/App'
import { EventsOn, EventsOff } from '@wails/runtime/runtime'
import { AgentChat } from '@wails/go/main/App'
import type { agent } from '@wails/go/models'

interface ScaffoldProgress {
    file: string
    index: number
    total: number
    done: boolean
    error?: string
}

const store = useWorkspace();
const chatHistoryContainer = ref<HTMLElement | null>(null);
const userInput = ref('');
const isAiThinking = ref(false);
const newSession = ref(false);
const currentSessionID = ref<string>("session-" + Math.random().toString(36).substring(2, 9));
const selectedProvider = ref<string>("gemini");

const messages = ref<any[]>([
    {
        role: 'assistant',
        content: 'Welcome back! I have context access to your workspace structure. Pick any file from the explorer to begin refactoring.',
    },
])

// ── Scaffold progress state ──────────────────────────────────────────────────

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

function resetScaffold() {
    scaffoldProgress.value = { active: false, file: '', index: 0, total: 0, files: [] }
}

const appendOrUpdateAssistant = (text: string) => {
    const lastMsg = messages.value[messages.value.length - 1];
    if (lastMsg && lastMsg.role === 'assistant') {
        lastMsg.content += text;
    } else {
        messages.value.push({ role: 'assistant', content: text });
    }
}

// ── Wails event listener ─────────────────────────────────────────────────────

onMounted(() => {
    EventsOn('scaffold:progress', (p: ScaffoldProgress) => {
        if (p.error) {
            isAiThinking.value = false
            messages.value.push({
                role: 'assistant',
                content: `❌ Scaffold failed:\n${p.error}`,
            })
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

        // Intermediate progress
        scaffoldProgress.value.active = true
        scaffoldProgress.value.file = p.file
        scaffoldProgress.value.index = p.index
        scaffoldProgress.value.total = p.total

        if (p.file && p.file !== 'Contacting LLM…' && p.file !== 'Parsing plan…') {
            scaffoldProgress.value.files.push(p.file)
        }

        scrollToBottom()
    })

    // In AiAssistant.vue — onMounted
    EventsOn('agent:message', (text: string) => {
        // append or stream into the last assistant bubble
        appendOrUpdateAssistant(text)
        scrollToBottom()
    })

    EventsOn('agent:tool', (call: { name: string; input: string }) => {
        // show a tool-call chip in the chat
        messages.value.push({
            role: 'tool',
            content: `🔧 ${call.name}`,
            // input: call.input,
        })
        scrollToBottom()
    })

    EventsOn('agent:done', () => {
        isAiThinking.value = false
        store.fetchSessions();
    })

    EventsOn('agent:error', (msg: string) => {
        isAiThinking.value = false
        messages.value.push({ role: 'assistant', content: `❌ ${msg}` })
    })
})

onUnmounted(() => {
    EventsOff('scaffold:progress')
    EventsOff('agent:message')
    EventsOff('agent:tool')
    EventsOff('agent:done')
    EventsOff('agent:error')
})

// ── Chat helpers ─────────────────────────────────────────────────────────────

const scrollToBottom = async () => {
    await nextTick()
    if (chatHistoryContainer.value) {
        chatHistoryContainer.value.scrollTop = chatHistoryContainer.value.scrollHeight
    }
}

const sendMessage = async () => {
    if (!userInput.value.trim()) return

    const text = userInput.value.trim()
    userInput.value = ''
    messages.value.push({ role: 'user', content: text })
    scrollToBottom()

    isAiThinking.value = true
    newSession.value = false;
    resetScaffold()

    try {
        await AgentChat(currentSessionID.value, text, selectedProvider.value)
    } catch (err: any) {
        isAiThinking.value = false
        messages.value.push({
            role: 'assistant',
            content: `❌ Error: ${err?.message ?? String(err)}`,
        })
        scrollToBottom()
    }
}

const triggerAiAction = (actionType: 'explain' | 'optimize') => {
    isAiThinking.value = true

    if (actionType === 'explain') {
        messages.value.push({ role: 'user', content: `Explain ${store.active!.activeFile!.name}` })
        scrollToBottom()
        setTimeout(() => {
            isAiThinking.value = false
            messages.value.push({
                role: 'assistant',
                content: `Analyzing \`${store.active!.activeFile!.name}\`...\n\nThis setup provides clean initialization patterns for structural layouts tagged under language group [${store.active!.activeFile!.lang}]. Code bounds look healthy.`,
            })
            scrollToBottom()
        }, 900)
    } else {
        messages.value.push({ role: 'user', content: `Refactor ${store.active!.activeFile!.name}` })
        scrollToBottom()
        setTimeout(() => {
            isAiThinking.value = false
            messages.value.push({
                role: 'assistant',
                content: `Code clean-up complete for \`${store.active!.activeFile!.name}\`. Structural spacing and configuration layouts adjusted according to industry standard guidelines.`,
            })
            scrollToBottom()
        }, 900)
    }
}
</script>