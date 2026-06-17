<template>
    <img
        :src="iconSrc"
        :alt="alt"
        :width="size"
        :height="size"
        class="shrink-0 object-contain select-none"
        draggable="false"
        @error="onError"
    />
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { getFileIcon, getFolderIcon, getLangIcon } from '@/scripts/fileIcons'

interface Props {
    /**
     * File extension with dot (e.g. '.ts'), full filename (e.g. 'Dockerfile'),
     * or language ID (e.g. 'typescript') — resolve priority: filename > ext > lang.
     */
    lang?: string
    /** 'file' | 'directory' */
    type?: 'file' | 'directory'
    /** Whether the folder is currently expanded (only used when type==='directory') */
    extended?: boolean
    /** Icon size in px. Default: 16 */
    size?: number
}

const props = withDefaults(defineProps<Props>(), {
    lang: '',
    type: 'file',
    extended: false,
    size: 16,
})

const errored = ref(false)

const alt = computed(() => `${props.lang || props.type} icon`)

const iconSrc = computed(() => {
    if (errored.value) return '/file-icons/default_file.svg'

    if (props.type === 'directory') {
        // Use the lang/name as the folder name
        return getFolderIcon(props.lang || '', props.extended)
    }

    // lang prop may be a raw extension without dot (e.g. "ts"), with dot (".ts"),
    // a full filename ("main.go", "Dockerfile"), or a language ID ("typescript")
    const raw = props.lang.trim()
    if (!raw) return '/file-icons/default_file.svg'

    // Normalise: if it looks like a bare extension (no slash, no dot prefix, short)
    // treat it as a file ".<ext>" so getFileIcon resolves it correctly.
    const asFilename = raw.startsWith('.') || raw.includes('/') || raw.includes('.')
        ? raw
        : `.${raw}`

    return getFileIcon(asFilename, raw)
})

const onError = () => {
    errored.value = true
}
</script>