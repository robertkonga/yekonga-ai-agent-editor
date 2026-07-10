<template>
    <div class="flex flex-col h-full">
        <!-- Search input -->
        <div class="p-2 border-b border-slate-800/60 space-y-2">
            <div class="flex items-center gap-1.5 bg-slate-800 rounded px-2 py-1.5">
                <i class="ye ye-magnifying-glass text-slate-500 ye-sm"></i>

                <input v-model="searchQuery" type="text" placeholder="Search files..."
                    class="flex-1 bg-transparent text-xs text-slate-200 outline-none placeholder:text-slate-500"
                    @keydown.enter="doSearch" />

                <button v-if="searchQuery" @click="searchQuery = ''; searchResults = []; searched = false;"
                    class="text-slate-500 hover:text-slate-300">
                    <i class="ye ye-xmark ye-xs"></i>
                </button>
            </div>
            <div class="flex items-center gap-1.5">
                <div class="flex items-center flex-1 gap-1 bg-slate-800 rounded px-2 py-1.5">
                    <i class="ye ye-arrows-rotate text-slate-500 ye-xs"></i>
                    <input v-model="replaceText" type="text" placeholder="Replace with..."
                        class="flex-1 bg-transparent text-xs text-slate-200 outline-none placeholder:text-slate-500" />
                </div>
                <button @click="doReplaceAll"
                    class="h-7 text-[12px] px-2 py-1.5 rounded bg-blue-600 hover:bg-blue-500 text-white font-medium disabled:opacity-40"
                    :disabled="!searchResults.length || !replaceText">
                    Replace
                </button>
            </div>
            <div v-if="searching" class="text-[10px] text-slate-500 flex items-center gap-1.5">
                <span
                    class="inline-block size-3 border-2 border-slate-500 border-t-transparent rounded-full animate-spin"></span>
                Searching...
            </div>
            <div v-else-if="searchResults.length" class="text-[10px] text-slate-400">
                Found {{ searchResults.length }} match{{ searchResults.length !== 1 ? 'es' : '' }}
                in {{ uniqueFiles }} file{{ uniqueFiles !== 1 ? 's' : '' }}
            </div>
            <div v-else-if="searched && !searching" class="text-[10px] text-slate-500">
                No results found
            </div>
        </div>

        <!-- Search results -->
        <div class="flex-1 overflow-y-auto">
            <div v-if="!searchResults.length && !searched"
                class="flex items-center justify-center h-full text-xs text-slate-600">
                <div class="text-center space-y-2">
                    <i class="ye ye-magnifying-glass text-2xl text-slate-700"></i>
                    <p>Type a query and press Enter to search</p>
                </div>
            </div>

            <!-- Grouped by file -->
            <div v-for="(group, filePath) in groupedResults" :key="filePath"
                class="border-b border-slate-800/40 last:border-b-0">
                <div @click="toggleFileGroup(filePath)"
                    class="flex items-center gap-1.5 px-2 py-1.5 text-[11px] text-slate-400 hover:bg-slate-800/50 cursor-pointer sticky top-0 bg-slate-900/95 backdrop-blur z-10">
                    <i class="ye ye-angle-right text-[10px] transition-transform"
                        :class="{ 'rotate-90': expandedFiles[filePath] }"></i>
                    <i class="ye ye-file text-blue-400"></i>
                    <span class="truncate flex-1 leading-none">{{ filePath }}</span>
                    <span class="text-[10px] text-slate-600">{{ group.length }}</span>
                </div>
                <template v-if="expandedFiles[filePath]">
                    <div v-for="(result, idx) in group" :key="idx" @click="openSearchFile(result)"
                        class="flex items-start gap-2 px-3 py-1.5 text-xs hover:bg-slate-800/30 cursor-pointer group">
                        <span class="text-[10px] text-slate-600 w-6 text-right shrink-0 mt-0.5 font-mono">{{ result.line
                            }}</span>
                        <span class="text-slate-300 truncate flex-1 font-mono text-[11px]">{{ result.content }}</span>
                    </div>
                </template>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue';
import { useWorkspace } from '@/store/workspace';
import { SearchInFiles, ReplaceInFiles } from '@wails/yekonga-builder/service.ts';

const store = useWorkspace();

// ── Search state ──────────────────────────────────────────────────────────
const searchQuery = ref('');
const replaceText = ref('');
const searchResults = ref<{ file: string; line: number; content: string }[]>([]);
const searched = ref(false);
const searching = ref(false);
const expandedFiles = reactive<Record<string, boolean>>({});

const uniqueFiles = computed(() => {
    const files = new Set(searchResults.value.map(r => r.file));
    return files.size;
});

const groupedResults = computed(() => {
    const groups: Record<string, { file: string; line: number; content: string }[]> = {};

    for (const r of searchResults.value) {
        if (!groups[r.file]) groups[r.file] = [];
        groups[r.file].push(r);
    }

    return groups;
});

function toggleFileGroup(filePath: string) {
    expandedFiles[filePath] = !expandedFiles[filePath];
}

async function doSearch() {
    const q = searchQuery.value.trim();
    if (!q || !store.active?.path) return;

    searching.value = true;
    searched.value = true;
    searchResults.value = [];

    try {
        const res = await SearchInFiles(store.active.path, q);
        searchResults.value = res || [];
        // Expand all file groups by default
        for (const r of searchResults.value) {
            expandedFiles[r.file] = true;
        }
    } catch (err: any) {
        console.error('Search failed:', err);
    } finally {
        searching.value = false;
    }
}

async function doReplaceAll() {
    const q = searchQuery.value.trim();
    const r = replaceText.value.trim();
    if (!q || !r || !store.active?.path) return;

    if (!confirm(`Replace all occurrences of "${q}" with "${r}" in the workspace?`)) return;

    try {
        const results = await ReplaceInFiles(store.active.path, q, r);
        const successCount = results!.filter((res: any) => res.success).length;
        const failCount = results!.filter((res: any) => !res.success).length;

        if (failCount > 0) {
            alert(`Replaced in ${successCount} files. ${failCount} file(s) failed.`);
        } else {
            alert(`Successfully replaced in ${successCount} file(s).`);
        }

        // Re-run search to update results
        await doSearch();
        // Refresh workspace files
        await store.fetchWorkspaceFiles();
    } catch (err: any) {
        console.error('Replace failed:', err);
        alert('Replace failed: ' + err);
    }
}

/**
 * Open a file from search results and scroll to the matching line.
 */
function openSearchFile(result: { file: string; line: number; content: string }) {
    if (!store.active?.path) return;

    const filePath = store.active.path + '/' + result.file;
    const fileNode = {
        id: filePath,
        name: result.file.split('/').pop() || result.file,
        path: filePath,
        type: 'file' as const,
    };

    // Set pending line number so the editor scrolls to this line after opening
    store.active.pendingLineNumber = result.line;

    store.openFile(fileNode, null);
}
</script>
