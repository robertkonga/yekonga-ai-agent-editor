<template>
    <div v-if="node.type === 'directory'" @click="switchFile(node)"
        class="px-2 py-0.5 text-xs font-medium text-slate-300 cursor-pointer transition relative group ">
        <div class="flex items-center space-x-1 z-20 relative">
            <svg :class="['size-4 p-px text-indigo-400/80 transition-transform', node.expanded ? 'rotate-90' : '']"
                fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="m8.25 4.5 7.5 7.5-7.5 7.5" />
            </svg>
            <svg class="h-4 w-4 text-amber-400/80" fill="none" viewBox="0 0 24 24" stroke-width="2"
                stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round"
                    d="M2.25 12.75V12A2.25 2.25 0 0 1 4.5 9.75h15A2.25 2.25 0 0 1 21.75 12v.75m-8.69-6.44-2.12-2.12a1.5 1.5 0 0 0-1.061-.44H4.5A2.25 2.25 0 0 0 2.25 6v12a2.25 2.25 0 0 0 2.25 2.25h15A2.25 2.25 0 0 0 21.75 18V9a2.25 2.25 0 0 0-2.25-2.25h-5.379a1.5 1.5 0 0 1-1.06-.44Z" />
            </svg>
            <span>{{ node.name }}</span>
        </div>
        <background :level="props.level"></background>
    </div>

    <div v-if="node.type === 'directory' && node.expanded"
        class="pl-0 space-y-0 border-l border-slate-800 ml-4">
        <template v-for="child in node.children" :key="child.id">
            <FileExplorerItem :item="child" :level="(props.level + 1)"></FileExplorerItem>
        </template>
    </div>

    <div v-if="node.type === 'file'" @click="switchFile(node)" :class="[
        ' px-2 py-0.5 text-xs cursor-pointer transition relative group',
    ]">
        <div class="flex items-center space-x-2 z-20 relative">
            <FileIcon 
                :lang="`${node.extension}`.substring(1)" 
                :type="node.type" 
                :extended="node.expanded"
                :class="['rounded-full', getFileColorClass(node.lang)]">
            </FileIcon>
            <span class="truncate">{{ node.name }}</span>
        </div>
        <background :level="props.level"></background>
    </div>
</template>
<script lang="ts" setup>
import { h, computed, defineComponent, inject, type Ref } from 'vue';
import { ReadFile } from '@wails/go/main/App';
import FileIcon from './FileIcon.vue';
import { useWorkspace, type FileNode } from '@/store/workspace.ts';

const props = defineProps<{
    item: FileNode;
    level: number;
}>()

const store = useWorkspace()
const node = computed(() => props.item)

// Initialize context pointer to first child element
const targetPath = inject<Ref<FileNode | null>>("targetPath") as Ref<FileNode>;

const switchFile = async (fileNode: FileNode) => {
    targetPath.value = fileNode;

    if (fileNode.type === 'directory') {
        fileNode.expanded = !fileNode.expanded;
    }

    if (store.active!.activeFile! && store.active!.activeFile!.id === fileNode.id) return;

    if (fileNode.type === 'file') {
        store.setActiveFile(fileNode)
    }

};
const getFileColorClass = inject("getFileColorClass") as (lang: any) => "bg-slate-400" | "bg-yellow-400" | "bg-blue-400" | "bg-emerald-400";

const background = defineComponent({
    name: 'FastButton',
    props: {
        level: { type: Number, required: true },
    },
    setup(subProps) {
        // Return a render function instead of an object
        return () => h(
            'div', 
            { 
                class: [
                    'absolute top-0 right-0 bottom-0 rounded-none text-white border z-10',
                    (targetPath.value && targetPath.value.id === node.value.id)
                        ? 'bg-blue-800/30 text-white border-blue-700/40'
                        : 'text-slate-400 border-transparent bg-transparent group-hover:bg-slate-700/50 group-hover:text-slate-300'
                ],
                style: {
                    left: `-${subProps.level * 1}rem`
                }
            }, 
            ''
        );
    }
});
</script>
