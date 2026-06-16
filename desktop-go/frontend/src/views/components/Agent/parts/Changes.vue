<template>
    <expandable title="Changes" :show="true">
        <template #extraHeader>
            <div class="flex items-center gap-2">
                <span v-if="store.active?.changedFiles.length" class="bg-blue-600 text-[10px] px-1.5 rounded-full text-white font-bold">
                    {{ store.active.changedFiles.length }}
                </span>
            </div>
        </template>
        <template #default>
            <div class="min-h-20 max-h-60 overflow-y-auto">
                <div v-if="!store.active?.changedFiles.length" class="text-xs text-gray-500 py-4 text-center">
                    No modified files
                </div>
                <div v-else class="flex flex-col gap-1">
                    <div v-for="file in store.active.changedFiles" :key="file.path" 
                         @click="store.openFile(file, null)"
                         class="flex items-center gap-2 text-xs text-gray-300 hover:bg-slate-800 p-1.5 rounded cursor-pointer group">
                        <i class="ye ye-file text-blue-400"></i>
                        <div class="truncate flex-grow">{{ file.name }}</div>
                        <button @click.stop="store.clearChange(file.path || '')" 
                                class="opacity-0 group-hover:opacity-100 text-gray-500 hover:text-red-400">
                            <i class="ye ye-xmark"></i>
                        </button>
                    </div>
                </div>
            </div>
        </template>
    </expandable>
</template>

<script setup lang="ts">
import { useWorkspace } from '@/store/workspace';
import Expandable from '@/components/Expandable.vue';

const store = useWorkspace();
</script>