<template>
    <div class="flex flex-col h-full">
        <!-- Header -->
        <div class="flex items-center justify-between px-3 py-2 border-b border-slate-800/60">
            <span class="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">Commits</span>
            <button @click="loadGitLog" class="text-slate-500 hover:text-slate-300" title="Refresh">
                <i class="ye ye-rotate ye-sm"></i>
            </button>
        </div>

        <!-- Loading -->
        <div v-if="loadingCommits" class="flex items-center justify-center py-8 text-xs text-slate-500">
            <span
                class="inline-block size-3 border-2 border-slate-500 border-t-transparent rounded-full animate-spin mr-2"></span>
            Loading commits...
        </div>

        <!-- Not a git repo -->
        <div v-else-if="notGitRepo"
            class="flex flex-col items-center justify-center py-8 text-xs text-slate-500 space-y-2">
            <i class="ye ye-git-alt text-2xl text-slate-700"></i>
            <p>Not a git repository</p>
        </div>

        <!-- Commit list -->
        <div v-else class="flex-1 overflow-y-auto">
            <div v-if="!commits.length" class="flex items-center justify-center h-full text-xs text-slate-600">
                No commits yet
            </div>
            <div v-for="commit in commits" :key="commit.hash" class="border-b border-slate-800/40 last:border-b-0">
                <div @click="loadCommitDetail(commit.hash)"
                    class="flex flex-col gap-1 px-3 py-2 cursor-pointer hover:bg-slate-800/30 transition-colors">
                    <div class="flex items-center gap-2">
                        <div
                            class="size-6 rounded-full bg-slate-700 flex items-center justify-center text-[10px] text-slate-300 font-semibold shrink-0">
                            {{ commit.author.charAt(0).toUpperCase() }}
                        </div>
                        <div class="flex-1 min-w-0">
                            <p class="text-xs text-slate-200 truncate font-medium">{{ commit.message }}</p>
                            <p class="text-[10px] text-slate-500">{{ commit.author }} · {{ formatDate(commit.date) }}
                            </p>
                        </div>
                        <span class="text-[9px] font-mono text-slate-600">{{ commit.hash.substring(0, 7) }}</span>
                    </div>
                </div>

                <!-- Expanded commit detail -->
                <div v-if="expandedCommit === commit.hash && commitDetail"
                    class="bg-slate-900/60 border-t border-slate-800/40">
                    <div v-if="loadingDetail" class="px-4 py-3 text-[10px] text-slate-500">Loading details...</div>
                    <template v-else>
                        <!-- Changed files list -->
                        <div class="px-2 py-1.5">
                            <div class="text-[10px] font-semibold text-slate-500 uppercase tracking-wider px-1 mb-1">
                                {{ commitDetail.files.length }} file{{ commitDetail.files.length !== 1 ? 's' : '' }}
                                changed
                            </div>
                            <div v-for="file in commitDetail.files" :key="file.path"
                                class="flex items-center gap-2 px-2 py-1.5 text-xs hover:bg-slate-800/40 cursor-pointer rounded"
                                @click="toggleFileDiff(file.path)">
                                <span class="text-[9px] font-semibold uppercase w-14 shrink-0" :class="{
                                    'text-emerald-400': file.status === 'added',
                                    'text-yellow-400': file.status === 'modified',
                                    'text-red-400': file.status === 'deleted',
                                    'text-blue-400': file.status === 'renamed' || file.status === 'copied'
                                }">
                                    {{ file.status }}
                                </span>
                                <span class="truncate text-slate-300">{{ file.path }}</span>
                                <i class="ye ye-angle-right text-[10px] text-slate-600 ml-auto transition-transform"
                                    :class="{ 'rotate-90': expandedFileDiffs[file.path] }"></i>
                            </div>

                            <!-- Diff content for each file -->
                            <div v-for="file in commitDetail.files" :key="'diff-' + file.path">
                                <div v-if="expandedFileDiffs[file.path] && file.patch"
                                    class="border-t border-slate-800/40">
                                    <div v-if="file.patch === '[Binary file]'"
                                        class="px-4 py-2 text-[10px] text-slate-500 italic">
                                        Binary file
                                    </div>
                                    <pre v-else
                                        class="text-[10px] font-mono leading-relaxed overflow-x-auto max-h-96 overflow-y-auto p-2 bg-slate-950/50">
                    <code>
