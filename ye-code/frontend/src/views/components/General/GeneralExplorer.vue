<template>
<div class="flex flex-col h-full">
    <div class="flex h-8 select-none items-center justify-between font-sans text-xs border-b border-slate-800/60">
        <div class="flex h-full min-w-0 flex-1 items-center overflow-hidden text-gray-400">
            <button type="button" title="Files"
                class="px-4 flex font-semibold h-full w-auto cursor-pointer items-center justify-center border-none bg-transparent p-0 outline-none transition-colors duration-120 hover:bg-white/10 focus-visible:outline focus-visible:outline-white/40 focus-visible:-outline-offset-2"
                :class="{ 'text-white bg-white/5': store.active?.exportView === 'FILES' }"
                @click="store.setViewExplore('FILES')">
                <i class="ye ye-files ye-lg"></i>
            </button>
            <button type="button" title="Search"
                class="px-4 flex font-semibold h-full w-auto cursor-pointer items-center justify-center border-none bg-transparent p-0 outline-none transition-colors duration-120 hover:bg-white/10 focus-visible:outline focus-visible:outline-white/40 focus-visible:-outline-offset-2"
                :class="{ 'text-white bg-white/5': store.active?.exportView === 'SEARCH' }"
                @click="store.setViewExplore('SEARCH')">
                <i class="ye ye-magnifying-glass ye-lg"></i>
            </button>
            <button type="button" title="Git History"
                class="px-4 flex font-semibold h-full w-auto cursor-pointer items-center justify-center border-none bg-transparent p-0 outline-none transition-colors duration-120 hover:bg-white/10 focus-visible:outline focus-visible:outline-white/40 focus-visible:-outline-offset-2"
                :class="{ 'text-white bg-white/5': store.active?.exportView === 'GIT' }"
                @click="store.setViewExplore('GIT')">
                <i class="ye ye-git-alt ye-lg"></i>
            </button>
        </div>
    </div>
    <div class="flex-1 overflow-hidden">
        <!-- SEARCH VIEW -->
        <SearchExplorer v-if="store.active?.exportView == 'SEARCH'" />

        <!-- GIT VIEW -->
        <GitExplorer v-else-if="store.active?.exportView == 'GIT'" />

        <!-- FILES VIEW -->
        <FileExplorer v-else :isEditor="true" />
    </div>
</div>
</template>

<script setup lang="ts">
import { useWorkspace } from '@/store/workspace';
import FileExplorer from './FileExplorer/FileExplorer.vue';
import SearchExplorer from './Search/SearchExplorer.vue';
import GitExplorer from './Git/GitExplorer.vue';

const store = useWorkspace();
</script>
