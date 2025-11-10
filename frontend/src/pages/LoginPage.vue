<template>
  <div class="login-page flex flex-col items-center justify-center h-screen bg-gray-100">
    <div class="p-6 bg-white rounded shadow-md w-80">
      <h2 class="text-lg font-bold mb-4 text-center">Login</h2>
      <form @submit.prevent="onSubmit">
        <input v-model="email" type="email" placeholder="Email" class="input mb-3" />
        <input v-model="password" type="password" placeholder="Password" class="input mb-3" />
        <button type="submit" class="btn w-full">Login</button>
      </form>
      <p v-if="error" class="text-red-500 text-sm mt-2 text-center">{{ error }}</p>
    </div>
  </div>
</template>

<style scoped>
.input {
  @apply border border-gray-300 rounded px-3 py-2 w-full;
}
.btn {
  @apply bg-blue-600 text-white py-2 rounded hover:bg-blue-700;
}
</style>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth';
import axios from 'axios';
import { ref } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();
const auth = useAuthStore();
const email = ref("");
const password = ref("");
const error = ref("");

async function onSubmit() {
    try {
        await auth.login(email.value, password.value);
        router.push("/dashboard");
    } catch (errors: unknown) {
        if (axios.isAxiosError(errors)) {
            const data = errors.response?.data as { error?: string } | undefined;
            error.value = data?.error ?? "Login failed";
        } else if (errors instanceof Error) {
            // Generic Error from JS/TS
            error.value = errors.message || "Login failed";
        } else {
            // Non-error thrown (rare)
            error.value = "Login failed";
        }
    }
}

</script>