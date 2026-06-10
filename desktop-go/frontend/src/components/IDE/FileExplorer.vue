<template>
    <div class="flex w-60 flex-col border-r border-slate-800/60 bg-slate-900/60 select-none">
        <div class="flex h-9 items-center justify-between px-4 border-b border-slate-800/60 bg-slate-900/40">
            <span class="text-xs font-bold uppercase tracking-wider text-slate-400">Workspace Explorer</span>
            <svg class="h-4 w-4 text-slate-500" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 9h16.5m-16.5 6.75h16.5" />
            </svg>
        </div>

        <div class="flex-1 overflow-y-auto pb-2">
            <template v-if="store.active">
                <div class="flex items-center justify-between bg-slate-700/30 text-sm px-3 py-1">
                    <span>{{ store.active.path.split('/').pop() }}</span>
                    <span @click="store.openWorkshop(null)">
                        <svg class="h-2 w-2 text-slate-500" viewBox="0 0 10 10">
                            <line x1="1" y1="1" x2="9" y2="9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" />
                            <line x1="9" y1="1" x2="1" y2="9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" />
                        </svg>
                    </span>
                </div>
                <div v-for="node in store.active?.workspaceFiles" :key="node.id">
                    <FileExplorerItem :item="node" :level="0"></FileExplorerItem>
                </div>
            </template>
            <template v-else>
                <div @click="openWorkspace" class="flex items-center justify-center h-full">
                    <div class="rounded-md border border-gray-500/30 p-3 ">Select Workspace</div>
                </div>
            </template>
        </div>
    </div>
</template>


<script setup lang="ts">
// Basic file explorer mockup

import { onMounted, ref } from 'vue';
import FileExplorerItem from './FileExplorerItem.vue';
import { OpenWorkspaceDialog } from '@wails/go/main/App';
import { useWorkspace, type FileNode } from '@/store/workspace.ts';

const store = useWorkspace();
// Virtual directory setup tracking structural file system tree

const openWorkspace = async () => {
    const targetPath = await OpenWorkspaceDialog();
    store.openWorkshop(targetPath);
}

if(!store.active) {
    // openWorkspace();
}

onMounted(async () => {
    store.fetchWorkspaceFiles()
})

</script>
