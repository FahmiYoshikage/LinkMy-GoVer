-- 001_initial_schema.sql
-- LinkMy v2.0 - Initial Database Schema

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_verified BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

-- Profiles (multi-profile per user)
CREATE TABLE profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    slug VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    title VARCHAR(100),
    bio TEXT,
    avatar VARCHAR(255) DEFAULT 'default-avatar.png',
    is_active BOOLEAN DEFAULT true,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

-- Themes (one per profile)
CREATE TABLE themes (
    id SERIAL PRIMARY KEY,
    profile_id INTEGER UNIQUE NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    bg_type VARCHAR(20) DEFAULT 'gradient',
    bg_value TEXT,
    button_style VARCHAR(20) DEFAULT 'rounded',
    button_color VARCHAR(20) DEFAULT '#667eea',
    text_color VARCHAR(20) DEFAULT '#333333',
    font VARCHAR(50) DEFAULT 'Inter',
    layout VARCHAR(20) DEFAULT 'centered',
    container_style VARCHAR(20) DEFAULT 'wide',
    enable_animations BOOLEAN DEFAULT true,
    enable_glass_effect BOOLEAN DEFAULT false,
    shadow_intensity VARCHAR(20) DEFAULT 'medium',
    -- Boxed layout (embedded)
    boxed_enabled BOOLEAN DEFAULT false,
    boxed_outer_bg_type VARCHAR(20) DEFAULT 'gradient',
    boxed_outer_bg_value TEXT,
    boxed_container_bg VARCHAR(20) DEFAULT '#ffffff',
    boxed_max_width INTEGER DEFAULT 480,
    boxed_radius INTEGER DEFAULT 30,
    boxed_shadow BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

-- Categories (link folders)
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    profile_id INTEGER NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    icon VARCHAR(50) DEFAULT 'bi-folder',
    color VARCHAR(20) DEFAULT '#667eea',
    position INTEGER DEFAULT 0,
    is_expanded BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Links
CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    profile_id INTEGER NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    title VARCHAR(100) NOT NULL,
    url VARCHAR(500) NOT NULL,
    icon VARCHAR(50) DEFAULT 'bi-link-45deg',
    position INTEGER DEFAULT 0,
    clicks INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

-- Click analytics
CREATE TABLE clicks (
    id BIGSERIAL PRIMARY KEY,
    link_id INTEGER NOT NULL REFERENCES links(id) ON DELETE CASCADE,
    ip INET,
    country VARCHAR(50),
    city VARCHAR(100),
    user_agent TEXT,
    referrer TEXT,
    clicked_at TIMESTAMPTZ DEFAULT NOW()
);

-- Sessions (JWT refresh tokens)
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token VARCHAR(255) UNIQUE NOT NULL,
    ip INET,
    user_agent VARCHAR(255),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Email verifications
CREATE TABLE email_verifications (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) NOT NULL,
    otp VARCHAR(6) NOT NULL,
    type VARCHAR(20) DEFAULT 'registration',
    ip INET,
    is_used BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL
);

-- Password resets
CREATE TABLE password_resets (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) NOT NULL,
    token VARCHAR(64) UNIQUE NOT NULL,
    ip INET,
    is_used BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL
);

-- Social icons reference
CREATE TABLE social_icons (
    id SERIAL PRIMARY KEY,
    platform_name VARCHAR(50) NOT NULL,
    icon_class VARCHAR(50) NOT NULL,
    icon_color VARCHAR(20),
    base_url VARCHAR(100)
);

-- Insert default social icons
INSERT INTO social_icons (platform_name, icon_class, icon_color, base_url) VALUES
('Instagram', 'bi-instagram', '#E4405F', 'https://instagram.com/'),
('Facebook', 'bi-facebook', '#1877F2', 'https://facebook.com/'),
('Twitter/X', 'bi-twitter-x', '#000000', 'https://twitter.com/'),
('LinkedIn', 'bi-linkedin', '#0A66C2', 'https://linkedin.com/in/'),
('GitHub', 'bi-github', '#181717', 'https://github.com/'),
('YouTube', 'bi-youtube', '#FF0000', 'https://youtube.com/'),
('TikTok', 'bi-tiktok', '#000000', 'https://tiktok.com/@'),
('WhatsApp', 'bi-whatsapp', '#25D366', 'https://wa.me/'),
('Telegram', 'bi-telegram', '#26A5E4', 'https://t.me/'),
('Discord', 'bi-discord', '#5865F2', 'https://discord.gg/'),
('Twitch', 'bi-twitch', '#9146FF', 'https://twitch.tv/'),
('Spotify', 'bi-spotify', '#1DB954', 'https://open.spotify.com/'),
('Email', 'bi-envelope-fill', '#EA4335', 'mailto:'),
('Website', 'bi-globe', '#667eea', 'https://'),
('Link', 'bi-link-45deg', '#6c757d', NULL);

-- Indexes
CREATE INDEX idx_profiles_user ON profiles(user_id);
CREATE INDEX idx_profiles_slug ON profiles(slug) WHERE is_active = true;
CREATE INDEX idx_links_profile ON links(profile_id) WHERE is_active = true;
CREATE INDEX idx_links_category ON links(category_id);
CREATE INDEX idx_categories_profile ON categories(profile_id);
CREATE INDEX idx_clicks_link ON clicks(link_id);
CREATE INDEX idx_clicks_date ON clicks(clicked_at);
CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_sessions_expires ON sessions(expires_at);
CREATE INDEX idx_sessions_token ON sessions(refresh_token);
CREATE INDEX idx_email_verifications_email ON email_verifications(email);
CREATE INDEX idx_password_resets_email ON password_resets(email);
CREATE INDEX idx_password_resets_token ON password_resets(token);
