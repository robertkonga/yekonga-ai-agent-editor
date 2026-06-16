<template>
    <div class="h-full bg-slate-900 text-gray-200 p-8 overflow-y-auto">
        <div class="max-w-2xl mx-auto">
            <h1 class="text-2xl font-bold mb-8 flex items-center gap-3">
                <i class="ye ye-gear"></i> Settings
            </h1>

            <section class="mb-10 bg-slate-800/50 p-6 rounded-xl border border-slate-700/50">
                <h2 class="text-lg font-semibold mb-4 flex items-center gap-2">
                    <i class="ye ye-key text-yellow-500"></i> API Configuration
                </h2>
                <div class="space-y-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-400 mb-1">Google Gemini API Key</label>
                        <div class="flex gap-2">
                            <input v-model="geminiKey" type="password" placeholder="Enter your Gemini API key" 
                                   class="flex-grow bg-slate-900 border border-slate-700 rounded px-3 py-2 text-sm focus:outline-none focus:border-blue-500 transition-colors" />
                            <button @click="saveKey('gemini')" :disabled="savingKey === 'gemini'"
                                    class="bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white px-4 py-2 rounded text-sm font-medium transition-colors">
                                {{ savingKey === 'gemini' ? 'Saving...' : 'Save' }}
                            </button>
                        </div>
                        <p class="mt-1 text-[10px] text-gray-500 italic">Used for Gemini 2.0 Flash models</p>
                    </div>

                    <div class="pt-2">
                        <label class="block text-sm font-medium text-gray-400 mb-1">Anthropic API Key</label>
                        <div class="flex gap-2">
                            <input v-model="anthropicKey" type="password" placeholder="Enter your Anthropic API key" 
                                   class="flex-grow bg-slate-900 border border-slate-700 rounded px-3 py-2 text-sm focus:outline-none focus:border-blue-500 transition-colors" />
                            <button @click="saveKey('anthropic')" :disabled="savingKey === 'anthropic'"
                                    class="bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white px-4 py-2 rounded text-sm font-medium transition-colors">
                                {{ savingKey === 'anthropic' ? 'Save' : 'Save' }}
                            </button>
                        </div>
                        <p class="mt-1 text-[10px] text-gray-500 italic">Used for Claude-3.5 Sonnet models</p>
                    </div>
                </div>
            </section>

            <section class="mb-10 bg-slate-800/50 p-6 rounded-xl border border-slate-700/50">
                <h2 class="text-lg font-semibold mb-4 flex items-center gap-2">
                    <i class="ye ye-robot text-blue-500"></i> Agent Defaults
                </h2>
                <div class="space-y-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-400 mb-1">Default LLM Provider</label>
                        <select v-model="defaultProvider" 
                                class="w-full bg-slate-900 border border-slate-700 rounded px-3 py-2 text-sm focus:outline-none focus:border-blue-500 transition-colors">
                            <option value="gemini">Google Gemini (Flash 2.0)</option>
                            <option value="anthropic">Anthropic (Claude 3.5 Sonnet)</option>
                        </select>
                    </div>
                </div>
            </section>

            <div class="flex justify-end gap-3 mt-8">
                <button @click="resetSettings" class="text-sm text-gray-400 hover:text-white transition-colors">Reset Defaults</button>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { SaveAPIKey, LoadAPIKey, SaveConfigToFile, LoadConfigFromFile } from '@wails/go/main/App';

const geminiKey = ref('');
const anthropicKey = ref('');
const defaultProvider = ref('gemini');
const savingKey = ref<string | null>(null);

const saveKey = async (type: string) => {
    savingKey.value = type;
    try {
        const key = type === 'gemini' ? geminiKey.value : anthropicKey.value;
        await SaveAPIKey(key);
        // Optionally store which key is which if they are different services
        await SaveConfigToFile(`APIKey_${type}`, key);
        console.log(`${type} key saved successfully`);
    } catch (error) {
        console.error("Failed to save key:", error);
    } finally {
        savingKey.value = null;
    }
}

const loadSettings = async () => {
    try {
        geminiKey.value = await LoadAPIKey(); // Default key
        // Maybe load specialized keys
        try { anthropicKey.value = await LoadConfigFromFile('APIKey_anthropic'); } catch (e) {}
        try { defaultProvider.value = await LoadConfigFromFile('DefaultProvider') || 'gemini'; } catch (e) {}
    } catch (error) {
        console.error("Failed to load settings:", error);
    }
}

const resetSettings = () => {
    if (confirm("Are you sure you want to reset all settings?")) {
        geminiKey.value = '';
        anthropicKey.value = '';
        defaultProvider.value = 'gemini';
    }
}

onMounted(() => {
    loadSettings();
})
</script>