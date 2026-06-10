import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAgentStore = defineStore('agent', () => {
  const ws = ref<WebSocket | null>(null)
  const isConnected = ref(false)
  
  const chatMessages = ref<{role: string, content: string}[]>([
    { role: 'assistant', content: 'Hello! I am your autonomous coding agent. What would you like me to build today?' }
  ])
  const terminalOutput = ref<string[]>([])
  
  const connect = () => {
    ws.value = new WebSocket('ws://localhost:8080/ws')
    
    ws.value.onopen = () => {
      isConnected.value = true
      terminalOutput.value.push('Connected to Agent Engine WS')
    }
    
    ws.value.onmessage = (event) => {
      const data = event.data
      
      // If the backend streams chunks, we append to the last message if it's from assistant
      if (data === '\n<done/>') {
        return // Task finished executing
      }
      
      const lastMsg = chatMessages.value[chatMessages.value.length - 1]
      if (lastMsg && lastMsg.role === 'assistant') {
        lastMsg.content += data
      } else {
        chatMessages.value.push({ role: 'assistant', content: data })
      }
      
      // Also write raw output to terminal for debugging
      terminalOutput.value.push(`[WS]: ${data}`)
    }
    
    ws.value.onclose = () => {
      isConnected.value = false
      terminalOutput.value.push('Disconnected from Agent Engine WS')
      setTimeout(connect, 3000) // Reconnect logic
    }
  }

  const sendTask = (task: string) => {
    if (ws.value && isConnected.value) {
      chatMessages.value.push({ role: 'user', content: task })
      ws.value.send(task)
      chatMessages.value.push({ role: 'assistant', content: '' }) // Placeholder for streaming response
    } else {
      console.error('WebSocket not connected')
    }
  }

  return {
    ws,
    isConnected,
    chatMessages,
    terminalOutput,
    connect,
    sendTask
  }
})
