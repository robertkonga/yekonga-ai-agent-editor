<template>
    <div class="flex flex-col h-screen w-full bg-slate-900 text-slate-200">
        <Layout></Layout>
    </div>
</template>

<script setup lang="ts">
import Layout from '@/layout/Main.vue';
import { useWorkspace } from './store/workspace.ts';

const store = useWorkspace();
store.loadIcons();
</script>
