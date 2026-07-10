<template>
    <div class="flex h-full w-full bg-slate-900 text-slate-200 leading-1"
        :class="{ 'select-none': store.activeSideView.isDragging }"
        @mouseup="store.onMouseUp"
        @mousemove="store.onMouseMove"
        @mouseleave="store.onMouseUp">

        <div class="h-full bg-slate-900/60 overflow-y-auto"
            :style="{ width: store.activeSideView.leftWidth + 'px' }">
            <GeneralExplorer :isEditor="true"></GeneralExplorer>
        </div>
        
        <!-- Left Divider -->
        <div v-if="store.active" class="relative flex items-center justify-center shrink-0 w-1.75 cursor-col-resize group"
            @mousedown.prevent="store.startDrag('left', $event)">
            <div class="h-full w-px  group-hover:bg-blue-500/50 group-hover:w-0.75 transition-all duration-150"></div>
            <div class="absolute inset-y-0 left-1/2 -translate-x-1/2 w-1.75 group-hover:bg-blue-500/10 rounded-full transition-colors duration-150"></div>
        </div>

        <!-- IDE Panel (Right pane) -->
        <div class="flex flex-1 flex-col h-full border-l border-slate-800/60">
            <div class="flex-1 h-full flex ">
                <template v-if="store.active">
                    <CodeEditor></CodeEditor>
                </template>
                <template>
                    <div class="flex items-center justify-center">Select Workspace</div>
                </template>
            </div>
            <div v-if="false" class="">
                <Terminal></Terminal>
            </div>
        </div>

        <!-- Right Divider -->
        <div v-if="store.active" class="relative flex items-center justify-center shrink-0 w-1.75 cursor-col-resize group"
            @mousedown.prevent="store.startDrag('right', $event)">
            <div class="h-full w-px bg-transparent group-hover:bg-blue-500/50 group-hover:w-0.75 transition-all duration-150"></div>
            <div class="absolute inset-y-0 left-1/2 -translate-x-1/2 w-1.75 group-hover:bg-blue-500/10 rounded-full transition-colors duration-150"></div>
        </div>

        <!-- Chat Panel (Left pane) -->
        <div v-if="store.active" class="h-full bg-slate-800 flex flex-col"
            :style="{ width: store.activeSideView.rightWidth + 'px' }">
            <ChatPanel :sessionId="store.active.sessionId!" :small="true" />
        </div>
    </div>
</template>

<script setup lang="ts">
import { provide, ref } from 'vue';
import { useWorkspace, type FileNode } from '@/store/workspace.ts';

import Terminal from './components/Editor/Terminal.vue';
import CodeEditor from './components/Editor/CodeEditor.vue';
import ChatPanel from './components/Chat/ChatPanel.vue';
import GeneralExplorer from './components/General/GeneralExplorer.vue';

const store = useWorkspace();
// Dynamic File Switching Handler
const targetPath = ref<FileNode | null>(null)

provide("targetPath", targetPath);
</script>


<style border>
/* Standard scrollbar formatting */
::-webkit-scrollbar {
    width: 4px;
    height: 4px;
}

::-webkit-scrollbar-track {
    background: transparent;
}

::-webkit-scrollbar-thumb {
    background: #1e293b;
    border-radius: 9999px;
}

::-webkit-scrollbar-thumb:hover {
    background: #334155;
}
</style>