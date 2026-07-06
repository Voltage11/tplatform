/**
 * Парсит сообщение об ошибке валидации от go-playground/validator
 * Возвращает объект: { имя_поля: "Читаемое сообщение" }
 */
export function parseValidationError(errorMessage: string): Record<string, string> {
  const fieldErrors: Record<string, string> = {};

  // Разбиваем на части по переводу строки или запятой (если несколько ошибок)
  const parts = errorMessage.split('\n');

  for (const part of parts) {
    // Ищем строку вида: Key: 'RoleCreateRequest.Name' Error:Field validation for 'Name' failed on the 'min' tag
    const match = part.match(/Field validation for '(\w+)' failed on the '(\w+)'(?: tag)?/);
    if (match) {
      const [, fieldName, rule] = match;
      fieldErrors[fieldName.toLowerCase()] = getHumanReadableError(fieldName, rule);
      continue;
    }

    // Альтернативный формат: Key: '...' Error:...
    const keyMatch = part.match(/Key: '[^']*\.(\w+)'.*Error:(.*)/);
    if (keyMatch) {
      const [, fieldName, errorText] = keyMatch;
      fieldErrors[fieldName.toLowerCase()] = errorText.trim();
    }
  }

  return fieldErrors;
}

function getHumanReadableError(field: string, rule: string): string {
  switch (rule) {
    case 'required':
      return 'Обязательное поле';
    case 'min':
      return 'Слишком короткое значение (мин. 3 символа)';
    case 'max':
      return 'Слишком длинное значение';
    case 'email':
      return 'Некорректный email';
    default:
      return `Ошибка валидации (${rule})`;
  }
}