<template>
<div class="flex flex-row h-full w-auto">
    <template v-for="(item, index) of store.active?.openTabs" :key="item.id">
        <div @click="store.setActiveFile(item)" class="flex h-full items-center border-t-2 border-r border-slate-800/60 text-xs cursor-pointer space-x-3 px-4 -mb-0.2"
            :class="[
                (store.active?.activeFile?.path == item.path)? 'bg-slate-950 border-t-blue-800':'bg-slate-900 border-t-transparent'
            ]">
            <FileIcon 
                :lang="`${item.lang}`" 
                :type="item.type" 
                :extended="item.expanded"
                :class="[getFileColorClass(item.lang || 'planetext')]">
            </FileIcon>
            <span>{{ item.name }}</span>
        </div>
    </template>
</div>
</template>


<script setup lang="ts">
import { useWorkspace, type FileNode } from '@/store/workspace';
import { getFileColorClass } from '../utils';

const store = useWorkspace()
</script>