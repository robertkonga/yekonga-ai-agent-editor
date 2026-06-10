<template>
    <div class="h-full flex w-full flex-col bg-slate-900/40 backdrop-blur-xl">
        <div class="flex h-9 items-center space-x-2 border-b border-slate-800/60 px-4 bg-slate-900/60">
            <div
                class="flex h-4 w-4 items-center justify-center rounded-md bg-indigo-600/20 text-indigo-400 border border-indigo-500/30">
                <svg class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round"
                        d="M9.813 15.904 9 21l8.982-11.795H14.18l.813-5.109L6 15.904h3.813Z" />
                </svg>
            </div>
            <span class="text-xs font-bold tracking-wider text-slate-200 uppercase">Context Assistant</span>
        </div>

        <div class="grid grid-cols-2 gap-2 p-3 border-b border-slate-800/60 bg-slate-900/20">
            <button @click="triggerAiAction('explain')"
                class="flex flex-col items-start rounded-lg border border-slate-800 bg-slate-900/60 p-2 text-left transition hover:border-slate-700 hover:bg-slate-850">
                <span class="text-xs font-semibold text-slate-200">Explain File</span>
                <span class="text-[10px] text-slate-500 mt-0.5">Deconstruct current tab</span>
            </button>
            <button @click="triggerAiAction('optimize')"
                class="flex flex-col items-start rounded-lg border border-slate-800 bg-slate-900/60 p-2 text-left transition hover:border-slate-700 hover:bg-slate-850">
                <span class="text-xs font-semibold text-indigo-400">Refactor Code</span>
                <span class="text-[10px] text-slate-500 mt-0.5">Optimize formatting</span>
            </button>
        </div>

        <div class="flex-1 overflow-y-auto p-4 space-y-4" ref="chatHistoryContainer">
            <div v-for="(msg, idx) in messages" :key="idx" :class="['flex flex-col max-w-[90%] rounded-xl p-3 text-xs leading-relaxed border',
                msg.role === 'user' ? 'ml-auto bg-slate-800 border-slate-700 text-slate-100' : 'mr-auto bg-slate-900/80 border-slate-800/60 text-slate-300'
            ]">
                <span class="font-bold text-[9px] uppercase tracking-wider mb-1 block"
                    :class="msg.role === 'user' ? 'text-indigo-300' : 'text-emerald-400'">
                    {{ msg.role === 'user' ? 'You' : 'Agent' }}
                </span>
                <p class="whitespace-pre-wrap text-[11px]">{{ msg.content }}</p>
            </div>

            <div v-if="isAiThinking"
                class="flex items-center space-x-1.5 mr-auto bg-slate-900/80 border border-slate-800/60 rounded-xl p-3 text-xs text-slate-500">
                <span class="h-1 w-1 animate-bounce rounded-full bg-slate-400" style="animation-delay: 0ms"></span>
                <span class="h-1 w-1 animate-bounce rounded-full bg-slate-400" style="animation-delay: 150ms"></span>
                <span class="h-1 w-1 animate-bounce rounded-full bg-slate-400" style="animation-delay: 300ms"></span>
            </div>
        </div>

        <div class="p-3 border-t border-slate-800/60 bg-slate-900/40">
            <form @submit.prevent="sendMessage" class="relative flex items-center">
                <input v-model="userInput" type="text" :placeholder="`Ask about ${store.active?.activeFile?.name}...`"
                    class="w-full rounded-lg border border-slate-800 bg-slate-950 py-2 pl-3 pr-8 text-xs text-slate-100 placeholder-slate-600 outline-none transition focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500" />
                <button type="submit"
                    class="absolute right-1.5 rounded-md p-1 text-slate-500 transition hover:text-slate-200">
                    <svg class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round"
                            d="M6 12 3.269 3.125A59.769 59.769 0 0 1 21.485 12 59.768 59.768 0 0 1 3.27 20.875L5.999 12Zm0 0h7.5" />
                    </svg>
                </button>
            </form>
        </div>
    </div>
</template>


<script setup lang="ts">
import { ref, onMounted, nextTick, inject, type Ref } from 'vue';
import { useWorkspace } from '@/store/workspace';

const store = useWorkspace();
const chatHistoryContainer = ref<HTMLElement | null>(null);
const userInput = ref('');
const isAiThinking = ref(false);

const messages = ref([
    { role: 'assistant', content: "Welcome back! I have context access to your workspace structure. Pick any file from the explorer to begin refactoring." }
]);

const scrollToBottom = async () => {
    await nextTick();
    if (chatHistoryContainer.value) {
        chatHistoryContainer.value!.scrollTop = chatHistoryContainer.value!.scrollHeight;
    }
};

const sendMessage = () => {
    if (!userInput.value.trim()) return;

    messages.value.push({ role: 'user', content: userInput.value });
    const textPayload = userInput.value;
    userInput.value = '';
    scrollToBottom();

    isAiThinking.value = true;
    setTimeout(() => {
        isAiThinking.value = false;

        messages.value.push({
            role: 'assistant',
            content: `Inspected instructions for "${store.active!.activeFile!.name}". \n\nQuery parsing matched context: "${textPayload}". Let me know if you would like me to output a script snippet to insert into this file.`
        });
        scrollToBottom();
    }, 1000);
};

const triggerAiAction = (actionType: any) => {
    // if (!editor) return;
    isAiThinking.value = true;

    if (actionType === 'explain') {
        messages.value.push({ role: 'user', content: `Explain ${store.active!.activeFile!.name}` });
        scrollToBottom();

        setTimeout(() => {
            isAiThinking.value = false;
            messages.value.push({
                role: 'assistant',
                content: `Analyzing file text content for \`${store.active!.activeFile!.name}\`... \n\nThis setup provides clean initialization patterns for structural layouts tagged under language group [${store.active!.activeFile!.lang}]. Code bounds look healthy.`
            });
            scrollToBottom();
        }, 900);
    } else if (actionType === 'optimize') {
        messages.value.push({ role: 'user', content: `Refactor ${store.active!.activeFile!.name}` });
        scrollToBottom();

        setTimeout(() => {
            isAiThinking.value = false;
            messages.value.push({
                role: 'assistant',
                content: `Code clean-up complete for \`${store.active!.activeFile!.name}\`. Structural spacing and configuration layouts adjusted according to industry standard guidelines.`
            });
            scrollToBottom();
        }, 900);
    }
};

const runCode = () => {
    alert(`Executing execution layer for file context: ${store.active!.activeFile!.name}`);
};
</script>
