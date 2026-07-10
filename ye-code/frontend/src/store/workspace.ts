import { defineStore } from 'pinia';
import { computed, nextTick, reactive, ref, type Ref } from 'vue';
import * as monaco from 'monaco-editor';
import { ReadDirectory, ReadFile, ListWorkspaceSessions, GetSession, SaveFile, CreateDirectory, MoveFile, ListIcons } from '@wails/yekonga-builder/service.ts';
import YekongaDatabase from '@/scripts/database';

const WORKSHOP_TABLE = "workshops"

export type VIEW_MODE = "EDITOR" | "AGENT" | "WORKSPACE" | "NONE";
export type EXPLORE_VIEW = "FILES" | "SEARCH" | "GIT";
export type CUSTOMIZATION_VIEW = "AGENT" | "DATA_SCHEMA" | "PROJECT";
export type DragSide = 'left' | 'right' | null;
export type SIDE_VIEW = 'editor' | 'agent';

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

export interface SideState {
    minLeft: number;
    maxLeft: number;
    minRight: number;
    maxRight: number; 

    leftWidth: number;
    rightWidth: number;
    isDragging: boolean;

    dragSide: DragSide;
    startX: number;
    startLeft: number;
    startRight: number;
}

export interface Workspace {
    id: string,
    name: string,
    path: string;
    sessionId: string | null;
    viewMode: VIEW_MODE;
    exportView: EXPLORE_VIEW;
    customizedView: CUSTOMIZATION_VIEW;
    workspaceFiles: FileNode[];
    changedFiles: FileNode[];
    sessions: Session[];
    openTabs: FileNode[];
    activeFile: FileNode | null;
    viewStates: Record<string, monaco.editor.ICodeEditorViewState | null>;
    isPinned?: boolean;
    lastOpened: Date;
    /** When set, the editor will navigate to this line after opening the next file */
    pendingLineNumber: number | null;

    sideView: SIDE_VIEW;
    sideStates: Record<SIDE_VIEW, SideState>;
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
        const isSimpleIcon = ref<boolean>(false);

        // ── Active session tracking ──────────────────────────────────────────
        const activeSessionMessages = ref<{ role: string; content: string }[]>([]);
        const activeSideView = computed<SideState>(() => {
            return active.value!.sideStates[active.value!.sideView];
        });

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

        const setViewMode = async (value: VIEW_MODE | null ) => {
            if (!active.value) return;
            if(value === null) {
                active.value = null;
            } else {
                active.value.viewMode = value;
    
                if(value === "AGENT") {
                    active.value.sideView = "agent";
                } else if(value === "EDITOR") {
                    active.value.sideView = "editor";
                }

                await saveLocally();
            }

        }

        const setSimpleIcon = async (value: boolean) => {
            isSimpleIcon.value = value;
        }

        const setViewExplore = async (value: EXPLORE_VIEW) => {
            if (!active.value) return;
            active.value.exportView = value;
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
                        sessionId: null,
                        exportView: "FILES",
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
                        pendingLineNumber: null,
                        sideView: "editor",
                        sideStates: {
                            "agent": {
                                minLeft: 200,
                                maxLeft: 500,
                                minRight: 200,
                                maxRight: 600,

                                leftWidth: 256,
                                rightWidth: 288,
                                isDragging: false,

                                dragSide: null,
                                startX: 0,
                                startLeft: 0,
                                startRight: 0
                            },
                            "editor": {
                                minLeft: 200,
                                maxLeft: 500,
                                minRight: 200,
                                maxRight: 600,

                                leftWidth: 256,
                                rightWidth: 288,
                                isDragging: false,

                                dragSide: null,
                                startX: 0,
                                startLeft: 0,
                                startRight: 0
                            }
                        }
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
                if(Array.isArray(sessions)) {
                    // Sort sessions by last_updated in descending order
                    sessions.sort((a, b) => new Date(b.last_updated).getTime() - new Date(a.last_updated).getTime());

                    active.value.sessions = sessions;
                } else {
                    active.value.sessions = [];
                }

                active.value.sessions = sessions || [];
                saveLocally();
            } catch (error) {
                console.error("Failed to fetch sessions:", error);
            }
        }

        // ── Load session history into the chat panel ──────────────────────
        const selectSession = async (sessionId: string) => {
            active.value!.sessionId = sessionId;
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
            active.value!.sessionId = id;
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
         * Opens a file and captures the previous file's layout state snapshot.
         * If the workspace has a pendingLineNumber, the editor will scroll to that line after opening.
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
                
                // Scroll to pending line number (set by search results or other navigation)
                const pendingLine = active.value!.pendingLineNumber;
                if (pendingLine !== null && pendingLine > 0) {
                    editorInstance.revealLineInCenter(pendingLine);
                    editorInstance.setPosition({ lineNumber: pendingLine, column: 1 });
                    editorInstance.focus();
                    active.value!.pendingLineNumber = null;
                } else {
                    editorInstance.focus();
                    if(savedState) {
                        editorInstance.restoreViewState(savedState)
                    }
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


        const startDrag = async (side: DragSide, e: MouseEvent) => {
            active.value!.sideStates[active.value!.sideView].isDragging = true;
            active.value!.sideStates[active.value!.sideView].dragSide = side;
            active.value!.sideStates[active.value!.sideView].startX = e.clientX;
            active.value!.sideStates[active.value!.sideView].startLeft = active.value!.sideStates[active.value!.sideView].leftWidth;
            active.value!.sideStates[active.value!.sideView].startRight = active.value!.sideStates[active.value!.sideView].rightWidth;

            await saveLocally();
        }

        const onMouseMove = async (e: MouseEvent) => {
            if (!active.value!.sideStates[active.value!.sideView].isDragging || !active.value!.sideStates[active.value!.sideView].dragSide) return

            const dx = e.clientX - active.value!.sideStates[active.value!.sideView].startX

            if (active.value!.sideStates[active.value!.sideView].dragSide === 'left') {
                const newWidth = Math.min(active.value!.sideStates[active.value!.sideView].maxLeft, Math.max(active.value!.sideStates[active.value!.sideView].minLeft, active.value!.sideStates[active.value!.sideView].startLeft + dx))
                active.value!.sideStates[active.value!.sideView].leftWidth = newWidth
            } else if (active.value!.sideStates[active.value!.sideView].dragSide === 'right') {
                const newWidth = Math.min(active.value!.sideStates[active.value!.sideView].maxRight, Math.max(active.value!.sideStates[active.value!.sideView].minRight, active.value!.sideStates[active.value!.sideView].startRight - dx))
                active.value!.sideStates[active.value!.sideView].rightWidth = newWidth
            }

            window.dispatchEvent(new Event("resize"));
            await saveLocally();
        }

        const onMouseUp = async () => {
            if (active.value!.sideStates[active.value!.sideView].isDragging) {
                active.value!.sideStates[active.value!.sideView].isDragging = false
                active.value!.sideStates[active.value!.sideView].dragSide = null
            }

            await saveLocally();
        }
    
        return {
            icons,
            isSimpleIcon,
            active,
            activePath,
            activeSideView,
            activeSessionMessages,
            workspaces,

            startDrag,
            onMouseMove,
            onMouseUp,

            loadIcons,
            setSimpleIcon,
            openFile,
            setViewMode,
            setViewExplore,
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

