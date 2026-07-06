<script setup lang="ts">
import { ref, watch } from 'vue';
import * as permissionApi from '@/api/permission';
import type { PermissionEntity } from '@/types/permission';

const props = defineProps<{
  visible: boolean;
  roleId: string | null;
  roleName?: string;
}>();

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void;
  (e: 'saved'): void;
}>();

const entities = ref<PermissionEntity[]>([]);
const isLoading = ref(false);
const isSaving = ref(false);
const error = ref<string | null>(null);

watch(() => props.visible, async (newVal) => {
  if (newVal && props.roleId) {
    isLoading.value = true;
    error.value = null;
    try {
      entities.value = await permissionApi.fetchPermissions(props.roleId);
    } catch (err: any) {
      error.value = err?.response?.data?.error || 'Ошибка загрузки прав';
    } finally {
      isLoading.value = false;
    }
  }
});

function toggleAction(entityIndex: number, actionIndex: number) {
  const entity = entities.value[entityIndex];
  if (!entity) return;
  const action = entity.actions[actionIndex];
  if (action) {
    action.is_active = !action.is_active;
  }
}

async function handleSave() {
  if (!props.roleId) return;
  isSaving.value = true;
  error.value = null;
  try {
    const selected = entities.value.flatMap(entity =>
      entity.actions
        .filter(a => a.is_active)
        .map(a => ({
          entity_name: entity.entity.name,
          action_name: a.name,
        }))
    );
    await permissionApi.updatePermissions(props.roleId, selected);
    emit('saved');
    emit('update:visible', false);
  } catch (err: any) {
    error.value = err?.response?.data?.error || 'Ошибка сохранения прав';
  } finally {
    isSaving.value = false;
  }
}

function handleCancel() {
  emit('update:visible', false);
}
</script>

<template>
  <div v-if="visible" class="modal-overlay" @click.self="handleCancel">
    <div class="modal-container">
      <h3>Права роли «{{ roleName }}»</h3>
      <div v-if="isLoading" class="loading">Загрузка...</div>
      <div v-else-if="error" class="error-message">{{ error }}</div>
      <div v-else class="permissions-table">
        <div class="entity-row header">
          <div class="entity-name">Сущность</div>
          <div class="actions-list">
            <span v-for="action in entities[0]?.actions" :key="action.name" class="action-header">
              {{ action.description }}
            </span>
          </div>
        </div>
        <div v-for="(entity, ei) in entities" :key="entity.entity.name" class="entity-row">
          <div class="entity-name">{{ entity.entity.description }}</div>
          <div class="actions-list">
            <label
              v-for="(action, ai) in entity.actions"
              :key="action.name"
              class="action-checkbox"
            >
              <input
                type="checkbox"
                :checked="action.is_active"
                @change="toggleAction(ei, ai)"
              />
            </label>
          </div>
        </div>
      </div>

      <div class="modal-actions">
        <button class="btn-cancel" @click="handleCancel" :disabled="isSaving">Отмена</button>
        <button class="btn-save" @click="handleSave" :disabled="isSaving || isLoading">
          {{ isSaving ? 'Сохранение...' : 'Сохранить' }}
        </button>
      </div>
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
  min-width: 600px;
  box-shadow: 0 8px 24px rgba(0,0,0,0.2);
}
h3 {
  margin-top: 0;
}
.loading,
.error-message {
  text-align: center;
  padding: 20px;
}
.error-message {
  color: #b91c1c;
  background: #fee2e2;
  border-radius: 8px;
}
.permissions-table {
  margin: 16px 0;
}
.entity-row {
  display: flex;
  align-items: center;
  border-bottom: 1px solid #e2e8f0;
  padding: 8px 0;
}
.entity-row.header {
  border-bottom: 2px solid #cbd5e1;
  font-weight: 600;
}
.entity-name {
  width: 200px;
  padding-right: 16px;
}
.actions-list {
  display: flex;
  gap: 16px;
  flex: 1;
}
.action-checkbox {
  width: 40px;
  text-align: center;
}
.action-header {
  width: 40px;
  text-align: center;
  font-size: 0.85rem;
}
.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 16px;
}
.btn-cancel,
.btn-save {
  padding: 8px 20px;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
}
.btn-cancel { background: #e2e8f0; color: #1e293b; }
.btn-save { background: #3b82f6; color: white; }
.btn-save:hover { background: #2563eb; }
.btn-cancel:hover { background: #cbd5e1; }
button:disabled { opacity: 0.6; cursor: not-allowed; }
</style>