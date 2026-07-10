<template>
    <div @contextmenu="onContextMenu($event)" class="h-full w-full">
        <input v-if="isEdit && isEditor" ref="editableElem" class="h-full bg-slate-800 flex-1 w-full" v-model="node.name" @blur="onBlur"/>
        <div v-else class="h-full flex items-center flex-1">
            <span class="w-full inline-block truncate leading-[1.2]">{{ node.name }}</span>
        </div>
    </div>

    <!--this is component mode of context-menu-->
    <context-menu v-model:show="showMenu" :options="optionsComponent">
        <context-menu-item v-if="node.type === 'directory'" label="New File" icon="ye-file-plus" @click="contextMenuItemClicked('new_file')" />
        <context-menu-item v-if="node.type === 'directory'" label="New Folder" icon="ye-folder-plus" @click="contextMenuItemClicked('new_folder')" />
        <!-- <context-menu-divider v-if="node.type === 'directory'" /> -->
        <context-menu-item label="Rename" icon="ye-text" @click="contextMenuItemClicked('rename')" />
        <context-menu-item label="Delete" icon="ye-trash" @click="contextMenuItemClicked('delete')" />
    </context-menu>
</template>

<script setup lang="ts">
import { useWorkspace, type FileNode } from '@/store/workspace';
import type { MenuOptions } from '@imengyu/vue3-context-menu';
import { DeleteFile, RenameFile } from '@wails/yekonga-builder/service.ts';
import { nextTick, reactive, ref } from 'vue';

const emit = defineEmits(["update:modelValue", "change"])
const props = defineProps<{
    modelValue: FileNode;
    isEditor: boolean;
}>()

const store = useWorkspace();
const node = ref<FileNode>(props.modelValue);
const isEdit = ref<boolean>(false);
const editableElem = ref<any>(null);
const count = ref<number>(0);
const showMenu = ref<boolean>(false);
const optionsComponent = reactive<MenuOptions>({
    iconFontClass: 'ye',
    customClass: 'bg-slate-800',
    zIndex: 99,
    minWidth: 150,
    x: 500,
    y: 200
})

const onBlur = () => {
    isEdit.value = false;
    RenameFile(node.value.name, node.value.path!)
}

const doubleClick = () => {
    count.value += 1;

    if(count.value >= 2) {
        isEdit.value = true;
    }

    setTimeout(()=>{
        if(isEdit.value && editableElem.value) {
            editableElem.value.focus();
        }
    }, 200)

    setTimeout(()=>{
        count.value = 0;
    }, 2000)
}

const onContextMenu = (e : MouseEvent) => {
    e.preventDefault();

    // Set the mouse position
    optionsComponent.x = e.x;
    optionsComponent.y = e.y;

    // Show menu
    showMenu.value = true;
}

const contextMenuItemClicked = async (name: string) =>{
    if(name == 'rename') {
        isEdit.value = true;

        setTimeout(()=>{
            if(isEdit.value && editableElem.value) {
                editableElem.value.focus();
            }
        }, 200)
    } else if(name == 'delete') {
        await DeleteFile(node.value.path!);
        store.fetchWorkspaceFiles();
    } else if(name == 'new_file') {
        const fileName = window.prompt("New File Name:");
        if (fileName) store.createNewFile(node.value.path!, fileName);
    } else if(name == 'new_folder') {
        const folderName = window.prompt("New Folder Name:");
        if (folderName) store.createNewFolder(node.value.path!, folderName);
    }
}
</script>