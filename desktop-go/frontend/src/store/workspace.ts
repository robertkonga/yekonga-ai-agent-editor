import { defineStore } from 'pinia';
import { computed, reactive, ref } from 'vue';
import * as monaco from 'monaco-editor';
import { ReadDirectory, ReadFile } from '@wails/go/main/App';
import YekongaDatabase from '@/scripts/database';

const WORKSHOP_TABLE = "workshops"

export interface FileNode {
    id: string;
    name: string;
    path?: string;
    type: 'file' | 'directory';
    lang?: string; // Optional: only files have languages
    content?: string; // Optional: only files have content strings
    expanded?: boolean; // Optional: only directories use this for UI state
    children?: FileNode[]; // Optional: only directories contain children arrays
    extension?: string;
}

export interface Workspace {
    id: string,
    name: string,
    path: string;
    workspaceFiles: FileNode[];
    openTabs: FileNode[];
    activeFile: FileNode | null;
    viewStates: Record<string, monaco.editor.ICodeEditorViewState | null>;
    isPinned?: boolean;
    lastOpened: Date;
}

export const useWorkspaceStore = (name: string) => {
    return defineStore(name, () => {
        const db = new YekongaDatabase({
            version: 1,
            tables: [
                { name: WORKSHOP_TABLE, key: "path" },
            ]
        })

        const workspaces = reactive<Record<string, Workspace>>({})
        const activePath = ref<string | null>(null);
        const active = computed<Workspace | null>(()=>{
            if(activePath.value && workspaces[activePath.value]) {
                return workspaces[activePath.value]
            }

            return null;
        })

        const saveLocally = async () => {
            await db.table(WORKSHOP_TABLE).create(window.copy(active.value));
        }

        const loadWorkshops = async () => {
            var res = await db.table(WORKSHOP_TABLE).find() as Workspace[];

            for (const e of res) {
                try {
                    workspaces[e.path] = e;
                } catch (error) {
                    console.log(error);
                }
            }
        }

        const fetchWorkspaceFiles = async () => {
            try {
                if(activePath.value) {
                    let list = await ReadDirectory(activePath.value);
            
                    if (list && Array.isArray(list.children)) {
                        active.value!.workspaceFiles = list.children as FileNode[];
                    }

                    saveLocally()
                }
            } catch (error) {
                
            }
        }

        const setActiveFile = (fileNode: FileNode) => {
            // Set as active file
            active.value!.activeFile = fileNode;
        }

        /**
         * Opens a file and captures the previous file's layout state snapshot
         */
        const openFile = async (fileNode: FileNode, editorInstance: monaco.editor.ICodeEditor | null ): Promise<string | monaco.editor.ICodeEditorViewState | void> => {
            if (!fileNode || fileNode.type === 'directory') return;

            // Add to tab list if it isn't already open
            const exists = active.value!.openTabs.some(tab => tab.id === fileNode.id);
            if (exists) {
                for (let i = 0; i < active.value!.openTabs.length; i++) {
                    let openTab = active.value!.openTabs[i];

                    if(openTab.id === fileNode.id) {
                        fileNode.content == openTab.content; break;
                    }
                }
            } else {
                try {
                    let content = await ReadFile(fileNode.path || "");
                    fileNode.content = content;
                } catch (error: any) {}
                
                active.value!.openTabs.push(fileNode);
            }

            if(editorInstance) {
                let savedState = active.value!.viewStates[fileNode.id];

                editorInstance.setValue(fileNode.content || "");
                monaco.editor.setModelLanguage(editorInstance!.getModel() as monaco.editor.ITextModel, fileNode.lang || "");
                
                editorInstance.focus();
                if(savedState) {
                    editorInstance.restoreViewState(savedState)
                }
            }

            saveLocally();
        }
    
        /**
         * Opens a file and captures the previous file's layout state snapshot
         */
        const storeFileState = (fileNode: FileNode, editorInstance: monaco.editor.ICodeEditor | null): void => {
            if (!fileNode || fileNode.type === 'directory') return;

            // Cache the previous file's view state before navigating away
            if (active.value!.activeFile && editorInstance) {
                fileNode.content = editorInstance.getValue();

                active.value!.viewStates[fileNode.id] = editorInstance.saveViewState();
            }

            // Add to tab list if it isn't already open
            const exists = active.value!.openTabs.some(tab => tab.id === fileNode.id);
            if (!exists) {
                active.value!.openTabs.push(fileNode);
            } else {
                for (let i = 0; i < active.value!.openTabs.length; i++) {
                    let id = active.value!.openTabs[i].id;

                    if(id === fileNode.id) {
                        active.value!.openTabs[i] = fileNode;
                        break;
                    }
                }
            }

            saveLocally();
        }
    
        /**
         * Restores the editor layout scroll and cursor positions
         */
        const restoreFileState = (fileId: string, editorInstance: monaco.editor.ICodeEditor | null): void => {
            if (!editorInstance) return;
    
            const savedState = active.value!.viewStates[fileId];
            if (savedState) {
                editorInstance.restoreViewState(savedState);
            }
            editorInstance.focus();
        }
    
        /**
         * Closes an open tab and flushes its state footprint from memory
         */
        const closeTab = (fileId: string, editorInstance: monaco.editor.ICodeEditor | null): void => {
            active.value!.openTabs = active.value!.openTabs.filter(tab => tab.id !== fileId);
            active.value!.viewStates[fileId] = null;
            delete active.value!.viewStates[fileId];
    
            if (active.value!.activeFile?.id === fileId) {
                if (active.value!.openTabs.length > 0) {
                    storeFileState(active.value!.openTabs[active.value!.openTabs.length - 1], editorInstance);
                } else {
                    active.value!.activeFile = null;
                }
            }
        }

        const openWorkshop = async (path: string | null) => {
            activePath.value = path;

            if(path) {
                let id = await generateID(path);
                let name = path.split("/").pop() || "";

                if(!workspaces[path]) {
                    workspaces[path] = {
                        id: id,
                        name: name,
                        path: path,
                        workspaceFiles: [],
                        openTabs: [],
                        activeFile: null,
                        isPinned: false,
                        viewStates: {},
                        lastOpened: new Date(),
                    }
                } else {
                    workspaces[path].id = id;
                    workspaces[path].name = name;
                    workspaces[path].lastOpened = new Date();
                }
            }
            
            await saveLocally();
            await fetchWorkspaceFiles();
        }
    
        return {
            workspaces,
            active,
            activePath,
            openFile,
            setActiveFile,
            closeTab,
            openWorkshop,
            storeFileState,
            restoreFileState,
            fetchWorkspaceFiles,
            loadWorkshops,
        };
    });
}

