<template>
    <div class="h-full flex-1 flex flex-col items-center justify-center pb-3">
        <div class="w-full h-full flex flex-col items-center justify-center bg-gray-400/5 border border-gray-400/20 rounded-xl overflow-hidden">
            <template v-if="store.active!.customizedView == 'DATA_SCHEMA'">
                <DataSchema></DataSchema>
            </template>
            <template v-else-if="store.active!.customizedView == 'PROJECT'">
                <Project></Project>
            </template>
            <template v-else>
                <ChatPanel></ChatPanel>
            </template>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useWorkspace, type CUSTOMIZATION_VIEW } from '@/store/workspace';
import ChatPanel from './parts/Customization/ChatPanel.vue';
import Project from './parts/Customization/Project.vue';
import DataSchema from './parts/Customization/DataSchema.vue';

const store = useWorkspace();
</script>
