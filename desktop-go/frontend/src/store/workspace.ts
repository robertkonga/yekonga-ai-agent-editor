import { defineStore } from 'pinia';
import { computed, nextTick, reactive, ref } from 'vue';
import * as monaco from 'monaco-editor';
import { ReadDirectory, ReadFile } from '@wails/go/main/App';
import YekongaDatabase from '@/scripts/database';

const WORKSHOP_TABLE = "workshops"

export type VIEW_MODE = "EDITOR" | "AGENT" | "WORKSPACE" | "NONE";
export type CUSTOMIZATION_VIEW = "AGENT" | "DATA_SCHEMA" | "PROJECT";

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
    lastUpdate?: Date;
}

export interface Workspace {
    id: string,
    name: string,
    path: string;
    customizedView: CUSTOMIZATION_VIEW;
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

        const viewMode = ref<VIEW_MODE>("EDITOR")
        const workspaces = reactive<Record<string, Workspace>>({})
        const activePath = ref<string | null>(null);
        const active = ref<Workspace | null>(null);

        const saveLocally = async () => {
            // await db.table(WORKSHOP_TABLE).create(window.copy(active.value));
            // =================================================== //

            if(activePath.value && active.value) {
                workspaces[activePath.value] = active.value;
            }

            let data = window.copy(workspaces);
            for (const key in data) {
                if (!Object.hasOwn(data, key)) continue;
                
                data[key].workspaceFiles = [];
            }

            window.savaLocalData(data)
        }

        const setViewMode = (value: VIEW_MODE) => {
            viewMode.value = value;
        }

        const loadWorkshops = async () => {
            let res = window.fetchLocalData();

            for (const key in res) {
                if (!Object.hasOwn(res, key)) continue;
                
                workspaces[key] = res[key];
            }

            // sortWorkshops();
        }

        const sortWorkshops = () => {
            const list = Object.fromEntries(
                Object.entries(window.copy(workspaces) as Record<string, Workspace>)
                    .sort(
                        ([, a], [, b]) => window.dateToNumber(b.lastOpened) - window.dateToNumber(a.lastOpened)
                    )
            );

            for (const key in workspaces) {
                delete workspaces[key];
            }

            for (const key in list) {
                if (!Object.hasOwn(list, key)) continue;
                workspaces[key] = list[key];
            }
        }

        const openWorkshop = async (path: string | null) => {
            viewMode.value = 'EDITOR';
            
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
                        customizedView: "AGENT",
                    }
                } else {
                    // workspaces[path].id = id;
                    // workspaces[path].name = name;
                    workspaces[path].lastOpened = new Date();
                }

                activePath.value = path;
                active.value = workspaces[path];
            } else {
                activePath.value = null;
                active.value = null;
            }
            
            await saveLocally();
            await fetchWorkspaceFiles();
        }

        const removeWorkshop = async (path: string) => {
            delete workspaces[path];
            await saveLocally();
        }

        const changeCustomizationView = async (name: CUSTOMIZATION_VIEW) => {
            if(active.value) {
                active.value!.customizedView = name;
            }
            
            await saveLocally();
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
                console.log(error);
            }
        }

        const setActiveFile = (fileNode: FileNode) => {
            fileNode = window.copy(fileNode);

            if(fileNode.children) {
                fileNode.children = [];
            }

            active.value!.activeFile = fileNode;
        }

        /**
         * Opens a file and captures the previous file's layout state snapshot
         */
        const openFile = async (fileNode: FileNode, editorInstance: monaco.editor.ICodeEditor | null ): Promise<string | void> => {
            fileNode = window.copy(fileNode);
            fileNode.children = [];
            
            if(fileNode.type === 'directory') {
                for (let i = 0; i < active.value!.workspaceFiles.length; i++) {
                    let file = active.value!.workspaceFiles[i];
    
                    if(file.id === fileNode.id) {
                        active.value!.workspaceFiles[i].expanded = !fileNode.expanded; 
                        break;
                    }
                }
            }
            
            if (!fileNode || fileNode.type === 'directory') return;

            // Add to tab list if it isn't already open
            const exists = active.value!.openTabs.some(tab => tab.id === fileNode.id);
            const driveContent = await ReadFile(fileNode.path || "");

            if (exists) {
                for (let i = 0; i < active.value!.openTabs.length; i++) {
                    const file = active.value!.openTabs[i];

                    if(file.id === fileNode.id) {
                        if(window.dateToNumber(driveContent["lastUpdate"]) > window.dateToNumber(fileNode.lastUpdate)) {
                            fileNode.content = driveContent.content; 
                            fileNode.lastUpdate = driveContent.lastUpdate;
                        } else {
                            fileNode.content = file.content;
                            fileNode.lastUpdate = file.lastUpdate;
                        }
                        
                        break;
                    }
                }
            } else {
                try {
                    fileNode.content = driveContent.content; 
                    fileNode.lastUpdate = driveContent.lastUpdate;
                } catch (error: any) {}
                
                active.value!.openTabs.push(fileNode);
            }

            active.value!.activeFile = fileNode;

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

            return driveContent.content;
        }
    
        /**
         * Opens a file and captures the previous file's layout state snapshot
         */
        const storeFileState = (fileNode: FileNode, editorInstance: monaco.editor.ICodeEditor | null): void => {
            if (!fileNode || fileNode.type === 'directory') return; 
            fileNode = window.copy(fileNode);
            fileNode.children = [];

            // Cache the previous file's view state before navigating away
            if (fileNode && editorInstance) {
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
    
            if (active.value!.openTabs.length > 0) {
                storeFileState(active.value!.openTabs[active.value!.openTabs.length - 1], editorInstance);
                if (active.value!.activeFile?.id === fileId) {
                    active.value!.activeFile = active.value!.openTabs[active.value!.openTabs.length - 1];
                }
            } else {
                active.value!.activeFile = null;
            }
            
        }
    
        return {
            viewMode,
            active,
            activePath,
            workspaces,
            openFile,
            setViewMode,
            setActiveFile,
            closeTab,
            saveLocally,
            openWorkshop,
            storeFileState,
            restoreFileState,
            fetchWorkspaceFiles,
            changeCustomizationView,
            loadWorkshops,
            removeWorkshop,
            sortWorkshops,
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