<span v-for="(line, li) in splitPatch(file.patch)" :key="li"
    :class="getDiffLineClass(line)"
    class="block whitespace-pre">{{ line }}</span>
                    </code>
                  </pre>
                                </div>
                                <div v-else-if="expandedFileDiffs[file.path] && !file.patch"
                                    class="px-4 py-2 text-[10px] text-slate-500 italic">
                                    No diff available
                                </div>
                            </div>
                        </div>
                    </template>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch, onMounted } from 'vue';
import { useWorkspace } from '@/store/workspace';
import { GetGitLog, GetGitCommitDetail } from '@wails/go/main/App';

const store = useWorkspace();

// ── Git history state ────────────────────────────────────────────────────
const commits = ref<{ hash: string; author: string; date: string; message: string }[]>([]);
const loadingCommits = ref(false);
const notGitRepo = ref(false);
const expandedCommit = ref<string | null>(null);
const commitDetail = ref<{
    hash: string;
    author: string;
    date: string;
    message: string;
    files: { path: string; status: string; patch: string }[];
} | null>(null);
const loadingDetail = ref(false);
const expandedFileDiffs = reactive<Record<string, boolean>>({});

function formatDate(dateStr: string): string {
    try {
        const d = new Date(dateStr);
        if (isNaN(d.getTime())) return dateStr;
        const now = new Date();
        const diff = now.getTime() - d.getTime();
        const mins = Math.floor(diff / 60000);
        const hours = Math.floor(diff / 3600000);
        const days = Math.floor(diff / 86400000);

        if (mins < 1) return 'just now';
        if (mins < 60) return `${mins}m ago`;
        if (hours < 24) return `${hours}h ago`;
        if (days < 7) return `${days}d ago`;
        return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
    } catch {
        return dateStr;
    }
}

async function loadGitLog() {
    if (!store.active?.path) return;

    loadingCommits.value = true;
    notGitRepo.value = false;
    commits.value = [];
    expandedCommit.value = null;
    commitDetail.value = null;

    try {
        const log = await GetGitLog(store.active.path);
        commits.value = log || [];
    } catch (err: any) {
        console.error('Git log failed:', err);
        notGitRepo.value = true;
    } finally {
        loadingCommits.value = false;
    }
}

async function loadCommitDetail(hash: string) {
    if (expandedCommit.value === hash) {
        expandedCommit.value = null;
        commitDetail.value = null;
        return;
    }

    if (!store.active?.path) return;

    expandedCommit.value = hash;
    loadingDetail.value = true;
    commitDetail.value = null;

    // Clear expanded diffs
    for (const key in expandedFileDiffs) {
        delete expandedFileDiffs[key];
    }

    try {
        const detail = await GetGitCommitDetail(store.active.path, hash);
        commitDetail.value = detail;
    } catch (err: any) {
        console.error('Failed to load commit detail:', err);
    } finally {
        loadingDetail.value = false;
    }
}

function toggleFileDiff(filePath: string) {
    expandedFileDiffs[filePath] = !expandedFileDiffs[filePath];
}

/**
 * Split a patch string into lines for rendering.
 */
function splitPatch(patch: string): string[] {
    return patch.split('\n');
}

/**
 * Get the CSS class for a diff line based on its content.
 */
function getDiffLineClass(line: string): string {
    if (line.startsWith('@@')) {
        return 'text-slate-600';
    }
    if (line.startsWith('-') && !line.startsWith('---')) {
        return 'text-red-400 bg-red-950/30';
    }
    if (line.startsWith('+') && !line.startsWith('+++')) {
        return 'text-emerald-400 bg-emerald-950/30';
    }
    if (line.startsWith('diff') || line.startsWith('index') ||
        line.startsWith('---') || line.startsWith('+++') ||
        line.startsWith('new') || line.startsWith('deleted') ||
        line.startsWith(' ')) {
        return 'text-slate-500';
    }
    return 'text-slate-500';
}

// Load git log when component is mounted
onMounted(() => {
    if (store.active?.path) {
        loadGitLog();
    }
});

// Reload when active path changes
watch(() => store.active?.path, () => {
    loadGitLog();
});
</script>
