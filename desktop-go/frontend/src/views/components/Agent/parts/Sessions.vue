<template>
    <expandable title="Sessions" :show="true">
        <template #extraHeader>
            <div class="flex items-center gap-2">
                <button @click="createNewSession" class="text-xs bg-slate-800 hover:bg-slate-700 px-2 py-1 rounded text-gray-400">
                    New <span class="text-[10px] ml-1 opacity-70">Ctrl+N</span>
                </button>
                <button @click="store.fetchSessions()" class="text-xs bg-slate-800 hover:bg-slate-700 px-2 py-1 rounded text-gray-400 hover:text-white">
                    <i class="ye ye-rotate ye-xs"></i>
                </button>
            </div>
        </template>
        <template #default>
            <div class="min-h-40 max-h-80 overflow-y-auto flex flex-col gap-2 py-2">
                <div v-if="!store.active?.sessions.length" class="text-xs text-gray-500 py-4 text-center">
                    No sessions found
                </div>
                <div v-for="session in store.active?.sessions" :key="session.id" 
                     @click="selectSession(session)"
                     class="group cursor-pointer p-2 hover:bg-slate-800 rounded transition-colors"
                     :class="{'bg-slate-800 border-l-2 border-blue-500': activeSessionId === session.id}">
                    <div class="flex items-center gap-2 text-xs text-gray-300">
                        <span class="w-1.5 h-1.5 rounded-full" :class="activeSessionId === session.id ? 'bg-blue-500' : 'bg-gray-500'"></span>
                        <div class="truncate font-medium">{{ session.id }}</div>
                    </div>
                    <div class="text-[10px] text-gray-500 ml-3 mt-1 flex items-center justify-between">
                        <div class="flex items-center gap-1 uppercase">
                            <i class="ye ye-robot"></i> {{ session.provider }}
                        </div>
                        <div>{{ formatDate(session.last_updated) }}</div>
                    </div>
                </div>  
            </div>
        </template>
    </expandable>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useWorkspace } from '@/store/workspace';
import Expandable from '@/components/Expandable.vue';

const store = useWorkspace();
const activeSessionId = ref<string | null>(null);

const createNewSession = () => {
    const id = "session-" + Math.random().toString(36).substring(2, 9);
    activeSessionId.value = id;
    // Notify parent or update state for new chat
}

const selectSession = (session: any) => {
    activeSessionId.value = session.id;
    // Load session history
}

const formatDate = (dateStr: string) => {
    if (!dateStr) return '';
    const date = new Date(dateStr);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

onMounted(() => {
    store.fetchSessions();
})
</script>