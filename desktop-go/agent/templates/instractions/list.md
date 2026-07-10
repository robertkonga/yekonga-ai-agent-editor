<!--
  LIST VIEW TEMPLATE — {{EntityPascalCase}}List.vue
  ================================================
  Replace ALL placeholders before use:

  {{EntityPascalCase}}      → e.g. GiftItemDistribution, Order, Outlet
  {{entityCamelCase}}       → e.g. giftItemDistribution, order, outlet
  {{entity-kebab-case}}     → e.g. gift-item-distribution, order, outlet
  {{collectionName}}        → e.g. giftItemDistributions, orders, outlets  (plural camelCase)
  {{ENTITY_LABEL}}          → e.g. Distribution, Order, Outlet  (human-readable)
  {{GRAPHQL_FIELDS}}        → GraphQL field selection  (id, name, ...)
  {{SEARCH_FIELDS}}         → Array of searchable field strings
  {{COLUMN_CONFIG}}         → Column definitions object
  {{PARENT_SCOPE_FIELD}}    → e.g. projectId (the FK used for scoping, or remove if not needed)
-->

<template>
    <div class="h-full flex flex-col">
        <div class="flex-grow overflow-hidden">
            <Paginator ref="paginatorRef"
                objectID="{{entity-kebab-case}}-list-paginator"
                :name="'{{collectionName}}'"
                :query="query"
                :limit="20"
                :options="options"
                :customClass="`grid grid-cols-1 ${paginatorRef?.displayMode === 'grid' ? 'md:grid-cols-2 lg:grid-cols-3' : ''} gap-3 py-4`">

                <template #addButton>
                    <button @click="openModal()" class="px-3 py-0.5 text-sm bg-primary-500 text-on-primary rounded-md hover:bg-primary-700 ms-4 transition-colors">
                        <i class="ye ye-circle-plus md:me-2"></i> 
                        <span class="hidden md:inline-block">{{ $t('general.add{{ENTITY_LABEL}}') }}</span>
                    </button>
                </template>

                <template #noData>
                    <div class="text-center py-5">
                        <div class="flex items-center justify-center opacity-60">
                            <img :src="`${$baseUrl}/core/image/6.svg`" alt="No Data" class="h-[150px] w-auto max-w-full " />
                        </div>
                        <div class="my-4 text-center">
                            <div class="text-2xl mb-3 text-light">{{ $t('general.noDataAvailable') }}</div>
                            <div class="p">{{ $t('general.theIsNoDataAvailableToShow') }}</div>
                            <div class="p">{{ $t('general.pleaseChooseDifferentFiltersAndTryAgain') }}</div>
                        </div>
                        <div class="pt-3">
                            <button @click="openModal()" class="px-5 py-3 bg-primary-500 text-on-primary rounded-md hover:bg-primary-700 ms-4 transition-colors me-2">
                                <i class="ye ye-plus-circle me-1"></i> {{ $t('general.addNew{{ENTITY_LABEL}}') }}
                            </button>
                        </div>
                    </div>
                </template>

                <!-- ========== CARD MODE (default) ========== -->
                <!-- Uses {{EntityPascalCase}}Card.vue — see .github/card-component-template.vue.txt -->
                <template #default="{ item, index }: { item: any; index: number }">
                    <{{EntityPascalCase}}Card :item="item" @view="openViewModal" @edit="openEditModal" />
                </template>

                <!-- ========== TABLE MODE SLOTS (uncomment if tableMode: true) ========== -->
                <!-- When using tableMode: true, remove #default slot above and use column slots instead: -->

                <!-- Primary clickable field -->
                <!-- <template #name="{ item, index }:{ item: any; index: number }">
                    <span @click="openViewModal(item.id)" class="text-primary-500 hover:underline cursor-pointer">
                        {{ item.name }}
                    </span>
                </template> -->

                <!-- Foreign key / relation display -->
                <!-- <template #relatedId="{ item, index }:{ item: any; index: number }">
                    {{ item.related?.name || '-' }}
                </template> -->

                <!-- Date display -->
                <!-- <template #createdAt="{ item, index }:{ item: any; index: number }">
                    {{ item.createdAt ? dayjs(item.createdAt).format('DD MMM, YYYY HH:mm') : '-' }}
                </template> -->

                <!-- Table actions column -->
                <!-- <template #actions="{ item }: { item: any }">
                     <div class="flex items-center">
                        <button @click="openViewModal(item.id)" class="btn btn-dark btn-circle btn-sm px-2 me-1 cursor-pointer" :title="$t('general.viewDetails')">
                             <i class="ye ye-eye ye-sm"></i>
                        </button>
                        <button @click="openEditModal(item.id)" class="btn btn-primary btn-circle btn-sm px-2 me-1"> 
                            <i class="ye ye-pencil ye-sm"></i> 
                        </button>
                        <button @click="onDelete(item.id, item.{{PRIMARY_FIELD}})" class="btn btn-red btn-circle text-white btn-sm px-2" :data-confirmation="'I you sure, you want to delete {{PRIMARY_FIELD}}?'">
                            <i class="ye ye-trash ye-sm"></i>
                        </button>
                    </div>
                </template> -->

                <!-- ========== END SLOTS ========== -->
            </Paginator>
        </div>
    </div>

    <!-- Add/Edit Modal -->
    <Modal ref="modalRef" :title="id ? 'Edit {{ENTITY_LABEL}}' : 'Add {{ENTITY_LABEL}}'" v-model="showModal" size="lg">
        <template #header> 
            {{ENTITY_LABEL}} Details
        </template>
        <{{EntityPascalCase}}Form :id="id" :isEdited="!!id" @success="onSuccessSaved" @cancel="showModal = false" />
    </Modal>

    <!-- View Modal -->
    <Modal ref="viewModalRef" title="{{ENTITY_LABEL}} Details" v-model="showViewModal" size="lg">
        <template #header> 
            {{ENTITY_LABEL}} Information 
        </template>
        <{{EntityPascalCase}}View :id="viewId" v-if="viewId" />
    </Modal>
