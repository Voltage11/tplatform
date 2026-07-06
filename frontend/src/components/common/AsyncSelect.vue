<script setup lang="ts">
import { ref, watch, onMounted, computed } from 'vue';

interface Option {
  id: string;
  name: string;
}

const props = defineProps<{
  modelValue: string | null;
  placeholder?: string;
  fetchOptions: (search: string) => Promise<Option[]>;
  initialOptions?: Option[]; // первые записи (без поиска)
  disabled?: boolean;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: string | null): void;
}>();

const options = ref<Option[]>(props.initialOptions || []);
const isLoading = ref(false);
const search = ref('');
const isOpen = ref(false);
const selectedOption = ref<Option | null>(null);
const inputRef = ref<HTMLInputElement | null>(null);

const filteredOptions = computed(() => {
  if (!search.value) return options.value;
  return options.value.filter(opt =>
    opt.name.toLowerCase().includes(search.value.toLowerCase())
  );
});

// При изменении modelValue или initialOptions подгружаем выбранное значение
watch(() => props.modelValue, async (newVal) => {
  if (newVal) {
    const option = options.value.find(o => o.id === newVal);
    if (option) {
      selectedOption.value = option;
    } else {
      // Если опция не в списке, попытаемся загрузить по id (можно дополнительно)
      // Пока просто сбросим
      selectedOption.value = null;
    }
  } else {
    selectedOption.value = null;
  }
}, { immediate: true });

watch(() => props.initialOptions, (newVal) => {
  if (newVal) options.value = newVal;
});

async function loadInitialOptions() {
  if (props.initialOptions?.length) {
    options.value = props.initialOptions;
    return;
  }
  isLoading.value = true;
  try {
    options.value = await props.fetchOptions('');
  } catch (e) {
    // ошибка
  } finally {
    isLoading.value = false;
  }
}

async function searchOptions() {
  isLoading.value = true;
  try {
    options.value = await props.fetchOptions(search.value);
  } catch (e) {
    // ошибка
  } finally {
    isLoading.value = false;
  }
}

let debounceTimer: ReturnType<typeof setTimeout> | null = null;
function onSearchInput() {
  if (debounceTimer) clearTimeout(debounceTimer);
  debounceTimer = setTimeout(searchOptions, 300);
}

function selectOption(option: Option) {
  selectedOption.value = option;
  emit('update:modelValue', option.id);
  search.value = '';
  isOpen.value = false;
}

function toggleOpen() {
  if (props.disabled) return;
  isOpen.value = !isOpen.value;
  if (isOpen.value) {
    if (options.value.length === 0) loadInitialOptions();
    // фокус на input внутри
    setTimeout(() => inputRef.value?.focus(), 50);
  }
}

function onBlur() {
  // Закрываем с задержкой, чтобы успел сработать клик по опции
  setTimeout(() => {
    isOpen.value = false;
    search.value = '';
  }, 200);
}

function clearSelection() {
  selectedOption.value = null;
  emit('update:modelValue', null);
  search.value = '';
}
</script>

<template>
  <div class="async-select" :class="{ 'is-disabled': disabled }">
    <div class="select-control" @click="toggleOpen">
      <span v-if="selectedOption" class="selected-text">{{ selectedOption.name }}</span>
      <span v-else class="placeholder">{{ placeholder || 'Выберите...' }}</span>
      <button v-if="selectedOption && !disabled" type="button" class="clear-btn" @click.stop="clearSelection">×</button>
      <span class="arrow">▼</span>
    </div>
    <div v-if="isOpen" class="dropdown">
      <input
        ref="inputRef"
        v-model="search"
        type="text"
        class="search-input"
        placeholder="Поиск..."
        @input="onSearchInput"
        @blur="onBlur"
      />
      <ul class="options-list">
        <li v-if="isLoading" class="loading">Загрузка...</li>
        <li
          v-for="option in filteredOptions"
          :key="option.id"
          class="option"
          :class="{ 'is-selected': option.id === selectedOption?.id }"
          @mousedown.prevent="selectOption(option)"
        >
          {{ option.name }}
        </li>
        <li v-if="!isLoading && filteredOptions.length === 0" class="no-options">Ничего не найдено</li>
      </ul>
    </div>
  </div>
</template>

<style scoped>
.async-select {
  position: relative;
  width: 100%;
}
.select-control {
  display: flex;
  align-items: center;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  padding: 8px 12px;
  cursor: pointer;
  background: white;
  min-height: 38px;
}
.is-disabled .select-control {
  background: #f1f5f9;
  cursor: not-allowed;
}
.selected-text {
  flex: 1;
}
.placeholder {
  color: #94a3b8;
}
.clear-btn {
  background: none;
  border: none;
  font-size: 1.2rem;
  cursor: pointer;
  padding: 0 4px;
}
.arrow {
  margin-left: 8px;
  font-size: 0.8rem;
  color: #64748b;
}
.dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: white;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  margin-top: 4px;
  z-index: 100;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}
.search-input {
  width: 100%;
  padding: 8px;
  border: none;
  border-bottom: 1px solid #e2e8f0;
  border-radius: 8px 8px 0 0;
  outline: none;
}
.options-list {
  list-style: none;
  margin: 0;
  padding: 0;
  max-height: 200px;
  overflow-y: auto;
}
.option {
  padding: 8px 12px;
  cursor: pointer;
}
.option:hover, .option.is-selected {
  background: #f1f5f9;
}
.loading, .no-options {
  padding: 8px 12px;
  color: #64748b;
}
</style>