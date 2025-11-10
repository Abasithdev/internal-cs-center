<template>
  <div class="min-h-screen bg-gray-50">
    <header class="bg-white shadow">
      <div class="max-w-7xl mx-auto px-4 py-4 flex justify-between items-center">
        <h1 class="text-2xl font-bold">Internal CS Center</h1>
        <div class="flex items-center gap-4">
          <span class="text-sm text-gray-600">{{ userEmail }}</span>
          <button @click="handleLogout" class="text-sm text-red-600 hover:text-red-800">
            Logout
          </button>
        </div>
      </div>
    </header>
    <main class="max-w-7xl mx-auto px-4 py-6">
      <slot />
    </main>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth';
import { useRouter } from 'vue-router';

const auth = useAuthStore();
const router = useRouter();
const userEmail = auth.email || 'User';

function handleLogout() {
  auth.logout();
  router.push('/login');
}
</script>
