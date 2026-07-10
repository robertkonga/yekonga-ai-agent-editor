import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ContextMenu from '@imengyu/vue3-context-menu'
import App from '@/App.vue'
import { useWorkspace } from '@/store/workspace.ts';
import Modal from '@/components/Modal.vue';
import ExpandableContent from '@/components/Expandable.vue';
// import CustomDropdown from '@/components/CustomDropdown.vue';

import '@/window.ts';
import '@imengyu/vue3-context-menu/lib/vue3-context-menu.css';
import '@/style.css';

const app = createApp(App);
app.use(ContextMenu);
app.use(createPinia());

app.config.globalProperties.$toDate = window.toDate;
app.config.globalProperties.$isImage = window.isImage;
app.config.globalProperties.$isVideo = window.isVideo;
app.config.globalProperties.$isAudio = window.isAudio;

app.component('modal', Modal);
app.component('expandable', ExpandableContent);
// app.component('custom-dropdown', CustomDropdown);

useWorkspace();

app.mount('#app');
