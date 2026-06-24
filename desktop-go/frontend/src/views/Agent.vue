<template>
    <div class="flex h-full gap-3 w-full bg-slate-900 text-gray-300 font-sans overflow-hidden">
        <aside class="w-64 flex flex-col">
            <LeftSidebar />
        </aside>

        <main class="flex-1 flex flex-col min-w-0 pt-3">
            <MainWorkspace />
        </main>

        <aside class="w-72 flex flex-col pe-3 pt-3">
            <RightSidebar />
        </aside>
    </div>
</template>

<script setup lang="ts">
import LeftSidebar from './components/Agent/LeftSidebar.vue'
import MainWorkspace from './components/Agent/MainWorkspace.vue'
import RightSidebar from './components/Agent/RightSidebar.vue'
</script>
