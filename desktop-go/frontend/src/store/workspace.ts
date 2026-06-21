import { defineStore } from 'pinia';
import { computed, nextTick, reactive, ref } from 'vue';
import * as monaco from 'monaco-editor';
import { ReadDirectory, ReadFile, ListWorkspaceSessions, GetSession, SaveFile, CreateDirectory, MoveFile, ListIcons } from '@wails/go/main/App';
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

export interface Session {
    id: string;
    provider: string;
    history: { role: string; content: any }[];
    last_updated: string;
    workspace: string;
}

export interface Workspace {
    id: string,
    name: string,
    path: string;
    viewMode: VIEW_MODE;
    customizedView: CUSTOMIZATION_VIEW;
    workspaceFiles: FileNode[];
    changedFiles: FileNode[];
    sessions: Session[];
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
        const active = ref<Workspace | null>(null);
        const icons: any[] = [];

        // ── Active session tracking ──────────────────────────────────────────
        const activeSessionId = ref<string | null>(null);
        const activeSessionMessages = ref<{ role: string; content: string }[]>([]);

        const saveLocally = async () => {
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

        const setViewMode = async (value: VIEW_MODE) => {
            if (!active.value) return;
            active.value.viewMode = value;
            await saveLocally();
        }

        const loadWorkshops = async () => {
            let res = window.fetchLocalData();

            for (const key in res) {
                if (!Object.hasOwn(res, key)) continue;
                
                workspaces[key] = res[key];
                if (!workspaces[key].changedFiles) workspaces[key].changedFiles = [];
                if (!workspaces[key].sessions) workspaces[key].sessions = [];
            }
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
            if(path) {
                let id = await generateID(path);
                let name = getFileName(path) || "";

                if(!workspaces[path]) {
                    workspaces[path] = {
                        id: id,
                        name: name,
                        path: path,
                        workspaceFiles: [],
                        changedFiles: [],
                        sessions: [],
                        openTabs: [],
                        activeFile: null,
                        isPinned: false,
                        viewStates: {},
                        lastOpened: new Date(),
                        customizedView: "AGENT",
                        viewMode: "EDITOR",
                    }
                } else {
                    workspaces[path].lastOpened = new Date();

                    if (!workspaces[path].id) workspaces[path].id = id;
                    if (!workspaces[path].name) workspaces[path].name = name;
                    if (!workspaces[path].changedFiles) workspaces[path].changedFiles = [];
                    if (!workspaces[path].sessions) workspaces[path].sessions = [];
                }

                activePath.value = path;
                active.value = workspaces[path];

                if(active.value && !["EDITOR","AGENT"].includes(active.value!.viewMode)) {
                    active.value!.viewMode = "AGENT";
                }

                // Load sessions for this workspace
                await fetchSessions();
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

        const trackChange = (file: FileNode) => {
            if (!active.value) return;
            const exists = active.value.changedFiles.some(f => f.path === file.path);
            if (!exists) {
                active.value.changedFiles.push(window.copy(file));
            }
            saveLocally();
        }

        const clearChange = (path: string) => {
            if (!active.value) return;
            active.value.changedFiles = active.value.changedFiles.filter(f => f.path !== path);
            saveLocally();
        }

        const createNewFile = async (parentPath: string, name: string) => {
            if (!active.value) return;
            const path = parentPath + "/" + name;
            try {
                await SaveFile("", path);
                await fetchWorkspaceFiles();
            } catch (error) {
                console.error("Failed to create file:", error);
            }
        }

        const createNewFolder = async (parentPath: string, name: string) => {
            if (!active.value) return;
            const path = parentPath + "/" + name;
            try {
                await CreateDirectory(path);
                await fetchWorkspaceFiles();
            } catch (error) {
                console.error("Failed to create folder:", error);
            }
        }

        const moveFile = async (sourcePath: string, destPath: string) => {
            if (!active.value) return;
            const fileName = sourcePath.split("/").pop();
            const targetPath = destPath + "/" + fileName;
            try {
                await MoveFile(sourcePath, targetPath);
                await fetchWorkspaceFiles();
            } catch (error) {
                console.error("Failed to move file:", error);
            }
        }

        const fetchSessions = async () => {
            if (!active.value) return;
            try {
                // Use workspace-scoped session listing
                const sessions = await ListWorkspaceSessions(activePath.value || "");
                active.value.sessions = sessions || [];
                saveLocally();
            } catch (error) {
                console.error("Failed to fetch sessions:", error);
            }
        }

        // ── Load session history into the chat panel ──────────────────────
        const selectSession = async (sessionId: string) => {
            activeSessionId.value = sessionId;
            activeSessionMessages.value = [];

            try {
                const session = await GetSession(sessionId);
                if (session && session.history) {
                    // Convert the history from the backend into plain message objects
                    const messages: { role: string; content: string }[] = [];
                    for (const msg of session.history) {
                        if (typeof msg.content === 'string') {
                            messages.push({ role: msg.role, content: msg.content });
                        } else if (Array.isArray(msg.content)) {
                            // ContentBlocks – extract text parts only
                            const textParts = msg.content
                                .filter((b: any) => b.type === 'text' && b.content)
                                .map((b: any) => b.content);
                            if (textParts.length > 0) {
                                messages.push({ role: msg.role, content: textParts.join('\n') });
                            }
                        }
                    }
                    activeSessionMessages.value = messages;
                }
            } catch (error) {
                console.error("Failed to load session:", error);
            }
        }

        // ── Create a new session (resets the chat) ──────────────────────
        const createNewSession = () => {
            const id = "session-" + Math.random().toString(36).substring(2, 9);
            activeSessionId.value = id;
            activeSessionMessages.value = [];
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

        const loadIcons = async () => {
            let list = await ListIcons()
            console.log(list)
        }
    
        return {
            icons,
            active,
            activePath,
            activeSessionId,
            activeSessionMessages,
            workspaces,
            loadIcons,
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
            trackChange,
            clearChange,
            fetchSessions,
            selectSession,
            createNewSession,
            createNewFile,
            createNewFolder,
            moveFile,
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

function getFileName(filePath: string) {
  // Normalize: split on either slash type, take the last non-empty segment
  const parts = filePath.split(/[\\/]/);
  return parts.pop() || parts.pop(); // handles trailing slash case
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