<!--
  PREVIEW VIEW TEMPLATE — {{EntityPascalCase}}View.vue
  ====================================================
  Replace ALL placeholders before use:

  {{EntityPascalCase}}      → e.g. GiftItemDistribution, Order, Outlet
  {{entityCamelCase}}       → e.g. giftItemDistribution, order, outlet
  {{SCHEMA_FIELDS}}         → Schema field definitions for display

  Schema Field Type Reference:
  ┌─────────────────────┬──────────────────────────────────────────────────────────┐
  │ type: "ID"          │ Hidden primary key.  hidden: true                       │
  │ type: "Title"       │ Section divider.  label: "Section Name"                 │
  │ type: "String"      │ Text display.  spans: 2, label: "Name"                 │
  │ type: "Text"        │ Long text.  spans: 4                                    │
  │ type: "Number"      │ Numeric.  spans: 2, label: "Qty"                       │
  │ type: "Datetime"    │ Date display.  spans: 2                                 │
  │ type: "Boolean"     │ Toggle.  actions: ['Activate', 'Deactivate']            │
  │ type: "Url"         │ Image/file preview.  spans: 4                           │
  │ type: "Location"    │ Lat/Long field.  spans: 6                               │
  │ type: "String"      │ With options → select display.  options: ["a","b"]      │
  │ type: "ID" + parent │ FK relation → expand parent fields for display          │
  │ type: "ID" +children│ One-to-many → inline table of child records             │
  └─────────────────────┴──────────────────────────────────────────────────────────┘
-->

<template>
  <div class="h-full flex flex-col pt-0 text-body">
    <DynamicView 
        :id="props.id"
        :schema="schema"
        :modal="`{{EntityPascalCase}}`"
        :fetch-query-name="`{{entityCamelCase}}`"
        :mutation-update-name="`update{{EntityPascalCase}}`"
    />
  </div>
</template> 
 
<script setup lang="ts">
import DynamicView from '@core/components/dynamic-view/DynamicView.vue';
import { SchemaType } from '@core/components/dynamic-form/types';

const props = defineProps<{
    id: string;
}>();

const schema: SchemaType = {
    id: { type: "ID", hidden: true },

    // ── Section Title ──
    // sectionName: { type: "Title", label: "Section Name" },

    // ── Basic String / Text ──
    // name: { type: "String", label: "Name", spans: 2 },
    // description: { type: "Text", label: "Description", spans: 4 },

    // ── String with options (select display) ──
    // status: { type: "String", label: "Status", options: ["active", "inactive", "banned"], spans: 2 },

    // ── Number ──
    // quantity: { type: "Number", label: "Quantity", spans: 2 },
    // totalAmount: { type: "Number", label: "Total Amount", spans: 2 },

    // ── Boolean (with toggle actions) ──
    // isActive: { type: "Boolean", actions: ['Activate', 'Deactivate'] },

    // ── Datetime ──
    // createdAt: { type: "Datetime", label: "Created At", spans: 2 },

    // ── Image / File URL ──
    // imageUrl: { type: "Url", label: "Photo", spans: 4 },

    // ── FK Relation (parent expand — shows resolved parent fields) ──
    // project: { 
    //     type: "ID", 
    //     label: "Project",
    //     parent: {
    //         title: { type: "String" },
    //         description: { type: "String" }
    //     },
    //     spans: 2 
    // },

    // ── One-to-Many Children (inline table of child records) ──
    // orderItems: { 
    //     type: "ID", 
    //     label: "Order Items",
    //     children: {
    //         productItem: {
    //             type: "ID",
    //             label: "Product",
    //             parent: {
    //                 name: { type: "String" },
    //             }
    //         },
    //         quantity: { type: "Number", label: "Quantity" },
    //         unitPrice: { type: "Number", label: "Unit Price" },
    //         totalPrice: { type: "Number", label: "Total Price" },
    //     },
    //     spans: 6 
    // },

    // ── Location ──
    // latitude: { type: "Location", label: "Latitude", spans: 6 },
    // longitude: { type: "Location", label: "Longitude", spans: 6 },

    {{SCHEMA_FIELDS}}
};
</script>
