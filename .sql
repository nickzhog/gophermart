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

