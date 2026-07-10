<template>
    <expandable title="Customizations" :show="true">
        <template #default>
            <ul class="space-y-1">
                <li @click="store.changeCustomizationView('AGENT')" class="flex items-center justify-between py-1.5 px-2 hover:bg-slate-800 rounded cursor-pointer text-sm"
                    :class="{'bg-slate-800':(store.active?.customizedView == 'AGENT')}"
                >
                    <div class="flex items-center gap-2 text-gray-400">
                        <span><i class="ye ye-user-tie"></i></span>
                        <span class="hover:text-gray-200">Agent</span>
                    </div>
                    <span class="text-xs text-gray-500"></span>
                </li>
                <li @click="store.changeCustomizationView('PROJECT')" class="flex items-center justify-between py-1.5 px-2 hover:bg-slate-800 rounded cursor-pointer text-sm"
                    :class="{'bg-slate-800':(store.active?.customizedView == 'PROJECT')}"
                >
                    <div class="flex items-center gap-2 text-gray-400">
                        <span><i class="ye ye-folder-open"></i></span>
                        <span class="hover:text-gray-200">Generate project</span>
                    </div>
                    <span class="text-xs text-gray-500"></span>
                </li>
                <li @click="store.changeCustomizationView('DATA_SCHEMA')" class="flex items-center justify-between py-1.5 px-2 hover:bg-slate-800 rounded cursor-pointer text-sm"
                    :class="{'bg-slate-800':(store.active?.customizedView == 'DATA_SCHEMA')}"
                >
                    <div class="flex items-center gap-2 text-gray-400">
                        <span><i class="ye ye-database"></i></span>
                        <span class="hover:text-gray-200">Generate data schema</span>
                    </div>
                    <span class="text-xs text-gray-500"></span>
                </li>
                <template v-if="false">
                    <li v-for="item in customizations" :key="item.name"
                        class="flex items-center justify-between py-1.5 px-2 hover:bg-slate-800 rounded cursor-pointer text-sm">
                        <div class="flex items-center gap-2 text-gray-400">
                            <span>{{ item.icon }}</span>
                            <span class="hover:text-gray-200">{{ item.name }}</span>
                        </div>
                        <span v-if="item.count !== null" class="text-xs text-gray-500">{{ item.count }}</span>
                    </li>
                </template>
            </ul>
        </template>
    </expandable>
</template>

<script setup lang="ts">
import { useWorkspace } from '@/store/workspace';
import { ref } from 'vue'

const store = useWorkspace();

const customizations = ref([
    { name: 'Agents', count: null, icon: '⧉' },
    { name: 'Skills', count: 12, icon: '💡' },
    { name: 'Instructions', count: null, icon: '📖' },
    { name: 'Hooks', count: null, icon: '⚡' },
    { name: 'MCP Servers', count: 1, icon: '🖥️' },
    { name: 'Plugins', count: null, icon: '🔌' },
])
</script>
