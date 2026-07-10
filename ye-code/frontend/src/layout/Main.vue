<template>
    <div class="flex flex-col h-screen w-full bg-slate-900 text-slate-200">
        <Toolbar></Toolbar>
        <!-- IDE Panel (Right pane) -->
        <div class="flex-1 h-full flex flex-row overflow-hidden relative">
            <template v-if="!store.active">
                <Workspace></Workspace>
            </template>
            <template v-else-if="store.active">
                <template v-if="store.active.viewMode == 'AGENT'">
                    <Agent></Agent>
                </template>
                <template v-if="store.active.viewMode == 'EDITOR'">
                    <Editor></Editor>
                </template>
            </template>
            <template v-else>
                <Workspace></Workspace>
            </template>
        </div>
    </div>
</template>

<script setup lang="ts">
import { provide, ref } from 'vue';
import Toolbar from '@/layout/Toolbar.vue';
import Editor from '@/views/Editor.vue';
import Agent from '@/views/Agent.vue';
import Settings from '@/views/Settings.vue';
import Workspace from '@/views/Workspace.vue';
import { useWorkspace } from '@/store/workspace.ts';

const store = useWorkspace();


</script>
