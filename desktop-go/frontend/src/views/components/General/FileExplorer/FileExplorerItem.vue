<template>
    <div v-if="node.type === 'directory'" @click="switchFile()"
        draggable="true" @dragstart="onDragStart" @dragover.prevent @drop.stop="onDrop"
        class="px-2 h-6 text-sm font-medium text-slate-300 cursor-pointer transition relative group ">
        <div class="flex h-full items-center space-x-1 z-20 relative">
            
            <span class="w-6 flex justify-center leading-tight">
                <template v-if="store.isSimpleIcon">
                    <svg :class="['size-4 p-px text-indigo-400/80 transition-transform', node.expanded ? 'rotate-90' : '']"
                        fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"> 
                        <path stroke-linecap="round" stroke-linejoin="round" d="m8.25 4.5 7.5 7.5-7.5 7.5" />
                    </svg>
                </template>
                <template v-else>
                    <FileIcon
                        :lang="node.name"
                        type="directory"
                        :extended="node.expanded"
                        :size="20"/>
                </template>
            </span>
            <span class="h-full flex-1 w-10/12"><FileName v-model="node" :isEditor="props.isEditor"></FileName></span>
        </div>
        <background v-if="props.isEditor" :level="props.level" :isEditor="props.isEditor"></background>
    </div>

    <div v-if="node.type === 'directory' && node.children && node.expanded"
        class="pl-0 space-y-0 border-l border-slate-800 ml-4">
        <template v-for="(child, i) of node.children" :key="child.id">
            <FileExplorerItem v-model="node.children[i]" @change="update" :level="(props.level + 1)" :isEditor="props.isEditor"></FileExplorerItem>
        </template>
    </div>

    <div v-if="node.type === 'file'" @click="switchFile()" 
        draggable="true" @dragstart="onDragStart"
        :class="[
        'px-2 h-6 text-sm cursor-pointer transition relative group',
    ]">
        <div class="flex h-full items-center space-x-1 z-20 relative">
            <span class="w-6 flex justify-center leading-tight">
                <template v-if="store.isSimpleIcon">
                    <FileIcon
                        :lang="'default'"
                        type="file"
                        :extended="false"
                        :size="16"
                    ></FileIcon>
                </template>
                <template v-else>
                    <FileIcon
                        :lang="node.extension || node.name"
                        :type="node.type"
                        :extended="node.expanded"
                        :size="16"
                    ></FileIcon>
                </template>
            </span>
            
            <span class="h-full flex-1 w-10/12"><FileName v-model="node" :isEditor="props.isEditor"></FileName></span>
        </div>
        <background :level="props.level" :isEditor="props.isEditor"></background>
    </div>
</template>
<script lang="ts" setup>
import { h, defineComponent, inject, type Ref, ref, watch } from 'vue';
import FileIcon from './FileIcon.vue';
import { useWorkspace, type FileNode } from '@/store/workspace.ts';
import FileName from './FileName.vue';

const emit = defineEmits(["update:modelValue", "change"])
const props = defineProps<{
    modelValue: FileNode;
    level: number;
    isEditor: boolean;
}>()

const store = useWorkspace();
const node = ref<FileNode>(props.modelValue);
 
// Initialize context pointer to first child element
const targetPath = inject<Ref<FileNode>>("targetPath") || ref<FileNode>();
// const targetPath = ref<FileNode>();

const update = () => {
    setTimeout(()=>{
        emit("change", node.value);
    }, 500)
}

const switchFile = async () => {
    targetPath.value = node.value;

    if (node.value.type === 'directory') {
        node.value.expanded = !node.value.expanded;
        
        emit("update:modelValue", node.value);
    } else {
        if (!(store.active!.activeFile! && store.active!.activeFile!.id === node.value.id)) {
            if (node.value.type === 'file') {
                store.setActiveFile(node.value);
            }
        }
    }

    setTimeout(()=>{
        store.saveLocally();
    }, 2000)
};

const onDragStart = (e: DragEvent) => {
    if (node.value.path) {
        e.dataTransfer?.setData("text/plain", node.value.path);
    }
}

const onDrop = (e: DragEvent) => {
    const sourcePath = e.dataTransfer?.getData("text/plain");
    if (sourcePath && node.value.path && node.value.type === 'directory' && sourcePath !== node.value.path) {
        store.moveFile(sourcePath, node.value.path);
    }
}

const background = defineComponent({
    name: 'FastButton',
    props: {
        level: { type: Number, required: true },
        isEditor: { type: Boolean, required: true },
    },
    setup(subProps) {
        
        // Return a render function instead of an object
        return () => h(
            'div', 
            { 
                class: [
                    'absolute top-0 right-0 bottom-0 rounded-none text-white border z-10',
                    (subProps.isEditor && targetPath && targetPath.value && targetPath.value.id === node.value.id)
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

watch(()=>props.modelValue, (v2, v1) => {
    node.value = v2;
}, { deep: true })

// watch(node, (v2, v1) => {
//     emit("update:modelValue", v2);
// }, { deep: true })
</script>
