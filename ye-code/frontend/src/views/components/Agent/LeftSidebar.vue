<template>
    <div class="flex flex-col h-full overflow-y-auto">
        <Sessions></Sessions>
        <!-- <Changes class=" border-t border-slate-800/60"></Changes> -->
        <Customization class=" border-t border-slate-800/60"></Customization>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Sessions from './parts/Sessions.vue';
import Changes from './parts/Changes.vue';
import Customization from './parts/Customization.vue';

</script>