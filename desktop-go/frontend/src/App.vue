<template>
    <div class="flex flex-col h-screen w-full bg-slate-900 text-slate-200">
        <Toolbar></Toolbar>
        <!-- IDE Panel (Right pane) -->
        <div class="flex-1 h-full flex flex-row overflow-hidden relative">
            <template v-if="store.active">
                <Editor></Editor>
            </template>
            <template v-else>
                <WorkspaceSelector></WorkspaceSelector>
            </template>
        </div>
    </div>
</template>

<script setup lang="ts">
import { provide, ref } from 'vue';
import Editor from './components/Editor.vue';
import Toolbar from './components/Toolbar.vue';
import WorkspaceSelector from './components/WorkspaceSelector.vue';

import { useWorkspace } from './store/workspace.ts';

const store = useWorkspace();
</script>
