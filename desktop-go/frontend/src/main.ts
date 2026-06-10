import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './style.css'
import App from './App.vue'
import { useWorkspace } from './store/workspace.ts';
import './window.ts'

const app = createApp(App);
app.use(createPinia());

useWorkspace();

app.mount('#app');
