
<template>
  <div class="flex h-full w-full bg-slate-900 text-slate-200">
    <!-- IDE Panel (Right pane) -->
    <div class="flex-1 h-full flex flex-row overflow-hidden relative">
      <IDEPanel />
    </div>
    
    <!-- Chat Panel (Left pane) -->
    <div v-if="store.active" class="h-full w-70 bg-slate-800 flex flex-col">
      <ChatPanel />
    </div>
  </div>
</template>

<script setup lang="ts">
import { provide, ref } from 'vue';
import { useWorkspace, type FileNode } from '@/store/workspace.ts';
import IDEPanel from './IDE/IDEPanel.vue';
import ChatPanel from './Chat/ChatPanel.vue';
import { getFileColorClass } from './IDE/utils.ts';

const store = useWorkspace();

// Dynamic File Switching Handler
const targetPath = ref<FileNode | null>(null)

provide("targetPath", targetPath);
provide("getFileColorClass", getFileColorClass);
</script>