</template>

<script setup lang="ts">
import { reactive, ref, computed, watch } from 'vue';
import { useRoute } from '@core/global';
import Paginator from '@core/components/paginator/Paginator.vue';
import { ConfigOptions } from '@core/components/paginator/state';
import Modal from '@core/components/modal/Modal.vue';
import {{EntityPascalCase}}Form from './{{EntityPascalCase}}Form.vue';
import {{EntityPascalCase}}View from './{{EntityPascalCase}}View.vue';
import {{EntityPascalCase}}Card from './{{EntityPascalCase}}Card.vue';
import dayjs from 'dayjs';

// ── Refs ──
const showModal = ref(false);
const showViewModal = ref(false);
const modalRef = ref<InstanceType<typeof Modal> | null>(null);
const viewModalRef = ref<InstanceType<typeof Modal> | null>(null);
const paginatorRef = ref<InstanceType<typeof Paginator> | null>(null);
const id = ref<string | null>(null);
const viewId = ref<string | null>(null);

// ── Route / Parent Scope ──
const route = useRoute();
const projectId = computed(() => (route && route.params && route.params.id as string));

// ── GraphQL Query ──
const query = ref<string>(`{
    {{GRAPHQL_FIELDS}}
}`);

// ── Modal Actions ──
const openModal = () => {
    id.value = null;
    showModal.value = true;
};

const openEditModal = (itemId: string) => {
    id.value = itemId;
    showModal.value = true;
};

const openViewModal = (itemId: string) => {
    viewId.value = itemId;
    showViewModal.value = true;
};

const onSuccessSaved = () => {
    showModal.value = false;
    paginatorRef.value?.reload();
};

const onDelete = async (itemId: string, title: string) => {
    const continueStatus = await window.customConfirmDelete(currentComp!.appContext, {
        title: $t('{{COLLECTION_NAME}}.deleteTitle', { title }),
        message: `Are you sure you want to delete "<strong>${title}</strong>"?`, 
    });

    if (continueStatus) {
        const mutation = `mutation { delete:delete{{ENTITY_PASCAL_CASE}}(where:{id:{equalTo:"${itemId}"}}){status,message} }`;
        const response = await window.$ajaxGraphql(mutation, {});
        if (response?.delete?.status) {
            paginatorRef.value?.reload();
            window.customAlert(currentComp!.appContext, 'success', $t('general.deleteSuccess'));
        } else {
            window.customAlert(currentComp!.appContext, 'danger', $t('general.deleteFail'));
        }
    }
};

// Add this back to your template file
const onDeleteMany = async (items: any[], title = '{{COLLECTION_NAME}}') => {
    const continueStatus = await window.customConfirmDelete(currentComp!.appContext, {
        title: $t(`{{COLLECTION_NAME}}.deleteMany`, { title }),
        message: `You are about to delete "<strong>${items.length} {{ENTITY_PASCAL_CASE}}</strong>"`, 
    });

    if (continueStatus) {
        const ids = items.map((i) => i.id);
        const mutation = `mutation { delete:delete{{ENTITY_PASCAL_CASE}}(where:{id:{in:${JSON.stringify(ids)}}}){status,message} }`;
        const response = await window.$ajaxGraphql(mutation, {});
        if (response?.delete?.status) {
            paginatorRef.value?.reload();
            window.customAlert(currentComp!.appContext, 'success', $t('general.deleteSuccess'));
        }
    }
};

const onCustomAction = (action: string, _title: string, _actionTitle: string) => {
    return async (items: any[], title = _title) => {
        const continueStatus = await window.customConfirm(currentComp!.appContext, $t(`{{COLLECTION_NAME}}.${action}ActionAlert`, { title, action: _actionTitle }));
        if (continueStatus) {
            const ids = items.map((i) => i.id);
            // NOTE: Replace 'genericAction' with the specific mutation name for this entity
            const mutation = `mutation { action:{{ENTITY_CAMEL_CASE}}Action(where:{id:{in:${JSON.stringify(ids)}}}, action:"${action}"){status,message} }`;
            const response = await window.$ajaxGraphql(mutation, {});
            if (response.action?.status) {
                paginatorRef.value?.reload();
                window.customAlert(currentComp!.appContext, 'success', $t(`general.actionSuccess`));
            }
        }
    }
};

// ── Paginator Config ──
const options = reactive<ConfigOptions>({
    name: '{{entityCamelCase}}Paginate',
    query: query.value,  
    showToolbar: true,
    showHeader: true,
    showFooter: true,
    tableMode: false,
    fixedSize: true,
    displayMode: 'list',
    modeOptions: ['grid', 'list'],
    argsRaw: `{{PARENT_SCOPE_FIELD}}:{equalTo:"${projectId.value}"}`,
    searchFields: [{{SEARCH_FIELDS}}],
    columns: {
        {{COLUMN_CONFIG}}
    },
    actions: [
        { 
            label: 'Delete',
            callback: onDeleteMany,
        },
    ],
});

// ── Re-scope on parent change ──
const resetWithinProject = () => {
    options.argsRaw = `{{PARENT_SCOPE_FIELD}}:{equalTo:"${projectId.value}"}`;
    paginatorRef.value?.reload();      
};

watch(() => projectId.value, () => {
    resetWithinProject();   
});
</script>
