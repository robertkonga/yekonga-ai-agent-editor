<template>
    <div class="flex flex-col h-full">
        <Sessions></Sessions>
        <Customization></Customization>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Sessions from './parts/Sessions.vue';
import Customization from './parts/Customization.vue';

</script>