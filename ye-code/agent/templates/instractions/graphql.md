# GraphQL Data Fetching Guide

## Query Pattern

Use `window.$ajaxGraphql(query, {})` to fetch data from GraphQL API.

**Basic Query Format:**
```
query{<modalName>(where:{<modelField>:{<condition>:"<value>"}}){<modelFields>}}
```

## Basic Query Parameters

- **modalName**: Model name in camelCase plural form
  - Examples: `users`, `events`, `userAccounts`, `messages`, `campaigns`, `transactions`, `wallets`, `batches`
  
- **modelField**: Model field in camelCase
  - Examples: `id`, `status`, `channel`, `createdAt`, `amount`, `type`
  
- **condition**: Filter condition
  - `equalTo`: Exact match (e.g., `status: {equalTo: "active"}`)
  - `notEqualTo`: Not equal
  - `greaterThan`: Greater than
  - `greaterThanOrEqualTo`: Greater than or equal
  - `lessThan`: Less than
  - `lessThanOrEqualTo`: Less than or equal
  - `contains`: String contains (case-insensitive)
  - `startsWith`: String starts with
  - `endsWith`: String ends with
  - `in`: Value in array (e.g., `status: {in: ["active", "pending"]}`)
  - `between`: Value between two numbers

- **modelFields**: Comma-separated fields to return
  - Examples: `id,name,status,createdAt`

## Basic Usage Examples

**Example 1: Fetch all active messages**
```javascript
const query = `query{messages(where:{status:{equalTo:"delivered"}}){id,recipient,channel,status,sentAt,cost}}`;
const data = await window.$ajaxGraphql(query, {});
```

**Example 2: Fetch campaigns by status**
```javascript
const query = `query{campaigns(where:{status:{equalTo:"running"}}){id,name,description,totalRecipients,status,createdAt}}`;
const data = await window.$ajaxGraphql(query, {});
```

**Example 3: Fetch recent transactions**
```javascript
const query = `query{transactions(where:{type:{in:["topup","deduction","refund"]}}){id,type,amount,reference,createdAt}}`;
const data = await window.$ajaxGraphql(query, {});
```

**Example 4: Fetch batches with counts**
```javascript
const query = `query{batches(where:{status:{notEqualTo:"pending"}}){id,title,totalCount,successCount,errorCount,status,createdAt}}`;
const data = await window.$ajaxGraphql(query, {});
```

## Implementation Pattern in Components

```typescript
import { ref, onMounted } from 'vue';

const data = ref([]);
const loading = ref(false);
const error = ref(null);

onMounted(async () => {
    try {
        loading.value = true;
        const query = `query{<modalName>(where:{<modelField>:{<condition>:"<value>"}}){<modelFields>}}`;
        const response = await window.$ajaxGraphql(query, {});
        data.value = response.data?.<modalName> || [];
    } catch (err) {
        error.value = err.message;
        console.error('GraphQL Error:', err);
    } finally {
        loading.value = false;
    }
});
```

## Error Handling

Always wrap GraphQL calls in try-catch blocks and handle:
- Network errors
- GraphQL validation errors
- Missing or null responses
- Type mismatches

---

# Advanced GraphQL: Aggregations & Grouping

## Aggregation Query Pattern

For analytics with grouping and summary functions:

```
query {
  <modalName>(groupBy:[<field1>,<field2>]) {
    <field1>
    <field2>
    <modelKey>Summary {
      count
      sum(targetKey:<numericField>)
      max(targetKey:<field>)
      average(targetKey:<numericField>)
      min(targetKey:<field>)
      
      countWithFilter: count(where:{<modelFieldCondition>})
      sumWithFilter: sum(targetKey:<numericField>, where:{<modelFieldCondition>})
    }
  }
}
```

## Important Syntax Rules

⚠️ **CRITICAL:**
- **`count`, `sum`, `min`, `max`, `average`** are COMPUTED functions, NOT static model fields
- **DO NOT use empty brackets `()`** - omit entirely if no parameters needed
- **Add `where` parameter only when filtering** which records are included in that aggregation
- Each aggregate function can have its own `where` to compute different metrics in one query
- **`<modelKey>Summary` naming**: Singular camelCase + "Summary" (e.g., `messageSummary`, `transactionSummary`)

### Correct vs Incorrect Syntax

✅ **CORRECT:**
```javascript
messageSummary {
  count
  sum(targetKey:cost)
  average(targetKey:cost)
}
```

❌ **WRONG:**
```javascript
messageSummary {
  count()
  sum(targetKey:cost, where:{})
  average(targetKey:cost, )
}
```

## Aggregation Usage Examples

**Example 1: Message statistics grouped by channel**
```javascript
const query = `query {
  messages(groupBy:[channel]) {
    channel
    messageSummary {
      total: count
      delivered: count(where:{status:{equalTo:"delivered"}})
      failed: count(where:{status:{equalTo:"failed"}})
      totalCost: sum(targetKey:cost)
      smsCost: sum(targetKey:cost, where:{channel:{equalTo:"sms"}})
      avgCost: average(targetKey:cost)
    }
  }
}`;
const data = await window.$ajaxGraphql(query, {});
```

