<template>
  <DefaultLayout>
    <div class="p-4">
      <h2 class="text-xl font-bold mb-4">Payments Dashboard</h2>
      <div class="flex gap-2 mb-4">
        <select v-model="status" class="border p-2 rounded">
          <option value="">All Status</option>
          <option value="completed">Completed</option>
          <option value="processing">Processing</option>
          <option value="failed">Failed</option>
        </select>
        <input v-model="search" placeholder="Search by ID" class="border p-2 rounded" />
        <button @click="fetchPayments" class="bg-blue-500 text-white px-4 rounded">Search</button>
      </div>
      <SummaryWidget :summary="summary" />
      <PaymentTable :payments="payments" :role="role" @review="onReview" />
      <Pagination :page="page" :totalPages="totalPages" @change="onPageChange" />
    </div>
  </DefaultLayout>
</template>

<script setup lang="ts">
import { getPayments, reviewPayment } from '@/api/paymentApi';
import { useAuthStore } from '@/stores/auth';
import { onMounted, ref } from 'vue';
import type { Payment, PaymentSummary } from '@/type/payment';
import DefaultLayout from '../layouts/DefaultLayout.vue';
import SummaryWidget from '@/components/SummaryWidget.vue';
import PaymentTable from '@/components/PaymentTable.vue';
import Pagination from '@/components/Pagination.vue';

const auth = useAuthStore();
const payments = ref<Payment[]>([]);
const summary = ref<PaymentSummary>()
const page = ref(1);
const totalPages = ref(1);
const status = ref("")
const search = ref("")
const role = auth.role;

async function fetchPayments() {
    const response = await getPayments({page: page.value, status: status.value, search: search.value})
    console.log(response)
    payments.value = response.meta.data;
    totalPages.value = response.meta.total_pages
    summary.value = response.summary
}

async function onReview(id: string) {
    if(role != "operational") return;
    await reviewPayment(id)
    await fetchPayments();
}

function onPageChange(p: number){
    page.value=p;
    fetchPayments();
}

onMounted(fetchPayments);
</script>