<script setup lang="ts">
import { ref, watch } from 'vue'
import { useThemeStore } from '@/stores/theme'
import { useFieldErrors } from '@/composables/useFieldErrors'
import type { Theme } from '@/types/theme'
import { uploadFile } from '@/api/upload'

const props = defineProps<{
  visible: boolean
  editingTheme?: Theme | null
}>()

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'saved'): void
}>()

const name = ref('')
const description = ref('')
const isActive = ref(true)
const dateBegin = ref('')
const dateEnd = ref('')
const checkPoint = ref(0)
const imgPath = ref('')
const isSubmitting = ref(false)

const { fieldErrors, generalError, setServerError, clearErrors } = useFieldErrors()

const imgFile = ref<File | null>(null)
const imgPreview = ref<string | null>(null)
const isUploadingImg = ref(false)
const imgError = ref<string | null>(null)

watch(
  () => props.visible,
  (newVal) => {
    if (newVal) {
      imgFile.value = null
      imgPreview.value = null
      imgError.value = null

      if (props.editingTheme) {
        const t = props.editingTheme
        name.value = t.name
        description.value = t.description || ''
        isActive.value = t.is_active
        dateBegin.value = t.date_begin ? t.date_begin.slice(0, 16) : '' // для input type="datetime-local"
        dateEnd.value = t.date_end ? t.date_end.slice(0, 16) : ''
        checkPoint.value = t.check_point
        imgPath.value = t.img_path || ''
      } else {
        name.value = ''
        description.value = ''
        isActive.value = true
        dateBegin.value = ''
        dateEnd.value = ''
        checkPoint.value = 0
        imgPath.value = ''
      }
      clearErrors()
    }
  },
)

async function handleSubmit() {
  if (!name.value.trim()) {
    fieldErrors.value = { name: 'Название обязательно' }
    return
  }

  isSubmitting.value = true
  try {
    const store = useThemeStore()
    const payload = {
      name: name.value,
      description: description.value || null,
      is_active: isActive.value,
      date_begin: dateBegin.value ? new Date(dateBegin.value).toISOString() : null,
      date_end: dateEnd.value ? new Date(dateEnd.value).toISOString() : null,
      check_point: checkPoint.value,
      img_path: imgPath.value || null,
    }

    if (props.editingTheme) {
      await store.updateTheme(props.editingTheme.id, payload)
    } else {
      await store.createTheme(payload)
    }
    emit('saved')
    emit('update:visible', false)
  } catch (err: any) {
    setServerError(err)
  } finally {
    isSubmitting.value = false
  }
}

function handleCancel() {
  emit('update:visible', false)
}

// Выбор файла
function onFileChange(event: Event) {
  const target = event.target as HTMLInputElement
  const files = target.files
  if (files && files.length > 0) {
    const file = files[0]
    // Валидация на клиенте
    const validTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/svg+xml']
    if (!validTypes.includes(file.type)) {
      imgError.value = 'Допустимы только JPEG, PNG, GIF, SVG'
      return
    }
    if (file.size > 5 * 1024 * 1024) {
      // 5 MB
      imgError.value = 'Размер файла не должен превышать 5 МБ'
      return
    }
    imgFile.value = file
    imgPreview.value = URL.createObjectURL(file)
    imgError.value = null
    // Автоматически загружаем на сервер
    uploadImage()
  }
}

async function uploadImage() {
  if (!imgFile.value) return
  isUploadingImg.value = true
  imgError.value = null
  try {
    const { url } = await uploadFile(imgFile.value)
    imgPath.value = url // записываем полученный URL
  } catch (err: any) {
    imgError.value = err?.response?.data?.error || 'Ошибка загрузки изображения'
    imgFile.value = null
    imgPreview.value = null
  } finally {
    isUploadingImg.value = false
  }
}

// Кнопка "Удалить изображение"
function removeImage() {
  imgFile.value = null
  imgPreview.value = null
  imgPath.value = ''
}
</script>

<template>
  <div v-if="visible" class="modal-overlay" @click.self="handleCancel">
    <div class="modal-container">
      <h3>{{ editingTheme ? 'Редактировать тему' : 'Создать тему' }}</h3>
      <form @submit.prevent="handleSubmit">
        <div class="form-group">
          <label>Изображение</label>
          <input
            type="file"
            accept="image/*"
            @change="onFileChange"
            :disabled="isSubmitting || isUploadingImg"
          />
          <div v-if="isUploadingImg" class="upload-progress">Загрузка изображения...</div>
          <div v-if="imgError" class="field-error">{{ imgError }}</div>
          <div v-if="imgPreview" class="image-preview">
            <img :src="imgPreview" alt="Предпросмотр" class="preview-img" />
            <button type="button" class="btn-remove" @click="removeImage">✕</button>
          </div>
          <!-- если редактирование и уже есть img_path -->
          <div v-else-if="editingTheme?.img_path">
            <img :src="editingTheme.img_path" alt="Текущее изображение" class="current-img" />
          </div>

          <label>Название *</label>
          <input
            v-model="name"
            type="text"
            :disabled="isSubmitting"
            :class="{ 'input-error': fieldErrors['name'] }"
          />
          <span v-if="fieldErrors['name']" class="field-error">{{ fieldErrors['name'] }}</span>
        </div>
        <div class="form-group">
          <label>Описание</label>
          <textarea v-model="description" :disabled="isSubmitting" rows="3"></textarea>
        </div>
        <div class="form-row">
          <div class="form-group">
            <label>Дата начала</label>
            <input v-model="dateBegin" type="datetime-local" :disabled="isSubmitting" />
          </div>
          <div class="form-group">
            <label>Дата окончания</label>
            <input v-model="dateEnd" type="datetime-local" :disabled="isSubmitting" />
          </div>
        </div>
        <div class="form-row">
          <div class="form-group">
            <label>Проходной балл</label>
            <input v-model.number="checkPoint" type="number" min="0" :disabled="isSubmitting" />
          </div>
          <div class="form-group checkbox-group">
            <label>
              <input v-model="isActive" type="checkbox" :disabled="isSubmitting" />
              Активна
            </label>
          </div>
        </div>
        <div class="form-group">
          <label>Ссылка на изображение</label>
          <input v-model="imgPath" type="text" :disabled="isSubmitting" placeholder="https://..." />
        </div>

        <div v-if="generalError" class="error-message">{{ generalError }}</div>

        <div class="modal-actions">
          <button type="button" class="btn-cancel" :disabled="isSubmitting" @click="handleCancel">
            Отмена
          </button>
          <button type="submit" class="btn-save" :disabled="isSubmitting">
            {{ isSubmitting ? 'Сохранение...' : editingTheme ? 'Обновить' : 'Создать' }}
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
  min-width: 360px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
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
input {
  padding: 8px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  font-size: 1rem;
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
.image-preview {
    position: relative;
    display: inline-block;
    margin-top: 8px;
}
.preview-img, .current-img {
    max-width: 200px;
    max-height: 150px;
    border-radius: 8px;
    border: 1px solid #ccc;
}
.btn-remove {
    position: absolute;
    top: -8px;
    right: -8px;
    background: #ef4444;
    color: white;
    border: none;
    border-radius: 50%;
    width: 24px;
    height: 24px;
    cursor: pointer;
}
.upload-progress {
    color: #3b82f6;
    font-size: 0.9rem;
}
</style>
