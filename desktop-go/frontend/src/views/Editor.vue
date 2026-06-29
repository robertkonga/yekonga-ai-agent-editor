<template>
    <div class="flex h-full w-full bg-slate-900 text-slate-200 leading-1"
        :class="{ 'select-none': isDragging }"
        @mouseup="onMouseUp"
        @mousemove="onMouseMove"
        @mouseleave="onMouseUp">

        <div class="h-full bg-slate-900/60 overflow-y-auto"
            :style="{ width: leftWidth + 'px' }">
            <FileExplorer :isEditor="true"></FileExplorer>
        </div>
        
        <!-- Left Divider -->
        <div v-if="store.active" class="relative flex items-center justify-center shrink-0 w-[7px] cursor-col-resize group"
            @mousedown.prevent="startDrag('left', $event)">
            <div class="h-full w-[1px] bg-gray-600/20 group-hover:bg-blue-500/50 group-hover:w-[3px] transition-all duration-150"></div>
            <div class="absolute inset-y-0 left-1/2 -translate-x-1/2 w-[7px] group-hover:bg-blue-500/10 rounded-full transition-colors duration-150"></div>
        </div>

        <!-- IDE Panel (Right pane) -->
        <div class="flex flex-1 flex-col h-full">
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
        <div v-if="store.active" class="relative flex items-center justify-center shrink-0 w-[7px] cursor-col-resize group"
            @mousedown.prevent="startDrag('right', $event)">
            <div class="h-full w-[1px] bg-gray-600/20 group-hover:bg-blue-500/50 group-hover:w-[3px] transition-all duration-150"></div>
            <div class="absolute inset-y-0 left-1/2 -translate-x-1/2 w-[7px] group-hover:bg-blue-500/10 rounded-full transition-colors duration-150"></div>
        </div>

        <!-- Chat Panel (Left pane) -->
        <div v-if="store.active" class="h-full bg-slate-800 flex flex-col"
            :style="{ width: rightWidth + 'px' }">
            <ChatPanel :sessionId="store.active?.sessionId!" :small="true" />
        </div>
    </div>
</template>

<script setup lang="ts">
import { provide, ref } from 'vue';
import { useWorkspace, type FileNode } from '@/store/workspace.ts';

import Terminal from './components/Editor/Terminal.vue';
import CodeEditor from './components/Editor/CodeEditor.vue';
import FileExplorer from './components/General/FileExplorer.vue';
import ChatPanel from './components/Chat/ChatPanel.vue';

import { getFileColorClass } from './components/Editor/utils.ts';

const store = useWorkspace();

// Dynamic File Switching Handler
const targetPath = ref<FileNode | null>(null)


const minLeft = 180
const maxLeft = 500
const minRight = 200
const maxRight = 600

const leftWidth = ref(256)
const rightWidth = ref(288)
const isDragging = ref(false)

type DragSide = 'left' | 'right' | null
const dragSide = ref<DragSide>(null)
const startX = ref(0)
const startLeft = ref(0)
const startRight = ref(0)

function startDrag(side: DragSide, e: MouseEvent) {
    isDragging.value = true
    dragSide.value = side
    startX.value = e.clientX
    startLeft.value = leftWidth.value
    startRight.value = rightWidth.value
}

function onMouseMove(e: MouseEvent) {
    if (!isDragging.value || !dragSide.value) return

    const dx = e.clientX - startX.value

    if (dragSide.value === 'left') {
        const newWidth = Math.min(maxLeft, Math.max(minLeft, startLeft.value + dx))
        leftWidth.value = newWidth
    } else if (dragSide.value === 'right') {
        const newWidth = Math.min(maxRight, Math.max(minRight, startRight.value - dx))
        rightWidth.value = newWidth
    }

    window.dispatchEvent(new Event("resize"));
}

function onMouseUp() {
    if (isDragging.value) {
        isDragging.value = false
        dragSide.value = null
    }
}

provide("targetPath", targetPath);
provide("getFileColorClass", getFileColorClass);
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