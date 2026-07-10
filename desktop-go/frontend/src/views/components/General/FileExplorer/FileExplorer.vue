<template>
    <div class="h-full flex flex-col select-none">
        <div class="flex h-9 items-center justify-between pl-3 pr-1 border-b border-slate-800/60 bg-slate-900/40">
            <span class="text-xs font-bold uppercase tracking-wider text-slate-400">Explorer</span>
            <div>
                <span @click="store.setSimpleIcon(!store.isSimpleIcon)" class="size-5 flex items-center justify-center cursor-pointer group">
                    <i class="ye ye-face-smile group-hover:text-white"
                        :class="!store.isSimpleIcon? 'text-blue-500': 'text-slate-500'"></i>
                </span>
            </div>
            <!-- <svg class="h-4 w-4 text-slate-500" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 9h16.5m-16.5 6.75h16.5" />
            </svg> -->
        </div>

        <div class="flex flex-col flex-1 pb-2 overflow-hidden">
            <template v-if="store.active">
                <div class="flex items-center justify-between text-sm pl-3 pr-1 py-1">
                    <span class="block truncate w-full" :title="store.active.path">{{ store.active.name }}</span>
                    <div class="flex gap-2">
                        <span @click="createFileAtRoot" class="size-5 flex items-center justify-center cursor-pointer group">
                            <i class="ye ye-file-plus text-slate-500 group-hover:text-white"></i>
                        </span>
                        <span @click="createFolderAtRoot" class="size-5 flex items-center justify-center cursor-pointer group">
                            <i class="ye ye-folder-plus text-slate-500 group-hover:text-white"></i>
                        </span>
                        <span @click="store.fetchWorkspaceFiles()" class="size-5 flex items-center justify-center cursor-pointer group">
                            <i class="ye ye-rotate text-slate-500 group-hover:text-white"></i>
                        </span>
                    </div>
                </div>
                <div class="flex-1 overflow-y-auto">
                    <div @dragover.prevent @drop="onDropRoot" v-for="(node, i) of store.active!.workspaceFiles" :key="node.id">
                        <FileExplorerItem v-model="store.active!.workspaceFiles[i]" :level="0" :isEditor="props.isEditor"></FileExplorerItem>
                    </div>
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

const props = defineProps<{
    isEditor: boolean;
}>()

const store = useWorkspace();
// Virtual directory setup tracking structural file system tree

const openWorkspace = async () => {
    const targetPath = await OpenWorkspaceDialog();
    store.openWorkshop(targetPath);
}

const createFileAtRoot = () => {
    if (!store.active) return;
    const name = window.prompt("New File Name:");
    if (name) store.createNewFile(store.active.path, name);
}

const createFolderAtRoot = () => {
    if (!store.active) return;
    const name = window.prompt("New Folder Name:");
    if (name) store.createNewFolder(store.active.path, name);
}

const onDropRoot = (e: DragEvent) => {
    if (!store.active) return;
    const sourcePath = e.dataTransfer?.getData("text/plain");
    if (sourcePath && sourcePath !== store.active.path) {
        store.moveFile(sourcePath, store.active.path);
    }
}

if(!store.active) {
    // openWorkspace();
}

onMounted(async () => {
    store.fetchWorkspaceFiles()
})

</script>
