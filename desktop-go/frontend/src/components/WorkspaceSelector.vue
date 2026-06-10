
<template>
    <div class="flex h-full min-h-0 w-full flex-col bg-slate-900 font-sans text-[#cccccc]">

        <!-- Header -->
        <div class="flex items-center border-b border-slate-800/60 px-5 py-4">
            <div class="flex items-center justify-between flex-1">
                <div class="flex items-center gap-3">
                    <!-- VS Code-style grid icon -->
                    <svg class="shrink-0 text-[#75beff]" width="18" height="18" viewBox="0 0 18 18" fill="none">
                        <rect x="1" y="1" width="7" height="7" rx="1" fill="currentColor" opacity="0.9" />
                        <rect x="10" y="1" width="7" height="7" rx="1" fill="currentColor" opacity="0.5" />
                        <rect x="1" y="10" width="7" height="7" rx="1" fill="currentColor" opacity="0.5" />
                        <rect x="10" y="10" width="7" height="7" rx="1" fill="currentColor" opacity="0.9" />
                    </svg>
                    <div class="min-w-0">
                        <h1 class="text-[13px] font-semibold leading-none tracking-wide text-[#e8e8e8]">Open Workspace</h1>
                        <p class="mt-1 text-[11px] text-[#858585]">Select a recent workspace or open a folder</p>
                    </div>
                </div>
                <button
                    class="flex items-center gap-1.5 rounded-full bg-blue-500/20 px-3 py-1.5 text-[12px] font-medium text-white transition-colors hover:bg-blue-600/30 focus-visible:outline focus-visible:outline-blue-500 focus-visible:outline-offset-1"
                    @click="openWorkspace">
                    <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
                        <path
                            d="M1 3.5C1 2.67 1.67 2 2.5 2H4.5L6 3.5H9.5C10.33 3.5 11 4.17 11 5V9C11 9.83 10.33 10.5 9.5 10.5H2.5C1.67 10.5 1 9.83 1 9V3.5Z"
                            stroke="currentColor" stroke-width="1.1" />
                    </svg>
                    Open Folder…
                </button>
            </div>
        </div>

        <div class="flex-1 overflow-y-auto">
            <div class="w-1/2 mx-auto ">
    
                <!-- Search -->
                <div class="px-4 py-3">
                    <div
                        class="flex items-center gap-2 rounded-md border border-slate-700 bg-slate-800 px-3 py-1.5 focus-within:border-[#007acc]">
                        <svg class="shrink-0 text-[#858585]" width="12" height="12" viewBox="0 0 12 12" fill="none">
                            <circle cx="5" cy="5" r="3.5" stroke="currentColor" stroke-width="1.2" />
                            <line x1="8" y1="8" x2="11" y2="11" stroke="currentColor" stroke-width="1.2"
                                stroke-linecap="round" />
                        </svg>
                        <input v-model="search" type="text" placeholder="Search workspaces..."
                            class="w-full bg-transparent text-[12px] text-[#cccccc] placeholder-[#555] outline-none" />
                        <button v-if="search" class="text-[#555] hover:text-[#999]" @click="search = ''"
                            aria-label="Clear search">
                            <svg width="10" height="10" viewBox="0 0 10 10">
                                <line x1="1" y1="1" x2="9" y2="9" stroke="currentColor" stroke-width="1.4"
                                    stroke-linecap="round" />
                                <line x1="9" y1="1" x2="1" y2="9" stroke="currentColor" stroke-width="1.4"
                                    stroke-linecap="round" />
                            </svg>
                        </button>
                    </div>
                </div>
    
                <!-- List -->
                <div class="min-h-0 px-2 pb-2">
                    <!-- Empty state -->
                    <div v-if="filtered.length === 0" class="flex flex-col items-center justify-center gap-2 py-10 text-center">
                        <svg class="text-slate-700/60" width="32" height="32" viewBox="0 0 32 32" fill="none">
                            <rect x="4" y="6" width="24" height="20" rx="2" stroke="currentColor" stroke-width="1.5" />
                            <line x1="10" y1="12" x2="22" y2="12" stroke="currentColor" stroke-width="1.5"
                                stroke-linecap="round" />
                            <line x1="10" y1="16" x2="18" y2="16" stroke="currentColor" stroke-width="1.5"
                                stroke-linecap="round" />
                        </svg>
                        <p class="text-[12px] text-[#555]">No workspaces match <span class="text-[#858585]">"{{ search
                                }}"</span></p>
                    </div>
    
                    <template v-else>
                        <!-- Pinned -->
                        <div v-if="pinned.length > 0" class="mb-1">
                            <div class="mb-1 px-2 pt-1 text-[10px] font-semibold uppercase tracking-widest text-[#555]">Pinned
                            </div>
                            <WorkspaceRow v-for="ws in pinned" :key="ws.id" :workspace="ws" :hovered="hoveredId === ws.id"
                                @mouseenter="hoveredId = ws.id" @mouseleave="hoveredId = null" @open="emit('open', ws)"
                                @remove="emit('removeRecent', ws.id)" />
                        </div>
    
                        <!-- Recent -->
                        <div v-if="recent.length > 0">
                            <div class="mb-1 px-2 pt-1 text-[10px] font-semibold uppercase tracking-widest text-[#555]">
                                {{ pinned.length ? 'Recent' : 'Recently Opened' }}
                            </div>
                            <WorkspaceRow v-for="ws in recent" :key="ws.id" :workspace="ws" :hovered="hoveredId === ws.id"
                                @mouseenter="hoveredId = ws.id" @mouseleave="hoveredId = null" @open="emit('open', ws)"
                                @remove="emit('removeRecent', ws.id)" />
                        </div>
                    </template>
                </div>
            </div>
        </div>

        <!-- Footer actions -->
        <div v-if="false" class="flex items-center gap-2 border-t border-slate-800/80 px-4 py-3">
            <button
                class="flex items-center gap-1.5 rounded-md bg-blue-500 px-3 py-1.5 text-[12px] font-medium text-white transition-colors hover:bg-blue-600 focus-visible:outline focus-visible:outline-blue-500 focus-visible:outline-offset-1"
                @click="openWorkspace">
                <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
                    <path
                        d="M1 3.5C1 2.67 1.67 2 2.5 2H4.5L6 3.5H9.5C10.33 3.5 11 4.17 11 5V9C11 9.83 10.33 10.5 9.5 10.5H2.5C1.67 10.5 1 9.83 1 9V3.5Z"
                        stroke="currentColor" stroke-width="1.1" />
                </svg>
                Open Folder…
            </button>
            <span class="text-[10px] text-slate-700/60">or drag a folder here</span>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import WorkspaceRow from './WorkspaceRow.vue'
