-- +goose Up
-- +goose StatementBegin
-- Создаём заглушки пользователей для всех user_id из токенов/блоклиста/лога, которых нет в таблице users
INSERT INTO users (id, username, first_name, photo_url, auth_date, tokentg)
SELECT DISTINCT src.user_id, '', '', '', NOW(), ''
FROM (
         SELECT user_id FROM refresh_tokens
         UNION
         SELECT user_id FROM token_blacklist
         UNION
         SELECT user_id FROM token_issuance_log
     ) AS src
LEFT JOIN users u ON u.id = src.user_id
WHERE u.id IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Удаляем только добавленные нами пустые заглушки
DELETE FROM users
WHERE username = ''
  AND first_name = ''
  AND photo_url = ''
  AND tokentg = '';
-- +goose StatementEnd
