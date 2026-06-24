<!--
  REPORT PAGE TEMPLATE — {{ReportPascalCase}}.vue
  ================================================
  Replace ALL placeholders before use:

  {{ReportPascalCase}}      → e.g. SalesAnalytics, FieldForcePerformance, GiftManagement
  {{REPORT_TITLE}}          → e.g. "Sales Analytics", "Field Force Performance"

  ─────────────────────────────────────────────────────────────────────────────
  DATA SCHEMA — derived from `database.json` in the parent server directory
  ─────────────────────────────────────────────────────────────────────────────

  `database.json` is the SINGLE SOURCE OF TRUTH for all report data.
  Every chart, KPI, table, and query MUST map to collections/fields defined there.

  How to read database.json for reports:
  ┌──────────────────────────────────────────────────────────────────────────┐
  │ database.json Key          →  What it means for reports                 │
  ├──────────────────────────────────────────────────────────────────────────┤
  │ Collection name (key)      →  model-key (camelCase singular)           │
  │   "Orders"                 →  model-key="order"                        │
  │   "FieldVisits"            →  model-key="fieldVisit"                   │
  │   "GiftItemDistributions"  →  model-key="giftItemDistribution"         │
  │   "ProductItems"           →  model-key="productItem"                  │
  │                                                                        │
  │ Collection name (plural)   →  GraphQL query name (camelCase plural)    │
  │   "Orders"                 →  orders(where: {...})                     │
  │   "FieldVisits"            →  fieldVisits(where: {...})                │
  │                                                                        │
  │ Field with type "Decimal"  →  target-key candidate (SUM/AVERAGE)       │
  │ or "Float" or "Number"     │                                           │
  │   totalAmount, quantity,   →  target-key="totalAmount"                 │
  │   price, unitPrice, value  →  target-key="quantity"                    │
  │                                                                        │
  │ Field with type "Datetime" →  period-key candidate (time axis)         │
  │   createdAt, distributedAt →  period-key="createdAt"                   │
  │   completedAt, startTime   →  period-key="distributedAt"              │
  │                                                                        │
  │ Field with type "String"   →  dimension candidate (group-by)           │
  │ + options array            │                                           │
  │   status, category,        →  dimension="status"                       │
  │   changeType, tradeType    →  dimension-breakdown="changeType"         │
  │                                                                        │
  │ Field with foreignKey      →  dimension for groupBy + parent expand    │
  │   fieldOfficerId → FK      →  groupBy: fieldOfficerId                  │
  │   outletId → FK            →  dimension="outletId"                     │
  │   giftItemId → FK          →  dimension="giftItemId"                   │
  │                                                                        │
  │ Summary naming convention  →  {{modelKey}}Summary                      │
  │   "Orders" collection      →  orderSummary { count, sum(...) }         │
  │   "FieldVisits" collection →  fieldVisitSummary { count }              │
  │   "OrderItems" child of    →  orderItemSummary { sum(targetKey:...) }  │
  │     "Orders"               │                                           │
  └──────────────────────────────────────────────────────────────────────────┘

  Mapping Rules:
  1. model-key    = Collection name → remove trailing "s"/"es" → camelCase
                    "Orders" → "order", "FieldVisits" → "fieldVisit"
                    "GiftItemDistributions" → "giftItemDistribution"
  2. target-key   = Any Decimal/Float/Number field you want to SUM or AVERAGE
                    If total-type="COUNT", target-key is NOT needed
  3. period-key   = The Datetime field used as the time axis
  4. dimension    = The field to group by (String with options, or FK ID)
  5. dimension-breakdown = Optional second grouping for stacked/grouped charts
  6. :where       = Always include { projectId: projectId } for scoping
  7. GraphQL groupBy = Use for leaderboard/table: collection(groupBy: fieldName)

  Example — Building a "Sales by Status" report from database.json:
    database.json → "Orders" collection has:
      - totalAmount: Decimal  → target-key="totalAmount"
      - createdAt: Datetime   → period-key="createdAt"
      - status: String opts   → dimension="status" or dimension-breakdown="status"
    Result:
      <ChartCard model-key="order" dimension="createdAt"
        dimension-breakdown="status" target-key="totalAmount"
        total-type="SUM" period-type="DAILY" type="LINE" />

  Example — Building a "Gift Distribution by Item" report:
    database.json → "GiftItemDistributions" collection has:
      - quantity: Number      → target-key="quantity"
      - distributedAt: DT     → period-key="distributedAt"
      - giftItemId: FK        → dimension="giftItemId"
    Result:
      <ChartCard model-key="giftItemDistribution" dimension="giftItemId"
        period-key="distributedAt" target-key="quantity"
        total-type="SUM" type="PIE" period-type="NONE" />

  Example — Building a custom leaderboard table:
    database.json → "FieldVisits" has fieldOfficerId (FK to FieldOfficer)
    GraphQL:
      fieldVisits(where:{projectId:{equalTo:"..."}}, groupBy:[fieldOfficerId]) {
        fieldOfficer { name, phone }
        summary: fieldVisitSummary {
          totalVisits: count
          thisWeek: count(where: { createdAt: { greaterThanOrEqualTo: "last week" } })
          lastWeek: count(where: { createdAt: { greaterThanOrEqualTo: "last 2 week", lessThanOrEqualTo: "last week" } })
        }
      }

  ─────────────────────────────────────────────────────────────────────────────
  REPORT BUILDING BLOCKS (copy/paste the sections you need)
  ─────────────────────────────────────────────────────────────────────────────

  1. ChartSimpleCard   — KPI summary card with sparkline (COUNT, SUM, AVERAGE)
  2. ChartCard         — Full chart panel (LINE, BAR, PIE, DOUGHNUT)
  3. Custom Table      — Manual data table with fetchData + GraphQL query
  4. Coverage Bars     — Progress bars for target tracking
  5. Google Map        — Geospatial view with markers/heatmap/polyline

  ChartSimpleCard Props Reference:
  ┌──────────────────────┬────────────────────────────────────────────────────┐
  │ :where               │ { projectId: projectId } — filter scope           │
  │ :dynamicTitle        │ true — auto-render computed value as title        │
  │ subtitle             │ Card label text                                   │
  │ model-key            │ Collection camelCase: "order","fieldVisit" etc.   │
  │ dimension            │ Field to group by: "createdAt","status" etc.      │
  │ dimension-breakdown  │ Optional sub-group: "fieldOfficerId","changeType" │
  │ type                 │ "LINE" | "BAR"                                    │
  │ period-type          │ "DAILY" | "WEEKLY" | "MONTHLY" | "NONE"          │
  │ total-type           │ "COUNT" | "SUM" | "AVERAGE"                      │
  │ target-key           │ Field to aggregate: "totalAmount","quantity"      │
  │ period-key           │ Date field: "createdAt","distributedAt"           │
  │ start-date           │ "30 days ago","This Month","last 2 weeks" etc.   │
  │ end-date             │ "now","last week" etc.                            │
  └──────────────────────┴────────────────────────────────────────────────────┘

  ChartCard Extra Props (in addition to above):
  ┌──────────────────────┬────────────────────────────────────────────────────┐
  │ type                 │ "LINE"|"BAR"|"PIE"|"DOUGHNUT"                     │
  │ show-summary         │ "total" | "mostActive" | "none"                  │
  │ :has-custom-legend   │ true — render legend outside chart               │
  │ :minVersion          │ true — compact card style                        │
  │ prefix / suffix      │ Label decorators: "Sales - " / " Gifts"         │
  │ chart-height         │ CSS height: "200px", "h-16"                      │
  └──────────────────────┴────────────────────────────────────────────────────┘

  GraphQL Summary/Aggregation Pattern:
    summary: {{modelKey}}Summary {
        totalCount: count
        thisWeek: count(where: { createdAt: { greaterThanOrEqualTo: "last week", lessThanOrEqualTo: "now" } })
        totalRevenue: sum(targetKey: totalAmount)
    }
