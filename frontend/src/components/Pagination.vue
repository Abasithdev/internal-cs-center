<template>
  <div class="mt-4 flex justify-center gap-2">
    <button :disabled="page === 1" @click="$emit('change', page - 1)">Prev</button>
    <span>Page {{ page }} / {{ totalPages }}</span>
    <button :disabled="page === totalPages" @click="$emit('change', page + 1)">Next</button>
  </div>
</template>

<script setup lang="ts">
defineProps<{ page: number; totalPages: number }>();
defineEmits(["change"]);

defineOptions({
  name: "AppPagination"
});
</script>
