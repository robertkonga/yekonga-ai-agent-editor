<template>
    <div class="flex flex-col h-full pb-3">
        <div class="w-full h-full flex flex-col bg-gray-400/5 border border-gray-400/20 rounded-xl">
            <div class="flex items-center justify-between p-3 border-b border-slate-800/60">
                <div class="flex space-x-2 text-sm">
                    <template v-for="tab in tabs" :key="tab.id">
                        <button @click="activeTab = tab.id" class="cursor-pointer text-xs px-3 rounded py-1"
                            :class="(activeTab === tab.id) ? 'text-white bg-slate-700/50 ' : 'text-gray-500 hover:text-gray-300'">
                            {{ tab.label }}
                        </button>
                    </template>
                </div>
                <div class="flex items-center gap-2 text-gray-400">
                    <button class="hover:text-white"><i class="ye ye-rotate ye-sm"></i></button>
                    <button class="hover:text-white"><i class="ye ye-binnacles ye-sm"></i></button>
                    <button class="hover:text-white"><i class="ye ye-file ye-sm"></i></button>
                </div>
            </div>

            <div class="flex-1 overflow-y-auto">
                <template v-if="activeTab == 'FILE'">
                    <FileExplorer :isEditor="false"></FileExplorer>
                </template>
                <template v-else>
                    <Changes></Changes>
                </template>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Changes from './parts/Changes.vue';
import FileExplorer from '@/views/components/General/FileExplorer/FileExplorer.vue';

type LABEL = "FILE" | "CHANGE";
const tabs:{
    id: LABEL;
    label: string;
}[] = [
    { id:"FILE", label: "Files" },
    { id:"CHANGE", label: "Changes" },
];
const activeTab = ref<LABEL>('FILE')

</script>