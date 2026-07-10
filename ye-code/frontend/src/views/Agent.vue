<template>
    <div class="flex h-full gap-0 w-full bg-slate-900 text-gray-300 font-sans overflow-hidden leading-1"
        :class="{ 'select-none': store.activeSideView.isDragging }"
        @mouseup="store.onMouseUp"
        @mousemove="store.onMouseMove"
        @mouseleave="store.onMouseUp">

        <!-- Left Sidebar -->
        <aside class="flex flex-col shrink-0 overflow-hidden"
            :style="{ width: store.activeSideView.leftWidth + 'px' }">
            <LeftSidebar />
        </aside>

        <!-- Left Divider -->
        <div class="relative flex items-center justify-center shrink-0 w-3 cursor-col-resize group"
            @mousedown.prevent="store.startDrag('left', $event)">
            <div class="h-full w-px group-hover:bg-blue-500/50 group-hover:w-0.75 transition-all duration-150"></div>
            <div class="absolute inset-y-0 left-1/2 -translate-x-1/2 w-3 group-hover:bg-blue-500/10 rounded-full transition-colors duration-150"></div>
        </div>

        <!-- Main Content -->
        <main class="flex-1 flex flex-col min-w-0 pt-3">
            <MainWorkspace />
        </main>

        <!-- Right Divider -->
        <div class="relative flex items-center justify-center shrink-0 w-3 cursor-col-resize group"
            @mousedown.prevent="store.startDrag('right', $event)">
            <div class="h-full w-px group-hover:bg-blue-500/50 group-hover:w-0.75 transition-all duration-150"></div>
            <div class="absolute inset-y-0 left-1/2 -translate-x-1/2 w-3 group-hover:bg-blue-500/10 rounded-full transition-colors duration-150"></div>
        </div>

        <!-- Right Sidebar -->
        <aside class="flex flex-col shrink-0 overflow-hidden pt-3 pe-3"
            :style="{ width: store.activeSideView.rightWidth + 'px' }">
            <RightSidebar />
        </aside>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import LeftSidebar from './components/Agent/LeftSidebar.vue'
import MainWorkspace from './components/Agent/MainWorkspace.vue'
import RightSidebar from './components/Agent/RightSidebar.vue'
import { useWorkspace } from '@/store/workspace.ts';

const store = useWorkspace();
</script>
