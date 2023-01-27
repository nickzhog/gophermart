CREATE TABLE IF NOT EXISTS public.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    login TEXT NOT NULL,
    password_hash TEXT NOT NULL, 
    balance TEXT NOT NULL DEFAULT '0'
);

CREATE TABLE IF NOT EXISTS public.orders (
    id TEXT PRIMARY KEY,
    user_id UUID NOT NULL,
    status TEXT NOT NULL DEFAULT 'NEW',
    accrual TEXT NOT NULL,
    sum TEXT NOT NULL,
    upload_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	constraint user_id FOREIGN KEY (user_id) REFERENCES public.users (id)
);

CREATE TABLE IF NOT EXISTS public.withdrawals (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    sum TEXT NOT NULL,
	processed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	constraint user_id FOREIGN KEY (user_id) REFERENCES public.users (id)
);

CREATE TABLE IF NOT EXISTS public.sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    create_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_agent TEXT NOT NULL,
    ip TEXT NOT NULL,
    is_active bool NOT NULL DEFAULT true
);


CREATE TABLE IF NOT EXISTS public.session_user (
    session_id UUID NOT NULL,
    user_id UUID NOT NULL,
    loginned_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_active bool NOT NULL DEFAULT true,
	constraint session_id FOREIGN KEY (session_id) REFERENCES public.sessions (id),
	constraint user_id FOREIGN KEY (user_id) REFERENCES public.users (id)
);