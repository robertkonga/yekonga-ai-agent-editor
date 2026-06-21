<template>
    <div
        class="h-50 bg-black/90 flex flex-col font-mono text-sm text-[#00ff00] shadow-inner border-[#30363d] border-t relative overflow-hidden">
        <div
            class="px-4 py-2 bg-[#161b22] border-b border-[#30363d] text-[#c9d1d9] text-xs font-semibold tracking-wider flex items-center gap-2">
            <svg class="w-4 h-4 text-[#8b949e]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M8 9l3 3-3 3m5 0h3M4 15h16a1 1 0 001-1V5a1 1 0 00-1-1H4a1 1 0 00-1 1v9a1 1 0 001 1z"></path>
            </svg>
            <span>Terminal</span>
        </div>
        <div class="flex-1 p-3 overflow-y-auto w-full max-w-full break-all leading-relaxed bg-[#0a0a0a]">
            <div class="opacity-80">root@sandbox:/workspace# python3 main.py</div>
            <div class="text-white mt-1">Starting computation...</div>
            <div class="text-white">Factorial of 5 is 120</div>
            <div class="opacity-80 mt-2">root@sandbox:/workspace# <span
                    class="animate-pulse bg-white/20 ml-1 w-2 h-4 inline-block align-middle rounded-sm"></span></div>
        </div>
    </div>
</template>

<script setup lang="ts">
// Simple terminal output view
</script>
