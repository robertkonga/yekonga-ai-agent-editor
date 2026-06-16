<template>
    <div class="h-full border-r border-slate-800/60 w-60 bg-slate-900/60 overflow-y-auto">
        <FileExplorer :isEditor="true"></FileExplorer>
    </div>

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
</template>

<script lang="ts" setup>
import { inject, type Ref } from 'vue';
import Terminal from './Terminal.vue';
import CodeEditor from './CodeEditor.vue';
import FileExplorer from '@/views/components/General/FileExplorer.vue';
import { useWorkspace } from '@/store/workspace.ts';

const store = useWorkspace();
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