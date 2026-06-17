
<template>
    <div class="cd-root" :class="{ 'cd-root--disabled': disabled }">
        <!-- Label -->
        <label v-if="label" :for="id" class="cd-label">{{ label }}</label>


        <!-- Trigger -->
        <template v-if="(!!slots.button)">
            <div :id="id" ref="triggerRef" type="button" role="combobox" 
                :aria-expanded="open" 
                :aria-haspopup="'listbox'"
                :aria-disabled="disabled" 
                :disabled="disabled" 
                @click="toggleMenu"
                @keydown="onTriggerKeydown">
                <slot name="button" v-bind:value="modelValue" v-bind:open="open" v-bind:label="selectedOption?.label ?? placeholder"></slot>
            </div>
        </template>
        <template v-else>
            <button :id="id" ref="triggerRef" type="button" role="combobox" :aria-expanded="open" :aria-haspopup="'listbox'"
                :aria-disabled="disabled" :disabled="disabled" class="cd-trigger"
                :class="[sizeClasses.trigger, { 'cd-trigger--open': open }]" @click="toggleMenu"
                @keydown="onTriggerKeydown">

                <!-- Selected value -->
                <span class="cd-trigger__value">
                    <span v-if="selectedOption?.icon" class="cd-icon" aria-hidden="true">
                        {{ selectedOption.icon }}
                    </span>
                    <span :class="selectedOption ? 'cd-trigger__label' : 'cd-trigger__placeholder'">
                        {{ selectedOption?.label ?? placeholder }}
                    </span>
                </span>
    
                <!-- Clear -->
                <span v-if="clearable && modelValue != null" class="cd-clear" aria-label="Clear selection" role="button"
                    tabindex="-1" @click="clear">
                    <svg width="10" height="10" viewBox="0 0 10 10" fill="none" aria-hidden="true">
                        <path d="M1 1L9 9M9 1L1 9" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
                    </svg>
                </span>
    
                <!-- Chevron -->
                <svg class="cd-chevron" :class="{ 'cd-chevron--open': open }" width="10" height="10" viewBox="0 0 10 10"
                    fill="none" aria-hidden="true">
                    <path d="M2 3.5L5 6.5L8 3.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"
                        stroke-linejoin="round" />
                </svg>
            </button>
        </template>

        <!-- Dropdown menu -->
        <Transition name="cd-menu">
            <div v-if="open" ref="menuRef" role="listbox" :aria-label="label ?? placeholder" class="cd-menu" :class="[
                sizeClasses.menu,
                resolvedPlacement === 'top' ? 'cd-menu--top' : 'cd-menu--bottom'
            ]" :style="{ maxHeight: `${maxHeight}px` }" @keydown="onMenuKeydown">
                <!-- Search -->
                <div v-if="searchable" class="cd-search">
                    <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true"
                        class="cd-search__icon">
                        <circle cx="5" cy="5" r="3.5" stroke="currentColor" stroke-width="1.25" />
                        <path d="M8 8L10.5 10.5" stroke="currentColor" stroke-width="1.25" stroke-linecap="round" />
                    </svg>
                    <input ref="searchRef" v-model="searchQuery" type="text" class="cd-search__input"
                        placeholder="Filter…" autocomplete="off" spellcheck="false" @keydown="onSearchKeydown" />
                </div>

                <div class="cd-menu__scroll">
                    <!-- Empty state -->
                    <div v-if="flatFiltered.length === 0" class="cd-empty">
                        No options match
                    </div>

                    <!-- Grouped options -->
                    <template v-for="[group, opts] in grouped" :key="group ?? '__ungrouped'">
                        <div v-if="group" class="cd-group-label">{{ group }}</div>
                        <div v-for="(opt, idx) in opts" :key="String(opt.value)"
                            :data-highlighted="flatFiltered.indexOf(opt) === highlightedIndex ? '' : undefined"
                            role="option" :aria-selected="opt.value === modelValue" :aria-disabled="opt.disabled"
                            class="cd-option" :class="[
                                sizeClasses.item,
                                {
                                    'cd-option--selected': opt.value === modelValue,
                                    'cd-option--highlighted': flatFiltered.indexOf(opt) === highlightedIndex,
                                    'cd-option--disabled': opt.disabled,
                                }
                            ]" @mouseenter="highlightedIndex = flatFiltered.indexOf(opt)" @mouseleave="highlightedIndex = -1"
                            @mousedown.prevent="select(opt)">
                            <span v-if="opt.icon" class="cd-icon" aria-hidden="true">{{ opt.icon }}</span>
                            <span class="cd-option__body">
                                <span class="cd-option__label">{{ opt.label }}</span>
                                <span v-if="opt.description" class="cd-option__desc">{{ opt.description }}</span>
                            </span>
                            <!-- Check mark -->
                            <svg v-if="opt.value === modelValue" class="cd-check" width="10" height="10"
                                viewBox="0 0 10 10" fill="none" aria-hidden="true">
                                <path d="M1.5 5L4 7.5L8.5 2.5" stroke="currentColor" stroke-width="1.5"
                                    stroke-linecap="round" stroke-linejoin="round" />
                            </svg>
                        </div>
                    </template>
                </div>
            </div>
        </Transition>
    </div>
