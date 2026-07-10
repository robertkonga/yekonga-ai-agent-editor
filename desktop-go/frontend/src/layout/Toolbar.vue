<template>
    <!-- [app-region:drag] is a Tailwind arbitrary property for Electron drag -->
    <div class="flex h-8 border-b select-none items-center justify-between font-sans text-xs [--wails-draggable:drag]"
        :class="[
            isMaximized ? 'bg-slate-900' : 'bg-slate-900',
            {
                'border-slate-700/20': (store.active && store.active.viewMode == 'EDITOR'),
                'border-transparent': (!store.active || (store.active && store.active.viewMode != 'EDITOR')),
            }
        ]" >
        <!-- Drag region + title -->
        <div class="flex h-full min-w-0 flex-1 items-center overflow-hidden pl-2 [app-region:drag] space-x-3">
            <img src="@/assets/appicon.png" class="h-4 w-auto"/>
            <span v-if="showTitle" class="truncate font-normal text-[#cccccc]/85">
                {{ title }}
            </span>
            <template v-if="store.active && store.active.viewMode != 'WORKSPACE'">
                <button type="button" title="Minimize" aria-label="Minimize window"
                    class="px-4 flex font-semibold h-full w-auto cursor-pointer items-center justify-center border-none bg-transparent p-0 text-blue-500 outline-none transition-colors duration-120 hover:bg-white/10 focus-visible:outline  focus-visible:outline-white/40 focus-visible:-outline-offset-2"
                    @click="store.setViewMode(null)">
                    Workspaces
                </button>
            </template>
        </div>

        <!-- Window controls — [app-region:no-drag] so buttons stay clickable in Electron -->
        <div class="flex h-full items-stretch [app-region:no-drag]" role="toolbar" aria-label="Window controls">
            <template v-if="store.active">
                <template v-if="store.active && store.active.viewMode == 'AGENT'">
                    <button type="button" title="Minimize" aria-label="Minimize window"
                        class="px-4 flex font-semibold h-full w-auto cursor-pointer items-center justify-center border-none bg-transparent p-0 text-[#cccccc] outline-none transition-colors duration-120 hover:bg-white/10 focus-visible:outline  focus-visible:outline-white/40 focus-visible:-outline-offset-2"
                        @click="store.setViewMode('EDITOR')">
                        Editor
                    </button>
                </template>
                <template v-else-if="store.active && store.active.viewMode == 'EDITOR'">
                    <button type="button" title="Minimize" aria-label="Minimize window"
                        class="px-4 flex font-semibold h-full w-auto cursor-pointer items-center justify-center border-none bg-indigo-900 p-0 text-[#cccccc] outline-none transition-colors duration-120 hover:bg-white/10 focus-visible:outline  focus-visible:outline-white/40 focus-visible:-outline-offset-2"
                        @click="store.setViewMode('AGENT')">
                        <!-- <i class="ye ye-microchip-ai ye-lg"></i> -->
                        <span>Agent</span>
                    </button>
                </template>
            </template>
            <button @click="showSetting = !showSetting" type="button" title="Minimize" aria-label="Minimize window"
                class="px-3 flex font-semibold h-full w-auto cursor-pointer items-center justify-center border-none bg-transparent p-0 text-[#cccccc] outline-none transition-colors duration-120 hover:bg-white/10 focus-visible:outline  focus-visible:outline-white/40 focus-visible:-outline-offset-2"
                >
                <i class="ye ye-gear ye-lg"></i>
            </button>
            <div class="border-r border-slate-800/60"></div>
            <!-- Minimize -->
            <button type="button" title="Minimize" aria-label="Minimize window"
                class="flex h-full w-11.5 cursor-pointer items-center justify-center border-none bg-transparent p-0 text-[#cccccc] outline-none transition-colors duration-120 hover:bg-white/10 focus-visible:outline  focus-visible:outline-white/40 focus-visible:-outline-offset-2"
                @click="handleMinimize">
                <svg width="10" height="1" viewBox="0 0 10 1" aria-hidden="true">
                    <rect width="10" height="1" fill="currentColor" />
                </svg>
            </button>

            <!-- Restore / Maximize -->
            <button type="button" :title="isMaximized ? 'Restore' : 'Maximize'"
                :aria-label="isMaximized ? 'Restore window' : 'Maximize window'"
                class="flex h-full w-11.5 cursor-pointer items-center justify-center border-none bg-transparent p-0 text-[#cccccc] outline-none transition-colors duration-120 hover:bg-white/10 focus-visible:outline  focus-visible:outline-white/40 focus-visible:-outline-offset-2"
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
                class="flex h-full w-11.5 cursor-pointer items-center justify-center border-none bg-transparent p-0 text-[#cccccc] outline-none transition-colors duration-120 hover:bg-[#e81123] hover:text-white focus-visible:outline  focus-visible:outline-white/40 focus-visible:-outline-offset-2"
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

<Modal title="Setting" v-model="showSetting">
    <Settings></Settings>
</Modal>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import Modal from '@/components/Modal.vue';
import { useWorkspace } from '@/store/workspace';
import Settings from '@/views/Settings.vue';
import { Quit, WindowMinimise, WindowToggleMaximise } from '@wails/runtime/runtime';

interface Props {
    title?: string
    isMaximized?: boolean
    showTitle?: boolean
}

const props = withDefaults(defineProps<Props>(), {
    title: 'YE-CODE',
    isMaximized: false,
    showTitle: true,
})

const store = useWorkspace();
const showSetting = ref<boolean>(false);
 
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
