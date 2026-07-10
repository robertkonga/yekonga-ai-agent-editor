<template>
    <teleport to="[main-modal-container]">
        <div v-if="modelValue"
            class="fixed inset-0 z-50 flex overflow-hidden w-full h-full backdrop-blur-xs"
            :class="{
                'justify-end': side === 'right',
                'justify-start': side === 'left',
                'justify-center items-center': side === 'center',
                'p-0 md:p-6': (size !== 'full' && side === 'center'),
            }" role="dialog" aria-modal="true" aria-labelledby="modal-title">

            <!-- Backdrop -->
            <div class="fixed inset-0 bg-slate-900/10 backdrop-blur-xs transition-opacity"
                @click="onToggleModal">
            </div>

            <!-- Modal Panel -->
            <div class="relative z-10 w-full flex flex-col shadow-xl transition-all duration-300 bg-slate-900 overflow-hidden"
                :class="[
                    // Size & Layout Logic
                    side === 'center' ? ' h-full max-h-screen' : 'h-full',
                    (side === 'center' && size != 'full')? 'rounded-none md:rounded-xl':'',
                    size === 'full' ? 'w-full': 'border border-slate-800/90',
                    size === 'lg' ? (side === 'center' ? 'max-w-4xl' : 'max-w-lg') : '',
                    size === 'sm' ? (side === 'center' ? 'max-w-sm' : 'max-w-xs') : '',
                    size === 'sm' ? (side === 'center' ? 'max-w-sm' : 'max-w-xs') : '',
                    (!size || size === 'md') ? (side === 'center' ? 'max-w-lg' : 'max-w-md') : '',

                    // Transparency
                    transparent ? 'bg-transparent shadow-none border-0' : ''
                ]">

                <!-- Header -->
                <div v-if="title"
                    class="flex items-center justify-start lg:justify-between pl-0 lg:pl-5 pr-3 min-h-12 bg-primary-500 dark:bg-primary-500 text-on-primary shrink-0"
                    :class="(side === 'center' && size != 'full') ? 'rounded-t-0 md:rounded-t-xl' : ''"> <!-- bg-primary -->
                    
                    <button type="button"
                        class="flex lg:hidden w-12 items-center justify-center text-on-primary/50 hover:text-on-primary transition-colors focus:outline-none rounded-lg p-1 hover:bg-white/10"
                        @click="onToggleModal">
                        
                        <i class="ye ye-arrow-left text-lg"></i>
                    </button>
                    <h5 id="modal-title" class="text-lg font-semibold tracking-wide">{{ title }}</h5>
                    <button type="button"
                        class="hidden lg:inline-block text-on-primary/50 hover:text-on-primary transition-colors focus:outline-none rounded-lg p-1 hover:bg-white/10"
                        @click="onToggleModal">
                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>

                <!-- Body -->
                <div class="flex-1 overflow-y-auto relative custom-scrollbar" :class="[
                    body ? 'p-6' : '',
                    transparent ? 'bg-transparent' : 'bg-slate-900'
                ]" :style="{ padding }">

                    <!-- Close button for headerless modals -->
                    <template v-if="!title && showCloseButton">
                        <div v-if="hasCloseSlot" @click="onToggleModal"
                            class="absolute top-4 right-4 z-20 cursor-pointer">
                            <slot name="close"></slot>
                        </div>
                        <button v-else type="button"
                            class="absolute top-4 right-4 z-20 text-gray-400 hover:text-gray-500 dark:hover:text-gray-200 transition-colors bg-white/50 dark:bg-black/20 rounded-full p-1"
                            @click="onToggleModal">
                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M6 18L18 6M6 6l12 12" />
                            </svg>
                        </button>
                    </template>

                    <slot name="default" v-bind:close="onToggleModal"></slot>
                </div>

                <!-- Footer -->
                <div v-if="hasFooterSlot"
                    class="px-3 py-3 bg-gray-50 dark:bg-gray-700/30 border-t border-gray-100 dark:border-gray-700 shrink-0 flex items-center justify-end"
                    :class="[
                        footerClass,
                        (side === 'center' && size != 'full') ? 'rounded-b-xl' : '',
                    ]">
                    <slot name="footer" v-bind="{ onClose: onToggleModal }"></slot>
                </div>
            </div>
        </div>
    </teleport>
</template>

<script setup lang="ts">
import { computed, ref, watch, onBeforeUnmount, getCurrentInstance } from 'vue'
import type { CSSProperties } from 'vue'

declare global {
    interface Window {
        closeModal: (value: string) => void;
        __modal_to_close__?: string[];
    }
}

const closeKey = '__modal_to_close__'

if (typeof window.closeModal !== 'function') {
    window.closeModal = (value: string) => {
        if (!Array.isArray(window[closeKey])) {
            window[closeKey] = [];
        }

        if (!window[closeKey].includes(value)) {
            window[closeKey].push(value);
        }
    }
}