-->

<template>
    <div class="text-body mb-6">
        
        <!-- ════════════════════════════════════════════════════════ -->
        <!-- SECTION 1: KPI Summary Cards (pick what you need)       -->
        <!-- ════════════════════════════════════════════════════════ -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-6">

            <!-- KPI Card — COUNT -->
            <ChartSimpleCard 
                :where="{ projectId: projectId }"
                :dynamicTitle="true"
                subtitle="{{KPI_LABEL}}" 
                model-key="{{modelKey}}"
                dimension="createdAt"
                type="LINE"
                period-type="DAILY"
                total-type="COUNT"
                period-key="createdAt"
                start-date="30 days ago"
                end-date="now">
                <template #action>
                    <div class="p-1.5 bg-purple-500/10 rounded-lg text-purple-500">
                        <i class="ye ye-box w-4 h-4"></i>
                    </div>
                </template>
            </ChartSimpleCard>

            <!-- KPI Card — SUM -->
            <!-- <ChartSimpleCard 
                :where="{ projectId: projectId }"
                :dynamicTitle="true"
                subtitle="Total Revenue" 
                model-key="order"
                dimension="createdAt"
                type="LINE"
                period-type="DAILY"
                total-type="SUM"
                period-key="createdAt"
                target-key="totalAmount"
                start-date="30 days ago"
                end-date="now">
            </ChartSimpleCard> -->

            <!-- KPI Card — AVERAGE -->
            <!-- <ChartSimpleCard 
                :where="{ projectId: projectId }"
                :dynamicTitle="true"
                subtitle="Avg Order Value" 
                model-key="order"
                dimension="createdAt"
                type="LINE"
                period-type="DAILY"
                total-type="AVERAGE"
                target-key="totalAmount"
                period-key="createdAt"
                start-date="30 days ago"
                end-date="now">
            </ChartSimpleCard> -->

            <!-- KPI Card — with dimension-breakdown (stacked/grouped) -->
            <!-- <ChartSimpleCard 
                :where="{ projectId: projectId }"
                :dynamicTitle="true"
                subtitle="Stock Movements" 
                model-key="giftStockMovement"
                dimension="createdAt"
                dimension-breakdown="changeType"
                target-key="quantity"
                type="BAR"
                period-type="DAILY"
                total-type="SUM"
                period-key="createdAt"
                start-date="30 days ago"
                end-date="now">
            </ChartSimpleCard> -->
        </div>

        <!-- ════════════════════════════════════════════════════════ -->
        <!-- SECTION 2: Charts + Sidebar (2/3 + 1/3 layout)         -->
        <!-- ════════════════════════════════════════════════════════ -->
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">

            <!-- LINE/BAR Chart (large, 2-col span) -->
            <ChartCard 
                :where="{ projectId: projectId }"
                class="lg:col-span-2"
                title="{{CHART_TITLE}}" 
                model-key="{{modelKey}}"
                dimension="createdAt"
                type="LINE"
                period-type="DAILY"
                total-type="COUNT"
                period-key="createdAt"
                start-date="This Month"
                end-date="now">
            </ChartCard>

            <!-- PIE/DOUGHNUT Chart (sidebar) -->
            <!-- <ChartCard 
                :where="{ projectId: projectId }"
                title="Distribution by Type" 
                model-key="giftItemDistribution"
                dimension="giftItemId" 
                period-key="distributedAt"
                type="PIE"
                period-type="NONE"
                target-key="quantity"
                total-type="SUM"
                show-summary="total"
                :has-custom-legend="true"
                :minVersion="true"
                suffix=" items"
                chart-height="200px"
                start-date="last 2 weeks"
                end-date="last week">
            </ChartCard> -->

            <!-- Sidebar: Custom Ranked List -->
            <!-- <div class="bg-card rounded-2xl p-6 border dark:border-gray-700/50 flex flex-col">
                <h3 class="font-bold text-lg mb-6">Top Items</h3>
                <div class="space-y-6 overflow-y-auto">
                    <div v-for="item in topItems" :key="item.id" class="flex items-center justify-between group">
                        <div class="flex items-center gap-4">
                             <div class="w-12 h-12 rounded-xl bg-gray-100 dark:bg-gray-800 flex items-center justify-center">
                                  <i class="ye ye-box w-6 h-6 text-gray-500"></i>
                             </div>
                            <div>
                                <p class="font-bold group-hover:text-primary-500 transition-colors">{{ item.name }}</p>
                                <p class="text-xs text-gray-500 font-mono">{{ item.subtitle }}</p>
                            </div>
                        </div>
                        <div class="text-right">
                            <p class="font-bold">{{ $filter.number(item.value) }}</p>
                            <p class="text-xs text-green-500 font-bold">{{ item.units }} units</p>
                        </div>
                    </div>
                </div>
                 <button class="w-full mt-auto pt-4 text-xs font-bold text-gray-400 hover:text-white transition-colors uppercase tracking-wider">
                    View All
                </button>
            </div> -->
        </div>

        <!-- ════════════════════════════════════════════════════════ -->
        <!-- SECTION 3: Data Table (optional — for leaderboards etc) -->
        <!-- ════════════════════════════════════════════════════════ -->
        <!-- <div class="bg-card rounded-2xl p-6 border dark:border-gray-700/50 mt-6">
            <div class="flex items-center justify-between mb-6">
                <h3 class="font-bold text-lg">{{TABLE_TITLE}}</h3>
            </div>
            <div class="overflow-x-auto">
                <table class="w-full text-sm text-left">
                    <thead class="text-xs text-gray-500 uppercase bg-gray-50 dark:bg-gray-800">
                        <tr>
                            <th class="px-4 py-3 rounded-l-lg">Name</th>
                            <th class="px-4 py-3">Value</th>
                            <th class="px-4 py-3 rounded-r-lg text-right">Trend</th>
                        </tr>
                    </thead>
                    <tbody>
                         <tr v-for="row in tableData" :key="row.id" class="border-b dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors">
                            <td class="px-4 py-3 font-medium flex items-center gap-3">
                                <div class="w-8 h-8 rounded-full bg-primary-100 dark:bg-primary-900/50 text-primary-600 flex items-center justify-center font-bold text-xs">
                                    {{ row.initials }}
                                </div>
                                {{ row.name }}
                            </td>
                            <td class="px-4 py-3 text-gray-600 dark:text-gray-300">{{ row.value }}</td>
                            <td class="px-4 py-3 text-right">
                                <p class="text-xs" :class="[
                                    (row.percent < 0) ? 'text-red-500' : 'text-green-500'
                                ]">{{ row.growth }}</p>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div> -->

        <!-- ════════════════════════════════════════════════════════ -->
        <!-- SECTION 4: Coverage/Progress Bars (optional)            -->
        <!-- ════════════════════════════════════════════════════════ -->
        <!-- <div class="bg-card rounded-2xl p-6 border dark:border-gray-700/50 mt-6">
            <h3 class="font-bold text-lg mb-6">Coverage Analysis</h3>
            <div class="space-y-6">
                <div v-for="(item, index) in progressItems" :key="item.id">
                     <div class="flex justify-between mb-2 text-sm font-medium">
                         <span>{{ item.label }}</span>
                         <span>{{ item.current }}/{{ item.target }}</span>
                     </div>
                     <div class="w-full bg-gray-100 dark:bg-gray-800 rounded-full h-2.5">
                        <div class="h-2.5 rounded-full" :style="{ 
                            width: item.percent + '%',
                            backgroundColor: window.getRandomColor(index)
                        }"></div>
                     </div>
                </div>
            </div>
        </div> -->

    </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import { useRoute } from '@core/global';
