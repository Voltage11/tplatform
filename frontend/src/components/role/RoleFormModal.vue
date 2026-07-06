<script setup lang="ts">
import { useRoleStore } from '@/stores/role';
import { ref, watch } from 'vue';
import { useFieldErrors } from '@/composables/useFieldErrors';

const props = defineProps<{
  visible: boolean;
  editingRole?: { id: string; name: string; description: string } | null;
}>();

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void;
  (e: 'saved'): void;
}>();

const name = ref('');
const description = ref('');
const isSubmitting = ref(false);

const { fieldErrors, generalError, setServerError, clearErrors } = useFieldErrors();

watch(() => props.visible, (newVal) => {
  if (newVal) {
    if (props.editingRole) {
      name.value = props.editingRole.name;
      description.value = props.editingRole.description;
    } else {
      name.value = '';
      description.value = '';
    }
    clearErrors();
  }
});

async function handleSubmit() {
  if (!name.value.trim()) {
    fieldErrors.value = { name: 'Название обязательно' };
    return;
  }
  isSubmitting.value = true;
  try {
    const store = useRoleStore();
    if (props.editingRole) {
      await store.updateRole(props.editingRole.id, name.value, description.value);
    } else {
      await store.createRole(name.value, description.value);
    }
    emit('saved');
    emit('update:visible', false);
  } catch (err: any) {
    setServerError(err);
  } finally {
    isSubmitting.value = false;
  }
}

function handleCancel() {
  emit('update:visible', false);
}
</script>

<template>
  <div v-if="visible" class="modal-overlay" @click.self="handleCancel">
    <div class="modal-container">
      <h3>{{ editingRole ? 'Редактировать роль' : 'Создать роль' }}</h3>
      <form @submit.prevent="handleSubmit">
        <div class="form-group">
          <label for="role-name">Название</label>
          <input
            id="role-name"
            v-model="name"
            type="text"
            placeholder="Название роли"
            :disabled="isSubmitting"
            :class="{ 'input-error': fieldErrors['name'] }"
          />
          <span v-if="fieldErrors['name']" class="field-error">{{ fieldErrors['name'] }}</span>
        </div>
        <div class="form-group">
          <label for="role-desc">Описание</label>
          <textarea
            id="role-desc"
            v-model="description"
            placeholder="Описание роли"
            :disabled="isSubmitting"
            rows="3"
          ></textarea>
        </div>

        <div v-if="generalError" class="error-message">{{ generalError }}</div>

        <div class="modal-actions">
          <button type="button" class="btn-cancel" :disabled="isSubmitting" @click="handleCancel">
            Отмена
          </button>
          <button type="submit" class="btn-save" :disabled="isSubmitting">
            {{ isSubmitting ? 'Сохранение...' : (editingRole ? 'Обновить' : 'Создать') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}
.modal-container {
  background: white;
  padding: 24px;
  border-radius: 12px;
  min-width: 400px;
  box-shadow: 0 8px 24px rgba(0,0,0,0.2);
}
h3 {
  margin-top: 0;
}
.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 12px;
}
label {
  font-weight: 500;
}
input, textarea {
  padding: 8px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  font-size: 1rem;
}
textarea {
  resize: vertical;
}
.input-error {
  border-color: #ef4444;
}
.field-error {
  color: #ef4444;
  font-size: 0.85rem;
  margin-top: 4px;
}
.error-message {
  background-color: #fee2e2;
  color: #b91c1c;
  padding: 8px;
  border-radius: 8px;
  margin-bottom: 12px;
}
.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 20px;
}
.btn-cancel,
.btn-save {
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
.btn-save {
  background: #3b82f6;
  color: white;
}
.btn-save:hover {
  background: #2563eb;
}
.btn-cancel:hover {
  background: #cbd5e1;
}
button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>