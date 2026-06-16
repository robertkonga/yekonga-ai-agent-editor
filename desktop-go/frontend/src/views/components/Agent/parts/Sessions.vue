<template>
    <expandable title="Sessions" :show="true">
        <template #extraHeader>
            <div class="flex items-center gap-2">
                <button class="text-xs bg-slate-800 hover:bg-slate-700 px-2 py-1 rounded text-gray-400">
                    New <span class="text-[10px] ml-1 opacity-70">Ctrl+N</span>
                </button>
                <button class="bg-slate-800 hover:bg-slate-700 px-2 py-0 rounded text-gray-400 hover:text-white">
                    <i class="ye ye-magnifying-glass ye-xs"></i>
                </button>
            </div>
        </template>
        <template #default>
            <div class="min-h-80">
                <div class="mb-2 text-xs text-gray-500">{{ store.active?.name }}</div>
            
                <div class="group cursor-pointer">
                    <div class="flex items-center gap-2 text-xs text-gray-300">
                        <span class="w-1.5 h-1.5 rounded-full bg-gray-500"></span>
                        <div class="truncate">Frontend project builder creation</div>
                    </div>
                    <div class="text-xs text-gray-500 ml-3 mt-1 flex items-center gap-1">
                        <span><i class="ye ye-folder-open"></i></span> 1 wk ago
                    </div>
                </div>  
            </div>
        </template>
    </expandable>

</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useWorkspace } from '@/store/workspace';

const store = useWorkspace();
</script>