
<template>
    <div class="group relative flex cursor-pointer items-center gap-3 rounded px-2 py-2 transition-colors hover:bg-slate-800"
        @click="emit('open')" @mouseenter="emit('mouseenter')" @mouseleave="emit('mouseleave')">
        <!-- Folder icon -->
        <svg class="shrink-0 text-[#c5a028]" width="16" height="16" viewBox="0 0 16 16" fill="none">
            <path
                d="M1.5 4C1.5 3.17 2.17 2.5 3 2.5H6L7.5 4H13C13.83 4 14.5 4.67 14.5 5.5V12C14.5 12.83 13.83 13.5 13 13.5H3C2.17 13.5 1.5 12.83 1.5 12V4Z"
                fill="currentColor" opacity="0.85" />
        </svg>

        <!-- Name + path -->
        <div @click="store.openWorkshop(workspace.path)" class="min-w-0 flex-1">
            <div class="flex items-center gap-1.5">
                <span class="truncate text-[12px] font-medium text-[#e8e8e8]">{{ workspace.name }}</span>
                <!-- Pin badge -->
                <svg v-if="workspace.isPinned" class="shrink-0 text-[#555]" width="8" height="8" viewBox="0 0 8 8">
                    <path d="M4 0L5 3H8L5.5 5L6.5 8L4 6L1.5 8L2.5 5L0 3H3L4 0Z" fill="currentColor" />
                </svg>
            </div>
            <div class="truncate text-[11px] text-[#6a6a6a]">{{ shortenPath(workspace.path) }}</div>
        </div>

        <!-- Last opened -->
        <span v-if="workspace.lastOpened" class="shrink-0 text-[10px] text-[#555] group-hover:opacity-0">
            {{ formatRelative(workspace.lastOpened) }}
        </span>

        <!-- Remove button (shown on hover) -->
        <button
            class="absolute right-2 hidden shrink-0 rounded p-0.5 text-[#555] hover:bg-[#3c3c3c] hover:text-[#cccccc] group-hover:flex"
            title="Remove from recent" aria-label="Remove from recent" @click.stop="emit('remove')">
            <svg width="10" height="10" viewBox="0 0 10 10">
                <line x1="1" y1="1" x2="9" y2="9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" />
                <line x1="9" y1="1" x2="1" y2="9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" />
            </svg>
        </button>
    </div>
</template>

<script setup lang="ts">
import { useWorkspace, type Workspace } from '@/store/workspace';


interface Props {
    workspace: Workspace
    hovered?: boolean
}

defineProps<Props>()

const emit = defineEmits<{
    open: []
    remove: []
    mouseenter: []
    mouseleave: []
}>()

const store = useWorkspace();

function formatRelative(date: any): string {
    if(typeof date == 'string') {
        try {
            date = new Date(Date.parse(date));
        } catch (error) {
            date = new Date();
        }
    }

    try {
        const diff = Date.now() - date.getTime()
        const mins = Math.floor(diff / 60000)
    
        if (mins < 1) return 'just now';
        if (mins < 60) return `${mins}m ago`;
    
        const hrs = Math.floor(mins / 60);
        if (hrs < 24) return `${hrs}h ago`;
    
        const days = Math.floor(hrs / 24);
        if (days < 7) return `${days}d ago`;
    
        return date.toLocaleDateString();
    } catch (error) {
        console.warn(error);
    }

    return (new Date()).toLocaleDateString();
}

function shortenPath(path: string): string {
    return path.replace(/^\/home\/[^/]+/, '~')
}
</script>
