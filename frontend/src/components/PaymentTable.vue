<template>
  <table class="w-full border-collapse border">
    <thead>
      <tr class="bg-gray-200">
        <th>ID</th><th>Merchant</th><th>Date</th><th>Amount</th><th>Status</th><th>Action</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="p in payments" :key="p.id">
        <td>{{ p.id.slice(0, 6) }}</td>
        <td>{{ p.merchant_name }}</td>
        <td>{{ new Date(p.date).toLocaleDateString() }}</td>
        <td>{{ p.amount }}</td>
        <td>{{ p.status }}</td>
        <td>
          <button
            v-if="role === 'operational'"
            @click="$emit('review', p.id)"
            class="text-blue-500 underline"
            :disabled="p.reviewed"
          >
            {{ p.reviewed ? 'Reviewed' : 'Mark as Reviewed' }}
          </button>
        </td>
      </tr>
    </tbody>
  </table>
</template>

<script setup lang="ts">
import type { Payment } from '@/type/payment';

defineProps<{ payments: Payment[]; role: string | null }>();
defineEmits(["review"]);
</script>
