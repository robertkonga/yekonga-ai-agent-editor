<template>
    <!-- [app-region:drag] is a Tailwind arbitrary property for Electron drag -->
    <div class="flex h-[30px] flex-shrink-0 select-none border-b border-slate-800/60 items-center justify-between font-sans text-xs [--wails-draggable:drag]"
        :class="isMaximized ? 'bg-slate-900' : 'bg-slate-900'" >
        <!-- Drag region + title -->
        <div class="flex min-w-0 flex-1 items-center overflow-hidden pl-2 [app-region:drag]">
            <span v-if="showTitle" class="truncate font-normal text-[#cccccc]/85">
                {{ title }}
            </span>
        </div>

        <!-- Window controls — [app-region:no-drag] so buttons stay clickable in Electron -->
        <div class="flex h-full items-stretch [app-region:no-drag]" role="toolbar" aria-label="Window controls">
            <!-- Minimize -->
            <button type="button" title="Minimize" aria-label="Minimize window"
                class="flex h-full w-[46px] cursor-pointer items-center justify-center border-none bg-transparent p-0 text-[#cccccc] outline-none transition-colors duration-[120ms] hover:bg-white/10 focus-visible:outline focus-visible:outline-1 focus-visible:outline-white/40 focus-visible:-outline-offset-2"
                @click="handleMinimize">
                <svg width="10" height="1" viewBox="0 0 10 1" aria-hidden="true">
                    <rect width="10" height="1" fill="currentColor" />
                </svg>
            </button>

            <!-- Restore / Maximize -->
            <button type="button" :title="isMaximized ? 'Restore' : 'Maximize'"
                :aria-label="isMaximized ? 'Restore window' : 'Maximize window'"
                class="flex h-full w-[46px] cursor-pointer items-center justify-center border-none bg-transparent p-0 text-[#cccccc] outline-none transition-colors duration-[120ms] hover:bg-white/10 focus-visible:outline focus-visible:outline-1 focus-visible:outline-white/40 focus-visible:-outline-offset-2"
                @click="handleRestore">
                <!-- Restore: two overlapping squares -->
                <svg v-if="isMaximized" width="10" height="10" viewBox="0 0 10 10" aria-hidden="true">
                    <path d="M3 0H10V7H8V2H3V0Z" fill="currentColor" />
                    <rect x="0" y="3" width="7" height="7" fill="none" stroke="currentColor" stroke-width="1" />
                </svg>
                <!-- Maximize: single square -->
                <svg v-else width="10" height="10" viewBox="0 0 10 10" aria-hidden="true">
                    <rect x="0.5" y="0.5" width="9" height="9" fill="none" stroke="currentColor" stroke-width="1" />
                </svg>
            </button>

            <!-- Close -->
            <button type="button" title="Close" aria-label="Close window"
                class="flex h-full w-[46px] cursor-pointer items-center justify-center border-none bg-transparent p-0 text-[#cccccc] outline-none transition-colors duration-[120ms] hover:bg-[#e81123] hover:text-white focus-visible:outline focus-visible:outline-1 focus-visible:outline-white/40 focus-visible:-outline-offset-2"
                @click="handleClose">
                <svg width="10" height="10" viewBox="0 0 10 10" aria-hidden="true">
                    <line x1="0" y1="0" x2="10" y2="10" stroke="currentColor" stroke-width="1.2"
                        stroke-linecap="round" />
                    <line x1="10" y1="0" x2="0" y2="10" stroke="currentColor" stroke-width="1.2"
                        stroke-linecap="round" />
                </svg>
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { Quit, WindowMinimise, WindowToggleMaximise } from '@wails/runtime/runtime';

interface Props {
    title?: string
    isMaximized?: boolean
    showTitle?: boolean
}

const props = withDefaults(defineProps<Props>(), {
    title: 'Yekonga Editor',
    isMaximized: false,
    showTitle: true,
})

 
function handleMinimize() {
    WindowMinimise()
}

function handleRestore() {
    WindowToggleMaximise()
}

function handleClose() {
    Quit()
}
</script>
