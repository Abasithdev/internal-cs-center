<template>
  <div v-if="summary">
    <h2 class="text-xl font-bold mb-4">Payments Summary</h2>
    <div class="flex gap-4 mb-4">
      <div>Total: {{ summary.total }}</div>
      <div>Completed: {{ summary.completed }}</div>
      <div>Processing: {{ summary.processing }}</div>
      <div>Failed: {{ summary.failed }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { PaymentSummary } from '@/type/payment';

defineProps<{ summary?: PaymentSummary }>();
</script>