</template>

<script setup lang="ts" generic="T extends string | number | object">
import { ref, computed, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'

export interface DropdownOption<V = string | number> {
    label: string
    value: V
    description?: string
    icon?: string       // emoji or SVG string
    disabled?: boolean
    group?: string
}

interface Props {
    modelValue?: T | null
    options: DropdownOption[]
    placeholder?: string
    disabled?: boolean
    searchable?: boolean
    clearable?: boolean
    size?: 'sm' | 'md' | 'lg'
    placement?: 'bottom' | 'top' | 'auto'
    maxHeight?: number
    id?: string
    name?: string
    label?: string
}

const props = withDefaults(defineProps<Props>(), {
    modelValue: null,
    placeholder: 'Select…',
    disabled: false,
    searchable: false,
    clearable: false,
    size: 'md',
    placement: 'auto',
    maxHeight: 240,
})

const emit = defineEmits<{
    'update:modelValue': [value: T | null]
    'change': [value: T | null, option: DropdownOption | null]
    'open': []
    'close': []
    'search': [query: string]
}>()


const slots = defineSlots<{
    button(props: { value: any, label: string, open?: boolean }): any
    item(props: { item: T }): any
    default(props: { options: DropdownOption[] }): any
}>();

// ── State ─────────────────────────────────────────────────────────────────

const open = ref(false)
const searchQuery = ref('')
const highlightedIndex = ref(-1)
const triggerRef = ref<HTMLElement | null>(null)
const menuRef = ref<HTMLElement | null>(null)
const searchRef = ref<HTMLInputElement | null>(null)
const resolvedPlacement = ref<'bottom' | 'top'>('bottom')

// ── Computed ──────────────────────────────────────────────────────────────

const selectedOption = computed(() =>
    props.options.find(o => o.value === props.modelValue) ?? null
)

const grouped = computed(() => {
    const enabledByGroup = new Map<string | undefined, DropdownOption[]>()
    const filtered = props.searchable && searchQuery.value
        ? props.options.filter(o =>
            o.label.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
            o.description?.toLowerCase().includes(searchQuery.value.toLowerCase()))
        : props.options

    for (const opt of filtered) {
        const key = opt.group
        if (!enabledByGroup.has(key)) enabledByGroup.set(key, [])
        enabledByGroup.get(key)!.push(opt)
    }
    return enabledByGroup
})

const flatFiltered = computed(() => {
    const all: DropdownOption[] = []
    for (const opts of grouped.value.values()) all.push(...opts)
    return all
})

const sizeClasses = computed(() => ({
    sm: { trigger: 'h-7 px-2.5 text-[11px]', menu: 'text-[11px]', item: 'px-2.5 py-1' },
    md: { trigger: 'h-8 px-3 text-[12px]', menu: 'text-[12px]', item: 'px-3 py-1.5' },
    lg: { trigger: 'h-9 px-3.5 text-[13px]', menu: 'text-[13px]', item: 'px-3.5 py-2' },
}[props.size]))

// ── Placement ─────────────────────────────────────────────────────────────

function calcPlacement() {
    if (props.placement !== 'auto') {
        resolvedPlacement.value = props.placement
        return
    }
    if (!triggerRef.value) return
    const rect = triggerRef.value.getBoundingClientRect()
    const spaceBelow = window.innerHeight - rect.bottom
    resolvedPlacement.value = spaceBelow < props.maxHeight + 8 && rect.top > props.maxHeight + 8
        ? 'top'
        : 'bottom'
}

// ── Open / Close ──────────────────────────────────────────────────────────

async function openMenu() {
    if (props.disabled) return
    calcPlacement()
    open.value = true
    searchQuery.value = ''
    highlightedIndex.value = selectedOption.value
        ? flatFiltered.value.findIndex(o => o.value === props.modelValue)
        : -1
    emit('open')
    await nextTick()
    if (props.searchable) searchRef.value?.focus()
    scrollToHighlighted()
}

function closeMenu() {
    open.value = false
    searchQuery.value = ''
    emit('close')
    triggerRef.value?.focus()
}

function toggleMenu() {
    open.value ? closeMenu() : openMenu()
}

// ── Selection ─────────────────────────────────────────────────────────────

function select(option: DropdownOption) {
    if (option.disabled) return
    emit('update:modelValue', option.value as T)
    emit('change', option.value as T, option)
    closeMenu()
}

function clear(e: MouseEvent) {
    e.stopPropagation()
    emit('update:modelValue', null)
    emit('change', null, null)
}

// ── Keyboard ──────────────────────────────────────────────────────────────

function onTriggerKeydown(e: KeyboardEvent) {
    if (['Enter', ' ', 'ArrowDown', 'ArrowUp'].includes(e.key)) {
        e.preventDefault()
        if (!open.value) openMenu()
        else if (e.key === 'ArrowDown') moveHighlight(1)
        else if (e.key === 'ArrowUp') moveHighlight(-1)
        else if (e.key === 'Enter') confirmHighlighted()
    }
    if (e.key === 'Escape') closeMenu()
}

function onMenuKeydown(e: KeyboardEvent) {
    if (e.key === 'ArrowDown') { e.preventDefault(); moveHighlight(1) }
    else if (e.key === 'ArrowUp') { e.preventDefault(); moveHighlight(-1) }
    else if (e.key === 'Enter') { e.preventDefault(); confirmHighlighted() }
    else if (e.key === 'Escape') closeMenu()
    else if (e.key === 'Tab') closeMenu()
}

function onSearchKeydown(e: KeyboardEvent) {
    if (e.key === 'ArrowDown') { e.preventDefault(); moveHighlight(1) }
    else if (e.key === 'ArrowUp') { e.preventDefault(); moveHighlight(-1) }
    else if (e.key === 'Enter') { e.preventDefault(); confirmHighlighted() }
    else if (e.key === 'Escape') closeMenu()
}

function moveHighlight(dir: 1 | -1) {
    const len = flatFiltered.value.length
    if (!len) return
    let idx = highlightedIndex.value + dir
    // Skip disabled
    let tries = 0
    while (tries < len) {
        if (idx < 0) idx = len - 1
        if (idx >= len) idx = 0
        if (!flatFiltered.value[idx].disabled) break
        idx += dir
        tries++
    }
    highlightedIndex.value = idx
    scrollToHighlighted()
}

function confirmHighlighted() {
    const opt = flatFiltered.value[highlightedIndex.value]
    if (opt) select(opt)
}

function scrollToHighlighted() {
    nextTick(() => {
        const el = menuRef.value?.querySelector('[data-highlighted]') as HTMLElement | null
        el?.scrollIntoView({ block: 'nearest' })
    })
}

// ── Search ────────────────────────────────────────────────────────────────

watch(searchQuery, q => {
    highlightedIndex.value = 0
    emit('search', q)
})

// ── Outside click ─────────────────────────────────────────────────────────

function onOutsideClick(e: MouseEvent) {
    if (!open.value) return
    const target = e.target as Node
    if (!triggerRef.value?.contains(target) && !menuRef.value?.contains(target)) {
        closeMenu()
    }
}

onMounted(() => document.addEventListener('mousedown', onOutsideClick, true))
onBeforeUnmount(() => document.removeEventListener('mousedown', onOutsideClick, true))
</script>


<style scoped>
/* ── Tokens ──────────────────────────────────────────────────────────────── */
.cd-root {
    --cd-bg: #1e1e1e;
    --cd-surface: #252526;
    --cd-border: #3c3c3c;
    --cd-border-hi: #007acc;
    --cd-text: #cccccc;
    --cd-muted: #6b6b6b;
    --cd-hover: #2a2d2e;
    --cd-selected: #094771;
    --cd-selected-text: #ffffff;
    --cd-highlight: #04395e;
    --cd-disabled: #4a4a4a;
    --cd-radius: 3px;
    --cd-shadow: 0 4px 16px rgba(0, 0, 0, .45), 0 1px 4px rgba(0, 0, 0, .3);
    --cd-font: 'Segoe UI', system-ui, sans-serif;
    --cd-mono: 'Cascadia Code', 'JetBrains Mono', monospace;

    position: relative;
    display: inline-flex;
    flex-direction: column;
    gap: 4px;
    font-family: var(--cd-font);
}

.cd-root--disabled {
    opacity: 0.45;
    pointer-events: none;
}

/* ── Label ───────────────────────────────────────────────────────────────── */
.cd-label {
    font-size: 11px;
    color: var(--cd-muted);
    letter-spacing: 0.04em;
    text-transform: uppercase;
    font-weight: 500;
    user-select: none;
}

/* ── Trigger ─────────────────────────────────────────────────────────────── */
.cd-trigger {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    width: 100%;
    border: 1px solid var(--cd-border);
    border-radius: var(--cd-radius);
    background: var(--cd-bg);
    color: var(--cd-text);
    cursor: pointer;
    font-family: var(--cd-font);
    font-size: inherit;
    outline: none;
    transition: border-color 120ms, background 120ms;
    user-select: none;
    white-space: nowrap;
}

.cd-trigger:hover,
.cd-trigger--open {
    border-color: var(--cd-border-hi);
    background: var(--cd-hover);
}

.cd-trigger:focus-visible {
    outline: 1px solid var(--cd-border-hi);
    outline-offset: 1px;
}

.cd-trigger__value {
    display: flex;
    align-items: center;
    gap: 6px;
    flex: 1;
    min-width: 0;
}

.cd-trigger__label {
    flex: 1;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    text-align: left;
    color: var(--cd-text);
}

.cd-trigger__placeholder {
    flex: 1;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    text-align: left;
    color: var(--cd-muted);
}

.cd-chevron {
    flex-shrink: 0;
    color: var(--cd-muted);
    transition: transform 160ms cubic-bezier(.4, 0, .2, 1);
}

.cd-chevron--open {
    transform: rotate(180deg);
}

.cd-clear {
    flex-shrink: 0;
    display: flex;
    align-items: center;
    color: var(--cd-muted);
    border-radius: 2px;
    padding: 1px;
    cursor: pointer;
    transition: color 100ms, background 100ms;
}

.cd-clear:hover {
    color: var(--cd-text);
    background: var(--cd-hover);
}

.cd-icon {
    font-size: 0.9em;
    line-height: 1;
    flex-shrink: 0;
}

/* ── Menu ─────────────────────────────────────────────────────────────────── */
.cd-menu {
    position: absolute;
    left: 0;
    right: 0;
    z-index: 9999;
    background: var(--cd-surface);
    border: 1px solid var(--cd-border);
    border-radius: var(--cd-radius);
    box-shadow: var(--cd-shadow);
    overflow: hidden;
    display: flex;
    flex-direction: column;
    font-family: var(--cd-font);
}

.cd-menu--bottom {
    top: calc(100% + 3px);
}

.cd-menu--top {
    bottom: calc(100% + 3px);
}

.cd-menu__scroll {
    overflow-y: auto;
    overflow-x: hidden;
    flex: 1;
    /* Custom scrollbar */
    scrollbar-width: thin;
    scrollbar-color: var(--cd-border) transparent;
}

.cd-menu__scroll::-webkit-scrollbar {
    width: 6px;
}

.cd-menu__scroll::-webkit-scrollbar-track {
    background: transparent;
}

.cd-menu__scroll::-webkit-scrollbar-thumb {
    background: var(--cd-border);
    border-radius: 3px;
}

/* ── Search ─────────────────────────────────────────────────────────────── */
.cd-search {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 6px 8px;
    border-bottom: 1px solid var(--cd-border);
    background: var(--cd-bg);
}

.cd-search__icon {
    color: var(--cd-muted);
    flex-shrink: 0;
}

.cd-search__input {
    flex: 1;
    background: transparent;
    border: none;
    outline: none;
    color: var(--cd-text);
    font-family: var(--cd-font);
    font-size: inherit;
    line-height: 1;
    padding: 0;
}

.cd-search__input::placeholder {
    color: var(--cd-muted);
}

/* ── Options ─────────────────────────────────────────────────────────────── */
.cd-option {
    display: flex;
    align-items: center;
    gap: 6px;
    cursor: pointer;
    color: var(--cd-text);
    border-left: 2px solid transparent;
    transition: background 80ms;
    user-select: none;
}

.cd-option--highlighted {
    background: var(--cd-hover);
    border-left-color: var(--cd-border-hi);
}

.cd-option--selected {
    background: var(--cd-selected);
    color: var(--cd-selected-text);
}

.cd-option--selected.cd-option--highlighted {
    background: var(--cd-highlight);
}

.cd-option--disabled {
    color: var(--cd-disabled);
    cursor: not-allowed;
}

.cd-option__body {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 1px;
}

.cd-option__label {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.cd-option__desc {
    font-size: 0.85em;
    color: var(--cd-muted);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.cd-option--selected .cd-option__desc {
    color: rgba(255, 255, 255, .65);
}

.cd-check {
    flex-shrink: 0;
    color: var(--cd-selected-text);
}

/* ── Group label ─────────────────────────────────────────────────────────── */
.cd-group-label {
    padding: 5px 10px 3px;
    font-size: 10px;
    font-weight: 600;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: var(--cd-muted);
    border-top: 1px solid var(--cd-border);
}

.cd-group-label:first-child {
    border-top: none;
}

/* ── Empty ───────────────────────────────────────────────────────────────── */
.cd-empty {
    padding: 10px 12px;
    color: var(--cd-muted);
    font-size: 11px;
    text-align: center;
    font-style: italic;
}

/* ── Transition ──────────────────────────────────────────────────────────── */
.cd-menu-enter-active {
    transition: opacity 120ms, transform 120ms cubic-bezier(.2, 0, 0, 1);
}

.cd-menu-leave-active {
    transition: opacity 100ms, transform 100ms cubic-bezier(.4, 0, 1, 1);
}

.cd-menu--bottom.cd-menu-enter-from {
    opacity: 0;
    transform: translateY(-4px);
}

.cd-menu--bottom.cd-menu-leave-to {
    opacity: 0;
    transform: translateY(-4px);
}

.cd-menu--top.cd-menu-enter-from {
    opacity: 0;
    transform: translateY(4px);
}

.cd-menu--top.cd-menu-leave-to {
    opacity: 0;
    transform: translateY(4px);
}

/* ── Reduced motion ──────────────────────────────────────────────────────── */
@media (prefers-reduced-motion: reduce) {

    .cd-menu-enter-active,
    .cd-menu-leave-active,
    .cd-chevron {
        transition: none;
    }
}
</style>