import { useWorkspace, type Workspace } from '@/store/workspace.ts'
import { OpenWorkspaceDialog } from '@wails/go/main/App.js'

interface Props {
    recentWorkspaces?: Workspace[]
}

const props = withDefaults(defineProps<Props>(), {
    recentWorkspaces: () => [
        // { id: '1', name: 'my-editor', path: '/home/user/projects/my-editor', lastOpened: new Date(Date.now() - 1000 * 60 * 5), isPinned: true },
        // { id: '2', name: 'vscode-ext', path: '/home/user/projects/vscode-ext', lastOpened: new Date(Date.now() - 1000 * 60 * 60 * 2) },
        // { id: '3', name: 'wails-app', path: '/home/user/projects/wails-app', lastOpened: new Date(Date.now() - 1000 * 60 * 60 * 24) },
        // { id: '4', name: 'rust-lsp', path: '/home/user/dev/rust-lsp', lastOpened: new Date(Date.now() - 1000 * 60 * 60 * 24 * 3) },
        // { id: '5', name: 'frontend-kit', path: '/home/user/projects/frontend-kit', lastOpened: new Date(Date.now() - 1000 * 60 * 60 * 24 * 7) },
    ],
})

const emit = defineEmits<{
    open: [workspace: Workspace]
    openFolder: []
    removeRecent: [id: string]
}>()


const store = useWorkspace();
// Virtual directory setup tracking structural file system tree

const recentWorkspaces = computed<Workspace[]>(()=>{
    let list:Workspace[] = [];

    for (const key in store.workspaces) {
        if (!Object.hasOwn(store.workspaces, key)) continue;
        const e = store.workspaces[key];
        
        list.push(e);
    }

    return list;
})

const openWorkspace = async () => {
    const targetPath = await OpenWorkspaceDialog();
    store.openWorkshop(targetPath);
}

const search = ref('')
const hoveredId = ref<string | null>(null)

const filtered = computed(() => {
    const q = search.value.trim().toLowerCase()
    if (!q) return recentWorkspaces.value
    return recentWorkspaces.value.filter(
        (w) => w.name.toLowerCase().includes(q) || w.path.toLowerCase().includes(q),
    )
})

const pinned = computed(() => filtered.value.filter((w) => w.isPinned))
const recent = computed(() => filtered.value.filter((w) => !w.isPinned))

</script>


<!-- ─── Inner row component defined in same file ─────────────────────────── -->
<script lang="ts">
// Re-export Workspace type for consumers
// export type { Workspace }
</script>