**Example 2: Transaction analysis by type with filters**
```javascript
const query = `query {
  transactions(groupBy:[type]) {
    type
    transactionSummary {
      all: count
      recent: count(where:{createdAt:{greaterThan:"2024-01-01"}})
      totalAmount: sum(targetKey:amount)
      largeTransactions: sum(targetKey:amount, where:{amount:{greaterThan:"100000"}})
      avgAmount: average(targetKey:amount)
      maxAmount: max(targetKey:amount)
    }
  }
}`;
const data = await window.$ajaxGraphql(query, {});
```

**Example 3: Campaign performance metrics by status**
```javascript
const query = `query {
  campaigns(groupBy:[status]) {
    status
    campaignSummary {
      totalCampaigns: count
      recentCampaigns: count(where:{createdAt:{greaterThan:"2024-01-01"}})
      totalRecipients: sum(targetKey:totalRecipients)
      smsCampaignRecipients: sum(targetKey:totalRecipients, where:{channel:{equalTo:"sms"}})
      avgRecipients: average(targetKey:totalRecipients)
      maxRecipients: max(targetKey:totalRecipients)
    }
  }
}`;
const data = await window.$ajaxGraphql(query, {});
```

**Example 4: Batch processing with completion filters**
```javascript
const query = `query {
  batches(groupBy:[status]) {
    status
    batchSummary {
      allBatches: count
      completedBatches: count(where:{status:{equalTo:"completed"}})
      totalMessages: sum(targetKey:totalCount)
      successMessages: sum(targetKey:successCount)
      errorMessages: sum(targetKey:errorCount, where:{errorCount:{greaterThan:"0"}})
      successRate: average(targetKey:successCount)
    }
  }
}`;
const data = await window.$ajaxGraphql(query, {});
```

**Example 5: Multi-level grouping (Messages by status, channel, campaign)**
```javascript
const query = `query {
  messages(groupBy:[status,channel,campaignId]) {
    status
    channel
    campaignId
    messageSummary {
      count
      sum(targetKey:cost)
      max(targetKey:cost)
      average(targetKey:cost)
      min(targetKey:cost)
    }
  }
}`;
const data = await window.$ajaxGraphql(query, {});
```

**Example 6: Comprehensive dashboard query with multiple filters**
```javascript
const query = `query {
  messages(groupBy:[channel]) {
    channel
    messageSummary {
      all: count
      delivered: count(where:{status:{equalTo:"delivered"}})
      sent: count(where:{status:{equalTo:"sent"}})
      failed: count(where:{status:{equalTo:"failed"}})
      totalCost: sum(targetKey:cost)
      smsCost: sum(targetKey:cost, where:{channel:{equalTo:"sms"}})
      whatsappCost: sum(targetKey:cost, where:{channel:{equalTo:"whatsapp"}})
      emailCost: sum(targetKey:cost, where:{channel:{equalTo:"email"}})
      avgCost: average(targetKey:cost)
      highCostMessages: count(where:{cost:{greaterThan:"1"}})
      maxCost: max(targetKey:cost)
      minCost: min(targetKey:cost, where:{cost:{greaterThan:"0"}})
    }
  }
}`;
const response = await window.$ajaxGraphql(query, {});
// Single query returns separate counts/sums for each condition
```

## Dashboard Implementation Pattern

```typescript
import { ref, onMounted } from 'vue';

const channelStats = ref([]);
const loading = ref(false);

onMounted(async () => {
    try {
        loading.value = true;
        const query = `query {
          messages(groupBy:[channel]) {
            channel
            messageSummary {
              total: count
              delivered: count(where:{status:{equalTo:"delivered"}})
              failed: count(where:{status:{equalTo:"failed"}})
              totalCost: sum(targetKey:cost)
              avgCost: average(targetKey:cost)
            }
          }
        }`;
        
        const response = await window.$ajaxGraphql(query, {});
        channelStats.value = response.data?.messages || [];
        
        // Returns: [
        //   { channel: 'sms', messageSummary: { total: 1250, delivered: 1200, failed: 50, totalCost: 625, avgCost: 0.5 } },
        //   { channel: 'whatsapp', messageSummary: { total: 850, delivered: 840, failed: 10, totalCost: 850, avgCost: 1.0 } },
        //   ...
        // ]
    } catch (err) {
        console.error('Aggregation query failed:', err);
    } finally {
        loading.value = false;
    }
});
```

## Benefits of Aggregation Queries

✅ Single query instead of multiple API calls  
✅ Database-level filtering (highly efficient)  
✅ Compute multiple metrics with different conditions simultaneously  
✅ Reduced network overhead  
✅ Perfect for dashboard KPIs and analytics  
✅ Clean, readable syntax without empty brackets  
✅ Flexible: each function can have independent `where` conditions  

## Aggregation Functions Available

- **count**: Total records in group (or matching `where` condition)
- **sum(targetKey:<field>)**: Sum numeric values
- **max(targetKey:<field>)**: Maximum value in group
- **min(targetKey:<field>)**: Minimum value in group
- **average(targetKey:<field>)**: Average of numeric values

All can optionally include `where:{<condition>}` parameter.

## Common Aggregation Patterns

| Use Case | Pattern | Benefit |
|----------|---------|---------|
| Channel breakdown | `groupBy:[channel]` | See performance by channel |
| Status analysis | `groupBy:[status]` | Understand distribution |
| Multi-dimension | `groupBy:[status,channel]` | Cross-tabulation analytics |
| Filtered metrics | `count(where:{...})` | Compare filtered vs total |
| Financial summary | `sum(targetKey:amount)` | Revenue/cost calculations |
| Performance KPIs | `average(targetKey:...)` | Average metrics per group |

