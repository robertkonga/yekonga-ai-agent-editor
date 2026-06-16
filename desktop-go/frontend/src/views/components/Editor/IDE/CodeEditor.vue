<template>
    <div class="flex flex-1 flex-col border-r border-slate-800/60 relative">
        <div class="flex h-9 items-center justify-between  bg-slate-900">
            <div class="block h-full w-full relative">
                <div class="absolute top-0 bottom-0 left-0 right-0 overflow-x-auto overflow-y-hidden">
                    <OpenTabs></OpenTabs>
                </div>
            </div>
            <div class="flex space-x-2 px-3 border-l border-slate-800/60">
                <button @click="saveFileContent()"
                    class="flex items-center space-x-1 rounded-md bg-indigo-900 px-2 py-1.5 cursor-pointer text-xs font-semibold text-white transition hover:bg-indigo-800 active:scale-95 shadow-md shadow-green-600/10">
                    <svg class="h-3.5 w-3.5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z" />
                        <polyline points="17 21 17 13 7 13 7 21" />
                        <polyline points="7 3 7 8 15 8" />
                    </svg>
                    <!-- <span>Save</span> -->
                </button>
                <button v-if="false" @click="runCode"
                    class="flex items-center space-x-1 rounded-md bg-indigo-600 px-3 py-1.5 text-xs font-semibold text-white transition hover:bg-indigo-500 active:scale-95 shadow-md shadow-indigo-600/10">
                    <svg class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round"
                            d="M5.25 5.653c0-.856.917-1.398 1.667-.986l11.54 6.347a1.125 1.125 0 0 1 0 1.972l-11.54 6.347c-.75.412-1.667-.13-1.667-.986V5.653Z" />
                    </svg>
                    <span>Execute</span>
                </button>
            </div>
        </div>

        <div class="relative flex-1 bg-slate-950">
            <template v-if="store.active?.activeFile">
                <div ref="editorContainer" class="h-full"></div>
            </template>
            <template v-else>
                <div class="flex items-center justify-center bg-slate-900 w-full h-full">
                    <div class="p-3 ">Select file</div>
                </div>
            </template>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref, inject, onMounted, onBeforeUnmount, watch, type Ref } from 'vue';
import * as monaco from 'monaco-editor';
import { SaveFile } from '@wails/go/main/App';
import { useWorkspace, type FileNode } from '@/store/workspace';
import OpenTabs from './parts/OpenTabs.vue';

const emit = defineEmits(['update:modelValue']);

// Props for two-way binding (v-model) and language config
const props = defineProps({
    modelValue: { type: String, default: '' },
    language: { type: String, default: 'javascript' },
    theme: { type: String, default: 'vs-dark' }
});

const store = useWorkspace();

const editorContainer = ref(null);
const getFileColorClass = inject("getFileColorClass") as (lang: any) => "bg-slate-400" | "bg-yellow-400" | "bg-blue-400" | "bg-emerald-400";
let editor: monaco.editor.IStandaloneCodeEditor | null = null;

const runCode = () => {
    alert(`Executing execution layer for file context: ${store.active!.activeFile!.name}`);
};

const saveFileContent = (content?: any) => {
    if(!content && editor) {
        content = editor.getValue() || "";
    }
    
    if (store.active && store.active.activeFile && (content !== store.active.activeFile.content)) {
        if(store.active.activeFile.path) {
            SaveFile(content, store.active.activeFile.path)
            store.trackChange(store.active.activeFile)
        }
    }
}

const setEditor = () => {
    if (editorContainer.value) {
        if(!editor) {
            editor = monaco.editor.create(editorContainer.value, {
                value: "",
                language: "text",
                theme: 'tailwind-slate',
                automaticLayout: true,
                fontSize: 13,
                fontFamily: 'JetBrains Mono, Fira Code, monospace',
                minimap: { enabled: false },
                padding: { top: 6 },
                lineHeight: 18,
            });
        }

        if(store.active && store.active.activeFile) {
            editor.setValue(store.active.activeFile.content || "")
            monaco.editor.setModelLanguage(editor!.getModel() as monaco.editor.ITextModel, store.active.activeFile.lang || "text");
        }

        // 2. Add the Ctrl + S / Cmd + S Keybinding
        editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyS, () => {
            let content = editor!.getValue() || "";
            
            saveFileContent(content);
        });

        // Save modifications back to file system references reactively
        editor.onDidChangeModelContent(() => {
            if(editor) { 
                if (store.active!.activeFile) {
                    setTimeout(()=>{
                        // console.log("editor!.getValue()", editor!.getValue());

                        store.storeFileState(store.active!.activeFile!, editor);
                    }, 1000)
                }
            }
        }); 
    }
}

const resetEditor = () => {
    destroyEditor()
    setEditor();
}

const destroyEditor = () => {
    if (editor) {
        editor.dispose();
        editor = null;
    }
}

const loadOpenContent = async () => {
    try {
        await store.openFile(store.active!.activeFile!, editor)
    } catch (error) {}

    setEditor()
}

onMounted(() => {
    // Custom VS-Dark Overrides for a dark slate look
    monaco.editor.defineTheme('tailwind-slate', {
        base: 'vs-dark',
        inherit: true,
        rules: [
            { token: 'comment', foreground: '64748b', fontStyle: 'italic' },
            // { token: 'keyword', foreground: '818cf8', fontWeight: 'bold' },
            { token: 'string', foreground: '34d399' },
            { token: 'number', foreground: 'fbbf24' },
        ],
        colors: {
            'editor.background': '#020617', // slate-950
            'editor.lineHighlightBackground': '#0f172a60',
            'editorLineNumber.foreground': '#475569',
            'editorLineNumber.activeForeground': '#cbd5e1',
        }
    });
    
    setEditor()
});

onBeforeUnmount(() => {
    if (editor) editor.dispose();
});

// Watch for external value changes to sync the editor
watch(() => props.modelValue, (newValue) => {
    if (editor && newValue !== editor.getValue()) {
        editor.setValue(newValue);
    }
});

// Watch for language changes dynamically
watch(() => props.language, (newLang) => {
    if (editor) {
        monaco.editor.setModelLanguage(editor!.getModel() as monaco.editor.ITextModel, newLang);
    }
});

watch(() => ((store.active && store.active.activeFile)? store.active!.activeFile.id || null: null), (v2, v1) => {
    // console.log("watch", 3); 

    // Safely update monaco model definitions
    if (v2 != v1) {
        // console.log("watch:loaded", v2);
        
        if(!v2){
            destroyEditor();
        } else {
            loadOpenContent();
        }
    }
}, {deep: true})

window.addEventListener("resize", () => {
    resetEditor();
})
</script>

<style scoped>
.editor-container {
    width: 100%;
    height: 500px;
    /* Adjust height based on your layout */
    border: 1px solid #ccc;
}
</style>