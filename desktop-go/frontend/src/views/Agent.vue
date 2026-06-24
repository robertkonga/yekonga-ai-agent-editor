<template>
    <div class="flex h-full gap-0 w-full bg-slate-900 text-gray-300 font-sans overflow-hidden leading-1"
        :class="{ 'select-none': isDragging }"
        @mouseup="onMouseUp"
        @mousemove="onMouseMove"
        @mouseleave="onMouseUp">

        <!-- Left Sidebar -->
        <aside class="flex flex-col shrink-0 overflow-hidden"
            :style="{ width: leftWidth + 'px' }">
            <LeftSidebar />
        </aside>

        <!-- Left Divider -->
        <div class="relative flex items-center justify-center shrink-0 w-3 cursor-col-resize group"
            @mousedown.prevent="startDrag('left', $event)">
            <div class="h-full w-[1px] group-hover:bg-blue-500/50 group-hover:w-[3px] transition-all duration-150"></div>
            <div class="absolute inset-y-0 left-1/2 -translate-x-1/2 w-3 group-hover:bg-blue-500/10 rounded-full transition-colors duration-150"></div>
        </div>

        <!-- Main Content -->
        <main class="flex-1 flex flex-col min-w-0 pt-3">
            <MainWorkspace />
        </main>

        <!-- Right Divider -->
        <div class="relative flex items-center justify-center shrink-0 w-3 cursor-col-resize group"
            @mousedown.prevent="startDrag('right', $event)">
            <div class="h-full w-[1px] group-hover:bg-blue-500/50 group-hover:w-[3px] transition-all duration-150"></div>
            <div class="absolute inset-y-0 left-1/2 -translate-x-1/2 w-3 group-hover:bg-blue-500/10 rounded-full transition-colors duration-150"></div>
        </div>

        <!-- Right Sidebar -->
        <aside class="flex flex-col shrink-0 overflow-hidden pt-3 pe-3"
            :style="{ width: rightWidth + 'px' }">
            <RightSidebar />
        </aside>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import LeftSidebar from './components/Agent/LeftSidebar.vue'
import MainWorkspace from './components/Agent/MainWorkspace.vue'
import RightSidebar from './components/Agent/RightSidebar.vue'

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
</script>