import {
    Chart as ChartJS,
    Title,
    Tooltip,
    Legend,
    BarElement,
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    ArcElement,
    Filler,
} from 'chart.js';
// Import only the chart types you use: Line, Bar, Doughnut, Pie
// import { Line, Bar, Doughnut } from 'vue-chartjs';
import ChartCard from '@core/components/cards/ChartCard.vue';
import ChartSimpleCard from '@core/components/cards/ChartSimpleCard.vue';
// import { DollarSign, Gift, TrendingUp, Filter, Package, MoreHorizontal } from 'lucide-vue-next';

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend, PointElement, LineElement, ArcElement, Filler);

// ── Route / Project Scope ──
const route = useRoute();
const projectId = computed(() => (route && route.params && route.params.id as string));

// ── Custom Data (for tables, ranked lists etc.) ──
// const tableData = ref<any[]>([]);
// const topItems = ref<any[]>([]);
// const progressItems = ref<any[]>([]);

// ── GraphQL Query (for custom data fetching) ──
// const query = computed(() => {
//     return `
// query Get{{ReportPascalCase}}Data {
//     {{collectionName}}(where: { projectId: { equalTo: "${projectId.value}" } }) {
//         id
//         name
//         summary: {{modelKey}}Summary {
//             totalCount: count
//             thisWeek: count(
//                 where: {
//                     createdAt: {
//                         greaterThanOrEqualTo: "last week"
//                         lessThanOrEqualTo: "now"
//                     }
//                 }
//             )
//             lastWeek: count(
//                 where: {
//                     createdAt: {
//                         greaterThanOrEqualTo: "last 2 week"
//                         lessThanOrEqualTo: "last week"
//                     }
//                 }
//             )
//         }
//     }
// }`;
// });

// ── Fetch & Transform ──
// const fetchData = async () => {
//     const res = await window.$ajaxGraphql(query.value, {});
//     if (res && res.{{collectionName}}) {
//         tableData.value = res.{{collectionName}}.map((item: any) => {
//             const total = item.summary.lastWeek || 1;
//             const percent = ((item.summary.thisWeek - item.summary.lastWeek) / total * 100);
//             return {
//                 id: item.id,
//                 name: item.name,
//                 value: item.summary.totalCount,
//                 initials: item.name?.substring(0, 2)?.toUpperCase() || '??',
//                 percent,
//                 growth: `${percent > 0 ? '▲' : '▼'} ${percent.toFixed(0)}% from last week`
//             };
//         });
//     }
// };
// fetchData();
</script>
