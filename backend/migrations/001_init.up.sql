-- ==================== TABLE: registrations ====================
CREATE TABLE registrations(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hashed VARCHAR(255) NOT NULL,
    ip_address VARCHAR(50) NOT NULL,
    user_agent TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    expired_at TIMESTAMPTZ NOT NULL,
    activated_at TIMESTAMPTZ,
    token VARCHAR(255) NOT NULL UNIQUE,
    verify_code varchar(5)
);
COMMENT ON TABLE registrations IS 'Заявки на регистрацию';
CREATE INDEX idx_registrations_email ON registrations (email);
CREATE INDEX idx_registrations_token ON registrations (token);

-- ==================== TABLE: users ====================
CREATE TABLE users
(
    id             SERIAL PRIMARY KEY,
    name           VARCHAR(255) NOT NULL,
    email          VARCHAR(255) NOT NULL UNIQUE,
    password_hashed  VARCHAR(255) NOT NULL,
    is_active      BOOLEAN      NOT NULL DEFAULT false,
    is_admin       BOOLEAN      NOT NULL DEFAULT false,
    last_login_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE users IS 'Пользователи системы';

CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_is_active ON users (is_active) WHERE is_active = true;

-- ==================== TABLE: sessions ====================
CREATE TABLE sessions
(
    id            SERIAL PRIMARY KEY,
    user_id       INTEGER             NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    refresh_token VARCHAR(1000) UNIQUE NOT NULL,
    user_agent    TEXT DEFAULT '',
    ip_address    VARCHAR(50),
    expired_at    TIMESTAMPTZ         NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE sessions IS 'Сессии пользователей с refresh tokens';

CREATE INDEX idx_sessions_refresh_token ON sessions (refresh_token);
CREATE INDEX idx_sessions_user_id ON sessions (user_id);
CREATE INDEX idx_sessions_expired_at ON sessions (expired_at);
CREATE INDEX idx_sessions_user_id_expired_at ON sessions (user_id, expired_at);

-- ===================== TABLE: topics ===================
CREATE TABLE topics(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,    
    name VARCHAR(150) NOT NULL,
    slug VARCHAR(150) NOT NULL UNIQUE,
    description TEXT DEFAULT '',
    is_active BOOLEAN DEFAULT TRUE,
    max_attempts INTEGER DEFAULT 1,
    date_end TIMESTAMPTZ,    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE topics IS 'Тема для тестирования, в которой содержатся вопросы';
COMMENT ON COLUMN topics.max_attempts IS 'Максимальное количество попыток, 0-неограничено';
CREATE INDEX idx_topics_user_id ON topics(user_id);
CREATE INDEX idx_topics_slug ON topics(slug);

-- ===================== TABLE: questions ===================
CREATE TABLE questions(
    id SERIAL PRIMARY KEY,
    topic_id INTEGER NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    name VARCHAR(500) NOT NULL,
    is_multiple BOOLEAN DEFAULT FALSE,
    position INTEGER DEFAULT 1,
    points INTEGER DEFAULT 1,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE questions IS 'Вопрос в теме для тестирования';    
CREATE INDEX idx_questions_topic_id ON questions(topic_id);

-- ===================== TABLE: answer ===================
CREATE TABLE answers (
    id SERIAL PRIMARY KEY,
    question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    name VARCHAR(500) NOT NULL,
    is_correct BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE answers IS 'Варианты ответа';
CREATE INDEX idx_answers_question_id ON answers(question_id);

-- ===================== TABLE: topic_results =============
CREATE TABLE topic_results(
    id SERIAL PRIMARY KEY,    
    topic_id INTEGER NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    begin_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    finish_at TIMESTAMPTZ,
    is_finished BOOLEAN DEFAULT FALSE,
    points INTEGER DEFAULT 0
);
COMMENT ON TABLE topic_results IS 'Результаты тестирования';
CREATE INDEX idx_topic_results_topic_id ON topic_results(topic_id);
CREATE INDEX idx_topic_results_user_id ON topic_results(user_id);
CREATE INDEX idx_topic_results_user_topic ON topic_results(user_id, topic_id);

-- ===================== TABLE: user_answers =============
CREATE TABLE user_answers (
    id SERIAL PRIMARY KEY,
    topic_result_id INTEGER NOT NULL REFERENCES topic_results(id) ON DELETE CASCADE,
    question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    answered_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    points INTEGER DEFAULT 0
);
COMMENT ON TABLE user_answers IS 'Ответы пользователя на конкретные вопросы';
CREATE INDEX idx_user_answers_topic_result_id ON user_answers(topic_result_id);
CREATE INDEX idx_user_answers_question_id ON user_answers(question_id);

-- ===================== TABLE: user_answer_details =============
-- Для хранения выбранных вариантов ответов (поддерживает multiple choice)
CREATE TABLE user_answer_details (
    id SERIAL PRIMARY KEY,
    user_answer_id INTEGER NOT NULL REFERENCES user_answers(id) ON DELETE CASCADE,
    answer_id INTEGER NOT NULL REFERENCES answers(id) ON DELETE CASCADE,    
    is_selected BOOLEAN DEFAULT TRUE,
    UNIQUE(user_answer_id, answer_id)
);
COMMENT ON TABLE user_answer_details IS 'Детализация выбранных вариантов ответов';
CREATE INDEX idx_user_answer_details_user_answer_id ON user_answer_details(user_answer_id);
CREATE INDEX idx_user_answer_details_answer_id ON user_answer_details(answer_id);
