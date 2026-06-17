import { EventsOn } from "@wails/runtime/runtime";


// ── Provider / model config ──────────────────────────────────────────────────

export interface ModelOption {
    label: string
    value: string
    description?: string
    disabled?: boolean
}

export interface ModelGroup {
    group: string
    models: ModelOption[]
}

export interface ProviderConfig {
    label: string
    value: string
    icon: string
    toolUse: boolean
    defaultModel: string
    modelGroups: ModelGroup[]
}

export const providerOptions: ProviderConfig[] = [
    {
        label: 'Ollama',
        value: 'ollama',
        icon: '🦙',
        toolUse: true,
        defaultModel: 'qwen3:8b',
        modelGroups: [
            {
                group: 'Qwen3',
                models: [
                    { label: 'qwen3:8b',  value: 'qwen3:8b',  description: '8B · fast' },
                    { label: 'qwen3:14b', value: 'qwen3:14b', description: '14B · balanced' },
                    { label: 'qwen3:32b', value: 'qwen3:32b', description: '32B · best quality' },
                ],
            },
            {
                group: 'DeepSeek',
                models: [
                    { label: 'deepseek-r1:7b',  value: 'deepseek-r1:7b',  description: '7B · reasoning' },
                    { label: 'deepseek-r1:14b', value: 'deepseek-r1:14b', description: '14B · reasoning' },
                    { label: 'deepseek-r1:32b', value: 'deepseek-r1:32b', description: '32B · reasoning' },
                ],
            },
        ],
    },
    {
        label: 'DeepSeek',
        value: 'deepseek',
        icon: '🐋',
        toolUse: true,
        defaultModel: 'deepseek-chat',
        modelGroups: [
            {
                group: '',
                models: [
                    { label: 'deepseek-chat',     value: 'deepseek-chat',     description: 'V3 · general purpose' },
                    { label: 'deepseek-reasoner', value: 'deepseek-reasoner', description: 'R1 · chain-of-thought' },
                    { label: 'DeepSeek-V4-Flash', value: 'deepseek-v4-flash', description: 'R1 · chain-of-thought' },
                    { label: 'DeepSeek-V4-Pro',   value: 'deepseek-v4-pro',   description: 'R1 · chain-of-thought' },
                ],
            },
        ],
    },
    {
        label: 'Gemini',
        value: 'gemini',
        icon: '✦',
        toolUse: true,
        defaultModel: 'gemini-2.5-flash',
        modelGroups: [
            {
                group: 'Gemini 2.5',
                models: [
                    { label: 'gemini-2.5-flash',      value: 'gemini-2.5-flash',      description: 'Fast · recommended' },
                    { label: 'gemini-2.5-flash-lite',  value: 'gemini-2.5-flash-lite', description: 'Lightest' },
                    { label: 'gemini-2.5-pro',        value: 'gemini-2.5-pro',        description: 'Most capable' },
                ],
            },
        ],
    },
    {
        label: 'Anthropic',
        value: 'anthropic',
        icon: '◆',
        toolUse: true,
        defaultModel: 'claude-sonnet-4-6',
        modelGroups: [
            {
                group: 'Claude 4',
                models: [
                    { label: 'claude-opus-4-6',    value: 'claude-opus-4-6',    description: 'Most capable' },
                    { label: 'claude-sonnet-4-6',  value: 'claude-sonnet-4-6',  description: 'Balanced · recommended' },
                ],
            },
            {
                group: 'Claude 3.5',
                models: [
                    { label: 'claude-haiku-4-5',   value: 'claude-haiku-4-5-20251001', description: 'Fastest · lowest cost' },
                ],
            },
        ],
    },
]