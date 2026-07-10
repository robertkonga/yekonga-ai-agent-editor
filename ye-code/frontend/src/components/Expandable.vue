<template>
    <div class="pl-4 pr-2 pt-3">
        <div class="flex items-center justify-between mb-3 cursor-pointer">
            <h2 @click="show = !show" class="h-full flex-1 text-sm font-semibold text-white">{{ props.title }}</h2>
            
            <slot v-if="!!(slots.extraHeader)" name="extraHeader" v-bind:show="show"></slot>
            <div v-else  @click="show = !show">
                <span class="size-4 flex items-center justify-center text-gray-500 transition-transform" :class="{'rotate-90':show}">
                    <i class="ye ye-angle-right"></i>
                </span>
            </div>
        </div>
        <div v-if="show" class="block" >
            <slot name="default" v-bind:show="show"></slot>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

const props = defineProps<{
    title: string;
    show?: boolean;
}>();

const slots = defineSlots<{
    extraHeader(props: { show: boolean }): any
    default(props: { show: boolean }): any
}>();

const show = ref<boolean>(props.show ? true: false);
</script>