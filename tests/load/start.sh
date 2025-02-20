#!/bin/bash

# Конфигурация
BASE_URL="http://localhost:8080"
RATE="100/s"        # RPS
DURATION="30s"      # Длительность теста

SSO_EMAIL="admin@example.com"
SSO_PASSWORD="superstrongpassword"

ROOT_DIR="./tests/load"
GET_TOKEN_PATH="$ROOT_DIR/scripts/get_token"
BODY_BASE_PATH="$ROOT_DIR/body"

TOKEN=$($GET_TOKEN_PATH -email $SSO_EMAIL -password $SSO_PASSWORD)
if [ -z "$TOKEN" ]; then
    echo "Не удалось получить токен"
    exit 1
fi

if ! command -v vegeta &> /dev/null; then
    echo "Vegeta не установлен. Установите с помощью: go install github.com/tsenart/vegeta/v12@latest"
    exit 1
fi

echo "Создаем файл целей..."
cat > tests/load/targets.txt <<EOF

# Список страниц
GET $BASE_URL/api/page
Authorization: Bearer $TOKEN

# Создание страницы
POST $BASE_URL/api/page
Authorization: Bearer $TOKEN
Content-Type: application/json
@$BODY_BASE_PATH/create_page.json

# Получение страницы по slug
GET $BASE_URL/api/page/sample-slug
Authorization: Bearer $TOKEN

# Обновление страницы по slug
PUT $BASE_URL/api/page/sample-slug
Authorization: Bearer $TOKEN
Content-Type: application/json
@$BODY_BASE_PATH/update_page.json

# Удаление страницы по slug
DELETE $BASE_URL/api/page/sample-slug
Authorization: Bearer $TOKEN

###########################################

# Создание SEO
POST $BASE_URL/api/seo
Authorization: Bearer $TOKEN
Content-Type: application/json
@$BODY_BASE_PATH/create_seo.json

# Получение SEO
GET $BASE_URL/api/seo/page/test-slug
Authorization: Bearer $TOKEN

# Обновление SEO
PUT $BASE_URL/api/seo/page/test-slug
Authorization: Bearer $TOKEN
Content-Type: application/json
@$BODY_BASE_PATH/update_seo.json

# Удаление SEO
DELETE $BASE_URL/api/seo/page/test-slug
Authorization: Bearer $TOKEN

EOF

echo "Запускаем нагрузочный тест..."
vegeta attack \
  -rate="$RATE" \
  -duration="$DURATION" \
  -targets=$ROOT_DIR/targets.txt \
  > $ROOT_DIR/results.bin

echo "Генерируем отчеты..."
vegeta report $ROOT_DIR/results.bin > $ROOT_DIR/report.txt

echo "Тестирование завершено!"
echo "Результаты:"
echo " - Текстовый отчет: report.txt"

rm -f $ROOT_DIR/targets.txt $ROOT_DIR/results.bin