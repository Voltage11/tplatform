<script setup lang="ts">
defineProps<{
  visible: boolean
  title?: string
  message?: string
  confirmText?: string
  cancelText?: string
  isLoading?: boolean
  errorMessage?: string | null
}>()

const emit = defineEmits<{
  (e: 'confirm'): void
  (e: 'cancel'): void
  (e: 'update:visible', value: boolean): void
}>()

function handleCancel() {
  emit('cancel')
  emit('update:visible', false)
}
</script>

<template>
  <div v-if="visible" class="modal-overlay" @click.self="handleCancel">
    <div class="modal-container">
      <h3 v-if="title">{{ title }}</h3>
      <p v-if="message">{{ message }}</p>
      <!-- Ошибка удаления -->
      <div v-if="errorMessage" class="error-message">{{ errorMessage }}</div>
      <div class="modal-actions">
        <button class="btn-cancel" :disabled="isLoading" @click="handleCancel">
          {{ cancelText || 'Отмена' }}
        </button>
        <button class="btn-confirm" :disabled="isLoading" @click="$emit('confirm')">
          {{ isLoading ? 'Выполнение...' : confirmText || 'Подтвердить' }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.error-message {
  background: #fee2e2;
  color: #b91c1c;
  padding: 8px;
  border-radius: 8px;
  margin: 12px 0;
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}
.modal-container {
  background: white;
  padding: 24px;
  border-radius: 12px;
  min-width: 320px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
}
h3 {
  margin-top: 0;
}
.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 20px;
}
.btn-cancel,
.btn-confirm {
  padding: 8px 20px;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
}
.btn-cancel {
  background: #e2e8f0;
  color: #1e293b;
}
.btn-confirm {
  background: #ef4444;
  color: white;
}
.btn-confirm:hover {
  background: #dc2626;
}
.btn-cancel:hover {
  background: #cbd5e1;
}
button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
