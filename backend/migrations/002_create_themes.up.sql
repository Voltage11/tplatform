-- 1. Темы
CREATE TABLE IF NOT EXISTS themes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    date_begin TIMESTAMP WITH TIME ZONE,
    date_end TIMESTAMP WITH TIME ZONE,
    max_point INT NOT NULL DEFAULT 0,
    check_point INT NOT NULL DEFAULT 0,
    img_path VARCHAR(255),
    correct_points INT NOT NULL DEFAULT 0,  
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT fk_themes_users FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE RESTRICT
);

-- 2. Вопросы
CREATE TABLE IF NOT EXISTS questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    theme_id UUID NOT NULL,
    question_type VARCHAR(50) NOT NULL,
    name TEXT NOT NULL,                
    point_correct INT NOT NULL DEFAULT 0,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_questions_themes FOREIGN KEY (theme_id) REFERENCES themes(id) ON DELETE CASCADE
);

-- 3. Варианты ответов
CREATE TABLE IF NOT EXISTS answers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_id UUID NOT NULL,
    name TEXT NOT NULL,                 
    is_correct BOOLEAN NOT NULL DEFAULT FALSE,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_answers_questions FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE
);

CREATE INDEX idx_questions_theme_id ON questions(theme_id);
CREATE INDEX idx_answers_question_id ON answers(question_id);
CREATE INDEX idx_themes_created_by ON themes(created_by_id);