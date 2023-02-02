CREATE TABLE IF NOT EXISTS public.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    login TEXT NOT NULL,
    password_hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS public.orders (
    id TEXT PRIMARY KEY,
    user_id UUID NOT NULL,
    status TEXT NOT NULL DEFAULT 'NEW',
    accrual TEXT NOT NULL,
    upload_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT user_id FOREIGN KEY (user_id) REFERENCES public.users (id)
);

CREATE TABLE IF NOT EXISTS public.withdrawals (
    id TEXT PRIMARY KEY,
    user_id UUID NOT NULL,
    sum TEXT NOT NULL,
	processed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT user_id FOREIGN KEY (user_id) REFERENCES public.users (id)
);

CREATE TABLE IF NOT EXISTS public.sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    create_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_active bool NOT NULL DEFAULT true,
	CONSTRAINT user_id FOREIGN KEY (user_id) REFERENCES public.users (id)
);