export const generateID = async function(absolutePath: string): Promise<string> {
    const standardizedPath = absolutePath.replace(/\\/g, '/')
    const encoded = new TextEncoder().encode(standardizedPath)
    const hashBuffer = await crypto.subtle.digest('SHA-256', encoded)
    const hex = Array.from(new Uint8Array(hashBuffer))
        .map((b) => b.toString(16).padStart(2, '0'))
        .join('')
    return hex.slice(0, 16)
}

export const useWorkspace = () => {
    let workspace = useWorkspaceStore("workspace")();
    workspace.loadWorkshops();

    return workspace;
}

const testWorkshopFiles:FileNode[] = [
    {
        id: 'src-dir',
        name: 'src',
        type: 'directory',
        expanded: true,
        children: [
            {
                id: 'src-dir-1',
                name: 'assets',
                type: 'directory',
                expanded: true,
                children: [
                    { id: 'app-js-1', name: 'generator.js', type: 'file', lang: 'javascript', content: `// Core logic entrypoint\nexport function initialize() {\n  console.log("App loaded smoothly.");\n}` },
                    { id: 'styles-css-1', name: 'generator.css', type: 'file', lang: 'css', content: `/* Core workspace presentation styling */\nbody {\n  background-color: #020617;\n  color: #f8fafc;\n}` }
                ]
            },
            { id: 'app-js', name: 'app.js', type: 'file', lang: 'javascript', content: `// Core logic entrypoint\nexport function initialize() {\n  console.log("App loaded smoothly.");\n}` },
            { id: 'styles-css', name: 'global.css', type: 'file', lang: 'css', content: `/* Core workspace presentation styling */\nbody {\n  background-color: #020617;\n  color: #f8fafc;\n}` }
        ]
    },
    {
        id: 'package-json',
        name: 'package.json',
        type: 'file',
        lang: 'json',
        content: `{\n  "name": "vue3-ai-editor",\n  "version": "1.0.0",\n  "private": true\n}`
    }
]