<template>
<div class="flex flex-row h-full w-auto">
    <template v-for="(item, index) of store.active?.openTabs" :key="item.id">
        <div class="flex h-full items-center border-t-2 border-r border-slate-800/60 text-xs cursor-pointer space-x-1 px-1 -mb-0.2"
            :class="[
                (store.active?.activeFile?.path == item.path)? 'bg-slate-950 border-t-blue-800':'bg-slate-900 border-t-transparent'
            ]">
            <div  @click="store.setActiveFile(item)" class="flex h-full items-center space-x-1 pl-2">
                <FileIcon
                    :lang="item.extension || item.name"
                    :type="item.type"
                    :extended="item.expanded"
                    :size="14"
                />
                <span class="whitespace-nowrap">{{ item.name }}</span>
            </div>
            <div @click="store.closeTab(item.id, null)" class="group size-5 flex items-center justify-center">
                <svg class="h-2 w-2 text-slate-500 group-hover:text-slate-300" viewBox="0 0 10 10">
                    <line x1="1" y1="1" x2="9" y2="9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" />
                    <line x1="9" y1="1" x2="1" y2="9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" />
                </svg>
            </div>
        </div>
    </template>
</div>
</template>


<script setup lang="ts">
import { useWorkspace } from '@/store/workspace';
import FileIcon from '@/views/components/General/FileIcon.vue';

const store = useWorkspace()
</script>