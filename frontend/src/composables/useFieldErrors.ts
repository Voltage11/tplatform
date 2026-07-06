import { ref, type Ref } from 'vue';
import { parseValidationError } from '@/utils/validation';

export function useFieldErrors() {
  const fieldErrors = ref<Record<string, string>>({});
  const generalError = ref<string | null>(null);

  function setServerError(error: any) {
    // Сбрасываем старые ошибки
    fieldErrors.value = {};
    generalError.value = null;

    // Проверяем, является ли ответ ошибкой валидации (400 и есть текст Field validation)
    if (
      error?.response?.status === 400 &&
      typeof error?.response?.data?.error === 'string' &&
      error.response.data.error.includes('Field validation')
    ) {
      const parsed = parseValidationError(error.response.data.error);
      fieldErrors.value = parsed;
      // Если не удалось распарсить ни одного поля, покажем общее сообщение
      if (Object.keys(parsed).length === 0) {
        generalError.value = error.response.data.error;
      }
    } else {
      // Любая другая ошибка (сеть, сервер, 500, 409 и т.д.)
      generalError.value =
        error?.response?.data?.error || error?.message || 'Произошла неизвестная ошибка';
    }
  }

  function clearErrors() {
    fieldErrors.value = {};
    generalError.value = null;
  }

  return { fieldErrors, generalError, setServerError, clearErrors };
}