interface Props {
    uuid?: string
    modelValue?: boolean
    routePath?: string | null
    modelModifiers?: Record<string, any>
    title?: string | boolean | null
    size?: 'lg' | 'sm' | 'md' | 'full'
    side?: 'right' | 'left' | 'center'
    padding?: string
    headerClass?: string
    bodyClass?: string
    footerClass?: string
    body?: boolean
    transparent?: boolean
    showCloseButton?: boolean
    closeBackground?: boolean
}

const props = withDefaults(defineProps<Props>(), {
    uuid: "",
    modelValue: false,
    routePath: null,
    modelModifiers: () => ({}),
    title: 'Info',
    size: 'lg',
    side: 'center',
    padding: '0rem',
    body: true,
    transparent: false,
    showCloseButton: false,
    closeBackground: false,
})

const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void
}>()

const slots = defineSlots<{
    default?: (props: { close: () => void }) => any
    footer?: (props: { onClose: () => void }) => any
    close?: () => any
}>()

const $this = getCurrentInstance();
const uuid = ref(props.uuid)
const modelValue = ref(props.modelValue)
const routePath = computed(() => props.routePath);
const transparent = ref(props.transparent);
const showCloseButton = ref(props.showCloseButton);
const closeBackground = ref(props.closeBackground);
const timer = ref<any>(null)

const hasFooterSlot = computed(() => !!slots.footer)
const hasCloseSlot = computed(() => !!slots.close)

const setTimer = () => {
    if (timer.value) return

    timer.value = setInterval(() => {
        if (uuid.value && window[closeKey]?.includes(uuid.value)) {
            onClose()
        }
    }, 50)
}

const onClose = () => {
    if (window[closeKey]?.includes(uuid.value)) {
        window[closeKey].splice(window[closeKey].indexOf(uuid.value), 1)
    }
    if (timer.value) {
        clearInterval(timer.value)
        timer.value = null
    }
    emit('update:modelValue', false)
    modelValue.value = false
}

const addClass = (elem: HTMLElement, clazz: string) => {
    const classes = elem.className.split(' ').filter(Boolean)
    if (!classes.includes(clazz)) {
        classes.push(clazz)
    }
    elem.className = classes.join(' ')
}

const removeClass = (elem: HTMLElement, clazz: string) => {
    const classes = elem.className.split(' ').filter(c => c !== clazz)
    elem.className = classes.join(' ')
}

const resetModal = () => {
    const classes = document.body.className.split(' ').filter(Boolean)
    if (!modelValue.value) {
        document.body.className = classes.filter(c => c !== 'modal-open').join(' ')
    } else if (!classes.includes('modal-open')) {
        classes.push('modal-open')
        document.body.className = classes.join(' ')
    }

    const modals = document.querySelectorAll('body > .modal')
    modals.forEach((e, i) => {
        if (i + 1 === modals.length) {
            removeClass(e as HTMLElement, 'd-none')
        } else {
            addClass(e as HTMLElement, 'd-none')
        }
    })
}

const onToggleModal = () => {
    const show = !modelValue.value
    emit('update:modelValue', show)
    modelValue.value = show
}

// Helper function to handle modal opening logic
function handleModalOpen(isOpen: any) {
    if (isOpen) {
        if ($this) {
            console.log('$isDesktop', $this.appContext.config.globalProperties.$isDesktop);

            // Fallback to normal modal behavior if the external function is not available
            modelValue.value = isOpen;
            setTimer();
        } else {
            modelValue.value = isOpen;
            setTimer();
        }
    } else {
        onClose();
    }
}

// Watch 1: Monitor props.modelValue
watch(
    () => props.modelValue,
    (newValue, oldValue) => {
        if (newValue !== oldValue) {
            handleModalOpen(newValue);
        }
    }
);

// Watch 2: Monitor modelValue
watch(
    modelValue,
    (newValue, oldValue) => {
        if (newValue !== oldValue) {
            resetModal();
        }
    }
);

// Watch 3: Monitor route.fullPath and window[closeKey]
watch(
    [
        () => window[closeKey]
    ],
    (
        [newCloseKey],
        [oldCloseKey]
    ) => {
        if ((newCloseKey && uuid.value && newCloseKey.includes(uuid.value))) {
            onClose();
        }
    }
);

if (modelValue.value) setTimer()

onBeforeUnmount(() => {
    onClose()
})
</script>

<style>
/* Global Helpers */
.modal-open {
    overflow: hidden !important;
}

/* Modal Mounting Container */
[main-modal-container] {
    position: relative;
    z-index: 9999;
    /* Ensure high z-index for the container */
}

/* Ensure only active modal is interactable if multiple (though teleport usually handles one stack) */
[main-modal-container]>div[role="dialog"] {
    display: none;
}

[main-modal-container]>div[role="dialog"]:last-child {
    display: flex !important;
}
</style>