<!--
  FORM VIEW TEMPLATE — {{EntityPascalCase}}Form.vue
  =================================================
  Replace ALL placeholders before use:

  {{EntityPascalCase}}      → e.g. GiftItemDistribution, Order, Outlet
  {{entityCamelCase}}       → e.g. giftItemDistribution, order, outlet
  {{SCHEMA_FIELDS}}         → Schema field definitions (see examples below)
-->

<template>
    <div class="px-4 text-body">
        <DynamicForm 
            :schema="schema"
            :id="props.id"
            :isEdited="props.isEdited"
            modal="{{EntityPascalCase}}"
            fetch-query-name="{{entityCamelCase}}"
            mutation-create-name="create{{EntityPascalCase}}"
            mutation-update-name="update{{EntityPascalCase}}"
            @success="emit('success')"
            @cancel="emit('cancel')"
        />
    </div>
</template> 
 
<script setup lang="ts">
import DynamicForm from '@core/components/dynamic-form/DynamicForm.vue';
import { SchemaType } from '@core/components/dynamic-form/types';

const emit = defineEmits(['success', 'cancel']);
const props = defineProps<{
    id?: string | null;
    isEdited?: boolean;
}>();

const schema: SchemaType = {
    id: { type: "ID" },

    // ── Section Title ──
    // sectionTitle: { type: "Title", label: "Section Name" },

    // ── Basic Fields ──
    // name: { type: "String", spans: 4 },
    // description: { type: "Text", spans: 4 },
    // isActive: { type: "Boolean", default: true },

    // ── ForeignKey (dropdown with search) ──
    // clientId: { 
    //     type: "ID", 
    //     label: "Client",
    //     foreignKey: {
    //         model: "Clients",       // Collection name
    //         label: "name",          // Display field
    //         key: "id",              // Value field
    //         search: ["name", "contactEmail"]  // Searchable fields
    //     }, 
    //     spans: 2  
    // },

    // ── Select / Options ──
    // status: { type: "String", options: ["active", "inactive"], default: "active", spans: 2 },

    // ── Number / Decimal ──
    // quantity: { type: "Number", spans: 2 },
    // price: { type: "String", spans: 2 },   // Use String for decimals in form inputs

    // ── Date ──
    // startTime: { type: "Datetime", default: null, spans: 2 },

    // ── File Upload ──
    // imageUrl: { type: "File", label: "Photo", spans: 2 },

    {{SCHEMA_FIELDS}}
};
</script>
